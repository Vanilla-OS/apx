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

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/roff"
)

func NewManCommand(a *App) *Command {
	c := &Command{}
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			doc := roff.NewDocument()
			a.rootDocs(doc)
			doc.Section(a.Name + " subcommands")
			for _, c := range a.RootCommand.Children() {
				if !c.Hidden {
					c.doc(doc)
				}
			}

			fmt.Println(doc.String())
			return nil

		},
	}
	c.Command = cmd
	return c

}
