package cmd

import (
	"fmt"
	"os"
	"path"

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
			err := setStorage()

			cobra.CheckErr(err)

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
	rootCmd.PersistentFlags().BoolVar(&alpine, "opensuse", false, "Install packages from the OpenSUSE repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "void", false, "Install packages from the Void Linux repository.")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Create or use custom container with this name.")

	rootCmd.PersistentFlags().BoolVar(&ubuntu, "apt", false, "(Deprecated) Install packages from the Ubuntu repository.")
	rootCmd.PersistentFlags().BoolVar(&fedora, "dnf", false, "(Deprecated) Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "apk", false, "(Deprecated) Install packages from the Alpine repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "zypper", false, "(Deprecated) Install packages from the OpenSUSE repository.")
	rootCmd.PersistentFlags().BoolVar(&alpine, "xbps", false, "(Deprecated) Install packages from the Void Linux repository.")
	rootCmd.PersistentFlags().MarkDeprecated("apt", "use --ubuntu instead")
	rootCmd.PersistentFlags().MarkDeprecated("dnf", "use --fedora instead")
	rootCmd.PersistentFlags().MarkDeprecated("apk", "use --alpine instead")
	rootCmd.PersistentFlags().MarkDeprecated("zypper", "use --opensuse instead")
	rootCmd.PersistentFlags().MarkDeprecated("xbps", "use --void instead")

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
	viper.BindPFlag("opensuse", rootCmd.PersistentFlags().Lookup("opensuse"))
	viper.BindPFlag("void", rootCmd.PersistentFlags().Lookup("void"))

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

func setStorage() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := path.Join(home, ".config", "containers", "storage.conf")
	_, err = os.Stat(configPath)
	if err == nil {
		// config exists
		return nil
	}
	// storage config path doesn't exist
	err = os.MkdirAll(path.Join(home, ".config", "containers"), 0755)
	if err != nil {
		return err
	}
	// storage config file doesn't exist
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(storageConf))

	return err
}

const storageConf = `[storage]
driver = "vfs"
`
