package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
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

	"github.com/spf13/viper"
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
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// dev paths
	viper.AddConfigPath("config/")

	// tests paths
	viper.AddConfigPath("../config/")

	// user paths
	viper.AddConfigPath(filepath.Join(userConfigDir, "apx/"))

	// prod paths
	viper.AddConfigPath("/etc/apx/")
	viper.AddConfigPath("/usr/share/apx/")

	viper.SetConfigName("apx")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Unable to read config file: \n\t%s\n", err)
		os.Exit(1)
	}

	// if viper.ConfigFileUsed() != "/etc/apx/apx.json" || viper.ConfigFileUsed() != "/usr/share/apx/apx.json" {
	// 	fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	// }

	distroboxPath := viper.GetString("distroboxPath")

	_, err = os.Stat(distroboxPath)
	if err != nil {
		if os.IsNotExist(err) {
			path, err := exec.LookPath("distrobox")
			if err != nil {
				fmt.Printf("Unable to find distrobox in PATH.\n")
			} else {
				distroboxPath = path
			}
		}
	}

	Cnf := NewApxConfig(
		viper.GetString("apxPath"),
		distroboxPath,
		viper.GetString("storageDriver"),
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
