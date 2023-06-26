package cmd

import (
	"embed"

	"github.com/vanilla-os/orchid/cmdr"
)

var apx *cmdr.App

func New(version string, fs embed.FS) *cmdr.App {
	apx = cmdr.NewApp("apx", version, fs)
	return apx
}

func NewRootCommand(version string) *cmdr.Command {
	root := cmdr.NewCommand(
		"apx",
		apx.Trans("apx.description"),
		apx.Trans("apx.description"),
		nil,
	)
	root.Version = version

	return root
}
