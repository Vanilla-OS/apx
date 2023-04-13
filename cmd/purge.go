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

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewPurgeCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("purge",
		apx.Trans("purge.long"),
		apx.Trans("purge.short"),
		purge)

	cmd.Example = "apx purge htop"
	cmd.Args = cobra.MinimumNArgs(1)
	return cmd
}

func purge(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return errors.New(apx.Trans("apx.notForNix"))

	}
	command := append([]string{}, container.GetPkgCommand("purge")...)
	command = append(command, args...)

	for _, pkg := range args {
        binaries, err := container.BinariesProvidedByPackage(pkg)
        if err != nil {
            return err
        }

        for _, binary := range binaries {
            container.RemoveDesktopEntry(binary)

            err := container.RemoveBinary(binary, false)
            if err != nil {
                cmdr.Error.Printf("Error unexporting binary: %s\n", err)
                return err
            }
        }
	}

	err := container.Run(command...)
	if err != nil {
		return err
	}

	return nil
}
