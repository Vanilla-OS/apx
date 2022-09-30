package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func autoRemoveUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Remove automatically all unused packages.

Usage:
  apx autoremove

Examples:
  apx autoremove
`)
	return nil
}

func NewAutoRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoremove",
		Short: "Remove automatically all unused packages",
		RunE:  autoRemove,
	}
	cmd.SetUsageFunc(autoRemoveUsage)
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().BoolP("on-persistent", "p", false, "used by systemd to enter in default mode only if the persistent mode is enabled")
	return cmd
}

func autoRemove(cmd *cobra.Command, args []string) error {
	fmt.Println("autoremove")

	// example of --sys flag
	if cmd.Flag("sys").Value.String() == "true" {
		if core.ImmutabilityStatus() {
			fmt.Println("immutable")
		}
		fmt.Println("system")
	}

	return nil
}
