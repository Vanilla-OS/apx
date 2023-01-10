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

func GetLegacyContainerNameAndDistro(pkgManager string) (string, string) {
	switch pkgManager {
	case "apt":
		return "apx_managed", "ubuntu"
	case "aur":
		return "apx_managed_aur", "archlinux"
	case "dnf":
		return "apx_managed_dnf", "fedora"
	case "apk":
		return "apx_managed_apk", "alpine"
	}
	log.Fatal(fmt.Errorf("unspecified container type"))
	return "", ""
}
