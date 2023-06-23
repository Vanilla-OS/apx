package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
)

type dbox struct {
	Engine       string
	EngineBinary string
	Version      string
}

type dboxContainer struct {
	ID        string
	CreatedAt string
	Status    string
	Labels    map[string]string
	Name      string
}

func NewDbox() (*dbox, error) {
	engineBinary, engine := getEngine()

	version, err := dboxGetVersion()
	if err != nil {
		return nil, err
	}

	return &dbox{
		Engine:       engine,
		EngineBinary: engineBinary,
		Version:      version,
	}, nil
}

func getEngine() (string, string) {
	podmanBinary, err := exec.LookPath("podman")
	if err == nil {
		return podmanBinary, "podman"
	}

	dockerBinary, err := exec.LookPath("docker")
	if err == nil {
		return dockerBinary, "docker"
	}

	log.Fatal("no container engine found. Please install Podman or Docker.")
	return "", ""
}

func dboxGetVersion() (version string, err error) {
	output, err := exec.Command(settings.Cnf.DistroboxPath, "version").Output()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(string(output), "distrobox: ")
	if len(splitted) != 2 {
		return "", errors.New("can't retrieve distrobox version")
	}

	return splitted[1], nil
}

func (d *dbox) RunCommand(command string, args []string, engineFlags []string, useEngine bool, capture_output bool) ([]byte, error) {
	entrypoint := settings.Cnf.DistroboxPath
	if useEngine {
		entrypoint = d.EngineBinary
	}

	// we need to append engineFlags after command, otherwise they will be
	// ignored in commands like "enter"
	finalArgs := []string{command}

	cmd := exec.Command(entrypoint, finalArgs...)

	if !capture_output {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = os.Environ()

	if d.Engine == "podman" {
		cmd.Env = append(cmd.Env, "CONTAINER_STORAGE_DRIVER="+settings.Cnf.StorageDriver)
		cmd.Env = append(cmd.Env, "XDG_DATA_HOME="+settings.Cnf.ApxStoragePath)

	} else if d.Engine == "docker" {
		cmd.Env = append(cmd.Env, "DOCKER_STORAGE_DRIVER="+settings.Cnf.StorageDriver) // TODO: check if this is correct
		cmd.Env = append(cmd.Env, "DOCKER_DATA_ROOT="+settings.Cnf.ApxStoragePath)     // TODO: check if this is correct
	}

	if len(engineFlags) > 0 {
		cmd.Args = append(cmd.Args, "--additional-flags")
		cmd.Args = append(cmd.Args, strings.Join(engineFlags, " "))
	}

	cmd.Args = append(cmd.Args, args...)
	// fmt.Println(cmd.String())

	if capture_output {
		output, err := cmd.Output()
		return output, err
	}

	err := cmd.Run()
	return nil, err
}

func (d *dbox) ListContainers() ([]dboxContainer, error) {
	output, err := d.RunCommand("ps", []string{
		"-a",
		"--format", "{{.ID}}|{{.CreatedAt}}|{{.Status}}|{{.Labels}}|{{.Names}}",
	}, []string{}, true, true)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(output), "\n")
	containers := []dboxContainer{}

	for _, row := range rows {
		if row == "" {
			continue
		}

		rowItems := strings.Split(row, "|")
		if len(rowItems) != 5 {
			continue
		}

		container := dboxContainer{
			ID:        strings.Trim(rowItems[0], "\""),
			CreatedAt: strings.Trim(rowItems[1], "\""),
			Status:    strings.Trim(rowItems[2], "\""),
			Name:      strings.Trim(rowItems[4], "\""),
			Labels:    map[string]string{},
		}

		// example labels: map[manager:apx name:alpine stack:alpine]
		labels := strings.ReplaceAll(rowItems[3], "map[", "")
		labels = strings.ReplaceAll(labels, "]", "")
		labelsItems := strings.Split(labels, " ")
		for _, label := range labelsItems {
			labelItems := strings.Split(label, ":")
			if len(labelItems) != 2 {
				continue
			}

			container.Labels[labelItems[0]] = labelItems[1]
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (d *dbox) CreateContainer(name string, image string, additionalPackages []string, labels map[string]string) error {
	args := []string{
		"--image", image,
		"--name", name,
		"--yes",
	}

	if len(additionalPackages) > 0 {
		args = append(args, "--additional-packages")
		args = append(args, strings.Join(additionalPackages, " "))
	}

	engineFlags := []string{}
	for key, value := range labels {
		engineFlags = append(engineFlags, fmt.Sprintf("--label=%s=%s", key, value))
	}
	engineFlags = append(engineFlags, "--label=manager=apx")

	_, err := d.RunCommand("create", args, engineFlags, false, true)
	// fmt.Println(string(out))
	return err
}

func (d *dbox) RunContainerCommand(name string, command []string) error {
	args := []string{
		"--name", name,
		"--",
	}

	args = append(args, command...)

	_, err := d.RunCommand("run", args, []string{}, false, false)
	return err
}

func (d *dbox) ContainerExec(name string, args ...string) error {
	finalArgs := []string{
		// "--verbose",
		name,
		"--",
	}

	finalArgs = append(finalArgs, args...)
	engineFlags := []string{}

	_, err := d.RunCommand("enter", finalArgs, engineFlags, false, false)
	return err
}

func (d *dbox) ContainerEnter(name string) error {
	finalArgs := []string{
		name,
	}

	engineFlags := []string{}

	_, err := d.RunCommand("enter", finalArgs, engineFlags, false, false)
	return err
}
