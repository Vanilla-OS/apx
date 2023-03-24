package cmdr

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Orchid is a cli helper for VanillaOS projects
*/

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewDocsCommand(a *App) *Command {
	c := &Command{}
	cmd := &cobra.Command{
		Use:                   "docs",
		Short:                 "Generates documentation for the cli application in the current directory.",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePrepender := func(filename string) string {
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), strings.Replace(base, "_", " ", -1))
			}
			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return strings.ToLower(base) + "/"
			}
			/*	err := doc.GenMarkdownTree(cmd.Root(), "./")
				if err != nil {
					log.Fatal(err)
				}
			*/
			err := doc.GenMarkdownTreeCustom(cmd.Root(), "./", filePrepender, linkHandler)

			if err != nil {
				log.Fatal(err)
			}
			return nil

		},
	}
	c.Command = cmd
	return c

}

const fmTemplate = `---
title: "%s"
description: %s

---
`
