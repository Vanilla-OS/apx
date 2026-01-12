package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Pietro di Caprio <pietro@fabricators.ltd>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vanilla-os/apx/v2/core"
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

func (c *PkgManagersListCmd) Run() error {
	pkgManagers := core.ListPkgManagers()

	if !c.Json {
		pkgManagersCount := len(pkgManagers)
		if pkgManagersCount == 0 {
			Apx.Log.Term.Info().Msg(Apx.LC.Get("pkgmanagers.list.info.noPkgManagers"))
			return nil
		}

		Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.list.info.foundPkgManagers"), pkgManagersCount)

		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{Apx.LC.Get("pkgmanagers.labels.name"), Apx.LC.Get("pkgmanagers.labels.builtIn")})

		for _, stack := range pkgManagers {
			builtIn := Apx.LC.Get("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = Apx.LC.Get("apx.terminal.yes")
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

func (c *PkgManagersShowCmd) Run() error {
	args := c.Args
	if len(args) == 0 {
		return fmt.Errorf("package manager name required")
	}
	pkgManagerName := args[0]
	pkgManager, err := core.LoadPkgManager(pkgManagerName)
	if err != nil {
		Apx.Log.Term.Error().Msg(err.Error())
		return nil
	}

	table := core.CreateApxTable(os.Stdout)
	table.Append([]string{Apx.LC.Get("pkgmanagers.labels.name"), pkgManager.Name})
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

func (c *PkgManagersNewCmd) Run() error {
	reader := bufio.NewReader(os.Stdin)

	if c.Name == "" {
		if c.NoPrompt {
			Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.new.error.noName"))
			return nil
		}

		Apx.Log.Term.Info().Msg(Apx.LC.Get("pkgmanagers.new.info.askName"))
		c.Name, _ = reader.ReadString('\n')
		c.Name = strings.ReplaceAll(c.Name, "\n", "")
		c.Name = strings.ReplaceAll(c.Name, " ", "")
		if c.Name == "" {
			Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.new.error.emptyName"))
			return nil
		}
	}

	if !c.NeedSudo && !c.NoPrompt {
		validChoice := false
		for !validChoice {
			choice, err := Apx.CLI.ConfirmAction(
				Apx.LC.Get("pkgmanagers.new.info.askSudo"),
				"y", "N",
				false,
			)
			if err != nil {
				return err
			}
			c.NeedSudo = choice
			validChoice = true
		}
	}

	cmdMap := map[string]*string{
		PkgManagerCmdAutoRemove: &c.AutoRemove,
		PkgManagerCmdClean:      &c.Clean,
		PkgManagerCmdInstall:    &c.Install,
		PkgManagerCmdList:       &c.List,
		PkgManagerCmdPurge:      &c.Purge,
		PkgManagerCmdRemove:     &c.Remove,
		PkgManagerCmdSearch:     &c.Search,
		PkgManagerCmdShow:       &c.Show,
		PkgManagerCmdUpdate:     &c.Update,
		PkgManagerCmdUpgrade:    &c.Upgrade,
	}

	for _, cmdName := range PkgManagerCmdSetOrder {
		cmd := cmdMap[cmdName]
		if *cmd == "" {
			if c.NoPrompt {
				Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.new.error.noCommand"), cmdName)
				return nil
			}
			if cmdName == PkgManagerCmdPurge || cmdName == PkgManagerCmdAutoRemove {
				Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.new.info.askCommandWithDefault"), cmdName, c.Remove)
				*cmd, _ = reader.ReadString('\n')
				*cmd = strings.ReplaceAll(*cmd, "\n", "")
				if *cmd == "" {
					*cmd = c.Remove
				}
				continue
			}

			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.new.info.askCommand"), cmdName)
			*cmd, _ = reader.ReadString('\n')
			*cmd = strings.ReplaceAll(*cmd, "\n", "")
			if *cmd == "" {
				Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.new.error.emptyCommand"), cmdName)
				return nil
			}
		}
	}

	if core.PkgManagerExists(c.Name) {
		if c.NoPrompt {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.new.error.alreadyExists"), c.Name)
			return nil
		}

		confirm, err := Apx.CLI.ConfirmAction(
			fmt.Sprintf(Apx.LC.Get("pkgmanagers.new.info.askOverwrite"), c.Name),
			"y", "N",
			false,
		)
		if err != nil {
			return err
		}

		if !confirm {
			Apx.Log.Term.Info().Msg(Apx.LC.Get("apx.info.aborting"))
			return nil
		}
	}

	pkgManager := core.NewPkgManager(c.Name, c.NeedSudo, c.AutoRemove, c.Clean, c.Install, c.List, c.Purge, c.Remove, c.Search, c.Show, c.Update, c.Upgrade, false)
	err := pkgManager.Save()
	if err != nil {
		Apx.Log.Term.Error().Msg(err.Error())
		return nil
	}

	Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.new.success"), c.Name)

	return nil
}

func (c *PkgManagersRmCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.rm.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	stacks := core.ListStackForPkgManager(pkgManager.Name)
	if len(stacks) > 0 {
		Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.rm.error.inUse"), len(stacks))
		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{Apx.LC.Get("pkgmanagers.labels.name"), "Base", "Packages", "PkgManager", Apx.LC.Get("pkgmanagers.labels.builtIn")})
		for _, stack := range stacks {
			builtIn := Apx.LC.Get("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = Apx.LC.Get("apx.terminal.yes")
			}
			table.Append([]string{stack.Name, stack.Base, strings.Join(stack.Packages, ", "), stack.PkgManager, builtIn})
		}
		table.Render()
		return nil
	}

	if !c.Force {
		confirm, err := Apx.CLI.ConfirmAction(
			fmt.Sprintf(Apx.LC.Get("pkgmanagers.rm.info.askConfirmation"), c.Name),
			"y", "N",
			false,
		)
		if err != nil {
			return err
		}
		if !confirm {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.rm.info.aborting"), c.Name)
			return nil
		}
	}

	error = pkgManager.Remove()
	if error != nil {
		return error
	}

	Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.rm.info.success"), c.Name)
	return nil
}

func (c *PkgManagersExportCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.export.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	if c.Output == "" {
		Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.export.error.noOutput"))
		return nil
	}

	error = pkgManager.Export(c.Output)
	if error != nil {
		return error
	}

	Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.export.info.success"), pkgManager.Name, c.Output)
	return nil
}

func (c *PkgManagersImportCmd) Run() error {
	if c.Input == "" {
		Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.import.error.noInput"))
		return nil
	}

	pkgmanager, error := core.LoadPkgManagerFromPath(c.Input)
	if error != nil {
		Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.import.error.cannotLoad"), c.Input)
	}

	error = pkgmanager.Save()
	if error != nil {
		return error
	}

	Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.import.info.success"), pkgmanager.Name)
	return nil
}

func (c *PkgManagersUpdateCmd) Run() error {
	args := c.Args
	if c.Name == "" {
		if len(args) != 1 || args[0] == "" {
			Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.update.error.noName"))
			return nil
		}

		c.Name = args[0]
	}

	pkgmanager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	if pkgmanager.BuiltIn {
		Apx.Log.Term.Error().Msg(Apx.LC.Get("pkgmanagers.update.error.builtIn"))
		os.Exit(126)
	}

	reader := bufio.NewReader(os.Stdin)

	if c.AutoRemove == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "autoRemove", pkgmanager.CmdAutoRemove)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.AutoRemove = pkgmanager.CmdAutoRemove
			} else {
				c.AutoRemove = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "autoRemove")
			return nil
		}
	}

	if c.Clean == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "clean", pkgmanager.CmdClean)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Clean = pkgmanager.CmdClean
			} else {
				c.Clean = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "clean")
			return nil
		}
	}

	if c.Install == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "install", pkgmanager.CmdInstall)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Install = pkgmanager.CmdInstall
			} else {
				c.Install = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "install")
			return nil
		}
	}

	if c.List == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "list", pkgmanager.CmdList)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.List = pkgmanager.CmdList
			} else {
				c.List = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "list")
			return nil
		}
	}

	if c.Purge == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "purge", pkgmanager.CmdPurge)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Purge = pkgmanager.CmdPurge
			} else {
				c.Purge = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "purge")
			return nil
		}
	}

	if c.Remove == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "remove", pkgmanager.CmdRemove)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Remove = pkgmanager.CmdRemove
			} else {
				c.Remove = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "remove")
			return nil
		}
	}

	if c.Search == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "search", pkgmanager.CmdSearch)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Search = pkgmanager.CmdSearch
			} else {
				c.Search = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "search")
			return nil
		}
	}

	if c.Show == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "show", pkgmanager.CmdShow)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Show = pkgmanager.CmdShow
			} else {
				c.Show = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "show")
			return nil
		}
	}

	if c.Update == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "update", pkgmanager.CmdUpdate)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Update = pkgmanager.CmdUpdate
			} else {
				c.Update = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "update")
			return nil
		}
	}

	if c.Upgrade == "" {
		if !c.NoPrompt {
			Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "upgrade", pkgmanager.CmdUpgrade)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				c.Upgrade = pkgmanager.CmdUpgrade
			} else {
				c.Upgrade = answer
			}
		} else {
			Apx.Log.Term.Error().Msgf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "upgrade")
			return nil
		}
	}

	pkgmanager.NeedSudo = c.NeedSudo
	pkgmanager.CmdAutoRemove = c.AutoRemove
	pkgmanager.CmdClean = c.Clean
	pkgmanager.CmdInstall = c.Install
	pkgmanager.CmdList = c.List
	pkgmanager.CmdPurge = c.Purge
	pkgmanager.CmdRemove = c.Remove
	pkgmanager.CmdSearch = c.Search
	pkgmanager.CmdShow = c.Show
	pkgmanager.CmdUpdate = c.Update
	pkgmanager.CmdUpgrade = c.Upgrade

	err := pkgmanager.Save()
	if err != nil {
		return err
	}

	Apx.Log.Term.Info().Msgf(Apx.LC.Get("pkgmanagers.new.success"), c.Name)

	return nil
}
