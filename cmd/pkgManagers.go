package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewPkgManagersCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"pkgmanagers",
		apx.Trans("pkgmanagers.description"),
		apx.Trans("pkgmanagers.description"),
		nil,
	)

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		apx.Trans("pkgmanagers.list.description"),
		apx.Trans("pkgmanagers.list.description"),
		listPkgManagers,
	)

	listCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"json",
			"j",
			apx.Trans("pkgmanagers.list.options.json.description"),
			false,
		),
	)

	// Show subcommand
	showCmd := cmdr.NewCommand(
		"show",
		apx.Trans("pkgmanagers.show.description"),
		apx.Trans("pkgmanagers.show.description"),
		showPkgManager,
	)
	showCmd.Args = cobra.MinimumNArgs(1)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		apx.Trans("pkgmanagers.new.description"),
		apx.Trans("pkgmanagers.new.description"),
		newPkgManager,
	)

	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"assume-yes",
			"y",
			apx.Trans("pkgmanagers.new.options.assumeYes.description"),
			false,
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("pkgmanagers.new.options.name.description"),
			"",
		),
	)
	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"need-sudo",
			"S",
			apx.Trans("pkgmanagers.new.options.needSudo.description"),
			false,
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"autoremove",
			"a",
			apx.Trans("pkgmanagers.new.options.autoremove.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"clean",
			"c",
			apx.Trans("pkgmanagers.new.options.clean.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"install",
			"i",
			apx.Trans("pkgmanagers.new.options.install.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"list",
			"l",
			apx.Trans("pkgmanagers.new.options.list.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"purge",
			"p",
			apx.Trans("pkgmanagers.new.options.purge.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"remove",
			"r",
			apx.Trans("pkgmanagers.new.options.remove.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"search",
			"s",
			apx.Trans("pkgmanagers.new.options.search.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"show",
			"w",
			apx.Trans("pkgmanagers.new.options.show.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"update",
			"u",
			apx.Trans("pkgmanagers.new.options.update.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"upgrade",
			"U",
			apx.Trans("pkgmanagers.new.options.upgrade.description"),
			"",
		),
	)

	// Rm subcommand
	rmCmd := cmdr.NewCommand(
		"rm",
		apx.Trans("pkgmanagers.rm.description"),
		apx.Trans("pkgmanagers.rm.description"),
		rmPkgManager,
	)

	rmCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("pkgmanagers.rm.options.name.description"),
			"",
		),
	)
	rmCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			apx.Trans("pkgmanagers.rm.options.force.description"),
			false,
		),
	)

	// Add subcommands to pkgmanagers
	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(rmCmd)

	return cmd
}

func listPkgManagers(cmd *cobra.Command, args []string) error {
	jsonFlag, _ := cmd.Flags().GetBool("json")

	pkgManagers := core.ListPkgManagers()

	if !jsonFlag {
		pkgManagersCount := len(pkgManagers)
		if pkgManagersCount == 0 {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.list.info.noPkgManagers"))
			return nil
		}

		fmt.Printf(apx.Trans("pkgmanagers.list.info.foundPkgManagers"), pkgManagersCount)

		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{apx.Trans("pkgmanagers.labels.name"), apx.Trans("pkgmanagers.labels.builtIn")})

		for _, stack := range pkgManagers {
			builtIn := apx.Trans("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = apx.Trans("apx.terminal.yes")
			}
			table.Append([]string{stack.Name, builtIn})
		}

		table.Render()
	} else {
		jsonPkgManagers, _ := json.MarshalIndent(pkgManagers, "", "  ")
		fmt.Println(string(jsonPkgManagers))
	}

	return nil
}

func showPkgManager(cmd *cobra.Command, args []string) error {
	pkgManagerName := args[0]
	pkgManager, err := core.LoadPkgManager(pkgManagerName)
	if err != nil {
		cmdr.Error.Println(err)
		return nil
	}

	table := core.CreateApxTable(os.Stdout)
	table.Append([]string{apx.Trans("pkgmanagers.labels.name"), pkgManager.Name})
	table.Append([]string{"NeedSudo", fmt.Sprintf("%t", pkgManager.NeedSudo)})
	table.Append([]string{"AutoRemove", pkgManager.CmdAutoRemove})
	table.Append([]string{"Clean", pkgManager.CmdClean})
	table.Append([]string{"Install", pkgManager.CmdInstall})
	table.Append([]string{"List", pkgManager.CmdList})
	table.Append([]string{"Purge", pkgManager.CmdPurge})
	table.Append([]string{"Remove", pkgManager.CmdRemove})
	table.Append([]string{"Search", pkgManager.CmdSearch})
	table.Append([]string{"Show", pkgManager.CmdShow})
	table.Append([]string{"Update", pkgManager.CmdUpdate})
	table.Append([]string{"Upgrade", pkgManager.CmdUpgrade})
	table.Render()

	return nil
}

func newPkgManager(cmd *cobra.Command, args []string) error {
	assumeYes, _ := cmd.Flags().GetBool("assume-yes")
	name, _ := cmd.Flags().GetString("name")
	needSudo, _ := cmd.Flags().GetBool("need-sudo")
	autoRemove, _ := cmd.Flags().GetString("autoremove")
	clean, _ := cmd.Flags().GetString("clean")
	install, _ := cmd.Flags().GetString("install")
	list, _ := cmd.Flags().GetString("list")
	purge, _ := cmd.Flags().GetString("purge")
	remove, _ := cmd.Flags().GetString("remove")
	search, _ := cmd.Flags().GetString("search")
	show, _ := cmd.Flags().GetString("show")
	update, _ := cmd.Flags().GetString("update")
	upgrade, _ := cmd.Flags().GetString("upgrade")

	reader := bufio.NewReader(os.Stdin)

	if name == "" {
		if assumeYes {
			cmdr.Error.Println(apx.Trans("pkgmanagers.new.error.noName"))
			return nil
		}

		cmdr.Info.Println(apx.Trans("pkgmanagers.new.info.askName"))
		name, _ = reader.ReadString('\n')
		name = strings.ReplaceAll(name, "\n", "")
		name = strings.ReplaceAll(name, " ", "")
		if name == "" {
			cmdr.Error.Println(apx.Trans("pkgmanagers.new.error.emptyName"))
			return nil
		}
	}

	if !needSudo && !assumeYes {
		validChoice := false
		for !validChoice {
			cmdr.Info.Println(apx.Trans("pkgmanagers.new.info.askSudo") + ` [y/N]`)
			answer, _ := reader.ReadString('\n')
			if answer == "\n" {
				answer = "n\n"
			}
			answer = strings.ToLower(strings.ReplaceAll(answer, " ", ""))
			switch answer {
			case "y\n":
				needSudo = true
				validChoice = true
			case "n\n":
				needSudo = false
				validChoice = true
			default:
				cmdr.Warning.Println(apx.Trans("apx.errors.invalidChoice"))
			}
		}
	}

	cmdMap := map[string]*string{
		"autoRemove": &autoRemove,
		"clean":      &clean,
		"install":    &install,
		"list":       &list,
		"purge":      &purge,
		"remove":     &remove,
		"search":     &search,
		"show":       &show,
		"update":     &update,
		"upgrade":    &upgrade,
	}

	for cmdName, cmd := range cmdMap {
		if *cmd == "" {
			if assumeYes {
				cmdr.Error.Printf(apx.Trans("pkgmanagers.new.error.noCommand"), cmdName)
				return nil
			}

			cmdr.Info.Printf(apx.Trans("pkgmanagers.new.info.askCommand"), cmdName)
			*cmd, _ = reader.ReadString('\n')
			*cmd = strings.ReplaceAll(*cmd, "\n", "")
			if *cmd == "" {
				cmdr.Error.Printf(apx.Trans("pkgmanagers.new.error.emptyCommand"), cmdName)
				return nil
			}
		}
	}

	if core.PkgManagerExists(name) {
		if assumeYes {
			cmdr.Error.Println(apx.Trans("pkgmanagers.new.error.alreadyExists"), name)
			return nil
		}

		cmdr.Info.Printf(apx.Trans("pkgmanagers.new.info.askOverwrite"), name)
		answer, _ := reader.ReadString('\n')
		answer = strings.ReplaceAll(answer, "\n", "")

		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			cmdr.Info.Println(apx.Trans("apx.info.aborting"))
			return nil
		}
	}

	pkgManager := core.NewPkgManager(name, needSudo, autoRemove, clean, install, list, purge, remove, search, show, update, upgrade, false)
	err := pkgManager.Save()
	if err != nil {
		cmdr.Error.Println(err)
		return nil
	}

	cmdr.Success.Printf(apx.Trans("pkgmanagers.new.success"), name)

	return nil
}

