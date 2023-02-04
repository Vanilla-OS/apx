package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewShowCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("snow", apx.Trans("show.long"), apx.Trans("show.short"), show)
	/*
			Example: "apx show htop",
			Use:     "show <package>",
			Short:   "Show details about a package",
			RunE:    show,
		}
		cmd.Flags().BoolP("isinstalled", "i", false, "Returns only whether package is installed")
	*/
	return cmd
}

func show(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Please specify a package name.")
	}

	if cmd.Flag("isinstalled").Value.String() == "true" {
		result, err := container.IsPackageInstalled(args[0])
		if err != nil {
			return err
		}

		if result {
			fmt.Printf("%s is installed", args[0])
			os.Exit(0)
		} else {
			fmt.Printf("%s is not installed", args[0])
			os.Exit(1)
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("show")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
