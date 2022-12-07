package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it work inside a container
	with support to installing packages from multiple sources without altering the root filesystem.
*/

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.2.4"
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
        autoremove  Remove all unused packages automatically
        clean       Clean the apx package manager cache
        enter       Enter in the container shell
        export      Export/Recreate a program's desktop entry from a managed container
        help        Show this help message and exit
        init        Initialize the managed container
        install     Install packages inside a managed container
        list        List installed packages
        log         Show logs
        purge       Purge packages inside a managed container
        run         Run a command inside a managed container
        remove      Remove packages inside a managed container
        search      Search for packages in a managed container
        show        Show details about a package
        unexport    Unexport/Remove a program's desktop entry from a managed container
        update      Update the list of available packages
        upgrade     Upgrade the system by installing/upgrading available packages
        version     Show version and exit`)
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
