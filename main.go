package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"log"

	"github.com/spf13/viper"
	"github.com/vanilla-os/apx/cmd"
)

var (
	Version = "1.5.0"
)

func init() {
	log.SetPrefix("\033[1m\033[34m⌬ Apx :: \033[0m")
	log.SetFlags(0)
	viper.SetEnvPrefix("apx") // will be uppercased automatically
	viper.AutomaticEnv()
}

func main() {
	rootCmd := cmd.NewApxCommand(Version)
	//	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
