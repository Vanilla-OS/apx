package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"github.com/spf13/viper"
)

type Config struct {
	Container  ContainerConfig  `json:"container"`
	PkgManager PkgManagerConfig `json:"pkg_manager"`
}

type ContainerConfig struct {
	Name string `json:"name"`
}

type PkgManagerConfig struct {
	Image     string `json:"pkg_manager_image,omitempty"`
	Path      string `json:"pkg_manager_path"`
	UpdateCmd string `json:"pkg_manager_update"`
}

var Cnf *Config

func init() {
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic("Config error!")
	}

	err = viper.Unmarshal(&Cnf)

	if err != nil {
		panic("Config error!\n" + err.Error())
	}
}
