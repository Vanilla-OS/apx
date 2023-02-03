package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func NewNixRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <pkg>",
		Short: "Remove a nix package",
		Long: `Remove the first matching package installed in the
default nix profile. Note that it is possible to have 
multiple installations of the same package. If you need 
more precision use 'nix profile list' and 'nix profile
remove'.`,

		RunE: removePackage,
		Args: cobra.ExactArgs(1),
	}

	return cmd
}
func removePackage(cmd *cobra.Command, args []string) error {
	err := core.NixRemovePackage(args[0])
	if err != nil {
		return err
	}
	log.Default().Printf("Package removal complete")
	return nil

}
