package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

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

	"github.com/spf13/cobra"
)

type UnitData struct {
	User string
}

func NewNixInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize nix repository",
		RunE:  initNix,
	}

	return cmd
}
func initNix(cmd *cobra.Command, args []string) error {
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
and some systemd units to mount that folder at /nix before running the installation
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
	nixDir := path.Join(homedir, ".nixu")
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
	log.Default().Printf("Enabling systemd units")

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
	log.Default().Printf("Installation complete. Reboot to use nix.")
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
