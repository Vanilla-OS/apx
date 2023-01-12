package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/apx/settings"
)

// package level variables for viper flags
var ubuntu, aur, fedora, alpine bool

// old packagemanager flags, deprecated in favor of the distribution flags above
var apt, dnf, apk bool

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
			var err error
			container, err = getContainer()

			// No commandline flag specified, get distro from config file or host
			if err != nil {
				container = core.NewContainer(settings.GetDistroFromSettings())
			}
		},
	}
	rootCmd.PersistentFlags().BoolVar(&ubuntu, "ubuntu", false, "Install packages from the Ubuntu repository.")
	rootCmd.PersistentFlags().BoolVar(&aur, "aur", false, "Install packages from the AUR (Arch User Repository).")
	rootCmd.PersistentFlags().BoolVar(&fedora, "fedora", false, "Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "alpine", false, "Install packages from the Alpine repository.")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Create or use custom container with this name.")

	rootCmd.PersistentFlags().BoolVar(&ubuntu, "apt", false, "(Deprecated) Install packages from the Ubuntu repository.")
	rootCmd.PersistentFlags().BoolVar(&fedora, "dnf", false, "(Deprecated) Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "apk", false, "(Deprecated) Install packages from the Alpine repository.")
	rootCmd.PersistentFlags().MarkDeprecated("apt", "use --ubuntu instead")
	rootCmd.PersistentFlags().MarkDeprecated("dnf", "use --fedora instead")
	rootCmd.PersistentFlags().MarkDeprecated("apk", "use --alpine instead")

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
	viper.BindPFlag("ubuntu", rootCmd.PersistentFlags().Lookup("ubuntu"))
	viper.BindPFlag("aur", rootCmd.PersistentFlags().Lookup("aur"))
	viper.BindPFlag("fedora", rootCmd.PersistentFlags().Lookup("fedora"))
	viper.BindPFlag("alpine", rootCmd.PersistentFlags().Lookup("alpine"))

	return rootCmd
}

func getContainer() (*core.Container, error) {
	// in the future these should be moved to
	// constants, and wrapped in package level calls
	ubuntu = viper.GetBool("ubuntu")
	aur = viper.GetBool("aur")
	fedora = viper.GetBool("fedora")
	alpine = viper.GetBool("alpine")

	var kind settings.DistroInfo
	if ubuntu {
		kind = settings.DistroUbuntu
	} else if aur {
		kind = settings.DistroArch
	} else if fedora {
		kind = settings.DistroFedora
	} else if alpine {
		kind = settings.DistroAlpine
	} else {
		return &core.Container{}, fmt.Errorf("No distribution passed as argument")
	}
	if len(name) > 0 {
		return core.NewNamedContainer(kind, name), nil
	} else {
		return core.NewContainer(kind), nil
	}
}
