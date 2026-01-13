package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/conf"
)

type Config struct {
	// Paths
	ApxPath       string `json:"apxPath"`
	DistroboxPath string `json:"distroboxPath"`
	StorageDriver string `json:"storageDriver"`

	// Virtual
	UserApxPath         string
	ApxStoragePath      string
	StacksPath          string
	UserStacksPath      string
	PkgManagersPath     string
	UserPkgManagersPath string
}

func GetApxDefaultConfig() (*Config, error) {
	config, err := conf.NewBuilder[Config]("apx").
		WithType("json").
		Build()
	if err != nil {
		fmt.Printf("Unable to read config file: \n\t%s\n", err)
		os.Exit(1)
	}

	distroboxPath := config.DistroboxPath

	err = TestFile(distroboxPath)
	if err != nil {
		if os.IsNotExist(err) {
			path, _ := LookPath("distrobox")
			distroboxPath = path
		}
	}

	Cnf := NewApxConfig(
		config.ApxPath,
		distroboxPath,
		config.StorageDriver,
	)
	return Cnf, nil
}

func NewApxConfig(apxPath, distroboxPath, storageDriver string) *Config {
	userDataDir, err := UserDataDir()
	if err != nil {
		panic(err)
	}

	Cnf := &Config{
		// Common
		ApxPath:       apxPath,
		DistroboxPath: distroboxPath,
		StorageDriver: storageDriver,

		// Virtual
		ApxStoragePath:      "",
		UserApxPath:         "",
		StacksPath:          "",
		UserStacksPath:      "",
		PkgManagersPath:     "",
		UserPkgManagersPath: "",
	}

	Cnf.UserApxPath = filepath.Join(userDataDir, "apx")
	Cnf.ApxStoragePath = filepath.Join(Cnf.UserApxPath, "storage")
	Cnf.StacksPath = filepath.Join(Cnf.ApxPath, "stacks")
	Cnf.UserStacksPath = filepath.Join(Cnf.UserApxPath, "stacks")
	Cnf.PkgManagersPath = filepath.Join(Cnf.ApxPath, "package-managers")
	Cnf.UserPkgManagersPath = filepath.Join(Cnf.UserApxPath, "package-managers")

	return Cnf
}

func UserDataDir() (string, error) {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir != "" {
		return dir, nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHome, ".local", "share"), nil
}

func IsFlatpak() bool {
	id := os.Getenv("FLATPAK_ID")
	return id != ""
}

func TestFile(path string) error {
	if IsFlatpak() {
		cmd := exec.Command("flatpak-spawn", "--host", "test", "-f", path)
		if err := cmd.Run(); err != nil {
			return os.ErrNotExist
		}
		return nil
	}
	_, err := os.Stat(path)
	return err
}

func LookPath(file string) (string, error) {
	if IsFlatpak() {
		cmd := exec.Command("flatpak-spawn", "--host", "which", file)
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}
	return exec.LookPath(file)
}
