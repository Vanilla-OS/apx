package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

type MetadataStruct struct {
	ImageName string `json:"image-name"`
	ImageID   string `json:"image-id"`
	Name      string `json:"name"`
	CreatedAt int    `json:"created-at"`
}

func NewListCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"list",
		apx.Trans("list.long"),
		apx.Trans("list.short"),
		list,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"containers",
			"c",
			apx.Trans("list.containers"),
			false,
		)).
		WithBoolFlag(
			cmdr.NewBoolFlag(
				"upgradable",
				"u",
				apx.Trans("list.upgradable"),
				false,
			)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"installed",
			"i",
			apx.Trans("list.installed"),
			false,
		))

	cmd.Flags().SetInterspersed(false)

	return cmd
}

func listContainers() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	containersJson := filepath.Join(homeDir, ".local/share/containers/storage/vfs-containers/containers.json")

	file, err := os.Open(containersJson)
	if err != nil {
		return err
	}

	defer file.Close()

	var data []interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	var metadata MetadataStruct
	for _, value := range data {
		obj := value.(map[string]interface{})
		metadataString := obj["metadata"].(string)

		json.Unmarshal([]byte(metadataString), &metadata)

		fmt.Println(metadata.Name)
	}

	return nil
}

func list(cmd *cobra.Command, args []string) error {
	if cmd.Flag("nix").Changed {
		return errors.New(apx.Trans("apx.notForNix"))

	}
	command := append([]string{}, container.GetPkgCommand("list")...)

	if cmd.Flag("containers").Changed {
		return listContainers()
	}
	if cmd.Flag("upgradable").Changed {
		command = append(command, "--upgradable")
	}
	if cmd.Flag("installed").Changed {
		command = append(command, "--installed")
	}

	command = append(command, args...)

	container.Run(command...)

	return nil
}
