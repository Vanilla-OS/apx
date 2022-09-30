package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.0.0"
)

func help(cmd *cobra.Command, args []string) {
	fmt.Print(`Usage: 
apx [options] [command] [arguments]

Options:
  -h, --help    Show this help message and exit
  -v, --version Show version and exit
  --sys        Perform operations on the system instead of the managed container

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
    version     Show version and exit
`)
}

func newApxCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "apx",
		Short:   "Apx is a package manager around apt which allows you to install packages in a container or system",
		Version: Version,
	}
}

func main() {
	rootCmd := newApxCommand()
	rootCmd.PersistentFlags().Bool("sys", false, "Perform operations on the system instead of the managed container")
	rootCmd.AddCommand(cmd.NewAutoRemoveCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
