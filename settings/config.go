package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"os"
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

var Cnf *Config

func init() {
	// dev paths
	viper.AddConfigPath("config/")

	// tests paths
	viper.AddConfigPath("../config/")

	// prod paths
	viper.AddConfigPath("/etc/apx/")
	viper.AddConfigPath("/usr/share/apx/")

	viper.SetConfigName("apx")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		return
	}

	// if viper.ConfigFileUsed() != "/etc/apx/apx.json" || viper.ConfigFileUsed() != "/usr/share/apx/apx.json" {
	// 	fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	// }

	Cnf = &Config{
		// Common
		ApxPath:       viper.GetString("apxPath"),
		DistroboxPath: viper.GetString("distroboxPath"),
		StorageDriver: viper.GetString("storageDriver"),

		// Virtual
		ApxStoragePath:      "",
		UserApxPath:         "",
		StacksPath:          "",
		UserStacksPath:      "",
		PkgManagersPath:     "",
		UserPkgManagersPath: "",
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	Cnf.UserApxPath = filepath.Join(userHome, ".local/share/apx")
	Cnf.ApxStoragePath = filepath.Join(Cnf.UserApxPath, "storage")
	Cnf.StacksPath = filepath.Join(Cnf.ApxPath, "stacks")
	Cnf.UserStacksPath = filepath.Join(Cnf.UserApxPath, "stacks")
	Cnf.PkgManagersPath = filepath.Join(Cnf.ApxPath, "package-managers")
	Cnf.UserPkgManagersPath = filepath.Join(Cnf.UserApxPath, "package-managers")
}
