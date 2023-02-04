package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"embed"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/cmd"
	"github.com/vanilla-os/orchid/cmdr"
)

var (
	Version = "1.5.0"
)

//go:embed locales/*.yml
var fs embed.FS
var apx *cmdr.App

/*
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
	    --zypper        Install packages from the openSUSE repository.
	    --xbps          Install packages from the Void (Linux) repository.
	-h, --help          help for apx
	-n, --name string   Create or use custom container with this name.
	-v, --version       version for apx

Use "apx [command] --help" for more information about a command.`)
}
*/
func main() {
	apx = cmd.New(Version, fs)

	// root command
	root := cmd.NewRootCommand(Version)
	// root command
	apx.CreateRootCommand(root)
	containerGroup := &cobra.Group{
		ID:    "container",
		Title: "Managed Container Commands",
	}
	root.AddGroup(containerGroup)
	nixGroup := &cobra.Group{
		ID:    "nix",
		Title: "Nix Commands",
	}
	root.AddGroup(nixGroup)

	autoremove := cmd.NewAutoRemoveCommand()
	autoremove.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(autoremove))

	clean := cmd.NewCleanCommand()
	clean.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(clean))

	enter := cmd.NewEnterCommand()
	enter.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(enter))

	export := cmd.NewExportCommand()
	export.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(export))

	initialize := cmd.NewInitializeCommand()
	initialize.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(initialize))

	install := cmd.NewInstallCommand()
	install.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(install))

	list := cmd.NewListCommand()
	list.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(list))

	purge := cmd.NewPurgeCommand()
	purge.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(purge))

	remove := cmd.NewRemoveCommand()
	remove.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(remove))

	run := cmd.NewRunCommand()
	run.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(run))

	search := cmd.NewSearchCommand()
	search.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(search))

	show := cmd.NewShowCommand()
	show.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(show))

	unexport := cmd.NewUnexportCommand()
	unexport.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(unexport))

	upgrade := cmd.NewUpgradeCommand()
	upgrade.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(upgrade))

	update := cmd.NewUpdateCommand()
	update.GroupID = containerGroup.ID
	root.AddCommand(cmd.AddContainerFlags(update))
	nix := cmd.NewNixCommand()
	nix.GroupID = nixGroup.ID
	root.AddCommand(nix)
	// run the app
	err := apx.Run()
	if err != nil {
		cmdr.Error.Println(err)
	}
}
