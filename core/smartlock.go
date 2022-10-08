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
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

func PkgManagerSmartLock() {
	_, err := exec.LookPath(settings.Cnf.PkgManager.Bin)
	if err != nil {
		panic("Config error! Package manager setting is not configured correctly!")
	}

	if settings.Cnf.PkgManager.Lock == settings.Cnf.PkgManager.Bin {
		panic("Bin/Lock collision!")
	}

	// here we skip the whole process if the package manager is already locked
	data, err := ioutil.ReadFile(settings.Cnf.PkgManager.Bin)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, l := range lines {
			if l == "# Apx::SmartLock" {
				fmt.Println("Apx::SmartLock is already active on this machine!")
				return
			}
		}
	}

	err = performSmartLock()
	if err != nil {
		panic("Something went wrong: " + err.Error())
	}
}

func performSmartLock() error {
	_, err := AlmostRun("sudo", "mv", settings.Cnf.PkgManager.Bin, settings.Cnf.PkgManager.Lock)
	if err != nil {
		return errors.New("can't lock the original binary")
	}

	// here we replace the original binary with a shell script which points to
	// apx so other programs won't break when they call it
	_, err = AlmostRun("sudo", "touch", settings.Cnf.PkgManager.Bin)
	if err != nil {
		return errors.New("can't recreate the locked binary")
	}

	_, err = AlmostRun("sudo", "chmod", "755", settings.Cnf.PkgManager.Bin)
	if err != nil {
		return errors.New("can't change permissions on " + settings.Cnf.PkgManager.Bin)
	}

	_, err = AlmostRun("sudo", "--",
		"sh", "-c", "echo '#!/bin/sh\\n# Apx::SmartLock\\napx --sys $@' >> "+settings.Cnf.PkgManager.Bin,
		"&&", "exit", "0")
	if err != nil {
		return errors.New("can't write the script binary file " + settings.Cnf.PkgManager.Bin)
	}

	return nil
}
