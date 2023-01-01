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

func GetPkgCommand(container string, command string) []string {
	switch container {
	case "apt":
		return GetAptPkgCommand(command)
	case "aur":
		return GetAurPkgCommand(command)
	case "dnf":
		return GetDnfPkgCommand(command)
	case "apk":
		return GetApkPkgCommand(command)
	case "default":
		return GetDefaultPkgCommand(command)
	default:
		return nil
	}
}

func GetDefaultPkgCommand(command string) []string {
	pkgmanager := settings.Cnf.PkgManager

	switch pkgmanager {
	case "apt":
		return GetAptPkgCommand(command)
	case "aur":
		return GetAurPkgCommand(command)
	case "dnf":
		return GetDnfPkgCommand(command)
	case "apk":
		return GetApkPkgCommand(command)
	default:
		return []string{"echo", pkgmanager + " is not implemented yet!"}
	}
}

func GetAptPkgCommand(command string) []string {
	bin := "apt"

	switch command {
	case "autoremove":
		return []string{"sudo", bin, "autoremove"}
	case "clean":
		return []string{"sudo", bin, "clean"}
	case "install":
		return []string{"sudo", bin, "install"}
	case "list":
		return []string{"sudo", bin, "list"}
	case "purge":
		return []string{"sudo", bin, "purge"}
	case "remove":
		return []string{"sudo", bin, "remove"}
	case "search":
		return []string{"sudo", bin, "search"}
	case "show":
		return []string{"sudo", bin, "show"}
	case "update":
		return []string{"sudo", bin, "update"}
	case "upgrade":
		return []string{"sudo", bin, "upgrade"}
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
			"bash", "-c", "cd ~/.local/src/yay  && tar -xvf yay.tar.gz && cd yay_*_x86_64* && sudo cp yay /usr/bin",
		}
	default:
		return nil
	}
}

func GetDnfPkgCommand(command string) []string {
	bin := "dnf"

	switch command {
	case "autoremove":
		return []string{"sudo", bin, "autoremove"}
	case "clean":
		return []string{"sudo", bin, "clean"}
	case "install":
		return []string{"sudo", bin, "install"}
	case "list":
		return []string{"sudo", bin, "list"}
	case "purge":
		return []string{"sudo", bin, "remove"}
	case "remove":
		return []string{"sudo", bin, "remove"}
	case "search":
		return []string{"sudo", bin, "search"}
	case "show":
		return []string{"sudo", bin, "info"}
	case "update":
		return []string{"sudo", bin, "upgrade", "--refresh"}
	case "upgrade":
		return []string{"sudo", bin, "upgrade"}
	default:
		return nil
	}
}

func GetApkPkgCommand(command string) []string {
	bin := "apk"

	switch command {
	case "autoremove":
		return []string{"echo", "Not implemented yet! "}
	case "clean":
		return []string{"echo", "Not implemented yet! "}
	case "install":
		return []string{"sudo", bin, "add"}
	case "list":
		return []string{"sudo", bin, "list"}
	case "purge":
		return []string{"sudo", bin, "del"}
	case "remove":
		return []string{"sudo", bin, "del"}
	case "search":
		return []string{"sudo", bin, "search"}
	case "show":
		return []string{"sudo", bin, "info"}
	case "update":
		return []string{"sudo", bin, "update"}
	case "upgrade":
		return []string{"sudo", bin, "upgrade", "--available"}
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
