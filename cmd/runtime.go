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

	handleFunc := func(subSystem *core.SubSystem, reqFunc func(*core.SubSystem, string, []string) error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			err := reqFunc(subSystem, cmd.Name(), args)
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

		commands = append(commands, subSystemCmd)
	}

	return commands
}

func runPkgCmd(subSystem *core.SubSystem, command string, args []string) error {
	if command != "enter" {
		if len(args) == 0 {
			return fmt.Errorf("no packages specified")
		}
	}

	if command != "run" && command != "enter" {
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
		err = subSystem.Exec(finalArgs...)
		if err != nil {
			return fmt.Errorf("error executing command: %s", err)
		}
	} else {
		if command == "run" {
			err := subSystem.Exec(args...)
			if err != nil {
				return fmt.Errorf("error executing command: %s", err)
			}
		} else {
			err := subSystem.Enter()
			if err != nil {
				return fmt.Errorf("error entering subsystem: %s", err)
			}
		}
	}

	return nil
}
