package main

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.4.1-1"
)

func init() {
	log.SetPrefix("\033[1m\033[34m⌬ Apx :: \033[0m")
	log.SetFlags(0)
	viper.SetEnvPrefix("apx") // will be uppercased automatically
	viper.AutomaticEnv()
}

func help(cmd *cobra.Command, args []string) {
	fmt.Println(`Usage:
  apx [command]

Available Commands:
  autoremove  Remove all unused packages automatically
  clean       Clean the apx package manager cache
  completion  Generate the autocompletion script for the specified shell
  enter       Enter in the container shell
  export      Export/Recreate a program's desktop entry from a managed container
  help        Help about any command
  init        Initialize the managed container
  install     Install packages inside a managed container
  list        List installed packages.
  purge       Purge packages inside a managed container
  remove      Remove packages inside a managed container.
  run         Run a program inside a managed container.
  search      Search for packages in a managed container.
  show        Show details about a package
  unexport    Unexport/Remove a program's desktop entry from a managed container
  update      Update the list of available packages
  upgrade     Upgrade the system by installing/upgrading available packages.

Global Flags:
      --apk           Install packages from the Alpine repository.
      --apt           Install packages from the Ubuntu repository.
      --aur           Install packages from the AUR (Arch User Repository).
      --dnf           Install packages from the Fedora's DNF (Dandified YUM) repository.
  -h, --help          help for apx
  -n, --name string   Create or use custom container with this name.
  -v, --version       version for apx

Use "apx [command] --help" for more information about a command.`)
}

func main() {
	rootCmd := cmd.NewApxCommand(Version)
	//	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
