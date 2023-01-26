package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/core"
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
			err := setStorage()

			cobra.CheckErr(err)

			container = getContainer()
		},
	}
	rootCmd.PersistentFlags().BoolVar(&apt, "apt", false, "Install packages from the Ubuntu repository.")
	rootCmd.PersistentFlags().BoolVar(&aur, "aur", false, "Install packages from the AUR (Arch User Repository).")
	rootCmd.PersistentFlags().BoolVar(&dnf, "dnf", false, "Install packages from the Fedora's DNF (Dandified YUM) repository.")
	rootCmd.PersistentFlags().BoolVar(&apk, "apk", false, "Install packages from the Alpine repository.")
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
	viper.BindPFlag("aur", rootCmd.PersistentFlags().Lookup("aur"))
	viper.BindPFlag("apt", rootCmd.PersistentFlags().Lookup("apt"))
	viper.BindPFlag("dnf", rootCmd.PersistentFlags().Lookup("dnf"))
	viper.BindPFlag("apk", rootCmd.PersistentFlags().Lookup("apk"))
	return rootCmd
}

func getContainer() *core.Container {
	var kind core.ContainerType = core.APT
	// in the future these should be moved to
	// constants, and wrapped in package level calls
	apt = viper.GetBool("apt")
	aur = viper.GetBool("aur")
	dnf = viper.GetBool("dnf")
	apk = viper.GetBool("apk")
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
	// storage config doesn't exist
	err = os.MkdirAll(path.Join(home, ".config", "containers"), 0755)
	if err != nil {
		return err
	}
	// storage config doesn't exist
	f, err := os.Create(configPath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(storageConf))

	return err
}

const storageConf = `[storage]
driver = "vfs"
`
