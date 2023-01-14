package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx export htop",
		Use:     "export <program>",
		Short:   "Export/Recreate a program's desktop entry from a managed container",
		RunE:    export,
	}
	cmd.Flags().String("bin", "", "Export a binary instead.")
	return cmd
}

func export(cmd *cobra.Command, args []string) error {
	if cmd.Flag("bin").Value.String() != "" {
		if err := container.ExportBinary(cmd.Flag("bin").Value.String()); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	} else {
		if len(args) == 0 {
			return errors.New("Please specify a program to export.")
		}
		container.ExportDesktopEntry(args[0])
	}
	return nil
}
