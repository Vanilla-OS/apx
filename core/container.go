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
	"log"
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

func GetHostImage() (img string, err error) {
	if settings.Cnf.PkgManager.Image != "" {
		return settings.Cnf.PkgManager.Image, nil
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

func RunContainer(args []string) error {
	log.Default().Printf("Running %v\n", args)

	// distrobox enter "$CONTAINER_NAME" -- "$@"
	cmd := exec.Command("distrobox", "enter",
		settings.Cnf.Container.Name, "--")
	cmd.Args = append(cmd.Args, args...)

	_, err := cmd.Output()
	return err
}

func CreateContainer() error {
	log.Default().Printf("Initializing container\n")

	host_image, err := GetHostImage()
	if err != nil {
		return err
	}

	// distrobox create --name "$CONTAINER_NAME" --image "$(host_image)" --yes
	// --additional-flags --label=manager=apx
	cmd := exec.Command("distrobox", "create", "--name", settings.Cnf.Container.Name,
		"--image", host_image, "--yes", "--no-entry", "--additional-flags", "--label=manager=apx")

	_, err = cmd.Output()
	return err
}

func StopContainer() error {
	log.Default().Printf("Stopping container\n")

	// distrobox stop "$CONTAINER_NAME" --yes
	cmd := exec.Command("distrobox", "stop",
		settings.Cnf.Container.Name, "--yes")

	_, err := cmd.Output()
	return err
}
