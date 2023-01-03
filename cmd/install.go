package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
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
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	no_export := cmd.Flag("no-export").Value.String() == "true"
	assume_yes := cmd.Flag("assume-yes").Value.String() == "true"
	fix_broken := cmd.Flag("fix-broken").Value.String() == "true"
	sideload := cmd.Flag("sideload").Value.String() == "true"

	command := append([]string{}, container.GetPkgCommand("install")...)

	if assume_yes {
		command = append(command, "-y")
	}
	if fix_broken {
		command = append(command, "-f")
	}

	if sideload {
		if len(args) != 1 {
			return fmt.Errorf("sideload requires the path to a local file")
		}
		path, err := core.MoveToUserTemp(args[0])
		if err != nil {
			return fmt.Errorf("can't move file to user temp: %s", err)
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
			container.ExportDesktopEntry(pkg)
		}
	}

	return nil
}
