package core

import (
	"fmt"

	"github.com/vanilla-os/apx/settings"
)

type ContainerLabels struct {
	Managed    bool
	Distro     string
	PkgManager settings.PackageManager
	Userid     int
}

func (l *ContainerLabels) ToArguments() []string {
	return []string{
		"--label=\"apx.managed=" + fmt.Sprint(l.Managed) + "\"",
		"--label=\"apx.distro=" + l.Distro + "\"",
		"--label=\"apx.pkgmanager=" + string(l.PkgManager) + "\"",
		"--label=\"apx.userid=" + fmt.Sprint(l.Userid) + "\"",
	}
}
