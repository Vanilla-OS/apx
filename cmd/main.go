//go:build !check_missing_strings

package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Pietro di Caprio <pietro@fabricators.ltd>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/vanilla-os/apx/v3/core"
	cmd "github.com/vanilla-os/apx/v3/internal/cli"
	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
)

var Version = "development"

//go:embed locales/*
var embeddedLocales embed.FS

func main() {
	var err error
	core.NewStandardApx()

	// Initialize SDK App
	subFS, err := fs.Sub(embeddedLocales, "locales")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Apx, err = app.NewApp(types.AppOptions{
		Name:          "apx",
		Version:       Version,
		RDNN:          "org.vanillaos.apx",
		LocalesFS:     subFS,
		DefaultLocale: "en",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize CLI
	rootCmdStruct := &cmd.RootCmd{
		Version: Version,
	}

	// Dynamic Commands (Runtime)
	subSystems, err := core.ListSubSystems(false, false)
	if err == nil {
		m := make(map[string]*cmd.SubsystemCmd)
		rootCmdStruct.DynamicSubsystems = &m
		for _, s := range subSystems {
			(*rootCmdStruct.DynamicSubsystems)[s.Name] = &cmd.SubsystemCmd{
				Name:       s.Name,
				Enter:      cmd.SubsystemEnterCmd{Name: s.Name},
				Run:        cmd.SubsystemRunCmd{Name: s.Name},
				Install:    cmd.SubsystemInstallCmd{Name: s.Name},
				Remove:     cmd.SubsystemRemoveCmd{Name: s.Name},
				Update:     cmd.SubsystemUpdateCmd{Name: s.Name},
				Upgrade:    cmd.SubsystemUpgradeCmd{Name: s.Name},
				List:       cmd.SubsystemListCmd{Name: s.Name},
				Search:     cmd.SubsystemSearchCmd{Name: s.Name},
				Show:       cmd.SubsystemShowCmd{Name: s.Name},
				Export:     cmd.SubsystemExportCmd{Name: s.Name},
				Unexport:   cmd.SubsystemUnexportCmd{Name: s.Name},
				Start:      cmd.SubsystemStartCmd{Name: s.Name},
				Stop:       cmd.SubsystemStopCmd{Name: s.Name},
				AutoRemove: cmd.SubsystemAutoRemoveCmd{Name: s.Name},
				Clean:      cmd.SubsystemCleanCmd{Name: s.Name},
				Purge:      cmd.SubsystemPurgeCmd{Name: s.Name},
			}
		}
	}

	err = cmd.Apx.WithCLI(rootCmdStruct)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Apx.CLI.Reload()
	cmd.Apx.CLI.SetName("apx")

	cmd.Apx.CLI.Reload()
	cmd.Apx.CLI.SetName("apx")
	cmd.Apx.CLI.SetTranslator(func(key string) string {
		return cmd.Apx.LC.Get(key)
	})

	err = cmd.Apx.CLI.Execute()
	if err != nil {
		cmd.Apx.Log.Error(err.Error())
		os.Exit(1)
	}
}
