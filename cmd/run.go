package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx run htop",
		Use:     "run <program>",
		Short:   "Run a program inside a managed container.",
		RunE:    run,
	}
	//cmd.SetUsageFunc(runUsage)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	container.Run(args...)

	return nil
}
