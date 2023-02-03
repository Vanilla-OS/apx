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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func NewNixInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize nix repository",
		Long: `Initializes a custom installation of nix by creating $HOME/.nix and
setting up some SystemD units to mount it as /nix.`,

		RunE: initNix,
	}

	return cmd
}
func initNix(cmd *cobra.Command, args []string) error {
	// prompt for confirmation
	log.Default().Printf(`This will create a ".nix" folder in your home directory
and set up some SystemD units to mount that folder at /nix before running the installation
Confirm 'y' to continue. [y/N] `)

	var proceed string
	fmt.Scanln(&proceed)
	proceed = strings.ToLower(proceed)

	if proceed != "y" {
		log.Default().Printf("operation canceled at user request")
		os.Exit(0)
	}
	err := core.NixInit()
	if err != nil {
		return err
	}
	log.Default().Printf("Installation complete. Reboot to use nix.")
	return nil

}
