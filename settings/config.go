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
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ContainerName string `json:"containername"`
	Image         string `json:"image"`
	PkgManager    string `json:"pkgmanager"`
}

var Cnf *Config

func init() {
	viper.AddConfigPath("/etc/apx/")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		image, pkgmanager, err := GetHostInfo()

		if err != nil {
			panic("Unsupported setup detected, set a default distro and package manager in the config file")
		}

		Cnf = &Config{ContainerName: "apx_managed", Image: image, PkgManager: pkgmanager}
	} else {
		err = viper.Unmarshal(&Cnf)

		if err != nil {
			panic("Config error!\n" + err.Error())
		}

		if Cnf.ContainerName == "" || Cnf.Image == "" || Cnf.PkgManager == "" {
			panic("Config error!\ncontainer_name, image and pkgmanager have to be set")
		}
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

func GetHostInfo() (img string, pkgmanager string, err error) {
	distro_raw, err := exec.Command("lsb_release", "-is").Output()
	if err != nil {
		return "", "", err
	}
	distro := strings.ToLower(strings.Trim(string(distro_raw), "\r\n"))

	release_raw, err := exec.Command("lsb_release", "-rs").Output()
	if err != nil {
		return "", "", err
	}
	release := strings.ToLower(strings.Trim(string(release_raw), "\r\n"))

	img = fmt.Sprintf("%v:%v", distro, release)

	pkgmanager = ""
	switch distro {
	case "ubuntu":
		pkgmanager = "apt"
	case "archlinux":
		pkgmanager = "yay"
	case "fedora":
		pkgmanager = "dnf"
	case "alpine":
		pkgmanager = "apk"
	default:
		return "", "", fmt.Errorf("Unknown package manager for this distro")
	}

	return img, pkgmanager, nil
}
