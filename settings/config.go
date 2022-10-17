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
	Name   string `json:"name"`
	Image  string `json:"container_image,omitempty"`
	Path   string `json:"container_path"`
	Update string `json:"container_update"`
}

type PkgManagerConfig struct {
	Bin           string `json:"pkgmanager_bin"`
	Sudo          bool   `json:"pkgmanager_sudo"`
	CmdAutoremove string `json:"pkgmanager_cmdAutoremove"`
	CmdClean      string `json:"pkgmanager_cmdClean"`
	CmdInstall    string `json:"pkgmanager_cmdInstall"`
	CmdList       string `json:"pkgmanager_cmdList"`
	CmdPurge      string `json:"pkgmanager_cmdPurge"`
	CmdRemove     string `json:"pkgmanager_cmdRemove"`
	CmdSearch     string `json:"pkgmanager_cmdSearch"`
	CmdShow       string `json:"pkgmanager_cmdShow"`
	CmdUpdate     string `json:"pkgmanager_cmdUpdate"`
	CmdUpgrade    string `json:"pkgmanager_cmdUpgrade"`
}

var Cnf *Config

func init() {
	viper.AddConfigPath("/etc/apx/")
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

	// fmt.Println("==========================")
	// fmt.Println("Config:")
	// fmt.Println("Container name:", Cnf.Container.Name)
	// fmt.Println("Container image:", Cnf.Container.Image)
	// fmt.Println("Container path:", Cnf.Container.Path)
	// fmt.Println("Container update command:", Cnf.Container.Update)
	// fmt.Println("Package manager bin:", Cnf.PkgManager.Bin)
	// fmt.Println("==========================")
}