func rmPkgManager(cmd *cobra.Command, args []string) error {
	pkgManagerName, _ := cmd.Flags().GetString("name")
	if pkgManagerName == "" {
		cmdr.Error.Println(apx.Trans("pkgmanagers.rm.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(pkgManagerName)
	if error != nil {
		return error
	}

	stacks := core.ListStackForPkgManager(pkgManager.Name)
	if len(stacks) > 0 {
		cmdr.Error.Printf(apx.Trans("pkgmanagers.rm.error.inUse"), len(stacks))
		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{apx.Trans("pkgmanagers.labels.name"), "Base", "Packages", "PkgManager", apx.Trans("pkgmanagers.labels.builtIn")})
		for _, stack := range stacks {
			builtIn := apx.Trans("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = apx.Trans("apx.terminal.yes")
			}
			table.Append([]string{stack.Name, stack.Base, strings.Join(stack.Packages, ", "), stack.PkgManager, builtIn})
		}
		table.Render()
		return nil
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		cmdr.Info.Printf(apx.Trans("pkgmanagers.rm.info.askConfirmation"), pkgManagerName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println(apx.Trans("pkgmanagers.rm.info.aborting"))
			return nil
		}
	}

	error = pkgManager.Remove()
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(apx.Trans("pkgmanagers.rm.info.success"), pkgManagerName)
	return nil
}
