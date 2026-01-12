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
			Apx.Log.Info(Apx.LC.Get("pkgmanagers.list.info.noPkgManagers"))
			return nil
		}

		Apx.Log.Infof(Apx.LC.Get("pkgmanagers.list.info.foundPkgManagers"), pkgManagersCount)

		headers := []string{Apx.LC.Get("pkgmanagers.labels.name"), Apx.LC.Get("pkgmanagers.labels.builtIn")}
		var data [][]string
		for _, stack := range pkgManagers {
			builtIn := Apx.LC.Get("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = Apx.LC.Get("apx.terminal.yes")
			}
			data = append(data, []string{stack.Name, builtIn})
		}

		err := Apx.CLI.Table(headers, data)
		if err != nil {
			return err
		}
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
		Apx.Log.Error(err.Error())
		return nil
	}

	headers := []string{"Property", "Value"}
	data := [][]string{
		{Apx.LC.Get("pkgmanagers.labels.name"), pkgManager.Name},
		{"NeedSudo", fmt.Sprintf("%t", pkgManager.NeedSudo)},
		{"AutoRemove", pkgManager.CmdAutoRemove},
		{"Clean", pkgManager.CmdClean},
		{"Install", pkgManager.CmdInstall},
		{"List", pkgManager.CmdList},
		{"Purge", pkgManager.CmdPurge},
		{"Remove", pkgManager.CmdRemove},
		{"Search", pkgManager.CmdSearch},
		{"Show", pkgManager.CmdShow},
		{"Update", pkgManager.CmdUpdate},
		{"Upgrade", pkgManager.CmdUpgrade},
	}

	err = Apx.CLI.Table(headers, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *PkgManagersNewCmd) Run() error {

	if c.Name == "" {
		if c.NoPrompt {
			Apx.Log.Error(Apx.LC.Get("pkgmanagers.new.error.noName"))
			return nil
		}

		name, err := Apx.CLI.PromptText(Apx.LC.Get("pkgmanagers.new.info.askName"), "")
		if err != nil {
			return err
		}
		c.Name = name
		c.Name = strings.ReplaceAll(c.Name, "\n", "")
		c.Name = strings.ReplaceAll(c.Name, " ", "")
		if c.Name == "" {
			Apx.Log.Error(Apx.LC.Get("pkgmanagers.new.error.emptyName"))
			return nil
		}
	}

	if !c.NeedSudo && !c.NoPrompt {
		choice, err := Apx.CLI.ConfirmAction(
			Apx.LC.Get("pkgmanagers.new.info.askSudo"),
			"y", "N",
			false,
		)
		if err != nil {
			return err
		}
		c.NeedSudo = choice
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
				Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.new.error.noCommand"), cmdName)
				return nil
			}
			if cmdName == PkgManagerCmdPurge || cmdName == PkgManagerCmdAutoRemove {
				answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.new.info.askCommandWithDefault"), cmdName, c.Remove), c.Remove)
				if err != nil {
					return err
				}
				*cmd = strings.TrimSpace(answer)
				if *cmd == "" {
					*cmd = c.Remove
				}
				continue
			}

			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.new.info.askCommand"), cmdName), "")
			if err != nil {
				return err
			}
			*cmd = strings.TrimSpace(answer)
			if *cmd == "" {
				Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.new.error.emptyCommand"), cmdName)
				return nil
			}
		}
	}

	if core.PkgManagerExists(c.Name) {
		if c.NoPrompt {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.new.error.alreadyExists"), c.Name)
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
			Apx.Log.Info(Apx.LC.Get("apx.info.aborting"))
			return nil
		}
	}

	pkgManager := core.NewPkgManager(c.Name, c.NeedSudo, c.AutoRemove, c.Clean, c.Install, c.List, c.Purge, c.Remove, c.Search, c.Show, c.Update, c.Upgrade, false)
	err := pkgManager.Save()
	if err != nil {
		Apx.Log.Error(err.Error())
		return nil
	}

	Apx.Log.Infof(Apx.LC.Get("pkgmanagers.new.success"), c.Name)

	return nil
}

func (c *PkgManagersRmCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("pkgmanagers.rm.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	stacks := core.ListStackForPkgManager(pkgManager.Name)
	if len(stacks) > 0 {
		Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.rm.error.inUse"), len(stacks))
		headers := []string{Apx.LC.Get("pkgmanagers.labels.name"), "Base", "Packages", "PkgManager", Apx.LC.Get("pkgmanagers.labels.builtIn")}
		var data [][]string
		for _, stack := range stacks {
			builtIn := Apx.LC.Get("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = Apx.LC.Get("apx.terminal.yes")
			}
			data = append(data, []string{stack.Name, stack.Base, strings.Join(stack.Packages, ", "), stack.PkgManager, builtIn})
		}
		Apx.CLI.Table(headers, data)
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
			Apx.Log.Infof(Apx.LC.Get("pkgmanagers.rm.info.aborting"), c.Name)
			return nil
		}
	}

	error = pkgManager.Remove()
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("pkgmanagers.rm.info.success"), c.Name)
	return nil
}

