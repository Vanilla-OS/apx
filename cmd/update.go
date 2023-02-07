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
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUpdateCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("update",
		apx.Trans("update.long"),
		apx.Trans("update.short"),
		update).WithBoolFlag(
		cmdr.NewBoolFlag(
			"all",
			"a",
			apx.Trans("apx.allFlag"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"assume-yes",
			"y",
			apx.Trans("apx.assumeYes"),
			false,
		))

	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return errors.New(apx.Trans("apx.notForNix"))

	}
	if cmd.Flag("all").Changed {
		var flags []string
		if cmd.Flag("assume-yes").Changed {
			flags = append(flags, "-y")
		}

		if err := core.ApplyForAll("update", flags); err != nil {
			return err
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("update")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Changed {
		command = append(command, "-y")
	}

	container.Run(command...)

	return nil
}
