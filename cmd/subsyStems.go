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
	"os"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewSubSystemsCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"subsystems",
		apx.Trans("subsystems.long"),
		apx.Trans("subsystems.short"),
		nil,
	)
	cmd.Example = "apx subsystems"

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		apx.Trans("listSubSystems.long"),
		apx.Trans("listSubSystems.short"),
		listSubSystems,
	)
	listCmd.Example = "apx subsystems list"

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		apx.Trans("newSubSystem.long"),
		apx.Trans("newSubSystem.short"),
		newSubSystem,
	)
	newCmd.Example = "apx subsystems new --name my-subsystem --stack my-stack"

	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"stack",
			"s",
			"The stack to be used for the subsystem",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			"The name of the subsystem",
			"",
		),
	)

	// Rm subcommand
	/*
		rmCmd := cmdr.NewCommand(
			"rm",
			apx.Trans("rmSubSystem.long"),
			apx.Trans("rmSubSystem.short"),
			rmSubSystem,
		)
		rmCmd.Example = "apx subsystems rm --name my-subsystem"
		rmCmd.Args = cobra.ExactArgs(1)
	*/

	// Add subcommands to subsystems
	cmd.AddCommand(listCmd)
	cmd.AddCommand(newCmd)
	// cmd.AddCommand(rmCmd, rmCmd)

	return cmd
}

func listSubSystems(cmd *cobra.Command, args []string) error {
	subSystems, err := core.ListSubSystems()
	if err != nil {
		return err
	}

	subSystemsCount := len(subSystems)
	if subSystemsCount == 0 {
		fmt.Println("No subsystems available. Create a new one with 'apx subsystems new' or contact the system administrator.")
		return nil
	}

	fmt.Printf("Found %d subsystems:\n", subSystemsCount)

	table := core.CreateApxTable(os.Stdout)
	table.SetHeader([]string{"Name", "Stack", "Status", "Pkgs"})

	for _, subSystem := range subSystems {
		table.Append([]string{
			subSystem.Name,
			subSystem.Stack.Name,
			subSystem.Status,
			fmt.Sprintf("%d", len(subSystem.Stack.Packages)),
		})
	}

	table.Render()

	return nil
}

func newSubSystem(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("stack")
	subSystemName, _ := cmd.Flags().GetString("name")

	if stackName == "" {
		cmdr.Info.Println("Please type a stack name:")
		fmt.Scanln(&stackName)
		if stackName == "" {
			cmdr.Error.Println("Stack name cannot be empty")
			return nil
		}
	}

	if subSystemName == "" {
		cmdr.Info.Println("Please type a subsystem name:")
		fmt.Scanln(&subSystemName)
		if subSystemName == "" {
			cmdr.Error.Println("Subsystem name cannot be empty")
			return nil
		}
	}

	stack, err := core.LoadStack(stackName)
	if err != nil {
		return err
	}

	subSystem, err := core.NewSubSystem(subSystemName, stack)
	if err != nil {
		return err
	}

	err = subSystem.Create()
	if err != nil {
		return err
	}

	cmdr.Success.Printf("Subsystem %s created successfully!\n", subSystemName)

	return nil
}

/*
func rmSubSystem(cmd *cobra.Command, args []string) error {
	subSystemName := args[0]

	subSystem, err := core.LoadSubSystem(subSystemName)
	if err != nil {
		return err
	}

	err = subSystem.Remove()
	if err != nil {
		return err
	}

	cmdr.Success.Printf("Subsystem %s removed successfully!\n", subSystemName)

	return nil
}
*/
