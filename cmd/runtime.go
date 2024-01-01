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
	"github.com/vanilla-os/apx/v2/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRuntimeCommands() []*cmdr.Command {
	var commands []*cmdr.Command

	subSystems, err := core.ListSubSystems(false, false)
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
			apx.Trans("runtimeCommand.description"),
			apx.Trans("runtimeCommand.description"),
			nil,
		)

		autoRemoveCmd := cmdr.NewCommand(
			"autoremove",
			apx.Trans("runtimeCommand.autoremove.description"),
			apx.Trans("runtimeCommand.autoremove.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		cleanCmd := cmdr.NewCommand(
			"clean",
			apx.Trans("runtimeCommand.clean.description"),
			apx.Trans("runtimeCommand.clean.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		installCmd := cmdr.NewCommand(
			"install",
			apx.Trans("runtimeCommand.install.description"),
			apx.Trans("runtimeCommand.install.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		installCmd.WithBoolFlag(
			cmdr.NewBoolFlag(
				"no-export",
				"n",
				apx.Trans("runtimeCommand.install.options.noExport.description"),
				false,
			),
		)

		listCmd := cmdr.NewCommand(
			"list",
			apx.Trans("runtimeCommand.list.description"),
			apx.Trans("runtimeCommand.list.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		purgeCmd := cmdr.NewCommand(
			"purge",
			apx.Trans("runtimeCommand.purge.description"),
			apx.Trans("runtimeCommand.purge.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		removeCmd := cmdr.NewCommand(
			"remove",
			apx.Trans("runtimeCommand.remove.description"),
			apx.Trans("runtimeCommand.remove.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		searchCmd := cmdr.NewCommand(
			"search",
			apx.Trans("runtimeCommand.search.description"),
			apx.Trans("runtimeCommand.search.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		showCmd := cmdr.NewCommand(
			"show",
			apx.Trans("runtimeCommand.show.description"),
			apx.Trans("runtimeCommand.show.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		updateCmd := cmdr.NewCommand(
			"update",
			apx.Trans("runtimeCommand.update.description"),
			apx.Trans("runtimeCommand.update.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		upgradeCmd := cmdr.NewCommand(
			"upgrade",
			apx.Trans("runtimeCommand.upgrade.description"),
			apx.Trans("runtimeCommand.upgrade.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		runCmd := cmdr.NewCommand(
			"run",
			apx.Trans("runtimeCommand.run.description"),
			apx.Trans("runtimeCommand.run.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		enterCmd := cmdr.NewCommand(
			"enter",
			apx.Trans("runtimeCommand.enter.description"),
			apx.Trans("runtimeCommand.enter.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		exportCmd := cmdr.NewCommand(
			"export",
			apx.Trans("runtimeCommand.export.description"),
			apx.Trans("runtimeCommand.export.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				apx.Trans("runtimeCommand.export.options.appName.description"),
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				apx.Trans("runtimeCommand.export.options.bin.description"),
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				apx.Trans("runtimeCommand.export.options.binOutput.description"),
				"",
			),
		)
		unexportCmd := cmdr.NewCommand(
			"unexport",
			apx.Trans("runtimeCommand.unexport.description"),
			apx.Trans("runtimeCommand.unexport.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				apx.Trans("runtimeCommand.unexport.options.appName.description"),
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				apx.Trans("runtimeCommand.unexport.options.bin.description"),
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				apx.Trans("runtimeCommand.unexport.options.binOutput.description"),
				"",
			),
		)

		startCmd := cmdr.NewCommand(
			"start",
			apx.Trans("subsystems.start.description"),
			apx.Trans("subsystems.start.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		stopCmd := cmdr.NewCommand(
			"stop",
			apx.Trans("subsystems.stop.description"),
			apx.Trans("subsystems.stop.description"),
			handleFunc(subSystem, runPkgCmd),
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
		subSystemCmd.AddCommand(startCmd)
		subSystemCmd.AddCommand(stopCmd)

		commands = append(commands, subSystemCmd)
	}

	return commands
}

func runPkgCmd(subSystem *core.SubSystem, command string, cmd *cobra.Command, args []string) error {
	if command != "run" && command != "enter" && command != "export" && command != "unexport" && command != "start" && command != "stop" {
		pkgManager, err := subSystem.Stack.GetPkgManager()
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.cantAccessPkgManager"), err)
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
			return fmt.Errorf(apx.Trans("apx.error.unknownCommand"), command)
		}

		if command == "remove" {
			exportedN, err := subSystem.UnexportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.unexportedApps"), exportedN)
			}
		}

		finalArgs := pkgManager.GenCmd(realCommand, args...)
		fmt.Println("finalArgs", finalArgs)
		_, err = subSystem.Exec(false, finalArgs...)
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.executingCommand"), err)
		}

		if command == "install" && !cmd.Flag("no-export").Changed {
			exportedN, err := subSystem.ExportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.exportedApps"), exportedN)
			}
		}

		return nil
	}

	if command == "run" {
		_, err := subSystem.Exec(false, args...)
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.executingCommand"), err)
		}

		return nil
	}

	if command == "enter" {
		err := subSystem.Enter()
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.enteringContainer"), err)
		}

		return nil
	}

	if command == "export" || command == "unexport" {
		appName, _ := cmd.Flags().GetString("app-name")
		bin, _ := cmd.Flags().GetString("bin")
		binOutput, _ := cmd.Flags().GetString("bin-output")

		if appName == "" && bin == "" {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.noAppNameOrBin"))
		}

		if appName != "" && bin != "" {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.sameAppOrBin"))
		}

		if command == "export" {
			if appName != "" {
				err := subSystem.ExportDesktopEntry(appName)
				if err != nil {
					return fmt.Errorf(apx.Trans("runtimeCommand.error.exportingApp"), err)
				}

				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.exportedApp"), appName)
			} else {
				err := subSystem.ExportBin(bin, binOutput)
				if err != nil {
					return fmt.Errorf(apx.Trans("runtimeCommand.error.exportingBin"), err)
				}

				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.exportedBin"), bin)
			}
		} else {
			if appName != "" {
				err := subSystem.UnexportDesktopEntry(appName)
				if err != nil {
					return fmt.Errorf(apx.Trans("runtimeCommand.error.unexportingApp"), err)
				}

				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.unexportedApp"), appName)
			} else {
				err := subSystem.UnexportBin(bin, binOutput)
				if err != nil {
					return fmt.Errorf(apx.Trans("runtimeCommand.error.unexportingBin"), err)
				}

				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.unexportedBin"), bin)
			}
		}
	}

	if command == "start" {
		cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.startingContainer"), subSystem.Name)
		err := subSystem.Start()
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.startingContainer"), err)
		}

		cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.startedContainer"), subSystem.Name)
	}

	if command == "stop" {
		cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.stoppingContainer"), subSystem.Name)
		err := subSystem.Stop()
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.stoppingContainer"), err)
		}

		cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.stoppedContainer"), subSystem.Name)
	}

	return nil
}
