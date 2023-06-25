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

func NewRuntimeCommands() []*cmdr.Command {
	var commands []*cmdr.Command

	subSystems, err := core.ListSubSystems()
	if err != nil {
		return []*cmdr.Command{}
	}

	handleFunc := func(subSystem *core.SubSystem, reqFunc func(*core.SubSystem, string, *cobra.Command, []string) error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			err := reqFunc(subSystem, cmd.Name(), cmd, args)
			return err
		}
	}

	for _, subSystem := range subSystems {
		subSystemCmd := cmdr.NewCommand(
			subSystem.Name,
			apx.Trans("runtimeCommand.long"),
			apx.Trans("runtimeCommand.short"),
			nil,
		)

		autoRemoveCmd := cmdr.NewCommand(
			"autoremove",
			"Remove unused packages from the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		cleanCmd := cmdr.NewCommand(
			"clean",
			"Clean the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)

		installCmd := cmdr.NewCommand(
			"install",
			"Install packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		installCmd.WithBoolFlag(
			cmdr.NewBoolFlag(
				"no-export",
				"n",
				"Do not export a desktop entry for the app",
				false,
			),
		)

		listCmd := cmdr.NewCommand(
			"list",
			"List packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		purgeCmd := cmdr.NewCommand(
			"purge",
			"Purge packages from the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		removeCmd := cmdr.NewCommand(
			"remove",
			"Remove packages from the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		searchCmd := cmdr.NewCommand(
			"search",
			"Search packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		showCmd := cmdr.NewCommand(
			"show",
			"Show packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		updateCmd := cmdr.NewCommand(
			"update",
			"Update packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		upgradeCmd := cmdr.NewCommand(
			"upgrade",
			"Upgrade packages in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		runCmd := cmdr.NewCommand(
			"run",
			"Run a command in the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		enterCmd := cmdr.NewCommand(
			"enter",
			"Enter the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)

		exportCmd := cmdr.NewCommand(
			"export",
			"Export a binary or an app from the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				"Name of the app to export",
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				"Path of the binary to export",
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				"Path of the binary output (default: ~/.local/bin/)",
				"",
			),
		)

		unexportCmd := cmdr.NewCommand(
			"unexport",
			"Unexport a binary or an app from the subsystem",
			"",
			handleFunc(subSystem, runPkgCmd),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				"Name of the app to unexport",
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				"Path of the binary to unexport",
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				"Path of where the binary was exported (default: ~/.local/bin/)",
				"",
			),
		)

		subSystemCmd.AddCommand(autoRemoveCmd)
		subSystemCmd.AddCommand(cleanCmd)
		subSystemCmd.AddCommand(installCmd)
		subSystemCmd.AddCommand(listCmd)
		subSystemCmd.AddCommand(purgeCmd)
		subSystemCmd.AddCommand(removeCmd)
		subSystemCmd.AddCommand(searchCmd)
		subSystemCmd.AddCommand(showCmd)
		subSystemCmd.AddCommand(updateCmd)
		subSystemCmd.AddCommand(upgradeCmd)
		subSystemCmd.AddCommand(runCmd)
		subSystemCmd.AddCommand(enterCmd)
		subSystemCmd.AddCommand(exportCmd)
		subSystemCmd.AddCommand(unexportCmd)

		commands = append(commands, subSystemCmd)
	}

	return commands
}

func runPkgCmd(subSystem *core.SubSystem, command string, cmd *cobra.Command, args []string) error {
	if command != "enter" && command != "export" && command != "unexport" {
		if len(args) == 0 {
			return fmt.Errorf("no packages specified")
		}
	}

	if command != "run" && command != "enter" && command != "export" && command != "unexport" {
		pkgManager, err := subSystem.Stack.GetPkgManager()
		if err != nil {
			return fmt.Errorf("error getting package manager: %s", err)
		}

		var realCommand string
		switch command {
		case "autoremove":
			realCommand = pkgManager.CmdAutoRemove
		case "clean":
			realCommand = pkgManager.CmdClean
		case "install":
			realCommand = pkgManager.CmdInstall
		case "list":
			realCommand = pkgManager.CmdList
		case "purge":
			realCommand = pkgManager.CmdPurge
		case "remove":
			realCommand = pkgManager.CmdRemove
		case "search":
			realCommand = pkgManager.CmdSearch
		case "show":
			realCommand = pkgManager.CmdShow
		case "update":
			realCommand = pkgManager.CmdUpdate
		case "upgrade":
			realCommand = pkgManager.CmdUpgrade
		default:
			return fmt.Errorf("unknown command: %s", command)
		}

		finalArgs := pkgManager.GenCmd(realCommand, args...)
		_, err = subSystem.Exec(false, finalArgs...)
		if err != nil {
			return fmt.Errorf("error executing command: %s", err)
		}

		if command == "install" && !cmd.Flag("no-export").Changed {
			exportedN, err := subSystem.ExportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printf("Exported %d desktop entries\n", exportedN)
			}
		}

		return nil
	}

	if command == "run" {
		_, err := subSystem.Exec(false, args...)
		if err != nil {
			return fmt.Errorf("error executing command: %s", err)
		}

		return nil
	}

	if command == "enter" {
		err := subSystem.Enter()
		if err != nil {
			return fmt.Errorf("error entering subsystem: %s", err)
		}

		return nil
	}

	if command == "export" || command == "unexport" {
		appName, _ := cmd.Flags().GetString("app-name")
		bin, _ := cmd.Flags().GetString("bin")
		binOutput, _ := cmd.Flags().GetString("bin-output")

		if appName == "" && bin == "" {
			return fmt.Errorf("app-name and bin cannot be both empty")
		}

		if appName != "" && bin != "" {
			return fmt.Errorf("app-name and bin cannot be both set")
		}

		if command == "export" {
			if appName != "" {
				err := subSystem.ExportDesktopEntry(appName)
				if err != nil {
					return fmt.Errorf("error exporting app: %s", err)
				}

				cmdr.Info.Printf("Exported app %s\n", appName)
			} else {
				err := subSystem.ExportBin(bin, binOutput)
				if err != nil {
					return fmt.Errorf("error exporting bin: %s", err)
				}

				cmdr.Info.Printf("Exported binary %s\n", bin)
			}
		} else {
			if appName != "" {
				err := subSystem.UnexportDesktopEntry(appName)
				if err != nil {
					return fmt.Errorf("error unexporting app: %s", err)
				}

				cmdr.Info.Printf("Unexported app %s\n", appName)
			} else {
				err := subSystem.UnexportBin(bin, binOutput)
				if err != nil {
					return fmt.Errorf("error unexporting bin: %s", err)
				}

				cmdr.Info.Printf("Unexported binary %s\n", bin)
			}
		}
	}

	return nil
}
