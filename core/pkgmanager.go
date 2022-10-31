package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

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
			"bash", "-c", "cd ~/.local/src/yay  && tar -xvf yay.tar.gz && cd yay_*_x86_64* && sudo cp yay /usr/bin",
		}
	default:
		return nil
	}
}

func GetLatestYay() string {
	url := "https://api.github.com/repos/Jguer/yay/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	assets_url := result["assets_url"].(string)
	resp, err = http.Get(assets_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var assets []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&assets)

	for _, asset := range assets {
		if strings.Contains(asset["name"].(string), "x86_64") {
			return asset["browser_download_url"].(string)
		}
	}
	panic("Failed to install yay. No asset found.")
}

func DownloadYay() {
	url := GetLatestYay()
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	yay_dir := fmt.Sprintf("%v/.local/src/yay", home)
	if _, err := os.Stat(yay_dir); os.IsNotExist(err) {
		err = os.MkdirAll(yay_dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	yay_file := fmt.Sprintf("%v/yay.tar.gz", yay_dir)
	out, err := os.Create(yay_file)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
