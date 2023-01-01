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
)

func runUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Run a program inside a managed container.

Usage:
  apx run <program> [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx run htop
`)
	return nil
}

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a program inside a managed container.",
		RunE:  run,
	}
	cmd.SetUsageFunc(runUsage)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	container.Run(args...)

	return nil
}
