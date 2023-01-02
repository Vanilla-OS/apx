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

func NewUnexportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx unexport code",
		Use:     "unexport <program>",
		Short:   "Unexport/Remove a program's desktop entry from a managed container",
		RunE:    unexport,
	}
	return cmd
}

func unexport(cmd *cobra.Command, args []string) error {

	return container.RemoveDesktopEntry(args[0])
}
