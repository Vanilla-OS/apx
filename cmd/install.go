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
	"log"

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
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	aur := cmd.Flag("aur").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	}

	command := append([]string{}, core.GetPkgCommand(sys, container, "install")...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}
	if cmd.Flag("fix-broken").Value.String() == "true" {
		command = append(command, "-f")
	}

	command = append(command, args...)

	if sys {
		log.Default().Println("Performing operations on the host system.")
		core.PkgManagerSmartLock()
		core.AlmostRun(false, command...)
		return nil
	}

	core.RunContainer(container, command...)

	for _, pkg := range args {
		core.ExportDesktopEntry(container, pkg)
	}

	return nil
}
