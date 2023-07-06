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

type MetadataStruct struct {
	ImageName string `json:"image-name"`
	ImageID   string `json:"image-id"`
	Name      string `json:"name"`
	CreatedAt int    `json:"created-at"`
}

func NewListCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"list",
		apx.Trans("list.long"),
		apx.Trans("list.short"),
		list,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"containers",
			"c",
			apx.Trans("list.containers"),
			false,
		)).
		WithBoolFlag(
			cmdr.NewBoolFlag(
				"upgradable",
				"u",
				apx.Trans("list.upgradable"),
				false,
			)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"installed",
			"i",
			apx.Trans("list.installed"),
			false,
		))

	cmd.Flags().SetInterspersed(false)

	cmd.SilenceUsage = true
	return cmd
}

func list(cmd *cobra.Command, args []string) error {
	if cmd.Flag("containers").Changed {
		command := "list"
		container.Run(command)
		return nil
	}

	if cmd.Flag("nix").Changed {
		return errors.New(apx.Trans("apx.notForNix"))

	}
	command := append([]string{}, container.GetPkgCommand("list")...)

	if cmd.Flag("upgradable").Changed {
		command = append(command, "--upgradable")
	}
	if cmd.Flag("installed").Changed {
		command = append(command, "--installed")
	}

	command = append(command, args...)

	container.Run(command...)

	return nil
}
