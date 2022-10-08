package core

import "github.com/vanilla-os/apx/settings"

func GetPkgManager(sys bool) []string {
	if sys {
		bin := settings.Cnf.PkgManager.Lock
		sudo := settings.Cnf.PkgManager.Sudo

		if sudo {
			return []string{"sudo", bin}
		}
		return []string{bin}
	}
	bin := settings.Cnf.PkgManager.Bin
	return []string{bin}
}
