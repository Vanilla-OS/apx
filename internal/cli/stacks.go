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

	"github.com/vanilla-os/apx/v3/core"
)

func (c *StacksListCmd) Run() error {
	stacks := core.ListStacks()

	if !c.Json {
		stacksCount := len(stacks)
		if stacksCount == 0 {
			fmt.Println(Apx.LC.Get("stacks.list.info.noStacks"))
			return nil
		}

		Apx.Log.Infof(Apx.LC.Get("stacks.list.info.foundStacks"), stacksCount)

		headers := []string{Apx.LC.Get("stacks.labels.name"), "Base", Apx.LC.Get("stacks.labels.builtIn"), "Pkgs", "Pkg manager"}
		var data [][]string
		for _, stack := range stacks {
			builtIn := Apx.LC.Get("apx.terminal.no")
			if stack.BuiltIn {
				builtIn = Apx.LC.Get("apx.terminal.yes")
			}
			data = append(data, []string{stack.Name, stack.Base, builtIn, fmt.Sprintf("%d", len(stack.Packages)), stack.PkgManager})
		}
		Apx.CLI.Table(headers, data)
	} else {
		jsonStacks, err := json.MarshalIndent(stacks, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonStacks))
	}

	return nil
}

func (c *StacksShowCmd) Run() error {
	args := c.Args
	if len(args) == 0 {
		return fmt.Errorf("stack name required")
	}
	stack, error := core.LoadStack(args[0])
	if error != nil {
		return error
	}

	headers := []string{"Property", "Value"}
	data := [][]string{
		{Apx.LC.Get("stacks.labels.name"), stack.Name},
		{"Base", stack.Base},
		{"Packages", strings.Join(stack.Packages, ", ")},
		{"Package manager", stack.PkgManager},
	}
	Apx.CLI.Table(headers, data)

	return nil
}

func (c *StacksNewCmd) Run() error {

	if c.Name == "" {
		if !c.NoPrompt {
			name, err := Apx.CLI.PromptText(Apx.LC.Get("stacks.new.info.askName"), "")
			if err != nil {
				return err
			}
			c.Name = name
			if c.Name == "" {
				Apx.Log.Error(Apx.LC.Get("stacks.new.error.emptyName"))
				return nil
			}
		} else {
			Apx.Log.Error(Apx.LC.Get("stacks.new.error.noName"))
			return nil
		}
	}

	ok := core.StackExists(c.Name)
	if ok {
		if ok {
			Apx.Log.Errorf(Apx.LC.Get("stacks.new.error.alreadyExists"), c.Name)
			return nil
		}
	}

	if c.BaseImage == "" {
		if !c.NoPrompt {
			base, err := Apx.CLI.PromptText(Apx.LC.Get("stacks.new.info.askBase"), "")
			if err != nil {
				return err
			}
			c.BaseImage = base
			if c.BaseImage == "" {
				Apx.Log.Error(Apx.LC.Get("stacks.new.error.emptyBase"))
				return nil
			}
		} else {
			Apx.Log.Error(Apx.LC.Get("stacks.new.error.noBase"))
			return nil
		}
	}

	if c.PkgManager == "" {
		pkgManagers := core.ListPkgManagers()
		if len(pkgManagers) == 0 {
			Apx.Log.Error(Apx.LC.Get("stacks.new.error.noPkgManagers"))
			return nil
		}

		var options []string
		for _, manager := range pkgManagers {
			options = append(options, manager.Name)
		}

		selected, err := Apx.CLI.SelectOption(Apx.LC.Get("stacks.new.info.selectPkgManager"), options)
		if err != nil {
			return err
		}
		c.PkgManager = selected
	}

	ok = core.PkgManagerExists(c.PkgManager)
	if !ok {
		Apx.Log.Error(Apx.LC.Get("stacks.new.error.pkgManagerDoesNotExist"))
		return nil
	}

	packagesArray := strings.Fields(c.Packages)
	if len(packagesArray) == 0 && !c.NoPrompt {

		confirm, _ := Apx.CLI.ConfirmAction(
			Apx.LC.Get("stacks.new.info.noPackages"),
			"y", "N",
			false,
		)
		if confirm {
			pkgs, _ := Apx.CLI.PromptText(Apx.LC.Get("stacks.new.info.askPackages"), "")
			packagesInput := strings.TrimSpace(pkgs)
			packagesArray = strings.Fields(packagesInput)
		} else {
			packagesArray = []string{}
		}
	}

	stack := core.NewStack(c.Name, c.BaseImage, packagesArray, c.PkgManager, false)

	err := stack.Save()
	if err != nil {
		return err
	}

	Apx.Log.Infof(Apx.LC.Get("stacks.new.info.success"), c.Name)

	return nil
}

