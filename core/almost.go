package core

import (
	"os"
	"os/exec"
	"strings"
)

func ImmutabilityStatus() bool {
	// if "almost" is not available, we assume that the system is mutable
	if _, err := os.Stat("/usr/bin/almost"); os.IsNotExist(err) {
		return false
	}

	// here we check if the system is immutable via "almost check"
	out, err := exec.Command("sudo", "almost", "check").Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(out), "ro") {
		return true
	}

	// if none of the above return true, we assume that the system is mutable
	// at this point will be at the user's risk
	return false
}

func AlmostRun(command ...string) (string, error) {
	var run *exec.Cmd

	// if the system is immutable, we run the command with "almost"
	// to obtain a temporary read-only state
	if ImmutabilityStatus() {
		cmd := []string{"sudo", "almost", "run"}
		cmd = append(cmd, command...)
		run = exec.Command(cmd[0], cmd[1:]...)
	} else {
		run = exec.Command(command[0], command[1:]...)
	}

	run.Env = os.Environ()
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	run.Stdin = os.Stdin
	run.Run()
	out, err := run.Output()
	return string(out), err
}

func AlmostEnterRw() (string, error) {
	var run *exec.Cmd

	run = exec.Command("sudo", "almost", "enter", "rw")
	out, err := run.Output()

	return string(out), err
}

func AlmostEnterRo() (string, error) {
	var run *exec.Cmd

	run = exec.Command("sudo", "almost", "enter", "ro")
	out, err := run.Output()

	return string(out), err
}
