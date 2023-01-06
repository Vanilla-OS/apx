package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"github.com/spf13/cobra"
)

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx export htop",
		Use:     "export <program>",
		Short:   "Export/Recreate a program's desktop entry from a managed container",
		RunE:    export,
	}
	return cmd
}

func export(cmd *cobra.Command, args []string) error {

	container.ExportDesktopEntry(args[0])
	return nil
}