func (c *StacksUpdateCmd) Run() error {
	args := c.Args
	if c.Name == "" {
		if len(args) != 1 || args[0] == "" {
			Apx.Log.Error(Apx.LC.Get("stacks.update.error.noName"))
			return nil
		}

		c.Name = args[0]
	}

	stack, error := core.LoadStack(c.Name)
	if error != nil {
		return error
	}

	if stack.BuiltIn {
		Apx.Log.Error(Apx.LC.Get("stacks.update.error.builtIn"))
		os.Exit(126)
	}

	if c.BaseImage == "" {
		if !c.NoPrompt {
			base, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("stacks.update.info.askBase"), stack.Base), stack.Base)
			if err != nil {
				return err
			}
			c.BaseImage = base
			if c.BaseImage == "" {
				c.BaseImage = stack.Base
			}
		} else {
			Apx.Log.Error(Apx.LC.Get("stacks.update.error.noBase"))
			return nil
		}
	}

	if c.PkgManager == "" {
		if !c.NoPrompt {
			manager, err := Apx.CLI.PromptText(fmt.Sprintf(Apx.LC.Get("stacks.update.info.askPkgManager"), stack.PkgManager), stack.PkgManager)
			if err != nil {
				return err
			}
			c.PkgManager = manager
			if c.PkgManager == "" {
				c.PkgManager = stack.PkgManager
			}
		} else {
			Apx.Log.Error(Apx.LC.Get("stacks.update.error.noPkgManager"))
			return nil
		}
	}

	ok := core.PkgManagerExists(c.PkgManager)
	if !ok {
		Apx.Log.Error(Apx.LC.Get("stacks.update.error.pkgManagerDoesNotExist"))
		return nil
	}

	if len(c.Packages) > 0 {
		stack.Packages = strings.Fields(c.Packages)
	} else if !c.NoPrompt {
		msg := Apx.LC.Get("stacks.update.info.noPackages")
		if len(stack.Packages) > 0 {
			msg = Apx.LC.Get("stacks.update.info.confirmPackages") + "\n\t -" + strings.Join(stack.Packages, "\n\t - ")
		}

		confirm, _ := Apx.CLI.ConfirmAction(msg, "y", "N", false)

		packagesArray := []string{}

		if confirm {
			pkgs, _ := Apx.CLI.PromptText(Apx.LC.Get("stacks.update.info.askPackages"), "")
			packagesInput := strings.TrimSpace(pkgs)
			packagesArray = strings.Fields(packagesInput)
			stack.Packages = packagesArray
		}
	}

	stack.Base = c.BaseImage
	stack.PkgManager = c.PkgManager

	err := stack.Save()
	if err != nil {
		return err
	}

	Apx.Log.Infof(Apx.LC.Get("stacks.update.info.success"), c.Name)

	return nil
}

func (c *StacksRmCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("stacks.rm.error.noName"))
		return nil
	}

	subSystems, _ := core.ListSubsystemForStack(c.Name)
	if len(subSystems) > 0 {
		Apx.Log.Errorf(Apx.LC.Get("stacks.rm.error.inUse"), len(subSystems))
		headers := []string{Apx.LC.Get("subsystems.labels.name"), "Stack", Apx.LC.Get("subsystems.labels.status"), "Pkgs"}
		var data [][]string
		for _, subSystem := range subSystems {
			data = append(data, []string{
				subSystem.Name,
				subSystem.Stack.Name,
				subSystem.Status,
				fmt.Sprintf("%d", len(subSystem.Stack.Packages)),
			})
		}
		Apx.CLI.Table(headers, data)
		return nil
	}

	if !c.Force {
		confirm, err := Apx.CLI.ConfirmAction(
			fmt.Sprintf(Apx.LC.Get("stacks.rm.info.askConfirmation"), c.Name),
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

	stack, error := core.LoadStack(c.Name)
	if error != nil {
		return error
	}

	error = stack.Remove()
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("stacks.rm.info.success"), c.Name)
	return nil
}

func (c *StacksExportCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("stacks.export.error.noName"))
		return nil
	}

	stack, error := core.LoadStack(c.Name)
	if error != nil {
		return error
	}

	if c.Output == "" {
		Apx.Log.Error(Apx.LC.Get("stacks.export.error.noOutput"))
		return nil
	}

	error = stack.Export(c.Output)
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("stacks.export.info.success"), stack.Name, c.Output)
	return nil
}

func (c *StacksImportCmd) Run() error {
	if c.Input == "" {
		Apx.Log.Error(Apx.LC.Get("stacks.import.error.noInput"))
		return nil
	}

	stack, error := core.LoadStackFromPath(c.Input)
	if error != nil {
		Apx.Log.Errorf(Apx.LC.Get("stacks.import.error.cannotLoad"), c.Input)
	}

	error = stack.Save()
	if error != nil {
		return error
	}

	Apx.Log.Infof(Apx.LC.Get("stacks.import.info.success"), stack.Name)
	return nil
}
