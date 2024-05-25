package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"
	"slices"

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
			return reqFunc(subSystem, cmd.Name(), cmd, args)
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
			apx.Trans("runtimeCommand.start.description"),
			apx.Trans("runtimeCommand.start.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		stopCmd := cmdr.NewCommand(
			"stop",
			apx.Trans("runtimeCommand.stop.description"),
			apx.Trans("runtimeCommand.stop.description"),
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

var baseCmds = []string{"run", "enter", "export", "unexport", "start", "stop"}

// isBaseCommand informs whether the command is a subsystem-base command
// (e.g. run, enter) instead of a subsystem-specific one (e.g. install, update)
func isBaseCommand(command string) bool {
	return slices.Contains(baseCmds, command)
}

// pkgManagerCommands maps command line arguments into package manager commands
func pkgManagerCommands(pkgManager *core.PkgManager, command string) (string, error) {
	switch command {
	case "autoremove":
		return pkgManager.CmdAutoRemove, nil
	case "clean":
		return pkgManager.CmdClean, nil
	case "install":
		return pkgManager.CmdInstall, nil
	case "list":
		return pkgManager.CmdList, nil
	case "purge":
		return pkgManager.CmdPurge, nil
	case "remove":
		return pkgManager.CmdRemove, nil
	case "search":
		return pkgManager.CmdSearch, nil
	case "show":
		return pkgManager.CmdShow, nil
	case "update":
		return pkgManager.CmdUpdate, nil
	case "upgrade":
		return pkgManager.CmdUpgrade, nil
	default:
		return "", fmt.Errorf(apx.Trans("apx.errors.unknownCommand"), command)
	}
}

func runPkgCmd(subSystem *core.SubSystem, command string, cmd *cobra.Command, args []string) error {
	if !isBaseCommand(command) {
		pkgManager, err := subSystem.Stack.GetPkgManager()
		if err != nil {
			return fmt.Errorf(apx.Trans("runtimeCommand.error.cantAccessPkgManager"), err)
		}

		realCommand, err := pkgManagerCommands(pkgManager, command)
		if err != nil {
			return err
		}

		if command == "remove" {
			exportedN, err := subSystem.UnexportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printfln(apx.Trans("runtimeCommand.info.unexportedApps"), exportedN)
			}
		}

		finalArgs := pkgManager.GenCmd(realCommand, args...)
		fmt.Println("finalArgs", finalArgs)
		_, err = subSystem.Exec(false, false, finalArgs...)
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
		_, err := subSystem.Exec(false, false, args...)
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

		return handleExport(subSystem, command, appName, bin, binOutput)
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

func handleExport(subSystem *core.SubSystem, command, appName, bin, binOutput string) error {
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

	return nil
}
