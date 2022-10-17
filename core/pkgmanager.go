package core

import "github.com/vanilla-os/apx/settings"

func GetPkgManager(sys bool) []string {
	sudo := settings.Cnf.PkgManager.Sudo

	if sys {
		bin := settings.Cnf.PkgManager.Bin

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

func GetPkgCommand(sys bool, container string, command string) []string {
	if sys {
		container = "default"
	}
	switch container {
	case "aur":
		return GetAurPkgCommand(command)
	case "default":
		return GetDefaultPkgCommand(sys, command)
	default:
		return nil
	}
}

func GetDefaultPkgCommand(sys bool, command string) []string {
	res := GetPkgManager(sys)
	switch command {
	case "autoremove":
		res = append(res, settings.Cnf.PkgManager.CmdAutoremove)
		break
	case "clean":
		res = append(res, settings.Cnf.PkgManager.CmdClean)
		break
	case "install":
		res = append(res, settings.Cnf.PkgManager.CmdInstall)
		break
	case "list":
		res = append(res, settings.Cnf.PkgManager.CmdList)
		break
	case "purge":
		res = append(res, settings.Cnf.PkgManager.CmdPurge)
		break
	case "remove":
		res = append(res, settings.Cnf.PkgManager.CmdRemove)
		break
	case "search":
		res = append(res, settings.Cnf.PkgManager.CmdSearch)
		break
	case "show":
		res = append(res, settings.Cnf.PkgManager.CmdShow)
		break
	case "update":
		res = append(res, settings.Cnf.PkgManager.CmdUpdate)
		break
	case "upgrade":
		res = append(res, settings.Cnf.PkgManager.CmdUpgrade)
		break
	default:
		return nil
	}
	return res

}

func GetAurPkgCommand(command string) []string {
	bin := "ame"

	switch command {
	// defaults
	case "autoremove":
		return []string{"echo", "Not implemented yet! "}
	case "clean":
		return []string{bin, "-Yc"}
	case "install":
		return []string{bin, "-S"}
	case "list":
		return []string{bin, "-Qm"}
	case "purge":
		return []string{bin, "-R"}
	case "remove":
		return []string{bin, "-Rs"}
	case "search":
		return []string{bin, "-Ss"}
	case "show":
		return []string{bin, "-Si"}
	case "update":
		return []string{bin, "-Syu"}
	case "upgrade":
		return []string{bin, "-Su"}

	// specials
	case "install-yay":
		return []string{
			"sh -c",
			`"sudo pacman -S --needed --noconfirm git base-devel &&
					rm -rf ~/.local/src/ame &&
					mkdir -p ~/.local/src/ame &&
					git clone https://aur.archlinux.org/ame.git ~/.local/src/ame &&
					cd ~/.local/src/ame && yes "" | /usr/sbin/makepkg -is --noconfirm &&
					cd -- ~ &&
					rm -rf ~/.local/src/ame"`,
		}
	default:
		return nil
	}
}
