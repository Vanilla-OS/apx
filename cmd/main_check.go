//go:build check_missing_strings

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

	cmd "github.com/vanilla-os/apx/v3/internal/cli"
	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
)

var Version = "development"

//go:embed locales/*
var embeddedLocales embed.FS

func main() {
	var err error

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
}
