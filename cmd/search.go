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

func NewSearchCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("search",
		apx.Trans("search.long"),
		apx.Trans("search.short"),
		search)

	cmd.Example = "apx search neovim"
	return cmd
}

func search(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		//return errors.New(apx.Trans("apx.notForNix"))
		return core.NixSearchPackage(args[0])

	}
	command := append([]string{}, container.GetPkgCommand("search")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
