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

func NewUnexportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx unexport code",
		Use:     "unexport <program>",
		Short:   "Unexport/Remove a program's desktop entry from a managed container",
		RunE:    unexport,
	}
	cmd.Flags().String("bin", "", "Unexport a previously exported binary.")
	return cmd
}

func unexport(cmd *cobra.Command, args []string) error {
	if cmd.Flag("bin").Value.String() != "" {
		bin_name := cmd.Flag("bin").Value.String()
		if err := container.RemoveBinary(bin_name, false); err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("Successfully removed exported binary `%s`.", bin_name)
		}
		return nil
	} else {
		if len(args) == 0 {
			return errors.New("Please specify a program to unexport.")
		}
		return container.RemoveDesktopEntry(args[0])
	}
}
