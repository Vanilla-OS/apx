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

	"github.com/vanilla-os/apx/v3/core"
)

func (c *SubsystemsListCmd) Run() error {
	subSystems, err := core.ListSubSystems(false, false)
	if err != nil {
		return err
	}

	if !c.Json {
		subSystemsCount := len(subSystems)
		if subSystemsCount == 0 {
			Apx.Log.Info(Apx.LC.Get("subsystems.list.info.noSubsystems"))
			return nil
		}

		Apx.Log.Infof(Apx.LC.Get("subsystems.list.info.foundSubsystems"), subSystemsCount)

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

		err := Apx.CLI.Table(headers, data)
		if err != nil {
			return err
		}
	} else {
		jsonSubSystems, err := json.MarshalIndent(subSystems, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonSubSystems))
	}

	return nil
}

func (c *SubsystemsNewCmd) Run() error {

	stacks := core.ListStacks()
	if len(stacks) == 0 {
		Apx.Log.Error(Apx.LC.Get("subsystems.new.error.noStacks"))
		return nil
	}

	if c.Name == "" {
		name, err := Apx.CLI.PromptText(Apx.LC.Get("subsystems.new.info.askName"), "")
		if err != nil {
			return err
		}
		c.Name = name
		if c.Name == "" {
			Apx.Log.Error(Apx.LC.Get("subsystems.new.error.emptyName"))
			return nil
		}
	}

	if c.Stack == "" {
		var options []string
		for _, stack := range stacks {
			options = append(options, stack.Name)
		}

		selected, err := Apx.CLI.SelectOption(Apx.LC.Get("subsystems.new.info.selectStack"), options)
		if err != nil {
			return err
		}
		c.Stack = selected
	}

	checkSubSystem, err := core.LoadSubSystem(c.Name, false)
	if err == nil {
		Apx.Log.Errorf(Apx.LC.Get("subsystems.new.error.alreadyExists"), checkSubSystem.Name)
		return nil
	}

	// Checking if name conflicts with existing commands.
	// In SDK declarative approach, we might need a way to check root commands.
	// Apx.CLI.GetRoot() returns the struct.
	// We can iterate simple tags or just use known conflicts.
	// For now we skip detailed conflict check as we are rebuilding.
	// Or we could check if c.Name is a valid command name if feasible.

	stack, err := core.LoadStack(c.Stack)
	if err != nil {
		return err
	}

	subSystem, err := core.NewSubSystem(c.Name, stack, c.Home, c.Init, false, false, false, true, "")
	if err != nil {
		return err
	}

	spinner := Apx.CLI.StartSpinner(fmt.Sprintf(Apx.LC.Get("subsystems.new.info.creatingSubsystem"), c.Name, c.Stack))

	err = subSystem.Create()
	if err != nil {
		spinner.Stop()
		return err
	}

	spinner.Stop()
	Apx.Log.Infof(Apx.LC.Get("subsystems.new.info.success"), c.Name)

	return nil
}

func (c *SubsystemsRmCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("subsystems.rm.error.noName"))
		return nil
	}

	if !c.Force {
		confirm, _ := Apx.CLI.ConfirmAction(
			fmt.Sprintf(Apx.LC.Get("subsystems.rm.info.askConfirmation"), c.Name),
			"y", "N",
			false,
		)
		if !confirm {
			Apx.Log.Term.Info().Msg(Apx.LC.Get("apx.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	err = subSystem.Remove()
	if err != nil {
		return err
	}

	Apx.Log.Infof(Apx.LC.Get("subsystems.rm.info.success"), c.Name)

	return nil
}

func (c *SubsystemsResetCmd) Run() error {
	if c.Name == "" {
		Apx.Log.Error(Apx.LC.Get("subsystems.reset.error.noName"))
		return nil
	}

	if !c.Force {
		confirm, _ := Apx.CLI.ConfirmAction(
			fmt.Sprintf(Apx.LC.Get("subsystems.reset.info.askConfirmation"), c.Name),
			"y", "N",
			false,
		)
		if !confirm {
			Apx.Log.Info(Apx.LC.Get("apx.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	err = subSystem.Reset()
	if err != nil {
		return err
	}

	Apx.Log.Infof(Apx.LC.Get("subsystems.reset.info.success"), c.Name)

	return nil
}
