package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/apx/settings"
)

// package level variables for viper flags
var apt, aur, dnf, apk bool

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

	rootCmd.PersistentFlags().BoolVar(&apt, string(core.APT), false, "Install packages from the Ubuntu repository.")
	rootCmd.PersistentFlags().BoolVar(&aur, string(core.AUR), false, "Install packages from the AUR (Arch User Repository).")
	rootCmd.PersistentFlags().BoolVar(&dnf, string(core.DNF), false, "Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&apk, string(core.APK), false, "Install packages from the Alpine repository.")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Create or use custom container with this name.")

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
	viper.BindPFlag(string(core.APT), rootCmd.PersistentFlags().Lookup(string(core.APT)))
	viper.BindPFlag(string(core.AUR), rootCmd.PersistentFlags().Lookup(string(core.AUR)))
	viper.BindPFlag(string(core.DNF), rootCmd.PersistentFlags().Lookup(string(core.DNF)))
	viper.BindPFlag(string(core.APK), rootCmd.PersistentFlags().Lookup(string(core.APK)))
	return rootCmd
}

func getContainer() *core.Container {
	var kind core.ContainerType = core.ContainerType(settings.Cnf.PkgManager)

	apt = viper.GetBool(string(core.APT))
	aur = viper.GetBool(string(core.AUR))
	dnf = viper.GetBool(string(core.DNF))
	apk = viper.GetBool(string(core.APK))

	if apt {
		kind = core.APT
	} else if aur {
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
