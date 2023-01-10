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

func (c *Container) GetLegacyContainerName() (name string) {
	var cn strings.Builder
	switch c.containerType {
	case APT:
		cn.WriteString("apx_managed")
	case AUR:
		cn.WriteString("apx_managed_aur")
	case DNF:
		cn.WriteString("apx_managed_dnf")
	case APK:
		cn.WriteString("apx_managed_apk")
	default:
		log.Fatal(fmt.Errorf("unspecified container type"))
	}
	if len(c.customName) > 0 {
		cn.WriteString("_")
		cn.WriteString(strings.Replace(c.customName, " ", "", -1))
	}
	return cn.String()
}
