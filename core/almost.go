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
