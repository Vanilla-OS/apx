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
		apx.Trans("apx.use"),
		apx.Trans("apx"),
		apx.Trans("apx"),
		nil,
	)
	root.Version = version

	return root
}
