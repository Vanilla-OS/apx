package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

func ContainerManager() string {
	docker := exec.Command("which", "docker")
	podman := exec.Command("which", "podman")

	if err := docker.Run(); err == nil {
		return "docker"
	} else if err := podman.Run(); err == nil {
		return "podman"
	}

	panic("No container engine found. Please install Podman or Docker.")
}

func GetHostImage() (img string, err error) {
	if settings.Cnf.Container.Image != "" {
		return settings.Cnf.Container.Image, nil
	}

	distro_raw, err := exec.Command("lsb_release", "-is").Output()
	if err != nil {
		return "", err
	}
	distro := strings.ToLower(strings.Trim(string(distro_raw), "\r\n"))

	release_raw, err := exec.Command("lsb_release", "-rs").Output()
	if err != nil {
		return "", err
	}
	release := strings.ToLower(strings.Trim(string(release_raw), "\r\n"))

	return fmt.Sprintf("%v:%v", distro, release), nil
}

func GetContainerImage(container string) (image string, err error) {
	switch container {
	case "default":
		return GetHostImage()
	case "aur":
		return "docker.io/library/archlinux", nil
	default:
		image = ""
		err = errors.New("Can't retrieve image for unknown container")
	}
	return image, err
}

func GetContainerName(container string) (name string) {
	switch container {
	case "default":
		name := "apx_managed"
		return name
	case "aur":
		name := "apx_managed_aur"
		return name
	default:
		panic("Unknown container not supported")
	}
}

func GetDistroboxVersion() (version string, err error) {
	output, err := exec.Command("/usr/lib/apx/distrobox", "version").Output()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(string(output), "distrobox: ")
	if len(splitted) != 2 {
		return "", errors.New("Can't retrieve distrobox version")
	}

	return splitted[1], nil
}

func RunContainer(container string, args ...string) error {
	if !ContainerExists(container) {
		log.Default().Printf("Managed container does not exist.\nTry: apx init")
		return errors.New("Managed container does not exist")
	}

	container_name := GetContainerName(container)

	cmd := exec.Command("/usr/lib/apx/distrobox", "enter", container_name, "--")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func EnterContainer(container string) error {
	if !ContainerExists(container) {
		log.Default().Printf("Managed container does not exist.\nTry: apx init")
		return errors.New("Managed container does not exist")
	}

	container_name := GetContainerName(container)

	cmd := exec.Command("/usr/lib/apx/distrobox", "enter", container_name)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return nil
}

func CreateContainer(container string) error {
	log.Default().Printf("Initializing container\n")

	container_image, err := GetContainerImage(container)
	container_name := GetContainerName(container)
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/lib/apx/distrobox", "create",
		"--name", container_name,
		"--image", container_image,
		"--yes",
		"--no-entry",
		"--additional-flags",
		"--label=manager=apx",
		"--yes")
	cmd.Env = os.Environ()
	// mute command output
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	//err = cmd.Run()

	_, err = cmd.Output()
	if err != nil {
		log.Panic(err)
	}

	if container == "aur" {
		RunContainer(container, GetAurPkgCommand("install-yay")...)
	}

	return err
}

func StopContainer(container string) error {
	log.Default().Printf("Stopping container\n")

	container_name := GetContainerName(container)

	cmd := exec.Command("/usr/lib/apx/distrobox", "stop", container_name, "--yes")
	_, err := cmd.Output()

	return err
}

func RemoveContainer(container string) error {
	log.Default().Printf("Removing container\n")

	container_name := GetContainerName(container)

	if !ContainerExists(container) {
		return nil
	}

	err := StopContainer(container)
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/lib/apx/distrobox", "rm", container_name, "--yes")
	_, err = cmd.Output()

	return err
}

func ExportDesktopEntry(container string, program string) error {
	log.Default().Printf("Exporting desktop entry %v\n", program)

	err := RunContainer(
		container,
		"distrobox-export", "--app", program,
		"--export-label", "◆", ">", "/dev/null")
	if err != nil {
		fmt.Printf("No desktop entry found for %v, nothing to export.\n", program)
		return err
	}

	log.Default().Printf("Desktop entry exported for %v\n", program)
	return nil
}

func RemoveDesktopEntry(container string, program string) error {
	log.Default().Printf("Removing desktop entry %v\n", program)

	container_name := GetContainerName(container)

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(home_dir + "/.local/share/applications")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()),
			strings.ToLower(container_name+"-"+program)) {
			log.Default().Printf("Removing desktop entry %v\n", file.Name())
			err := os.Remove(home_dir + "/.local/share/applications/" + file.Name())
			if err != nil {
				return err
			}
		}
	}
	log.Default().Printf("Desktop entry %v not found.\n", program)
	return nil
}

func ContainerExists(container string) bool {
	container_name := GetContainerName(container)
	manager := ContainerManager()

	cmd := exec.Command(manager, "ps", "-a", "-q", "-f", "name="+container_name)
	output, _ := cmd.Output()

	return len(output) > 0
}
