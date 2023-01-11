package core

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

var legacy_containers_ids = []string(nil)

func GetLegacyContainersIds() []string {
	if legacy_containers_ids != nil {
		return legacy_containers_ids
	}

	manager := ContainerManager()

	cmd := exec.Command(manager, "ps", "-aqf", "label=manager=apx")
	output, _ := cmd.Output()

	containers_ids := strings.Split(string(output), "\n")

	cmd = exec.Command(manager, "ps", "-aqf", "label=apx.managed="+fmt.Sprint(true))
	output, _ = cmd.Output()

	migrated_containers_ids := strings.Split(string(output), "\n")

	legacy_containers_ids = []string{}
	for _, id := range containers_ids {
		if !slices.Contains(migrated_containers_ids, id) {
			legacy_containers_ids = append(legacy_containers_ids, id)
		}
	}

	return legacy_containers_ids
}

func (c *Container) GetLegacyContainerNameAndDistro(pkgManager string) (string, string) {
	var cn strings.Builder
	var distro string

	switch pkgManager {
	case "apt":
		cn.WriteString("apx_managed")
		distro = "ubuntu"
	case "aur":
		cn.WriteString("apx_managed_aur")
		distro = "archlinux"
	case "dnf":
		cn.WriteString("apx_managed_dnf")
		distro = "fedora"
	case "apk":
		cn.WriteString("apx_managed_apk")
		distro = "alpine"
	default:
		log.Fatal(fmt.Errorf("unspecified container type"))
	}

	if len(c.customName) > 0 {
		cn.WriteString("_")
		cn.WriteString(strings.Replace(c.customName, " ", "", -1))
	}

	return cn.String(), distro
}
