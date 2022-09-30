package core

import (
	"fmt"
	"os"
)

var apxDir = "/etc/apx"

func init() {
	if !RootCheck(false) {
		return
	}
	if _, err := os.Stat(apxDir); os.IsNotExist(err) {
		if err := os.Mkdir(apxDir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func RootCheck(display bool) bool {
	if os.Geteuid() != 0 {
		if display {
			fmt.Println("You must be root to run this command")
		}
		return false
	}
	return true
}
