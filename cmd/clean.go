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
	"github.com/vanilla-os/apx/settings"
)

func cleanUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Clean the package manager cache.

Usage:
  apx clean

Examples:
  apx clean
`)
	return nil
}

func NewCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean the package manager cache",
		RunE:  clean,
	}
	cmd.SetUsageFunc(cleanUsage)
	return cmd
}

func clean(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	command := append([]string{}, core.GetPkgManager(sys)...)
	command = append(command, settings.Cnf.PkgManager.CmdClean)
	command = append(command, args...)

	if sys {
		log.Default().Println("Performing operations on the host system.")
		core.PkgManagerSmartLock()
		core.AlmostRun(false, command...)
		return nil
	}

	core.RunContainer(command...)
	return nil
}
