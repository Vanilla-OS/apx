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

func NewNixRemoveCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"remove <pkg>",
		apx.Trans("nixremove.long"),
		apx.Trans("nixremove.short"),
		removePackage,
	)

	cmd.Args = cobra.ExactArgs(1)
	cmd.Example = "apx nix remove jq"
	return cmd
}
func removePackage(cmd *cobra.Command, args []string) error {
	err := core.NixRemovePackage(args[0])
	if err != nil {
		return err
	}
	cmdr.Success.Println(apx.Trans("nixremove.success"))
	return nil

}