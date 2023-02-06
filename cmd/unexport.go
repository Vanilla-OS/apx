package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUnexportCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("unexport",
		apx.Trans("unexport.long"),
		apx.Trans("unexport.short"),
		unexport).WithBoolFlag(
		cmdr.NewBoolFlag(
			"bin",
			"",
			apx.Trans("unexport.binFlag"),
			false,
		),
	)

	cmd.Args = cobra.ExactArgs(1)
	cmd.Example = "apx unexport code"
	return cmd
}

func unexport(cmd *cobra.Command, args []string) error {
	if cmd.Flag("bin").Changed {
		bin_name := args[0]
		if err := container.RemoveBinary(bin_name, false); err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			cmdr.Success.Println(apx.Trans("unexport.success", bin_name))
		}
		return nil
	} else {

		return container.RemoveDesktopEntry(args[0])
	}
}
