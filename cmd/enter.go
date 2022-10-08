package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
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

Examples:
  apx enter
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
	if err := core.EnterContainer(); err != nil {
		log.Default().Fatal("Failed to enter container: ", err)
		return err
	}

	fmt.Print("You are now outside the container.\n")
	return nil
}
