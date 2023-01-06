package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	err := CheckContainerTools()
	if err != nil {
		fmt.Println(`One or more core components are not available. 
Please refer to our documentation at https://documentation.vanillaos.org/`)
		log.Fatal(err)
	}
}

func CheckContainerTools() error {
	_, distrobox := os.Stat("/usr/lib/apx/distrobox")
	docker := exec.Command("which", "docker")
	podman := exec.Command("which", "podman")

	if distrobox != nil {
		return errors.New(`distrobox is not installed`)
	}

	if err := docker.Run(); err != nil {
		if err := podman.Run(); err != nil {
			return errors.New(`container engine (docker or podman) not found`)
		}
	}

	return nil
}

func IsVM() bool {
	_, err := exec.Command("systemd-detect-virt").Output()
	return err == nil
}

func ExitIfVM() {
	if IsVM() {
		log.Default().Printf("Apx does not work inside a VM.")
		os.Exit(1)
	}
}

func IsOverlayTypeFS() bool {
	out, err := exec.Command("df", "-T", "/").Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), "overlay")
}

func ExitIfOverlayTypeFS() {
	if IsOverlayTypeFS() {
		log.Default().Printf("Apx does not work with overlay type filesystem.")
		os.Exit(1)
	}
}
