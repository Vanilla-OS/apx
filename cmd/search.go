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

func searchUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Search in package descriptions.

Usage:
  apx search <packages>

Examples:
  apx search htop
`)
	return nil
}

func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search in package descriptions.",
		RunE:  search,
	}
	cmd.SetUsageFunc(searchUsage)
	return cmd
}

func search(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(sys, container, "search")...)
	command = append(command, args...)

	if sys {
		log.Default().Println("Performing operations on the host system.")
		core.AlmostRun(false, command...)
		return nil
	}

	core.RunContainer(container, command...)

	return nil
}
