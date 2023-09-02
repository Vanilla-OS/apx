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
	output, err := exec.Command(apx.Cnf.DistroboxPath, "version").Output()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(string(output), "distrobox: ")
	if len(splitted) != 2 {
		return "", errors.New("can't retrieve distrobox version")
	}

	return splitted[1], nil
}

func (d *dbox) RunCommand(command string, args []string, engineFlags []string, useEngine bool, captureOutput bool, muteOutput bool, rootFull bool) ([]byte, error) {
	entrypoint := apx.Cnf.DistroboxPath
	if useEngine {
		entrypoint = d.EngineBinary
	}

	// we need to append engineFlags after command, otherwise they will be
	// ignored in commands like "enter"
	finalArgs := []string{command}

	// NOTE: for engine-specific commands, we need to use pkexec for rootfull
	//		 containers, since podman does not offer a dedicated flag for this.
	if rootFull && useEngine {
		entrypoint = "pkexec"
		finalArgs = []string{d.EngineBinary, command}
	}

	cmd := exec.Command(entrypoint, finalArgs...)

	if !captureOutput && !muteOutput {
		cmd.Stdout = os.Stdout
	}
	if !muteOutput {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	cmd.Env = os.Environ()

	// NOTE: the custom storage is not being used since it prevent other
	//		 utilities, like VSCode, to access the container.
	if d.Engine == "podman" {
		cmd.Env = append(cmd.Env, "CONTAINER_STORAGE_DRIVER="+apx.Cnf.StorageDriver)
		// cmd.Env = append(cmd.Env, "XDG_DATA_HOME="+apx.Cnf.ApxStoragePath)
	} else if d.Engine == "docker" {
		cmd.Env = append(cmd.Env, "DOCKER_STORAGE_DRIVER="+apx.Cnf.StorageDriver)
		// cmd.Env = append(cmd.Env, "DOCKER_DATA_ROOT="+apx.Cnf.ApxStoragePath)
	}

	if len(engineFlags) > 0 {
		cmd.Args = append(cmd.Args, "--additional-flags")
		cmd.Args = append(cmd.Args, strings.Join(engineFlags, " "))
	}

	// NOTE: the root flag is not being used by the Apx CLI, but it's useful
	//		 for those using Apx as a library, e.g. VSO.
	if rootFull && !useEngine {
		cmd.Args = append(cmd.Args, "--root")
	}

	cmd.Args = append(cmd.Args, args...)

	if os.Getenv("APX_VERBOSE") == "1" {
		fmt.Println("running command:")
		fmt.Println(cmd.String())
	}

	if captureOutput {
		output, err := cmd.Output()
		return output, err
	}

	err := cmd.Run()
	return nil, err
}

func (d *dbox) ListContainers(rootFull bool) ([]dboxContainer, error) {
	output, err := d.RunCommand("ps", []string{
		"-a",
		"--format", "{{.ID}}|{{.CreatedAt}}|{{.Status}}|{{.Labels}}|{{.Names}}",
	}, []string{}, true, true, false, rootFull)
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

func (d *dbox) GetContainer(name string, rootFull bool) (*dboxContainer, error) {
	containers, err := d.ListContainers(rootFull)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		// fmt.Println("found container", container.Name, "requested", name)
		if container.Name == name {
			return &container, nil
		}
	}

	return nil, errors.New("container not found")
}

func (d *dbox) ContainerDelete(name string, rootFull bool) error {
	_, err := d.RunCommand("rm", []string{
		"--force",
		name,
	}, []string{}, false, false, true, rootFull)
	return err
}

func (d *dbox) CreateContainer(name string, image string, additionalPackages []string, labels map[string]string, withInit bool, rootFull bool, unshared bool) error {
	args := []string{
		"--image", image,
		"--name", name,
		"--no-entry",
		"--yes",
	}

	if hasNvidiaGPU() {
		args = append(args, "--nvidia")
	}

	if withInit {
		args = append(args, "--init")
	}

	if unshared {
		args = append(args, "--unshare-all")
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

	_, err := d.RunCommand("create", args, engineFlags, false, true, true, rootFull)
	// fmt.Println(string(out))
	return err
}

func (d *dbox) RunContainerCommand(name string, command []string, rootFull bool) error {
	args := []string{
		"--name", name,
		"--",
	}

	args = append(args, command...)

	_, err := d.RunCommand("run", args, []string{}, false, false, false, rootFull)
	return err
}

func (d *dbox) ContainerExec(name string, captureOutput bool, muteOutput bool, rootFull bool, args ...string) (string, error) {
	finalArgs := []string{
		// "--verbose",
		name,
		"--",
	}

	finalArgs = append(finalArgs, args...)
	engineFlags := []string{}

	out, err := d.RunCommand("enter", finalArgs, engineFlags, false, captureOutput, muteOutput, rootFull)
	return string(out), err
}

func (d *dbox) ContainerEnter(name string, rootFull bool) error {
	finalArgs := []string{
		name,
	}

	engineFlags := []string{}

	_, err := d.RunCommand("enter", finalArgs, engineFlags, false, false, false, rootFull)
	return err
}

func (d *dbox) ContainerExport(name string, delete bool, rootFull bool, args ...string) error {
	finalArgs := []string{"distrobox-export"}

	if delete {
		finalArgs = append(finalArgs, "--delete")
	}

	finalArgs = append(finalArgs, args...)

	_, err := d.ContainerExec(name, false, true, rootFull, finalArgs...)
	return err
}

func (d *dbox) ContainerExportDesktopEntry(containerName string, appName string, label string, rootFull bool) error {
	args := []string{"--app", appName, "--export-label", label}
	return d.ContainerExport(containerName, false, rootFull, args...)
}

func (d *dbox) ContainerUnexportDesktopEntry(containerName string, appName string, rootFull bool) error {
	args := []string{"--app", appName}
	return d.ContainerExport(containerName, true, rootFull, args...)
}

func (d *dbox) ContainerExportBin(containerName string, binary string, exportPath string, rootFull bool) error {
	args := []string{"--bin", binary, "--export-path", exportPath}
	return d.ContainerExport(containerName, false, rootFull, args...)
}

func (d *dbox) ContainerUnexportBin(containerName string, binary string, exportPath string, rootFull bool) error {
	args := []string{"--bin", binary, "--export-path", exportPath}
	return d.ContainerExport(containerName, true, rootFull, args...)
}
