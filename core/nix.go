package core

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/vanilla-os/orchid/cmdr"
)

type UnitData struct {
	User string
}

func NixInstallPackage(pkg string, unfree bool) error {
	cmd := []string{}
	cmd = append(cmd, "nix", "profile", "install")
	if unfree {
		cmd = append(cmd, "--impure")
	}
	cmd = append(cmd, "nixpkgs#"+pkg)
	install := exec.Command(cmd[0], cmd[1:]...)
	install.Env = append(install.Env, "NIXPKGS_ALLOW_UNFREE=1")
	install.Stderr = os.Stderr
	install.Stdin = os.Stdin
	install.Stdout = os.Stdout

	err := install.Run()
	if err != nil {
		cmdr.Error.Println("error installing package")
		cmdr.Error.Println("have you run the `init` command yet?")
		return err
	}

	return nil

}
func NixRemovePackage(pkg string) error {
	list := exec.Command("nix", "profile", "list")
	bb, err := list.Output()
	if err != nil {
		log.Default().Println("error getting installed packaged")
		log.Default().Println("have you run the `init` command yet?")
		return err
	}
	lines := bytes.Split(bb, []byte("\n"))
	needle := []byte("." + pkg)
	var pkgNumber string
	// output:
	//5 flake:nixpkgs#legacyPackages.x86_64-linux.go github:NixOS/nixpkgs/79feedf38536de2a27d13fe2eaf200a9c05193ba#legacyPackages.x86_64-linux.go /nix/store/v6i0a6bfx3707airawpc2589pbbl465r-go-1.19.5
	if len(lines) > 0 {
		for _, line := range lines {
			// split the line by fields, field[0] is the package number
			// field[1] has the full package name
			pieces := bytes.Split(line, []byte(" "))
			if len(pieces) > 1 {
				if bytes.Contains(pieces[1], needle) {
					// this is our package
					pkgNumber = string(pieces[0])
					break
				}
			}
		}
		if pkgNumber == "" {
			return errors.New("package not found")
		}
		remove := exec.Command("nix", "profile", "remove", pkgNumber)
		err = remove.Run()
		return err

	}
	return errors.New("no packages installed")

}
func NixInit() error {
	// get user name for the systemd units
	user := os.Getenv("USER")
	if user == "" {
		return errors.New("can't get current user")
	}
	// make sure this isn't being run as root
	if user == "root" {
		return errors.New("init must not be run as root user")
	}
	unitData := UnitData{User: user}

	// prompt for confirmation
	log.Default().Printf(`This will create a ".nix" folder in your home directory
and set up some SystemD units to mount that folder at /nix before running the installation
Confirm 'y' to continue. [y/N] `)

	var proceed string
	fmt.Scanln(&proceed)
	proceed = strings.ToLower(proceed)

	if proceed != "y" {
		log.Default().Printf("operation canceled at user request")
		os.Exit(0)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Default().Printf("unable to get home directory")
		return err
	}
	nixDir := path.Join(homedir, ".nix")
	fi, err := os.Stat(nixDir)
	if err != nil {
		// it's ok if the directory doesn't exist
		// but not ok if there is some other error
		if !os.IsNotExist(err) {
			log.Default().Printf(err.Error())
			return err
		}
	}
	if fi != nil {
		log.Default().Printf("$HOME/.nix already exists, refusing to overwrite")
		os.Exit(0)
	}
	// create local nix store
	log.Default().Printf("Creating $HOME/.nix")

	err = os.MkdirAll(nixDir, 0755)
	if err != nil {
		log.Default().Printf("error creating $HOME/.nix")
		return err
	}
	log.Default().Printf("Creating systemd units to mount /nix")

	err = makeUnit(unitData, "/etc/systemd/system/nix.mount", mountTemplate)
	if err != nil {
		log.Default().Printf("error creating directory mount unit")
		return err
	}
	err = makeUnit(unitData, "/etc/systemd/system/ensure-nix-dir.service", ensureTemplate)
	if err != nil {
		log.Default().Printf("error creating ensure directory unit")
		return err
	}
	err = makeUnit(unitData, "/etc/systemd/system/ensure-nix-own.service", ownerTemplate)
	if err != nil {
		log.Default().Printf("error creating directory ownership unit")
		return err
	}
	err = makeUnit(u, "/etc/profile.d/xXDG.sh", xdgConfig)
	if err != nil {
		log.Default().Printf("error creating directory ownership unit")
		return err
	}

	log.Default().Printf("Enabling systemd units")
	reload := exec.Command("sudo", "systemctl", "daemon-reload")
	reload.Stderr = os.Stderr
	reload.Stdin = os.Stdin
	reload.Stdout = os.Stdout

	err = reload.Run()
	if err != nil {
		log.Default().Printf("error reloading systemd daemon")
		return err
	}
	// enable the mount unit, which depends on the others
	enable := exec.Command("sudo", "systemctl", "enable", "--now", "/etc/systemd/system/nix.mount")
	enable.Stderr = os.Stderr
	enable.Stdin = os.Stdin
	enable.Stdout = os.Stdout

	err = enable.Run()
	if err != nil {
		log.Default().Printf("error enabling nix mount")
		return err
	}
	// chown now so we can install
	chown := exec.Command("sudo", "chown", "-R", user+":root", "/nix")
	chown.Stderr = os.Stderr
	chown.Stdin = os.Stdin
	chown.Stdout = os.Stdout

	err = chown.Run()
	if err != nil {
		log.Default().Printf("error changing ownership of /nix")
		return err
	}

	nix := exec.Command("bash", "-c", singleUserCommand)
	nix.Stderr = os.Stderr
	nix.Stdin = os.Stdin
	nix.Stdout = os.Stdout

	err = nix.Run()
	if err != nil {
		log.Default().Printf("error installing nix")
		return err
	}
	// nix conf
	log.Default().Printf("Creating nix configuration file to enable flakes and nix command")

	nixConfDir := path.Join(homedir, ".config", "nix")

	err = os.MkdirAll(nixConfDir, 0755)
	if err != nil {
		log.Default().Printf("error creating nix configuration file")
		return err
	}
	nixConfFile := path.Join(nixConfDir, "nix.conf")
	conf, err := os.Create(nixConfFile)
	if err != nil {
		return err
	}
	_, err = conf.Write([]byte(nixConf))
	if err != nil {
		return err
	}
	return nil

}

