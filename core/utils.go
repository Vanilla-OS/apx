package core

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var apxDir = "/etc/apx"
var ProcessPath string

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

	ex, err := os.Executable()
	if err != nil {
		panic("Can't get process path!\n" + err.Error())
	}
	ProcessPath = filepath.Dir(ex)
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

func AskConfirmation(s string) bool {
	var response string
	fmt.Print(s + " [y/N]: ")
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		return true
	}
	return false
}

func MoveToUserTemp(path string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(user.HomeDir, ".cache", "apx")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return "", err
		}
	}

	fileName := filepath.Base(path)
	newPath := filepath.Join(cacheDir, fileName)
	if err := os.Rename(path, newPath); err != nil {
		return "", err
	}

	return newPath, nil
}