func (c *PkgManagersExportCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("pkgmanagers.export.error.noName"))
		return nil
	}

	pkgManager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	if c.Output == "" {
		Apx.Log.Error(Apx.LC.Get("pkgmanagers.export.error.noOutput"))
		return nil
	}

	error = pkgManager.Export(c.Output)
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("pkgmanagers.export.info.success"), pkgManager.Name, c.Output)
	return nil
}

func (c *PkgManagersImportCmd) Run() error {
	if c.Input == "" {
		Apx.Log.Error(Apx.LC.Get("pkgmanagers.import.error.noInput"))
		return nil
	}

	pkgmanager, error := core.LoadPkgManagerFromPath(c.Input)
	if error != nil {
		Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.import.error.cannotLoad"), c.Input)
	}

	error = pkgmanager.Save()
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("pkgmanagers.import.info.success"), pkgmanager.Name)
	return nil
}

func (c *PkgManagersUpdateCmd) Run() error {
	args := c.Args
	if c.Name == "" {
		if len(args) != 1 || args[0] == "" {
			Apx.Log.Error(Apx.LC.Get("pkgmanagers.update.error.noName"))
			return nil
		}

		c.Name = args[0]
	}

	pkgmanager, error := core.LoadPkgManager(c.Name)
	if error != nil {
		return error
	}

	if pkgmanager.BuiltIn {
		Apx.Log.Error(Apx.LC.Get("pkgmanagers.update.error.builtIn"))
		os.Exit(126)
	}

	if c.AutoRemove == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "autoRemove", pkgmanager.CmdAutoRemove), pkgmanager.CmdAutoRemove)
			if err != nil {
				return err
			}
			c.AutoRemove = strings.TrimSpace(answer)
			if c.AutoRemove == "" {
				c.AutoRemove = pkgmanager.CmdAutoRemove
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "autoRemove")
			return nil
		}
	}

	if c.Clean == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "clean", pkgmanager.CmdClean), pkgmanager.CmdClean)
			if err != nil {
				return err
			}
			c.Clean = strings.TrimSpace(answer)
			if c.Clean == "" {
				c.Clean = pkgmanager.CmdClean
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "clean")
			return nil
		}
	}

	if c.Install == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "install", pkgmanager.CmdInstall), pkgmanager.CmdInstall)
			if err != nil {
				return err
			}
			c.Install = strings.TrimSpace(answer)
			if c.Install == "" {
				c.Install = pkgmanager.CmdInstall
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "install")
			return nil
		}
	}

	if c.List == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "list", pkgmanager.CmdList), pkgmanager.CmdList)
			if err != nil {
				return err
			}
			c.List = strings.TrimSpace(answer)
			if c.List == "" {
				c.List = pkgmanager.CmdList
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "list")
			return nil
		}
	}

	if c.Purge == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "purge", pkgmanager.CmdPurge), pkgmanager.CmdPurge)
			if err != nil {
				return err
			}
			c.Purge = strings.TrimSpace(answer)
			if c.Purge == "" {
				c.Purge = pkgmanager.CmdPurge
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "purge")
			return nil
		}
	}

	if c.Remove == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "remove", pkgmanager.CmdRemove), pkgmanager.CmdRemove)
			if err != nil {
				return err
			}
			c.Remove = strings.TrimSpace(answer)
			if c.Remove == "" {
				c.Remove = pkgmanager.CmdRemove
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "remove")
			return nil
		}
	}

	if c.Search == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "search", pkgmanager.CmdSearch), pkgmanager.CmdSearch)
			if err != nil {
				return err
			}
			c.Search = strings.TrimSpace(answer)
			if c.Search == "" {
				c.Search = pkgmanager.CmdSearch
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "search")
			return nil
		}
	}

	if c.Show == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "show", pkgmanager.CmdShow), pkgmanager.CmdShow)
			if err != nil {
				return err
			}
			c.Show = strings.TrimSpace(answer)
			if c.Show == "" {
				c.Show = pkgmanager.CmdShow
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "show")
			return nil
		}
	}

	if c.Update == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "update", pkgmanager.CmdUpdate), pkgmanager.CmdUpdate)
			if err != nil {
				return err
			}
			c.Update = strings.TrimSpace(answer)
			if c.Update == "" {
				c.Update = pkgmanager.CmdUpdate
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "update")
			return nil
		}
	}

	if c.Upgrade == "" {
		if !c.NoPrompt {
			answer, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("pkgmanagers.update.info.askNewCommand"), "upgrade", pkgmanager.CmdUpgrade), pkgmanager.CmdUpgrade)
			if err != nil {
				return err
			}
			c.Upgrade = strings.TrimSpace(answer)
			if c.Upgrade == "" {
				c.Upgrade = pkgmanager.CmdUpgrade
			}
		} else {
			Apx.Log.Errorf(Apx.LC.Get("pkgmanagers.update.error.missingCommand"), "upgrade")
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
