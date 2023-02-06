package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewShowCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("show <package>",
		apx.Trans("show.long"),
		apx.Trans("show.short"),
		show).WithBoolFlag(
		cmdr.NewBoolFlag(
			"isinstalled",
			"i",
			apx.Trans("show.isInstalled"),
			false,
		),
	)

	cmd.Example = "apx show htop\napx show -i neovim"
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func show(cmd *cobra.Command, args []string) error {

	if cmd.Flag("isinstalled").Changed {
		result, err := container.IsPackageInstalled(args[0])
		if err != nil {
			return err
		}

		if result {
			cmdr.Info.Printf(apx.Trans("show.found", args[0]))
			os.Exit(0)
		} else {
			cmdr.Info.Printf(apx.Trans("show.notFound", args[0]))
			os.Exit(1)
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("show")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
