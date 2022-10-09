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
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

func PkgManagerSmartLock() {
	if settings.Cnf.PkgManager.Lock == settings.Cnf.PkgManager.Bin {
		panic("Bin/Lock collision!")
	}

	// here we skip the whole process if the package manager is already locked
	_, err := exec.LookPath(settings.Cnf.PkgManager.Bin)
	bin_exists := err == nil

	if bin_exists {
		data, err := ioutil.ReadFile(settings.Cnf.PkgManager.Bin)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, l := range lines {
				if l == "# Apx::SmartLock" {
					return
				}
			}
		}
	}

	err = performSmartLock(bin_exists)
	if err != nil {
		panic("Something went wrong: " + err.Error())
	}
}

func performSmartLock(bin_exists bool) error {
	if bin_exists {
		_, err := AlmostRun(false, "sudo", "mv", settings.Cnf.PkgManager.Bin, settings.Cnf.PkgManager.Lock)
		if err != nil {
			fmt.Println(err.Error())
			return errors.New("can't lock the original binary")
		}
	} else {
		log.Default().Println("There original package manager binary wasn't found. Assuming this was already locked. Skipping this step.")
	}

	// here we replace the original binary with a shell script which points to
	// apx so other programs won't break when they call it
	_, err := AlmostRun(false, "sudo", "touch", settings.Cnf.PkgManager.Bin)
	if err != nil {
		return errors.New("can't recreate the locked binary")
	}

	_, err = AlmostRun(false, "sudo", "chmod", "755", settings.Cnf.PkgManager.Bin)
	if err != nil {
		return errors.New("can't change permissions on " + settings.Cnf.PkgManager.Bin)
	}

	_, err = AlmostRun(false, "sudo", "--",
		"sh", "-c", "echo '#!/bin/sh\\n# Apx::SmartLock\\napx --sys $@' >> "+settings.Cnf.PkgManager.Bin,
		"&&", "exit", "0")
	if err != nil {
		return errors.New("can't write the script binary file " + settings.Cnf.PkgManager.Bin)
	}

	return nil
}
