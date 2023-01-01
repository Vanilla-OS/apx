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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func initializeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Initialize the managed container.

Usage:
  apx init [options]

Options:
  -h, --help            Show this help message and exit
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

	if container.Exists() {
		log.Default().Printf(`Container already exists. Do you want to re-initialize it?\ 
This operation will remove everything, including your files in the container. [y/N] `)

		var proceed string
		fmt.Scanln(&proceed)
		proceed = strings.ToLower(proceed)

		if proceed != "y" {
			os.Exit(0)
		}
	}

	if err := container.Remove(); err != nil {
		panic(err)
	}
	if err := container.Create(); err != nil {
		panic(err)
	}

	return nil
}
