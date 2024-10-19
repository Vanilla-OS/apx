package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/v2/core"
	"github.com/vanilla-os/orchid/cmdr"
)

const (
	PkgManagerCmdAutoRemove = "autoRemove"
	PkgManagerCmdClean      = "clean"
	PkgManagerCmdInstall    = "install"
	PkgManagerCmdList       = "list"
	PkgManagerCmdPurge      = "purge"
	PkgManagerCmdRemove     = "remove"
	PkgManagerCmdSearch     = "search"
	PkgManagerCmdShow       = "show"
	PkgManagerCmdUpdate     = "update"
	PkgManagerCmdUpgrade    = "upgrade"
)

var PkgManagerCmdSetOrder = []string{
	PkgManagerCmdInstall,
	PkgManagerCmdUpdate,
	PkgManagerCmdRemove,
	PkgManagerCmdPurge,
	PkgManagerCmdAutoRemove,
	PkgManagerCmdClean,
	PkgManagerCmdList,
	PkgManagerCmdSearch,
	PkgManagerCmdShow,
	PkgManagerCmdUpgrade,
}

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
			"no-prompt",
			"y",
			apx.Trans("pkgmanagers.new.options.noPrompt.description"),
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

	// Export subcommand
	exportCmd := cmdr.NewCommand(
		"export",
		apx.Trans("pkgmanagers.export.description"),
		apx.Trans("pkgmanagers.export.description"),
		exportPkgmanager,
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("pkgmanagers.export.options.name.description"),
			"",
		),
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"output",
			"o",
			apx.Trans("stacks.export.options.output.description"),
			"",
		),
	)

	// Import subcommand
	importCmd := cmdr.NewCommand(
		"import",
		apx.Trans("pkgmanagers.import.description"),
		apx.Trans("pkgmanagers.import.description"),
		importPkgmanager,
	)
	importCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"input",
			"i",
			apx.Trans("pkgmanagers.import.options.input.description"),
			"",
		),
	)

	// Update subcommand
	updateCmd := cmdr.NewCommand(
		"update",
		apx.Trans("pkgmanagers.update.description"),
		apx.Trans("pkgmanagers.update.description"),
		updatePkgManager,
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("pkgmanagers.new.options.name.description"),
			"",
		),
	)
	updateCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"need-sudo",
			"S",
			apx.Trans("pkgmanagers.new.options.needSudo.description"),
			false,
		),
	)
	updateCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"no-prompt",
			"y",
			apx.Trans("pkgmanagers.new.options.noPrompt.description"),
			false,
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"autoremove",
			"a",
			apx.Trans("pkgmanagers.new.options.autoremove.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"clean",
			"c",
			apx.Trans("pkgmanagers.new.options.clean.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"install",
			"i",
			apx.Trans("pkgmanagers.new.options.install.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"list",
			"l",
			apx.Trans("pkgmanagers.new.options.list.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"purge",
			"p",
			apx.Trans("pkgmanagers.new.options.purge.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"remove",
			"r",
			apx.Trans("pkgmanagers.new.options.remove.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"search",
			"s",
			apx.Trans("pkgmanagers.new.options.search.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"show",
			"w",
			apx.Trans("pkgmanagers.new.options.show.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"update",
			"u",
			apx.Trans("pkgmanagers.new.options.update.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"upgrade",
			"U",
			apx.Trans("pkgmanagers.new.options.upgrade.description"),
			"",
		),
	)

	// Add subcommands to pkgmanagers
	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(exportCmd)
	cmd.AddCommand(importCmd)
	cmd.AddCommand(updateCmd)

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

		cmdr.Info.Printfln(apx.Trans("pkgmanagers.list.info.foundPkgManagers"), pkgManagersCount)

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
	noPrompt, _ := cmd.Flags().GetBool("no-prompt")
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
		if noPrompt {
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

	if !needSudo && !noPrompt {
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
		PkgManagerCmdAutoRemove: &autoRemove,
		PkgManagerCmdClean:      &clean,
		PkgManagerCmdInstall:    &install,
		PkgManagerCmdList:       &list,
		PkgManagerCmdPurge:      &purge,
		PkgManagerCmdRemove:     &remove,
		PkgManagerCmdSearch:     &search,
		PkgManagerCmdShow:       &show,
		PkgManagerCmdUpdate:     &update,
		PkgManagerCmdUpgrade:    &upgrade,
	}

	for _, cmdName := range PkgManagerCmdSetOrder {
		cmd := cmdMap[cmdName]
		if *cmd == "" {
			if noPrompt {
				cmdr.Error.Printf(apx.Trans("pkgmanagers.new.error.noCommand"), cmdName)
				return nil
			}
			if cmdName == PkgManagerCmdPurge || cmdName == PkgManagerCmdAutoRemove {
				cmdr.Info.Printfln(apx.Trans("pkgmanagers.new.info.askCommandWithDefault"), cmdName, remove)
				*cmd, _ = reader.ReadString('\n')
				*cmd = strings.ReplaceAll(*cmd, "\n", "")
				if *cmd == "" {
					*cmd = remove
				}
				continue
			}

			cmdr.Info.Printfln(apx.Trans("pkgmanagers.new.info.askCommand"), cmdName)
			*cmd, _ = reader.ReadString('\n')
			*cmd = strings.ReplaceAll(*cmd, "\n", "")
			if *cmd == "" {
				cmdr.Error.Printf(apx.Trans("pkgmanagers.new.error.emptyCommand"), cmdName)
				return nil
			}
		}
	}

	if core.PkgManagerExists(name) {
		if noPrompt {
			cmdr.Error.Println(apx.Trans("pkgmanagers.new.error.alreadyExists"), name)
			return nil
		}

		cmdr.Info.Printfln(apx.Trans("pkgmanagers.new.info.askOverwrite"), name)
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

	cmdr.Success.Printfln(apx.Trans("pkgmanagers.new.success"), name)

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
		reader := bufio.NewReader(os.Stdin)
		validChoice := false
		for !validChoice {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.rm.info.askConfirmation")+` [y/N]`, pkgManagerName)
			answer, _ := reader.ReadString('\n')
			if answer == "\n" {
				answer = "n\n"
			}
			answer = strings.ToLower(strings.ReplaceAll(answer, " ", ""))
			switch answer {
			case "y\n":
				validChoice = true
				force = true
			case "n\n":
				validChoice = true
			default:
				cmdr.Warning.Println(apx.Trans("apx.errors.invalidChoice"))
			}
		}
	}

	if !force {
		cmdr.Info.Printfln(apx.Trans("pkgmanagers.rm.info.aborting"), pkgManagerName)
		return nil
	}

	error = pkgManager.Remove()
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(apx.Trans("pkgmanagers.rm.info.success"), pkgManagerName)
	return nil
}

func exportPkgmanager(cmd *cobra.Command, args []string) error {
	pkgManagerName, _ := cmd.Flags().GetString("name")
	if pkgManagerName == "" {
		cmdr.Error.Println(apx.Trans("pkgmanagers.export.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(pkgManagerName)
	if error != nil {
		return error
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		cmdr.Error.Println(apx.Trans("pkgmanagers.export.error.noOutput"))
		return nil
	}

	error = pkgManager.Export(output)
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(apx.Trans("pkgmanagers.export.info.success"), pkgManager.Name, output)
	return nil
}

func importPkgmanager(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	if input == "" {
		cmdr.Error.Println(apx.Trans("pkgmanagers.import.error.noInput"))
		return nil
	}

	pkgmanager, error := core.LoadPkgManagerFromPath(input)
	if error != nil {
		cmdr.Error.Printf(apx.Trans("pkgmanagers.import.error.cannotLoad"), input)
	}

	error = pkgmanager.Save()
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(apx.Trans("pkgmanagers.import.info.success"), pkgmanager.Name)
	return nil
}

func updatePkgManager(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	needSudo, _ := cmd.Flags().GetBool("need-sudo")
	noPrompt, _ := cmd.Flags().GetBool("no-prompt")
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

	if name == "" {
		if len(args) != 1 || args[0] == "" {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.noName"))
			return nil
		}

		cmd.Flags().Set("name", args[0])
		name = args[0]
	}

	pkgmanager, error := core.LoadPkgManager(name)
	if error != nil {
		return error
	}

	if pkgmanager.BuiltIn {
		cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.builtIn"))
		os.Exit(126)
	}

	reader := bufio.NewReader(os.Stdin)

	if autoRemove == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "autoRemove", pkgmanager.CmdAutoRemove)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				autoRemove = pkgmanager.CmdAutoRemove
			} else {
				autoRemove = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "autoRemove")
			return nil
		}
	}

	if clean == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "clean", pkgmanager.CmdClean)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				clean = pkgmanager.CmdClean
			} else {
				clean = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "clean")
			return nil
		}
	}

	if install == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "install", pkgmanager.CmdInstall)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				install = pkgmanager.CmdInstall
			} else {
				install = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "install")
			return nil
		}
	}

	if list == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "list", pkgmanager.CmdList)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				list = pkgmanager.CmdList
			} else {
				list = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "list")
			return nil
		}
	}

	if purge == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "purge", pkgmanager.CmdPurge)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				purge = pkgmanager.CmdPurge
			} else {
				purge = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "purge")
			return nil
		}
	}

	if remove == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "remove", pkgmanager.CmdRemove)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				remove = pkgmanager.CmdRemove
			} else {
				remove = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "remove")
			return nil
		}
	}

	if search == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "search", pkgmanager.CmdSearch)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				search = pkgmanager.CmdSearch
			} else {
				search = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "search")
			return nil
		}
	}

	if show == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "show", pkgmanager.CmdShow)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				show = pkgmanager.CmdShow
			} else {
				show = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "show")
			return nil
		}
	}

	if update == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "update", pkgmanager.CmdUpdate)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				update = pkgmanager.CmdUpdate
			} else {
				update = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "update")
			return nil
		}
	}

	if upgrade == "" {
		if !noPrompt {
			cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.askNewCommand"), "upgrade", pkgmanager.CmdUpgrade)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				upgrade = pkgmanager.CmdUpgrade
			} else {
				upgrade = answer
			}
		} else {
			cmdr.Error.Println(apx.Trans("pkgmanagers.update.error.missingCommand"), "upgrade")
			return nil
		}
	}

	pkgmanager.NeedSudo = needSudo
	pkgmanager.CmdAutoRemove = autoRemove
	pkgmanager.CmdClean = clean
	pkgmanager.CmdInstall = install
	pkgmanager.CmdList = list
	pkgmanager.CmdPurge = purge
	pkgmanager.CmdRemove = remove
	pkgmanager.CmdSearch = search
	pkgmanager.CmdShow = show
	pkgmanager.CmdUpdate = update
	pkgmanager.CmdUpgrade = upgrade

	err := pkgmanager.Save()
	if err != nil {
		return err
	}

	cmdr.Info.Printfln(apx.Trans("pkgmanagers.update.info.success"), name)

	return nil
}
