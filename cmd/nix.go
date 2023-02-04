package cmd

import "github.com/vanilla-os/orchid/cmdr"

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

func NewNixCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"nix [command]",
		apx.Trans("nix.long"),
		apx.Trans("nix.short"),
		nil,
	)

	cmd.AddCommand(NewNixInitCommand())
	cmd.AddCommand(NewNixInstallCommand())
	cmd.AddCommand(NewNixRemoveCommand())

	return cmd
}
