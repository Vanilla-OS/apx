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

func GetPkgCommand(sys bool, container string, command string) []string {
	if sys {
		container = "default"
	}
	switch container {
	case "aur":
		return GetAurPkgCommand(command)
	case "default":
		return GetDefaultPkgCommand(command)
	default:
		return nil
	}
}

func GetDefaultPkgCommand(command string) []string {
	switch command {
	case "autoremove":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdAutoremove,
		}
	case "clean":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdClean,
		}
	case "install":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdInstall,
		}
	case "list":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdList,
		}
	case "purge":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdPurge,
		}
	case "remove":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdRemove,
		}
	case "search":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdSearch,
		}
	case "show":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdShow,
		}
	case "update":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdUpdate,
		}
	case "upgrade":
		return []string{
			settings.Cnf.PkgManager.Bin,
			settings.Cnf.PkgManager.CmdUpgrade,
		}
	default:
		return nil
	}
}

func GetAurPkgCommand(command string) []string {
	bin := "yay"

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
					rm -rf ~/.local/src/yay &&
					mkdir -p ~/.local/src/yay &&
					git clone https://aur.archlinux.org/yay.git ~/.local/src/yay &&
					cd ~/.local/src/yay && yes "" | /usr/sbin/makepkg -is &&
					cd -- ~ &&
					rm -rf ~/.local/src/yay"`,
		}
	default:
		return nil
	}
}
