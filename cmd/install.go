package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it work inside a container
	with support to installing packages from multiple sources without altering the root filesystem.
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
  apx install <packages>
  apx --aur install <packages>

Examples:
  apx install htop git
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
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "install")...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}
	if cmd.Flag("fix-broken").Value.String() == "true" {
		command = append(command, "-f")
	}

	command = append(command, args...)

	err := core.RunContainer(container, command...)
	if err != nil {
		return err
	}

	if cmd.Flag("no-export").Value.String() == "true" {
		return nil
	}

	for _, pkg := range args {
		core.ExportDesktopEntry(container, pkg)
	}

	return nil
}
