package core

import (
	"fmt"
	"os"

	"github.com/vanilla-os/apx/settings"
)

type ContainerLabels struct {
	Id           string
	Managed      bool
	Distro       string
	PkgManager   settings.PackageManager
	Userid       int
	CustomName   string
	MigratedFrom string
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
		"--label=apx.id=" + l.Id,
		"--label=apx.managed=" + fmt.Sprint(l.Managed),
		"--label=apx.distro=" + l.Distro,
		"--label=apx.pkgmanager=" + string(l.PkgManager),
		"--label=apx.userid=" + fmt.Sprint(l.Userid),
		"--label=apx.customname=" + l.CustomName,
		"--label=apx.migratedfrom=" + l.MigratedFrom,
	}
}

func (l *ContainerLabels) ToFilters() []string {
	return []string{
		"--filter", "label=apx.managed=" + fmt.Sprint(l.Managed),
		"--filter", "label=apx.distro=" + l.Distro,
		"--filter", "label=apx.pkgmanager=" + string(l.PkgManager),
		"--filter", "label=apx.userid=" + fmt.Sprint(l.Userid),
		"--filter", "label=apx.customname=" + l.CustomName,
	}
}
