package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewSubSystemsCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"subsystems",
		apx.Trans("subsystems.description"),
		apx.Trans("subsystems.description"),
		nil,
	)

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		apx.Trans("subsystems.list.description"),
		apx.Trans("subsystems.list.description"),
		listSubSystems,
	)

	listCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"json",
			"j",
			apx.Trans("subsystems.list.options.json"),
			false,
		),
	)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		apx.Trans("subsystems.new.description"),
		apx.Trans("subsystems.new.description"),
		newSubSystem,
	)

	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"stack",
			"s",
			apx.Trans("subsystems.new.options.stack"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("subsystems.new.options.name"),
			"",
		),
	)

	// Rm subcommand
	rmCmd := cmdr.NewCommand(
		"rm",
		apx.Trans("subsystems.rm.description"),
		apx.Trans("subsystems.rm.description"),
		rmSubSystem,
	)

	rmCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("subsystems.rm.options.name"),
			"",
		),
	)
	rmCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			apx.Trans("subsystems.rm.options.force"),
			false,
		),
	)

	// Reset subcommand
	resetCmd := cmdr.NewCommand(
		"reset",
		apx.Trans("subsystems.reset.description"),
		apx.Trans("subsystems.reset.description"),
		resetSubSystem,
	)

	resetCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			apx.Trans("subsystems.reset.options.name"),
			"",
		),
	)
	resetCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			apx.Trans("subsystems.reset.options.force"),
			false,
		),
	)

	// Add subcommands to subsystems
	cmd.AddCommand(listCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(resetCmd)

	return cmd
}

func listSubSystems(cmd *cobra.Command, args []string) error {
	jsonFlag, _ := cmd.Flags().GetBool("json")

	subSystems, err := core.ListSubSystems()
	if err != nil {
		return err
	}

	if !jsonFlag {
		subSystemsCount := len(subSystems)
		if subSystemsCount == 0 {
			cmdr.Info.Println(apx.Trans("subsystems.list.info.noSubsystems"))
			return nil
		}

		fmt.Printf(apx.Trans("subsystems.list.info.foundSubsystems"), subSystemsCount)

		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{apx.Trans("subsystems.labels.name"), "Stack", apx.Trans("subsystems.labels.status"), "Pkgs"})

		for _, subSystem := range subSystems {
			table.Append([]string{
				subSystem.Name,
				subSystem.Stack.Name,
				subSystem.Status,
				fmt.Sprintf("%d", len(subSystem.Stack.Packages)),
			})
		}

		table.Render()
	} else {
		jsonSubSystems, err := json.MarshalIndent(subSystems, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonSubSystems))
	}

	return nil
}

func newSubSystem(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("stack")
	subSystemName, _ := cmd.Flags().GetString("name")

	stacks := core.ListStacks()
	if len(stacks) == 0 {
		cmdr.Error.Println(apx.Trans("subsystems.new.error.noStacks"))
		return nil
	}

	if subSystemName == "" {
		cmdr.Info.Println(apx.Trans("subsystems.new.info.askName"))
		fmt.Scanln(&subSystemName)
		if subSystemName == "" {
			cmdr.Error.Println(apx.Trans("apx.error.noName"))
			return nil
		}
	}

	if stackName == "" {
		cmdr.Info.Println(apx.Trans("subsystems.new.info.askStack"))
		for i, stack := range stacks {
			fmt.Printf("%d. %s\n", i+1, stack.Name)
		}
		fmt.Printf(apx.Trans("subsystems.new.info.selectStack"), len(stacks))

		var stackIndex int
		_, err := fmt.Scanln(&stackIndex)
		if err != nil {
			cmdr.Error.Println(apx.Trans("apx.error.invalidInput"))
			return nil
		}

		if stackIndex < 1 || stackIndex > len(stacks) {
			cmdr.Error.Println(apx.Trans("apx.error.invalidInput"))
			return nil
		}

		stackName = stacks[stackIndex-1].Name
	}

	checkSubSystem, err := core.LoadSubSystem(subSystemName)
	if err == nil {
		cmdr.Error.Printf(apx.Trans("subsystems.new.error.alreadyExists"), checkSubSystem.Name)
		return nil
	}

	stack, err := core.LoadStack(stackName)
	if err != nil {
		return err
	}

	subSystem, err := core.NewSubSystem(subSystemName, stack)
	if err != nil {
		return err
	}

	cmdr.Info.Printf(apx.Trans("subsystems.new.info.creatingSubsystem"), subSystemName, stackName)
	err = subSystem.Create()
	if err != nil {
		return err
	}

	cmdr.Success.Printf(apx.Trans("subsystems.new.info.success"), subSystemName)

	return nil
}

func rmSubSystem(cmd *cobra.Command, args []string) error {
	subSystemName, _ := cmd.Flags().GetString("name")
	forceFlag, _ := cmd.Flags().GetBool("force")

	if subSystemName == "" {
		cmdr.Error.Println(apx.Trans("subsystems.rm.error.noName"))
		return nil
	}

	if !forceFlag {
		cmdr.Info.Printf(apx.Trans("subsystems.rm.info.askConfirmation"), subSystemName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println(apx.Trans("apx.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(subSystemName)
	if err != nil {
		return err
	}

	err = subSystem.Remove()
	if err != nil {
		return err
	}

	cmdr.Success.Printf(apx.Trans("subsystems.rm.info.success"), subSystemName)

	return nil
}

func resetSubSystem(cmd *cobra.Command, args []string) error {
	subSystemName, _ := cmd.Flags().GetString("name")
	if subSystemName == "" {
		cmdr.Error.Println(apx.Trans("subsystems.reset.error.noName"))
		return nil
	}

	forceFlag, _ := cmd.Flags().GetBool("force")

	if !forceFlag {
		cmdr.Info.Printf(apx.Trans("subsystems.reset.info.askConfirmation"), subSystemName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println(apx.Trans("apx.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(subSystemName)
	if err != nil {
		return err
	}

	err = subSystem.Reset()
	if err != nil {
		return err
	}

	cmdr.Success.Printf(apx.Trans("subsystems.reset.info.success"), subSystemName)

	return nil
}
