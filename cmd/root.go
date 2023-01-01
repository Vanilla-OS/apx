package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

// package level variables for viper flags
var aur, dnf, apk bool

// package level variable for container name,
// set in root command's PersistentPreRun function
var container *core.Container
var name string

func NewApxCommand(Version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "apx",
		Short:   "Apx is a package manager with support for multiple sources allowing you to install packages in a managed container.",
		Version: Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			container = getContainer()
		},
	}
	rootCmd.PersistentFlags().BoolVar(&aur, "aur", false, "Install packages from the AUR (Arch User Repository).")
	rootCmd.PersistentFlags().BoolVar(&dnf, "dnf", false, "Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&apk, "apk", false, "Install packages from the Alpine repository.")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Create named container with this name.")

	rootCmd.AddCommand(NewInitializeCommand())
	rootCmd.AddCommand(NewAutoRemoveCommand())
	rootCmd.AddCommand(NewInstallCommand())
	rootCmd.AddCommand(NewCleanCommand())
	rootCmd.AddCommand(NewEnterCommand())
	rootCmd.AddCommand(NewExportCommand())
	rootCmd.AddCommand(NewListCommand())
	rootCmd.AddCommand(NewPurgeCommand())
	rootCmd.AddCommand(NewRemoveCommand())
	rootCmd.AddCommand(NewRunCommand())
	rootCmd.AddCommand(NewSearchCommand())
	rootCmd.AddCommand(NewShowCommand())
	rootCmd.AddCommand(NewUnexportCommand())
	rootCmd.AddCommand(NewUpdateCommand())
	rootCmd.AddCommand(NewUpgradeCommand())
	return rootCmd
}

func getContainer() *core.Container {
	var kind core.ContainerType = core.DEFAULT
	if aur {
		kind = core.AUR
	} else if dnf {
		kind = core.DNF
	} else if apk {
		kind = core.APK
	}
	if len(name) > 0 {
		return core.NewNamedContainer(kind, name)
	} else {
		return core.NewContainer(kind)
	}

}
