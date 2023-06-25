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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewStacksCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"stacks",
		apx.Trans("stacks"),
		apx.Trans("stacks"),
		nil,
	)
	cmd.Example = "apx stacks"

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		apx.Trans("listStacks"),
		apx.Trans("listStacks"),
		listStacks,
	)
	listCmd.Example = "apx stacks list"

	// Show subcommand
	showCmd := cmdr.NewCommand(
		"show",
		apx.Trans("showStack"),
		apx.Trans("showStack"),
		showStack,
	)
	showCmd.Example = "apx stacks show my-stack"
	showCmd.Args = cobra.MinimumNArgs(1)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		apx.Trans("newStack"),
		apx.Trans("newStack"),
		newStack,
	)
	newCmd.Example = "apx stacks new --name my-stack --base vanillaos:pico --packages nano,git --pkg-manager apt"

	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"assume-yes",
			"y",
			"Assume yes; assume that the answer to any question which would be asked is yes",
			false,
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			"Name of the stack",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"base",
			"b",
			"Base image to use for the stack",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"packages",
			"p",
			"List of packages to install in the stack",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"pkg-manager",
			"k",
			"Package manager to use in the subsystem",
			"",
		),
	)

	// Update subcommand
	updateCmd := cmdr.NewCommand(
		"update",
		apx.Trans("updateStack"),
		apx.Trans("updateStack"),
		updateStack,
	)
	updateCmd.Example = "apx stacks update --base vanillaos:pico --packages nano,git --pkg-manager apt"

	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"base",
			"b",
			"Base image to use for the stack",
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"packages",
			"p",
			"List of packages to install in the stack",
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"pkg-manager",
			"k",
			"Package manager to use in the subsystem",
			"",
		),
	)

	// Rm subcommand
	rmStackCmd := cmdr.NewCommand(
		"rm",
		apx.Trans("rmStack"),
		apx.Trans("rmStack"),
		removeStack,
	)
	rmStackCmd.Example = "apx stacks rm my-stack"

	rmStackCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			"Name of the stack",
			"",
		),
	)
	rmStackCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			"Force the removal of the stack",
			false,
		),
	)

	// Export subcommand
	exportCmd := cmdr.NewCommand(
		"export",
		apx.Trans("exportStack"),
		apx.Trans("exportStack"),
		exportStack,
	)
	exportCmd.Example = "apx stacks export my-stack"

	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			"Name of the stack",
			"",
		),
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"output",
			"o",
			"Output location",
			"",
		),
	)

	// Import subcommand
	importCmd := cmdr.NewCommand(
		"import",
		apx.Trans("importStack"),
		apx.Trans("importStack"),
		importStack,
	)
	importCmd.Example = "apx stacks import my-stack"

	importCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"input",
			"i",
			"Input file",
			"",
		),
	)

	// Add subcommands to stacks
	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(updateCmd)
	cmd.AddCommand(rmStackCmd)
	cmd.AddCommand(exportCmd)
	cmd.AddCommand(importCmd)

	return cmd
}

func listStacks(cmd *cobra.Command, args []string) error {
	stacks := core.ListStacks()
	stacksCount := len(stacks)
	if stacksCount == 0 {
		fmt.Println("No stacks available. Create a new one with 'apx stacks new' or contact the system administrator.")
		return nil
	}

	fmt.Printf("Found %d stacks:\n", stacksCount)

	table := core.CreateApxTable(os.Stdout)
	table.SetHeader([]string{"Name", "Base", "Built-in", "Pkgs", "Pkg manager"})

	for _, stack := range stacks {
		builtIn := "No"
		if stack.BuiltIn {
			builtIn = "Yes"
		}
		table.Append([]string{stack.Name, stack.Base, builtIn, fmt.Sprintf("%d", len(stack.Packages)), stack.PkgManager})
	}

	table.Render()

	return nil
}

func showStack(cmd *cobra.Command, args []string) error {
	stack, error := core.LoadStack(args[0])
	if error != nil {
		return error
	}

	table := core.CreateApxTable(os.Stdout)
	table.Append([]string{"Name", stack.Name})
	table.Append([]string{"Base", stack.Base})
	table.Append([]string{"Packages", strings.Join(stack.Packages, ", ")})
	table.Append([]string{"Package manager", stack.PkgManager})
	table.Render()

	return nil
}

