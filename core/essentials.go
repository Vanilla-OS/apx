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

	"github.com/vanilla-os/apx/settings"
)

func init() {
	err := CheckContainerTools()
	if err != nil {
		fmt.Println(`One or more core components are not available. 
Please refer to our documentation at https://documentation.vanillaos.org/`)
		log.Fatal(err)
	}

	err = CheckAndCreateUserStacksDirectory()
	if err != nil {
		log.Fatal(err)
	}

	err = CheckAndCreateApxStorageDirectory()
	if err != nil {
		log.Fatal(err)
	}

	err = CheckAndCreateApxUserPkgManagersDirectory()
	if err != nil {
		log.Fatal(err)
	}

}

func CheckContainerTools() error {
	_, err := os.Stat(settings.Cnf.DistroboxPath)
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

func CheckAndCreateUserStacksDirectory() error {
	_, err := os.Stat(settings.Cnf.UserStacksPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(settings.Cnf.UserStacksPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create stacks directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check stacks directory: %w", err)
		}
	}

	return nil
}

func CheckAndCreateApxStorageDirectory() error {
	_, err := os.Stat(settings.Cnf.ApxStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(settings.Cnf.ApxStoragePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create apx storage directory: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check apx storage directory: %w", err)
		}
	}

	return nil
}

func CheckAndCreateApxUserPkgManagersDirectory() error {
	_, err := os.Stat(settings.Cnf.UserPkgManagersPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(settings.Cnf.UserPkgManagersPath, 0755)
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
