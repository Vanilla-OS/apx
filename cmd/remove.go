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

func removeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Remove packages inside a managed container.

Usage:
  apx remove <packages> [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx remove htop
`)
	return nil
}

func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove packages inside a managed container.",
		RunE:  remove,
	}
	cmd.SetUsageFunc(removeUsage)
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")

	return cmd
}

func remove(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"
	apk := cmd.Flag("apk").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	} else if apk {
		container = "apk"
	}

	command := append([]string{}, core.GetPkgCommand(container, "remove")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	err := core.RunContainer(container, command...)
	if err != nil {
		return err
	}

	for _, pkg := range args {
		core.RemoveDesktopEntry(container, pkg)
	}

	return nil
}
