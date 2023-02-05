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
	"github.com/vanilla-os/orchid/cmdr"
)

func NewInitializeCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"init",
		apx.Trans("init.long"),
		apx.Trans("init.short"),
		initialize,
	)
	/*
			Example: "apx init",
			Use:     "init",
			Short:   "Initialize the managed container",
			RunE:    initialize,
		}
	*/
	cmd.Example = "apx init"
	return cmd
}

func initialize(cmd *cobra.Command, args []string) error {

	if container.Exists() {

		b, err := cmdr.Confirm.Show(apx.Trans("init.confirm"))

		if err != nil {
			return err
		}

		if !b {
			cmdr.Info.Println("Canceled operation at user request")
			return nil
		}
	}

	if err := container.Remove(); err != nil {
		cmdr.Error.Printf("error removing container: %v", err)
		return err
	}
	if err := container.Create(); err != nil {
		cmdr.Error.Printf("error creating container: %v", err)
		return err
	}

	return nil
}
