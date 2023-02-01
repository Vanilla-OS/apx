package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DistroId      string `json:"distroid"`
	ContainerName string `json:"containername"`
	Image         string `json:"image"`
	PkgManager    string `json:"pkgmanager"`
	DistroboxPath string `json:"distroboxpath"`
}

var Cnf *Config

func init() {
	viper.AddConfigPath("/etc/apx/")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic("Failed to read the config file")
	}

	err = viper.Unmarshal(&Cnf)

	if err != nil {
		log.Fatalf("Config error: %v\n" + err.Error())
	}

	if Cnf.DistroboxPath == "" {
		panic("Config error: distroboxpath needs to be set")
	}
}

func GetDistroFromSettings() DistroInfo {
	if Cnf.DistroId == "" || Cnf.Image == "" || Cnf.PkgManager == "" {
		distro, err := GetHostInfo()

		if err != nil {
			log.Fatal("Unsupported setup detected, set a default distro and package manager in the config file")
		}

		Cnf = &Config{ContainerName: "apx_managed", Image: distro.Image, PkgManager: string(distro.Pkgmanager), DistroboxPath: Cnf.DistroboxPath}
		return distro
	} else {
		if Cnf.ContainerName == "" {
			log.Fatal("Config error!\ncontainername has to be set")
		}

		var supportedPackageManager = false
		switch Cnf.PkgManager {
		case Apt, Yay, Dnf, Apk:
			supportedPackageManager = true
		}

		if !supportedPackageManager {
			log.Fatal("Config error!\nInvalid package manager")
		}

		return DistroInfo{Id: Cnf.DistroId, Image: Cnf.Image, Pkgmanager: PackageManager(Cnf.PkgManager), ContainerName: Cnf.ContainerName}
	}
}

func GetHostInfo() (distro DistroInfo, err error) {
	file, err := os.Open("/etc/os-release")

	if err != nil {
		return DistroInfo{}, fmt.Errorf("Failed to open /etc/os-release: " + err.Error())
	}

	var distro_id string
	var distro_version string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_content := scanner.Text()
		if strings.HasPrefix(line_content, "ID=") {
			distro_id = strings.Trim(line_content, "ID=")
		} else if strings.HasPrefix(line_content, "VERSION_ID=") {
			distro_version = strings.Trim(line_content, "VERSION_ID=")
		}
	}

	if err := scanner.Err(); err != nil {
		return DistroInfo{}, fmt.Errorf("Failure while reading /etc/os-release: " + err.Error())
	}

	defer file.Close()

	if len(distro_id) == 0 {
		return DistroInfo{}, fmt.Errorf("Failed to read distro type")
	}
	if len(distro_version) == 0 {
		return DistroInfo{}, fmt.Errorf("Failed to read distro version")
	}

	switch distro_id {
	case "ubuntu":
		return DistroUbuntu, nil
	case "arch":
		return DistroArch, nil
	case "fedora":
		return DistroFedora, nil
	case "alpine":
		return DistroAlpine, nil
	case "opensuse":
		return DistroOpensuse, nil
	case "void":
		return DistroVoid, nil

	default:
		return DistroInfo{}, fmt.Errorf("Unsupported distro detected")
	}
}