func newStack(cmd *cobra.Command, args []string) error {
	assumeYes, _ := cmd.Flags().GetBool("assume-yes")
	name, _ := cmd.Flags().GetString("name")
	base, _ := cmd.Flags().GetString("base")
	packages, _ := cmd.Flags().GetStringArray("packages")
	pkgManager, _ := cmd.Flags().GetString("pkg-manager")

	if name == "" {
		if !assumeYes {
			cmdr.Info.Println("Please type a name for the stack:")
			fmt.Scanln(&name)
			if name == "" {
				cmdr.Error.Println("The name cannot be empty.")
				return nil
			}
		} else {
			cmdr.Error.Println("Please provide a name for the stack.")
			return nil
		}
	}

	ok := core.StackExists(name)
	if ok {
		cmdr.Error.Println("A stack with the same name already exists.")
		return nil
	}

	if base == "" {
		if !assumeYes {
			cmdr.Info.Println("Please type a base image for the stack, (e.g. vanillaos/pico):")
			fmt.Scanln(&base)
			if base == "" {
				cmdr.Error.Println("The base image cannot be empty.")
				return nil
			}
		} else {
			cmdr.Error.Println("Please provide a base image for the stack")
			return nil
		}
	}

	if pkgManager == "" {
		pkgManagers := core.ListPkgManagers()
		if len(pkgManagers) == 0 {
			cmdr.Error.Println("No package managers available")
			return nil
		}

		cmdr.Info.Println("Please select a package manager:")
		for i, manager := range pkgManagers {
			fmt.Printf("%d. %s\n", i+1, manager.Name)
		}
		fmt.Printf("Select a package manager [1-%d]: ", len(pkgManagers))

		var pkgManagerIndex int
		_, err := fmt.Scanln(&pkgManagerIndex)
		if err != nil {
			cmdr.Error.Println("Invalid input")
			return nil
		}

		if pkgManagerIndex < 1 || pkgManagerIndex > len(pkgManagers) {
			cmdr.Error.Println("Invalid package manager selection")
			return nil
		}

		pkgManager = pkgManagers[pkgManagerIndex-1].Name
	}

	ok = core.PkgManagerExists(pkgManager)
	if !ok {
		cmdr.Error.Println("The specified package manager does not exist. Create it with 'apx pkgmanagers new' or contact the system administrator.")
		return nil
	}

	if len(packages) == 0 && !assumeYes {
		cmdr.Info.Println("You have not provided any package to install in the stack. Do you want to add some now? [y/N]")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			cmdr.Info.Println("Please type the packages you want to install in the stack, separated by a space:")
			packagesInput, _ := reader.ReadString('\n')
			packagesInput = strings.TrimSpace(packagesInput)
			packages = strings.Fields(packagesInput)
		} else {
			packages = []string{}
		}
	}

	stack := core.NewStack(name, base, packages, pkgManager, false)

	err := stack.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Stack %s created successfully!\n", name)

	return nil
}

func updateStack(cmd *cobra.Command, args []string) error {
	assumeYes, _ := cmd.Flags().GetBool("assume-yes")
	name, _ := cmd.Flags().GetString("name")
	base, _ := cmd.Flags().GetString("base")
	packages, _ := cmd.Flags().GetStringArray("packages")
	pkgManager, _ := cmd.Flags().GetString("pkg-manager")

	if name == "" {
		if len(args) != 1 || args[0] == "" {
			cmdr.Error.Println("Please provide the name of the stack to update.")
			return nil
		}

		cmd.Flags().Set("name", args[0])
		name = args[0]
	}

	stack, error := core.LoadStack(name)
	if error != nil {
		return error
	}

	if base == "" {
		if !assumeYes {
			cmdr.Info.Printf("Please type a new base image for the stack or confirm the current one (%s):", stack.Base)
			fmt.Scanln(&base)
			if base == "" {
				base = stack.Base
			}
		} else {
			cmdr.Error.Println("Please provide a base image for the stack")
			return nil
		}
	}

	if pkgManager == "" {
		if !assumeYes {
			cmdr.Info.Printf("Please type a new package manager for the stack or confirm the current one (%s):", stack.PkgManager)
			fmt.Scanln(&pkgManager)
			if pkgManager == "" {
				pkgManager = stack.PkgManager
			}
		} else {
			cmdr.Error.Println("Please provide a package manager for the stack")
			return nil
		}
	}

	ok := core.PkgManagerExists(pkgManager)
	if !ok {
		cmdr.Error.Println("The specified package manager does not exist. Create it with 'apx pkgmanagers new' or contact the system administrator.")
		return nil
	}

	if len(packages) == 0 && !assumeYes {
		if len(stack.Packages) == 0 {
			cmdr.Info.Println("You have not provided any package to install in the stack. Do you want to add some now? [y/N]")
		} else {
			cmdr.Info.Println("Do you want to confirm the current packages list?\n" + strings.Join(stack.Packages, "\n\t - ") + "\n[y/N]")
		}

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if answer == "y" || answer == "Y" {
			if len(stack.Packages) > 0 {
				packages = stack.Packages
			} else {
				cmdr.Info.Println("Please type the packages you want to install in the stack, separated by a space:")
				packagesInput, _ := reader.ReadString('\n')
				packagesInput = strings.TrimSpace(packagesInput)
				packages = strings.Fields(packagesInput)
			}
		} else {
			packages = []string{}
		}
	}

	stack.Base = base
	stack.Packages = packages
	stack.PkgManager = pkgManager

	err := stack.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Stack %s updated successfully!\n", name)

	return nil
}

func removeStack(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("name")
	if stackName == "" {
		cmdr.Error.Println("Please provide the name of the stack to remove.")
		return nil
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		cmdr.Info.Printf("Are you sure you want to remove the stack %s? [y/N]\n", stackName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println("Aborting...")
			return nil
		}
	}

	stack, error := core.LoadStack(stackName)
	if error != nil {
		return error
	}

	error = stack.Remove()
	if error != nil {
		return error
	}

	fmt.Printf("Stack %s removed successfully\n", stack.Name)
	return nil
}

func exportStack(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("name")
	if stackName == "" {
		cmdr.Error.Println("Please provide the name of the stack to export.")
		return nil
	}

	stack, error := core.LoadStack(stackName)
	if error != nil {
		return error
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		cmdr.Error.Println("Please provide the output location.")
		return nil
	}

	error = stack.Export(output)
	if error != nil {
		return error
	}

	fmt.Printf("Stack %s exported successfully to %s\n", stack.Name, output)
	return nil
}

func importStack(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	if input == "" {
		cmdr.Error.Println("Please provide the input file.")
		return nil
	}

	stack, error := core.LoadStackFromPath(input)
	if error != nil {
		cmdr.Error.Printf("Error while loading the stack from %s: %s\n", input, error)
	}

	error = stack.Save()
	if error != nil {
		return error
	}

	fmt.Printf("Stack %s imported successfully\n", stack.Name)
	return nil
}
