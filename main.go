package main

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
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.3.2"
)

func init() {
	log.SetPrefix("\033[1m\033[34m⌬ Apx :: \033[0m")
	log.SetFlags(0)
}

func help(cmd *cobra.Command, args []string) {
	fmt.Println(`Usage:
	apx [options] [command] [arguments]

Options:
	-h, --help    Show this help message and exit
	-v, --version Show version and exit
	--aur	    Install packages from the AUR repository
	--dnf	    Install packages from the Fedora repository

Commands:
	autoremove  Remove all unused packages
	clean       Clean the apx package manager cache
	enter       Enter the container shell
	export      Export/Recreate a program's desktop entry from the container
	help        Show this help message and exit
	init        Initialize a managed container
	install     Install packages inside the container
	list        List installed packages
	purge       Purge packages from the container
	run         Run a command inside the container
	remove      Remove packages from the container
	search      Search for packages
	show        Show details about a package
	unexport    Unexport/Remove a program's desktop entry
	update      Update the list of available packages
	upgrade     Upgrade the system by installing/upgrading available packages`)
}

func newApxCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "apx",
		Short:   "Apx is a package manager with support for multiple sources allowing you to install packages in a managed container.",
		Version: Version,
	}
}

func main() {
	rootCmd := newApxCommand()
	rootCmd.PersistentFlags().Bool("aur", false, "Install packages from the AUR (Arch User Repository).")
	rootCmd.PersistentFlags().Bool("dnf", false, "Install packages from the Fedora's DNF (Dandified YUM) repository.")

	rootCmd.AddCommand(cmd.NewInitializeCommand())
	rootCmd.AddCommand(cmd.NewAutoRemoveCommand())
	rootCmd.AddCommand(cmd.NewInstallCommand())
	rootCmd.AddCommand(cmd.NewCleanCommand())
	rootCmd.AddCommand(cmd.NewEnterCommand())
	rootCmd.AddCommand(cmd.NewExportCommand())
	rootCmd.AddCommand(cmd.NewListCommand())
	rootCmd.AddCommand(cmd.NewPurgeCommand())
	rootCmd.AddCommand(cmd.NewRemoveCommand())
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewSearchCommand())
	rootCmd.AddCommand(cmd.NewShowCommand())
	rootCmd.AddCommand(cmd.NewUnexportCommand())
	rootCmd.AddCommand(cmd.NewUpdateCommand())
	rootCmd.AddCommand(cmd.NewUpgradeCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
