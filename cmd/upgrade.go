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

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx upgrade",
		Use:     "upgrade",
		Short:   "Upgrade the system by installing/upgrading available packages.",
		RunE:    upgrade,
	}
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	return cmd
}

func upgrade(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("upgrade")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	container.Run(command...)

	return nil
}
