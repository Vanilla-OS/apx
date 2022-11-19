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

func runUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Run a program inside the container.

Usage:
  apx run <program>

Examples:
  apx run htop
`)
	return nil
}

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a program inside the container.",
		RunE:  run,
	}
	cmd.SetUsageFunc(runUsage)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}
	core.RunContainer(container, args...)

	return nil
}
