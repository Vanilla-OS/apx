package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewNixInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "Install nix package",
		Example: "apx nix install jq",
		RunE:    installPackage,
		Args:    cobra.ExactArgs(1),
	}

	return cmd
}
func installPackage(cmd *cobra.Command, args []string) error {
	install := exec.Command("nix", "profile", "install", "nixpkgs#"+args[0])
	install.Stderr = os.Stderr
	install.Stdin = os.Stdin
	install.Stdout = os.Stdout

	err := install.Run()
	if err != nil {
		log.Default().Printf("error installing package")
		return err
	}

	log.Default().Printf("Package installation complete")
	return nil

}
