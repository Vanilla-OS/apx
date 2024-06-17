package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"embed"
	"os"

	"github.com/vanilla-os/apx/v2/cmd"
	"github.com/vanilla-os/apx/v2/core"
	"github.com/vanilla-os/orchid/cmdr"
)

var Version = "2.4.1"

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
	apx.CreateRootCommand(root, apx.Trans("apx.msg.help"), apx.Trans("apx.msg.version"))

	msgs := cmdr.UsageStrings{
		Usage:                apx.Trans("apx.msg.usage"),
		Aliases:              apx.Trans("apx.msg.aliases"),
		Examples:             apx.Trans("apx.msg.examples"),
		AvailableCommands:    apx.Trans("apx.msg.availableCommands"),
		AdditionalCommands:   apx.Trans("apx.msg.additionalCommands"),
		Flags:                apx.Trans("apx.msg.flags"),
		GlobalFlags:          apx.Trans("apx.msg.globalFlags"),
		AdditionalHelpTopics: apx.Trans("apx.msg.additionalHelpTopics"),
		MoreInfo:             apx.Trans("apx.msg.moreInfo"),
	}
	apx.SetUsageStrings(msgs)

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
