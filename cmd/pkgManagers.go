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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/apx/settings"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewPkgManagersCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"pkgmanagers",
		apx.Trans("pkgmanagers"),
		apx.Trans("pkgmanagers"),
		nil,
	)
	cmd.Example = "apx pkgmanagers"

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		apx.Trans("pkgmanagers"),
		apx.Trans("pkgmanagers"),
		listPkgManagers,
	)
	listCmd.Example = "apx pkgmanagers list"

	// Show subcommand
	showCmd := cmdr.NewCommand(
		"show",
		apx.Trans("showPkgManager"),
		apx.Trans("showPkgManager"),
		showPkgManager,
	)
	showCmd.Example = "apx pkgmanagers show myPkgManager"
	showCmd.Args = cobra.MinimumNArgs(1)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		apx.Trans("newPkgManager"),
		apx.Trans("newPkgManager"),
		newPkgManager,
	)
	newCmd.Example = "apx pkgmanagers new -n myPkgManager -a \"autoremove\" -c \"clean\" -i \"install\" -l \"list\" -p \"purge\" -r \"remove\" -s \"search\" -w \"show\" -u \"update\" -U \"upgrade\""

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
			"Name of the package manager",
			"",
		),
	)
	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"need-sudo",
			"S",
			"Whether the package manager needs sudo to run",
			false,
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"autoremove",
			"a",
			"Command to autoremove packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"clean",
			"c",
			"Command to clean the package manager",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"install",
			"i",
			"Command to install packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"list",
			"l",
			"Command to list installed packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"purge",
			"p",
			"Command to purge packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"remove",
			"r",
			"Command to remove packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"search",
			"s",
			"Command to search packages",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"show",
			"w",
			"Command to show package info",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"update",
			"u",
			"Command to update the package manager",
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"upgrade",
			"U",
			"Command to upgrade packages",
			"",
		),
	)

	// Rm subcommand
	rmCmd := cmdr.NewCommand(
		"rm",
		apx.Trans("rmPkgManager"),
		apx.Trans("rmPkgManager"),
		rmPkgManager,
	)
	rmCmd.Example = "apx pkgmanagers rm myPkgManager"

	rmCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			"Name of the package manager",
			"",
		),
	)
	rmCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			"Force removal of the package manager",
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
	pkgManagers := core.ListPkgManagers()
	pkgManagersCount := len(pkgManagers)
	if pkgManagersCount == 0 {
		fmt.Println("No package managers available. Create a new one with 'apx pkgmanagers new' or contact the system administrator.")
		return nil
	}

	fmt.Printf("Found %d pkgManagers:\n", pkgManagersCount)

	table := core.CreateApxTable(os.Stdout)
	table.SetHeader([]string{"Name", "Built-in"})

	for _, stack := range pkgManagers {
		builtIn := "No"
		if stack.BuiltIn {
			builtIn = "Yes"
		}
		table.Append([]string{stack.Name, builtIn})
	}

	table.Render()

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
	table.Append([]string{"Name", pkgManager.Name})
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

	if name == "" {
		if !assumeYes {
			cmdr.Info.Println("Please type a name for the package manager:")
			fmt.Scanln(&name)
			if name == "" {
				cmdr.Error.Println("The name cannot be empty.")
				return nil
			}
		} else {
			cmdr.Error.Println("Please provide a name for the package manager.")
			return nil
		}
	}

	if !needSudo && !assumeYes {
		cmdr.Info.Println("Does the package manager need sudo to run? [y/N]")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(answer)) == "y" {
			needSudo = true
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
			if !assumeYes {
				cmdr.Info.Printf("Please type the command for %s:\n", cmdName)
				fmt.Scanln(cmd)
				if *cmd == "" {
					cmdr.Error.Printf("The command for %s cannot be empty.\n", cmdName)
					return nil
				}
			} else {
				cmdr.Error.Printf("Please provide the command for %s.\n", cmdName)
				return nil
			}
		}
	}

	if _, err := os.Stat(filepath.Join(settings.Cnf.PkgManagersPath, name+".yml")); err == nil {
		if !assumeYes {
			cmdr.Info.Printf("A package manager with the name %s already exists. Do you want to overwrite it? [y/N]\n", name)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(answer)) != "y" {
				cmdr.Info.Println("Aborting.")
				return nil
			}
		} else {
			cmdr.Error.Println("A package manager with the same name already exists.")
			return nil
		}
	}

	pkgManager := core.NewPkgManager(name, needSudo, autoRemove, clean, install, list, purge, remove, search, show, update, upgrade, false)
	err := pkgManager.Save()
	if err != nil {
		cmdr.Error.Println(err)
		return nil
	}

	cmdr.Success.Printf("Package manager %s created successfully!\n", name)

	return nil
}

func rmPkgManager(cmd *cobra.Command, args []string) error {
	pkgManagerName, _ := cmd.Flags().GetString("name")
	if pkgManagerName == "" {
		cmdr.Error.Println("Please specify the name of the package manager you want to remove.")
		return nil
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		cmdr.Info.Printf("Are you sure you want to remove the package manager %s? [y/N]\n", pkgManagerName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println("Aborting...")
			return nil
		}
	}

	pkgManager, error := core.LoadPkgManager(pkgManagerName)
	if error != nil {
		return error
	}

	stacks := core.ListStackForPkgManager(pkgManager.Name)
	if len(stacks) > 0 {
		fmt.Printf("The package manager %s is used by %d stacks:\n", pkgManager.Name, len(stacks))
		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{"Name", "Base", "Packages", "PkgManager", "Built-in"})
		for _, stack := range stacks {
			builtIn := "No"
			if stack.BuiltIn {
				builtIn = "Yes"
			}
			table.Append([]string{stack.Name, stack.Base, strings.Join(stack.Packages, ", "), stack.PkgManager, builtIn})
		}
		table.Render()
	}

	error = pkgManager.Remove()
	if error != nil {
		return error
	}

	fmt.Printf("Package manager %s removed successfully\n", pkgManager.Name)
	return nil
}