func makeUnit(u UnitData, location, t string) error {
	tmpl, err := template.New("unit").Parse(t)
	if err != nil {
		return err
	}
	var bb bytes.Buffer
	err = tmpl.Execute(&bb, u)
	if err != nil {
		return err
	}
	// create a file in /tmp
	file, err := os.Create("/tmp/tmpUnit")
	if err != nil {
		return err
	}
	defer os.Remove("/tmp/tmpUnit")
	_, err = file.Write(bb.Bytes())
	if err != nil {
		return err
	}
	// move it to the system unit directory
	cmd := exec.Command("sudo", "mv", "/tmp/tmpUnit", location)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("sudo", "chown", "root:root", location)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()

}

var mountTemplate = `[Unit]
Description=Mount /nix from ~/.nix
After=local-fs.target var-home.mount ensure-nix-dir.service
Wants=ensure-nix-dir.service
[Mount]
Options=bind,nofail
What=/home/{{.User}}/.nix
Where=/nix
[Install]
WantedBy=multi-user.target
`
var ensureTemplate = `[Unit]
Description=Ensure /nix is present
[Service]
Type=oneshot
ExecStart=mkdir -p -m 0755 /nix
`
var ownerTemplate = `[Unit]
Description=Ensure /nix ownership is correct
Wants=ensure-nix-dir.service
[Service]
Type=oneshot
ExecStart=chown -R {{.User}}:root /nix
`
var singleUserCommand = "sh <(curl -L https://nixos.org/nix/install) --no-daemon"

var nixConf = "experimental-features = nix-command flakes"

var xdgConfig = "XDG_DATA_DIRS=$HOME/.nix-profile/share:$XDG_DATA_DIRS"
