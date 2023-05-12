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

func NewInitializeCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"init",
		apx.Trans("init.long"),
		apx.Trans("init.short"),
		initialize,
	)

	cmd.Example = "apx init"
	cmd.SilenceUsage = true
	return cmd
}

func initialize(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return initNix(cmd, args)
	}
	if container.Exists() {

		b, err := cmdr.Confirm.Show(apx.Trans("init.confirm"))

		if err != nil {
			return err
		}

		if !b {
			cmdr.Info.Println(apx.Trans("apx.cxl"))
			return nil
		}
	}

	if err := container.Remove(); err != nil {
		cmdr.Error.Printf(apx.Trans("init.remove"), err)
		return err
	}
	if err := container.Create(); err != nil {
		cmdr.Error.Printf(apx.Trans("init.create"), err)
		return err
	}

	return nil
}
func initNix(cmd *cobra.Command, args []string) error {
	// prompt for confirmation

	b, err := cmdr.Confirm.Show(apx.Trans("nixinit.confirm"))
	if err != nil {
		return err
	}

	if !b {
		cmdr.Info.Println(apx.Trans("apx.cxl"))
		return nil
	}

	unfree, err := cmdr.Confirm.Show(apx.Trans("nixinit.unfree"))
	if err != nil {
		return err
	}
	insecure, err := cmdr.Confirm.Show(apx.Trans("nixinit.insecure"))
	if err != nil {
		return err
	}
	err = core.NixInit(unfree, insecure)
	if err != nil {
		return err
	}
	cmdr.Success.Println(apx.Trans("nixinit.success"))
	return nil

}
