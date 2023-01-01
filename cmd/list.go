package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	List installed packages.

Usage:
  apx list [options]

Options:
  -h, --help            Show this help message and exit
  -u, --upgradable      Show only upgradable packages
  -i, --installed       Show only installed packages

`)
	return nil
}

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed packages.",
		RunE:  list,
	}
	//cmd.SetUsageFunc(listUsage)
	cmd.Flags().SetInterspersed(false)
	cmd.Flags().BoolP("upgradable", "u", false, "List only upgradable packages.")
	cmd.Flags().BoolP("installed", "i", false, "List only installed packages.")

	return cmd
}

func list(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("list")...)

	if cmd.Flag("upgradable").Value.String() == "true" {
		command = append(command, "--upgradable")
	}
	if cmd.Flag("installed").Value.String() == "true" {
		command = append(command, "--installed")
	}

	command = append(command, args...)

	container.Run(command...)

	return nil
}
