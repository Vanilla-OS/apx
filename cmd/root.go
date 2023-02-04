package cmd

import (
	"embed"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

var apx *cmdr.App

// package level variables for viper flags
var apt, aur, dnf, apk, zypper, xbps bool

// package level variable for container name,
// set in root command's PersistentPreRun function
var container *core.Container
var name string

const (
	verboseFlag string = "verbose"
)

func New(version string, fs embed.FS) *cmdr.App {
	apx = cmdr.NewApp("apx", version, fs)
	return apx
}
func NewRootCommand(version string) *cmdr.Command {
	root := cmdr.NewCommand(
		apx.Trans("apx.use"),
		apx.Trans("apx.long"),
		apx.Trans("apx.short"),
		nil).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				verboseFlag,
				"v",
				apx.Trans("apx.verboseFlag"),
				false))
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		err := setStorage()
		cobra.CheckErr(err)
		container = getContainer()
	}

	/*
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

		rootCmd.AddCommand(addContainerFlags(NewInitializeCommand()))
		rootCmd.AddCommand(addContainerFlags(NewAutoRemoveCommand()))
		rootCmd.AddCommand(addContainerFlags(NewInstallCommand()))
		rootCmd.AddCommand(addContainerFlags(NewCleanCommand()))
		rootCmd.AddCommand(addContainerFlags(NewEnterCommand()))
		rootCmd.AddCommand(addContainerFlags(NewExportCommand()))
		rootCmd.AddCommand(addContainerFlags(NewListCommand()))
		rootCmd.AddCommand(addContainerFlags(NewPurgeCommand()))
		rootCmd.AddCommand(addContainerFlags(NewRemoveCommand()))
		rootCmd.AddCommand(addContainerFlags(NewRunCommand()))
		rootCmd.AddCommand(addContainerFlags(NewSearchCommand()))
		rootCmd.AddCommand(addContainerFlags(NewShowCommand()))
		rootCmd.AddCommand(addContainerFlags(NewUnexportCommand()))
		rootCmd.AddCommand(addContainerFlags(NewUpdateCommand()))
		rootCmd.AddCommand(addContainerFlags(NewUpgradeCommand()))
		rootCmd.AddCommand(NewNixCommand())

		return rootCmd
	*/
	root.Version = version
	return root
}

// AddContainerFlags applies flags that are only relevant
// to apx commands that pertain to containers. e.g. not nix
func AddContainerFlags(cmd *cmdr.Command) *cmdr.Command {
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false
	cmd.PersistentFlags().BoolVar(&apt, "apt", false, apx.Trans("flags.apt"))
	cmd.PersistentFlags().BoolVar(&aur, "aur", false, apx.Trans("flags.aur"))
	cmd.PersistentFlags().BoolVar(&dnf, "dnf", false, apx.Trans("flags.dnf"))
	cmd.PersistentFlags().BoolVar(&apk, "apk", false, apx.Trans("flags.apk"))
	cmd.PersistentFlags().BoolVar(&zypper, "zypper", false, apx.Trans("flags.zypper"))
	cmd.PersistentFlags().BoolVar(&xbps, "xbps", false, apx.Trans("flags.xbps"))
	cmd.PersistentFlags().StringVarP(&name, "name", "n", "", apx.Trans("flags.name"))
	viper.BindPFlag("aur", cmd.PersistentFlags().Lookup("aur"))
	viper.BindPFlag("apt", cmd.PersistentFlags().Lookup("apt"))
	viper.BindPFlag("dnf", cmd.PersistentFlags().Lookup("dnf"))
	viper.BindPFlag("apk", cmd.PersistentFlags().Lookup("apk"))
	viper.BindPFlag("zypper", cmd.PersistentFlags().Lookup("zypper"))
	viper.BindPFlag("xbps", cmd.PersistentFlags().Lookup("xbps"))
	return cmd
}

func getContainer() *core.Container {
	var kind core.ContainerType = core.APT
	// in the future these should be moved to
	// constants, and wrapped in package level calls
	apt = viper.GetBool("apt")
	aur = viper.GetBool("aur")
	dnf = viper.GetBool("dnf")
	apk = viper.GetBool("apk")
	zypper = viper.GetBool("zypper")
	xbps = viper.GetBool("xbps")
	if aur {
		kind = core.AUR
	} else if dnf {
		kind = core.DNF
	} else if apk {
		kind = core.APK
	} else if zypper {
		kind = core.ZYPPER
	} else if xbps {
		kind = core.XBPS
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
