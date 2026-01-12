package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Pietro di Caprio <pietro@fabricators.ltd>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var ProcessPath string

func RootCheck(display bool) bool {
	if os.Geteuid() != 0 {
		if display {
			fmt.Println("You must be root to run this command")
		}
		return false
	}
	return true
}

func CopyToUserTemp(path string) (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(userCacheDir, "apx")
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

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}

func SelectYamlFile(basePath string, name string) string {
	const (
		YML  string = ".yml"
		YAML string = ".yaml"
	)

	yamlFile := filepath.Join(basePath, fmt.Sprintf("%s%s", name, YAML))
	ymlFile := filepath.Join(basePath, fmt.Sprintf("%s%s", name, YML))

	if _, err := os.Stat(yamlFile); errors.Is(err, os.ErrNotExist) {
		return ymlFile
	}

	return yamlFile
}
