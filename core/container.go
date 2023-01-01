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
	"time"

	"github.com/briandowns/spinner"
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

func GetContainerImage(container string) (image string, err error) {
	switch container {
	case "default":
		return settings.Cnf.Image, nil
	case "apt":
		return "docker.io/library/ubuntu", nil
	case "aur":
		return "docker.io/library/archlinux", nil
	case "dnf":
		return "docker.io/library/fedora", nil
	case "apk":
		return "docker.io/library/alpine", nil
	default:
		image = ""
		err = errors.New("can't retrieve image for unknown container")
	}
	return image, err
}

func GetContainerName(container string) (name string) {
	switch container {
	case "default":
		return settings.Cnf.ContainerName
	case "apt":
		return "apx_managed_apt"
	case "aur":
		return "apx_managed_aur"
	case "dnf":
		return "apx_managed_dnf"
	case "apk":
		return "apx_managed_apk"
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
	ExitIfOverlayTypeFS()

	if !ContainerExists(container) {
		err := CreateContainer(container)
		if err != nil {
			log.Default().Println("Failed to initialize the container. Try manually with `apx init`.")
			return err
		}
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
	ExitIfOverlayTypeFS()

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
		if err.Error() != "exit status 130" {
			return err
		}
	}

	return nil
}

func CreateContainer(container string) error {
	ExitIfOverlayTypeFS()

	if !CheckConnection() {
		log.Default().Println("No internet connection. Please connect to the internet and try again.")
		return errors.New("failed to create container")
	}

	container_image, err := GetContainerImage(container)
	if err != nil {
		return err
	}

	container_name := GetContainerName(container)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Creating container..."

	spinner.Start()

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

	spinner.Stop()

	if container == "aur" {
		DownloadYay()
		RunContainer(container, GetAurPkgCommand("install-yay")...)
	}

	log.Default().Println("Container created")

	return err
}

func StopContainer(container string) error {
	ExitIfOverlayTypeFS()

	container_name := GetContainerName(container)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Stopping container..."

	spinner.Start()

	cmd := exec.Command("/usr/lib/apx/distrobox", "stop", container_name, "--yes")
	_, err := cmd.Output()

	spinner.Stop()

	if err != nil {
		log.Panic(err)
	}

	log.Default().Println("Container stopped")

	return err
}

func RemoveContainer(container string) error {
	ExitIfOverlayTypeFS()

	container_name := GetContainerName(container)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Removing container..."

	if !ContainerExists(container) {
		return nil
	}

	err := StopContainer(container)
	if err != nil {
		return err
	}

	spinner.Start()

	cmd := exec.Command("/usr/lib/apx/distrobox", "rm", container_name, "--yes")
	_, err = cmd.Output()

	spinner.Stop()

	log.Default().Println("Container removed")

	return err
}

func ExportDesktopEntry(container string, program string) {
	RunContainer(container, "sh", "-c", "distrobox-export --app "+program+" 2>/dev/null || true")
}

func RemoveDesktopEntry(container string, program string) error {
	container_name := GetContainerName(container)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = fmt.Sprintf("Removing desktop entry: %v\n", program)

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	spinner.Start()

	files, err := ioutil.ReadDir(home_dir + "/.local/share/applications")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()),
			strings.ToLower(container_name+"-"+program)) {
			spinner.Stop()
			err := os.Remove(home_dir + "/.local/share/applications/" + file.Name())
			if err != nil {
				return err
			}
		}
	}

	spinner.Stop()

	log.Default().Printf("Desktop entry %v not found.\n", program)
	return nil
}

func ContainerExists(container string) bool {
	container_name := GetContainerName(container)
	manager := ContainerManager()

	cmd := exec.Command(manager, "ps", "-a", "-q", "-f", "name="+container_name+"$")
	output, _ := cmd.Output()

	// fmt.Println("container_name: ", container_name)
	// fmt.Println("command: ", cmd.String())
	// fmt.Println("output: ", string(output))

	return len(output) > 0
}
