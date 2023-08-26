package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (a *Apx) EssentialChecks() error {
	err := a.CheckContainerTools()
	if err != nil {
		fmt.Println(`One or more core components are not available. 
Please refer to our documentation at https://documentation.vanillaos.org/`)
		return err
	}

	err = a.CheckAndCreateUserStacksDirectory()
	if err != nil {
		fmt.Println(`Failed to create stacks directory.`)
		return err
	}

	err = a.CheckAndCreateApxStorageDirectory()
	if err != nil {
		fmt.Println(`Failed to create apx storage directory.`)
		return err
	}

	err = a.CheckAndCreateApxUserPkgManagersDirectory()
	if err != nil {
		fmt.Println(`Failed to create apx user pkg managers directory.`)
		return err
	}

	return nil
}

func (a *Apx) CheckContainerTools() error {
	_, err := os.Stat(a.Cnf.DistroboxPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("distrobox is not installed")
		}
		return err
	}

	if _, err := exec.LookPath("docker"); err != nil {
		if _, err := exec.LookPath("podman"); err != nil {
			return errors.New("container engine (docker or podman) not found")
		}
	}

	return nil
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

func (a *Apx) CheckAndCreateUserStacksDirectory() error {
	_, err := os.Stat(a.Cnf.UserStacksPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(a.Cnf.UserStacksPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create stacks directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check stacks directory: %w", err)
		}
	}

	return nil
}

func (a *Apx) CheckAndCreateApxStorageDirectory() error {
	_, err := os.Stat(a.Cnf.ApxStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(a.Cnf.ApxStoragePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create apx storage directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check apx storage directory: %w", err)
		}
	}

	return nil
}

func (a *Apx) CheckAndCreateApxUserPkgManagersDirectory() error {
	_, err := os.Stat(a.Cnf.UserPkgManagersPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(a.Cnf.UserPkgManagersPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create apx user pkg managers directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check apx user pkg managers directory: %w", err)
		}
	}

	return nil
}

func hasNvidiaGPU() bool {
	_, err := os.Stat("/dev/nvidia0")
	return err == nil
}
