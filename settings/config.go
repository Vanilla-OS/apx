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
		image, pkgmanager, err := GetHostInfo()

		if err != nil {
			log.Fatal("Unsupported setup detected, set a default distro and package manager in the config file")
		}

		Cnf = &Config{ContainerName: "apx_managed", Image: image, PkgManager: pkgmanager, DistroboxPath: "/usr/lib/apx/distrobox"}
	} else {
		err = viper.Unmarshal(&Cnf)

		if err != nil {
			log.Fatalf("Config error: %v\n" + err.Error())
		}

		if Cnf.ContainerName == "" || Cnf.Image == "" || Cnf.PkgManager == "" {
			log.Fatal("Config error!\ncontainer_name, image and pkgmanager have to be set")
		}
	}
}

func GetHostInfo() (img string, pkgmanager string, err error) {
	file, err := os.Open("/etc/os-release")

	if err != nil {
		return "", "", fmt.Errorf("Failed to open /etc/os-release: " + err.Error())
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
		return "", "", fmt.Errorf("Failure while reading /etc/os-release: " + err.Error())
	}

	defer file.Close()

	if len(distro_id) == 0 {
		return "", "", fmt.Errorf("Failed to read distro type")
	}
	if len(distro_version) == 0 {
		return "", "", fmt.Errorf("Failed to read distro version")
	}

	switch distro_id {
	case "ubuntu":
		return "docker.io/library/ubuntu:" + distro_version, "apt", nil
	case "arch":
		return "docker.io/library/archlinux:" + distro_version, "yay", nil
	case "fedora":
		return "docker.io/library/fedora:" + distro_version, "dnf", nil
	case "alpine":
		return "docker.io/library/alpine:" + distro_version, "apk", nil
	case "zypper":
		return "registry.opensuse.org/opensuse/tumbleweed:latest" + distro_version, "zypper", nil
	case "xbps":
		return "ghcr.io/void-linux/void-linux:latest-full-x86_64" + distro_version, "xbps", nil
	default:
		return "", "", fmt.Errorf("Unsupported distro detected")
	}
}
