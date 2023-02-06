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

func NewNixInstallCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"install <pkg>",
		apx.Trans("nixinstall.long"),
		apx.Trans("nixinstall.short"),
		installPackage,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"allow-unfree",
			"u",
			apx.Trans("nixinstall.allowUnfree"),
			false,
		),
	)
	/*
			Use:     "install <pkg>",
			Short:   "Install nix package",
			Long:    `Install a package from the nixpkgs repository as a flake.`,
			Example: "apx nix install jq",
			RunE:    installPackage,
			Args:    cobra.ExactArgs(1),
		}
	*/
	cmd.Args = cobra.ExactArgs(1)
	cmd.Example = "apx nix install jq"
	return cmd
}
func installPackage(cmd *cobra.Command, args []string) error {
	allowUnfree := false
	if cmd.Flags().Changed("allow-unfree") {
		allowUnfree = true
	}
	err := core.NixInstallPackage(args[0], allowUnfree)
	if err != nil {
		return err
	}
	cmdr.Success.Println("Package installation complete")
	return nil

}
