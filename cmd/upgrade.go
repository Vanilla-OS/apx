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
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUpgradeCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"upgrade",
		apx.Trans("upgrade.long"),
		apx.Trans("upgrade.short"),
		upgrade).WithBoolFlag(
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
	cmd.Example = "apx upgrade"

	return cmd
}

func upgrade(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return upgradePackage(cmd, args)

	}
	if cmd.Flag("all").Changed {
		var flags []string
		if cmd.Flag("assume-yes").Changed {
			flags = append(flags, "-y")
		}

		if err := core.ApplyForAll("upgrade", flags); err != nil {
			return err
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("upgrade")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Changed {
		command = append(command, "-y")
	}

	container.Run(command...)

	return nil
}

func upgradePackage(cmd *cobra.Command, args []string) error {

	err := core.NixUpgradePackage(args[0])
	if err != nil {
		return err
	}
	cmdr.Success.Println(apx.Trans("nixupgrade.success"))
	return nil

}
