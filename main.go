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
	"os"

	"github.com/vanilla-os/apx/cmd"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/orchid/cmdr"
)

var (
	Version = "2.0.0"
)

//go:embed locales/*.yml
var fs embed.FS
var apx *cmdr.App

func main() {
	core.NewStandardApx()

	apx = cmd.New(Version, fs)

	// check if root, exit if so
	if core.RootCheck(false) {
		cmdr.Error.Println(apx.Trans("apx.errors.noRoot"))
		os.Exit(1)
	}

	// root command
	root := cmd.NewRootCommand(Version)
	apx.CreateRootCommand(root)

	// commands
	stacks := cmd.NewStacksCommand()
	root.AddCommand(stacks)

	subsystems := cmd.NewSubSystemsCommand()
	root.AddCommand(subsystems)

	pkgManagers := cmd.NewPkgManagersCommand()
	root.AddCommand(pkgManagers)

	runtimeCmds := cmd.NewRuntimeCommands()
	root.AddCommand(runtimeCmds...)

	// run the app
	err := apx.Run()
	if err != nil {
		cmdr.Error.Println(err)
	}
}
