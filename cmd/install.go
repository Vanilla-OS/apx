package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewInstallCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"install",
		apx.Trans("install.long"),
		apx.Trans("install.short"),
		install,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"assume-yes",
			"y",
			apx.Trans("install.assumeYes"),
			false,
		),
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"fix-broken",
			"f",
			apx.Trans("install.fixBroken"),
			false,
		),
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"no-export",
			"",
			apx.Trans("install.noExport"),
			false,
		),
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"sideload",
			"",
			apx.Trans("install.sideload"),
			false,
		),
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"allow-unfree",
			"",
			apx.Trans("nixinstall.allowUnfree"),
			false,
		),
	)
	/*
				Example: `apx install htop git
		apx install --sideload /path/to/file.deb`,
				Use:   "install <packages>",
				Short: "Install packages inside a managed container",
				RunE:  install,
			}
			cmd.Flags().SetInterspersed(false)
			cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
			cmd.Flags().BoolP("fix-broken", "f", false, "Fix broken deps before installing.")
			cmd.Flags().Bool("no-export", false, "Do not export a desktop entry after the installation.")
			cmd.Flags().Bool("sideload", false, "Install a package from a local file.")
	*/
	cmd.Example = "apx install htop git"
	cmd.Flags().SetInterspersed(false)
	cmd.Args = cobra.MinimumNArgs(1)
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return installPackage(cmd, args)
	}

	no_export := cmd.Flag("no-export").Changed
	assume_yes := cmd.Flag("assume-yes").Changed
	fix_broken := cmd.Flag("fix-broken").Changed
	sideload := cmd.Flag("sideload").Changed

	command := append([]string{}, container.GetPkgCommand("install")...)

	if assume_yes {
		command = append(command, "-y")
	}
	if fix_broken {
		command = append(command, "-f")
	}

	if sideload {
		if len(args) != 1 {

			return fmt.Errorf(apx.Trans("install.sideArgs"))
		}
		path, err := core.MoveToUserTemp(args[0])
		if err != nil {
			return fmt.Errorf(apx.Trans("install.sideUserTemp", err))
		}
		command = append(command, path)
	} else {
		command = append(command, args...)
	}

	err := container.Run(command...)
	if err != nil {
		return err
	}

	if no_export {
		return nil
	}

	if !sideload {
		for _, pkg := range args {
			result, err := cmdr.Confirm.Show(apx.Trans("install.exportBinPrompt"))
			if err != nil {
				return err
			}

			if result {
				err := container.ExportBinary(pkg)
				if err != nil {
					cmdr.Error.Printf("Error exporting binary: %s\n", err)
					return err
				}
			}

			container.ExportDesktopEntry(pkg)
		}
	}

	return nil
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
	cmdr.Success.Println(apx.Trans("nixinstall.success"))
	return nil

}
