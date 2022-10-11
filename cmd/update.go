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

func updateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Update list of available packages.

Usage:
  apx update

Examples:
  apx update
`)
	return nil
}

func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update list of available packages.",
		RunE:  update,
	}
	cmd.SetUsageFunc(updateUsage)
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	aur := cmd.Flag("aur").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	}

	command := append([]string{}, core.GetPkgManager(sys)...)
	command = append(command, core.GetPkgCommand(sys, container, "update")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	if sys {
		log.Default().Println("Performing operations on the host system.")
		core.PkgManagerSmartLock()
		core.AlmostRun(false, command...)
		return nil
	}

	core.RunContainer(container, command...)

	return nil
}
