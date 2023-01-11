package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
	"github.com/vanilla-os/apx/settings"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx migrate",
		Use:     "migrate",
		Short:   "Migrate legacy containers to newer format",
		RunE:    migrate,
	}
	return cmd
}

func migrate(cmd *cobra.Command, args []string) error {
	legacy_containers_ids := core.GetLegacyContainersIds()
	fmt.Println(legacy_containers_ids)

	manager := core.ContainerManager()

	for _, id := range legacy_containers_ids {
		out, err := exec.Command(manager, "stop", id).Output()
		if err != nil {
			log.Fatal(string(out))
		}

		cmd := exec.Command(manager, "inspect", "-f", "'{{ .Name }}'", id)
		output, _ := cmd.Output()
		name := string(output[2 : len(output)-2])

		cmd_args := []string{
			"create",
			"--clone", id,
			"--name", container.GenerateNewContainerName(),
		}

		labels := core.ContainerLabels{
			Managed: true,
			Userid:  os.Geteuid(),
		}

		found := false
		for _, pm := range []string{"apt", "aur", "dnf", "apk"} {
			legacy_name, distro := core.GetLegacyContainerNameAndDistro(pm)
			if name == legacy_name {
				found = true
				labels.PkgManager = settings.PackageManager(pm)
				labels.Distro = distro
				break
			}
		}

		if !found {
			// TODO Implement migration for named containers
			log.Fatal("Unknown container type for container: " + name)
		}

		cmd_args = append(cmd_args, "--additional-flags")
		cmd_args = append(cmd_args, strings.Join(labels.ToArguments(), " "))
		fmt.Println(cmd_args)

		out, err = exec.Command(settings.Cnf.DistroboxPath, cmd_args...).Output()
		if err != nil {
			log.Fatal(string(out))
		}

		out, err = exec.Command(manager, "rm", "-v", id).Output()
		if err != nil {
			log.Fatal(string(out))
		}
	}

	return nil
}
