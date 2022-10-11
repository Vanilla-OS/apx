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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func initializeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Initialize the managed container.

Usage:
  apx init

Examples:
  apx init
`)
	return nil
}

func NewInitializeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the managed container",
		RunE:  initialize,
	}
	cmd.SetUsageFunc(initializeUsage)
	return cmd
}

func initialize(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	}

	if core.ContainerExists(container) {
		log.Default().Printf(`Container already exists. Do you want to re-initialize it?\ 
This operation will remove everything, of course your files as well. [y/N] `)

		var proceed string
		fmt.Scanln(&proceed)
		proceed = strings.ToLower(proceed)

		if proceed != "y" {
			os.Exit(0)
		}
	}

	if err := core.RemoveContainer(container); err != nil {
		panic(err)
	}
	if err := core.CreateContainer(container); err != nil {
		panic(err)
	}

	log.Default().Println("Container created")

	return nil
}
