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

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewExportCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"export <application/binary>",
		apx.Trans("export.long"),
		apx.Trans("export.short"),
		export,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"bin",
			"",
			apx.Trans("export.binFlag"),
			false,
		),
	)
	cmd.Example = "apx export htop\napx export --bin fzf"
	cmd.Args = cobra.ExactArgs(1)
	cmd.SilenceUsage = true
	return cmd
}

func export(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return errors.New(apx.Trans("apx.notForNix"))

	}
	if cmd.Flag("bin").Changed {
		err := container.ExportBinary(args[0])
		if err != nil {
			cmdr.Error.Printf("Error exporting binary: %s\n", err)
			return err
		}
	} else {
		container.ExportDesktopEntry(args[0])
	}
	return nil
}
