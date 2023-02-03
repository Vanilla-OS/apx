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
)

func NewNixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nix",
		Short: "Manage nix installation",
		Long:  "Manage a custom installation of nix in your $HOME directory.",
	}
	cmd.AddCommand(NewNixInitCommand())
	cmd.AddCommand(NewNixInstallCommand())
	cmd.AddCommand(NewNixRemoveCommand())

	return cmd
}
