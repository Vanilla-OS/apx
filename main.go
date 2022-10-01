package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/cmd"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/apx/lang"
)

var (
	Version = "1.0.0"
)

func help(cmd *cobra.Command, args []string) {
	text := lang.GetText("en", "cmd_help")
	fmt.Println(text)
}

func newApxCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "apx",
		Short:   "Apx is a package manager around apt which allows you to install packages in a container or system",
		Version: Version,
	}
}

func main() {
	rootCmd := newApxCommand()
	rootCmd.PersistentFlags().Bool("sys", false, "Perform operations on the system instead of the managed container")
	rootCmd.AddCommand(cmd.NewAutoRemoveCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()

	fmt.Println(" --------- ")
	image, _ := core.GetHostImage()
	fmt.Println(image)
	fmt.Println(core.GetDistroboxVersion())
	fmt.Println(" --------- ")
}
