package core

import (
	"fmt"
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

func AlmostRun(output bool, command ...string) (string, error) {
	var run *exec.Cmd
	var immutable_need_restore bool

	// if the system is immutable, we need to enter rw mode and
	// then exit rw mode after the command is executed
	if ImmutabilityStatus() {
		immutable_need_restore = true
		err := exec.Command("sudo", "almost", "enter", "rw").Run()
		if err != nil {
			panic(err)
		}
	}

	run = exec.Command(command[0], command[1:]...)

	run.Env = os.Environ()
	if !output {
		if run.Stdout == nil {
			run.Stdout = os.Stdout
		}
		if run.Stderr == nil {
			run.Stderr = os.Stderr
		}
		if run.Stdin == nil {
			run.Stdin = os.Stdin
		}
		err := run.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

	out, err := run.Output()

	// restore the system to immutable mode due to previous check
	if immutable_need_restore {
		err := exec.Command("sudo", "almost", "enter", "ro").Run()
		if err != nil {
			panic(err)
		}
	}

	return string(out), err
}

func AlmostEnterRw() (string, error) {
	run := exec.Command("sudo", "almost", "enter", "rw")
	out, err := run.Output()

	return string(out), err
}

func AlmostEnterRo() (string, error) {
	run := exec.Command("sudo", "almost", "enter", "ro")
	out, err := run.Output()

	return string(out), err
}
