package core

import "github.com/vanilla-os/apx/settings"

func GetPkgManager(sys bool) []string {
	sudo := settings.Cnf.PkgManager.Sudo

	if sys {
		bin := settings.Cnf.PkgManager.Lock

		if sudo {
			return []string{"sudo", bin}
		}
		return []string{bin}
	}

	bin := settings.Cnf.PkgManager.Bin

	if sudo {
		return []string{"sudo", bin}
	}
	return []string{bin}
}
