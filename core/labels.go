package core

import (
	"fmt"
	"os"

	"github.com/vanilla-os/apx/settings"
)

type ContainerLabels struct {
	Managed    bool
	Distro     string
	PkgManager settings.PackageManager
	Userid     int
	CustomName string
}

func CreateLabels(distro string, pkgManager settings.PackageManager, name string) *ContainerLabels {
	return &ContainerLabels{
		Managed:    true,
		Distro:     distro,
		PkgManager: pkgManager,
		Userid:     os.Geteuid(),
		CustomName: name,
	}
}

func (l *ContainerLabels) ToArguments() []string {
	return []string{
		"--label=apx.managed=" + fmt.Sprint(l.Managed),
		"--label=apx.distro=" + l.Distro,
		"--label=apx.pkgmanager=" + string(l.PkgManager),
		"--label=apx.userid=" + fmt.Sprint(l.Userid),
		"--label=apx.customname=" + l.CustomName,
	}
}
