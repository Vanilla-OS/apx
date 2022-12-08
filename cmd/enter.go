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
	"log"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func enterUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Enter in the container shell.

Usage:
  apx enter
  apx --aur enter
  apx --dnf enter
`)
	return nil
}

func NewEnterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enter",
		Short: "Enter in the container shell",
		RunE:  enter,
	}
	cmd.SetUsageFunc(enterUsage)
	return cmd
}

func enter(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	if err := core.EnterContainer(container); err != nil {
		log.Default().Fatal("Failed to enter container: ", err)
		return err
	}

	fmt.Print("You are now outside the container.\n")
	return nil
}
