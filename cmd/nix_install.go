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

func NewNixInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install <pkg>",
		Short:   "Install nix package",
		Long:    `Install a package from the nixpkgs repository as a flake.`,
		Example: "apx nix install jq",
		RunE:    installPackage,
		Args:    cobra.ExactArgs(1),
	}

	return cmd
}
func installPackage(cmd *cobra.Command, args []string) error {
	err := core.NixInstallPackage(args[0])
	if err != nil {
		return err
	}
	log.Default().Printf("Package installation complete")
	return nil

}
