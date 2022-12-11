package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func installUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Install packages inside a managed container.

Usage:
  apx install [options] <packages>

Options:
  -h, --help            Show this help message and exit
  -y, --assume-yes      Proceed without manual confirmation.
  -f, --fix-broken      Fix broken deps before installing.
  --no-export           Do not export a desktop entry after the installation.
  --sideload [path]     Install a package from a local file.

Examples:
  apx install htop git
  apx install --sideload /path/to/file.deb
  apx --aur install htop
  apx --dnf install htop
`)
	return nil
}

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages inside a managed container",
		RunE:  install,
	}
	cmd.SetUsageFunc(installUsage)
	cmd.Flags().SetInterspersed(false)
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	cmd.Flags().BoolP("fix-broken", "f", false, "Fix broken deps before installing.")
	cmd.Flags().Bool("no-export", false, "Do not export a desktop entry after the installation.")
	cmd.Flags().Bool("sideload", false, "Install a package from a local file.")
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"
	no_export := cmd.Flag("no-export").Value.String() == "true"
	assume_yes := cmd.Flag("assume-yes").Value.String() == "true"
	fix_broken := cmd.Flag("fix-broken").Value.String() == "true"
	sideload := cmd.Flag("sideload").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "install")...)

	if assume_yes {
		command = append(command, "-y")
	}
	if fix_broken {
		command = append(command, "-f")
	}

	if sideload {
		if len(args) != 1 {
			return fmt.Errorf("sideload requires the path to a local file")
		}
		path, err := core.MoveToUserTemp(args[0])
		if err != nil {
			return fmt.Errorf("can't move file to user temp: %s", err)
		}
		command = append(command, path)
	} else {
		command = append(command, args...)
	}

	err := core.RunContainer(container, command...)
	if err != nil {
		return err
	}

	if no_export {
		return nil
	}

	for _, pkg := range args {
		core.ExportDesktopEntry(container, pkg)
	}

	return nil
}
