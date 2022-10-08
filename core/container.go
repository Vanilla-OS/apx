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

func GetDistroboxVersion() (version string, err error) {
	output, err := exec.Command("distrobox", "version").Output()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(string(output), "distrobox: ")
	if len(splitted) != 2 {
		return "", errors.New("Can't retrieve distrobox version")
	}

	return splitted[1], nil
}

func RunContainer(args ...string) error {
	log.Default().Printf("Running %v\n", args)

	cmd := exec.Command("distrobox", "enter",
		settings.Cnf.Container.Name, "--")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func EnterContainer() error {
	cmd := exec.Command("distrobox", "enter", settings.Cnf.Container.Name)
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

func CreateContainer() error {
	log.Default().Printf("Initializing container\n")

	host_image, err := GetHostImage()
	if err != nil {
		return err
	}

	cmd := exec.Command("distrobox", "create", "--name", settings.Cnf.Container.Name,
		"--image", host_image, "--yes", "--no-entry", "--additional-flags", "--label=manager=apx")
	_, err = cmd.Output()

	return err
}

func StopContainer() error {
	log.Default().Printf("Stopping container\n")

	cmd := exec.Command("distrobox", "stop",
		settings.Cnf.Container.Name, "--yes")

	_, err := cmd.Output()
	return err
}

func RemoveContainer() error {
	log.Default().Printf("Removing container\n")

	err := StopContainer()
	if err != nil {
		return err
	}

	cmd := exec.Command("distrobox", "rm",
		settings.Cnf.Container.Name, "--yes")

	_, err = cmd.Output()
	return err
}

func ExportDesktopEntry(program string) error {
	log.Default().Printf("Exporting desktop entry %v\n", program)

	err := RunContainer(
		"distrobox=export", "--app", program,
		"--export-label", "◆", ">", "/dev/null")
	if err != nil {
		fmt.Println("No desktop entry found for %w, nothing to export.\n", program)
		return err
	}

	return nil
}

func RemoveDesktopEntry(program string) error {
	log.Default().Printf("Removing desktop entry %v\n", program)

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
			strings.ToLower(settings.Cnf.Container.Name+"-"+program)) {
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

func ContainerExists() bool {
	manager := ContainerManager()
	cmd := exec.Command(manager, "ps", "-a", "-q", "-f", "name="+settings.Cnf.Container.Name)
	output, _ := cmd.Output()
	return len(output) > 0
}
