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
	if RootCheck(false) {
		fmt.Println("Do not run Apx as root!")
		os.Exit(1)
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

func AskConfirmation(s string) bool {
	var response string
	fmt.Print(s + " [y/N]: ")
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		return true
	}
	return false
}

func CopyToUserTemp(path string) (string, error) {
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

	pathContents, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer pathContents.Close()

	newPathContents, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer newPathContents.Close()

	_, err = newPathContents.ReadFrom(pathContents)
	if err != nil {
		return "", err
	}

	return newPath, nil
}
