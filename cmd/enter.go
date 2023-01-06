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
	"log"

	"github.com/spf13/cobra"
)

func NewEnterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx enter",
		Use:     "enter",
		Short:   "Enter in the container shell",
		RunE:    enter,
	}
	return cmd
}

func enter(cmd *cobra.Command, args []string) error {

	if err := container.Enter(); err != nil {
		log.Default().Fatal("Failed to enter container: ", err)
		return err
	}

	fmt.Print("You are now outside the container.\n")
	return nil
}
