package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func installUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Install packages.

Usage:
  apx install <packages>

Examples:
  apx install htop git
`)
	return nil
}

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages",
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
