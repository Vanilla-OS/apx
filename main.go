package main

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
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.1.5"
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
	--sys         Perform operations on the system instead of the managed container
	--aur	    Install packages from the AUR repository

Commands:
	autoremove  Remove automatically all unused packages
	clean       Clean the apt cache
	enter       Enter the container
	help        Show this help message and exit
	init        Initialize the container
	install     Install packages
	list        List packages based on package names
	log         Show logs
	purge       Purge packages
	run         Run a command inside the container
	remove      Remove packages
	search      Search in package descriptions
	show        Show package details
	update      Update list of available packages
	upgrade     Upgrade the system by installing/upgrading packages
	version     Show version and exit`)
}

func newApxCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "apx",
		Short:   "Apx is a package manager around apt which allows you to install packages in a container or host system.",
		Version: Version,
	}
}

func main() {
	rootCmd := newApxCommand()
	rootCmd.PersistentFlags().Bool("sys", false, "Perform operations on the system host rather than in the container.")
	rootCmd.PersistentFlags().Bool("aur", false, "Install packages from the AUR repository.")

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
