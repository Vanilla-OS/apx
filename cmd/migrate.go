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

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx migrate",
		Use:     "migrate",
		Short:   "Migrate legacy containers to newer format",
		RunE:    migrate,
	}
	return cmd
}

func migrate(cmd *cobra.Command, args []string) error {
	legacy_containers_ids := core.GetLegacyContainersIds()
	fmt.Println(legacy_containers_ids)

	return nil
}
