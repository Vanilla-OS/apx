package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

type SubSystem struct {
	InternalName string
	Name         string
	Stack        *Stack
	Status       string
}

func NewSubSystem(name string, stack *Stack) (*SubSystem, error) {
	return &SubSystem{
		InternalName: genInternalName(name),
		Name:         name,
		Stack:        stack,
	}, nil
}

func genInternalName(name string) string {
	return fmt.Sprintf("apx-%s", strings.ReplaceAll(strings.ToLower(name), " ", "-"))
}

func (s *SubSystem) Create() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	err = dbox.CreateContainer(
		s.InternalName,
		s.Stack.Base,
		[]string{},
		map[string]string{
			"stack": s.Stack.Name,
			"name":  s.Name,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func LoadSubSystem(name string) (*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
		return nil, err
	}

	container, err := dbox.GetContainer(fmt.Sprintf("apx-%s", name))
	if err != nil {
		return nil, err
	}

	stack, err := LoadStack(container.Labels["stack"])
	if err != nil {
		return nil, err
	}
	return &SubSystem{
		InternalName: genInternalName(name),
		Name:         container.Labels["name"],
		Stack:        stack,
		Status:       container.Status,
	}, nil
}

func ListSubSystems() ([]*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
		return nil, err
	}

	containers, err := dbox.ListContainers()
	if err != nil {
		return nil, err
	}

	subsystems := []*SubSystem{}
	for _, container := range containers {
		if _, ok := container.Labels["name"]; !ok {
			log.Printf("Container %s has no name label, skipping", container.Name)
			continue
		}

		stack, err := LoadStack(container.Labels["stack"])
		if err != nil {
			log.Printf("Error loading stack %s: %s", container.Labels["stack"], err)
			continue
		}

		subsystem := &SubSystem{
			Name:   container.Labels["name"],
			Stack:  stack,
			Status: container.Status,
		}

		subsystems = append(subsystems, subsystem)
	}

	return subsystems, nil
}

func (s *SubSystem) Exec(captureOutput bool, args ...string) (string, error) {
	dbox, err := NewDbox()
	if err != nil {
		return "", err
	}

	out, err := dbox.ContainerExec(s.InternalName, captureOutput, args...)
	if err != nil {
		return "", err
	}

	if captureOutput {
		return out, nil
	}

	return "", nil
}

func (s *SubSystem) Enter() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerEnter(s.InternalName)
}

func (s *SubSystem) Remove() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerDelete(s.InternalName)
}

func (s *SubSystem) Reset() error {
	err := s.Remove()
	if err != nil {
		return err
	}

	return s.Create()
}

func (s *SubSystem) ExportDesktopEntry(appName string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerExportDesktopEntry(s.InternalName, appName, s.Name)
}

func (s *SubSystem) ExportBin(binary string, exportPath string) error {
	if !strings.HasPrefix(binary, "/") {
		binaryPath, err := s.Exec(true, "which", binary)
		if err != nil {
			return err
		}

		binary = binaryPath
		binary = strings.TrimSuffix(binary, "\r\n")
	}

	binaryName := filepath.Base(binary)

	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if exportPath == "" {
		exportPath = fmt.Sprintf("%s/%s", homeDir, ".local/bin")
	}

	joinedPath := filepath.Join(exportPath, binaryName)
	if _, err := os.Stat(joinedPath); err == nil {
		tmpExportPath := fmt.Sprintf("/tmp/%s", uuid.New().String())
		err = os.MkdirAll(tmpExportPath, 0755)
		if err != nil {
			return err
		}

		err = dbox.ContainerExportBin(s.InternalName, binary, tmpExportPath)
		if err != nil {
			return err
		}

		err = CopyFile(filepath.Join(tmpExportPath, binaryName), filepath.Join(exportPath, fmt.Sprintf("%s-%s", binaryName, s.InternalName)))
		if err != nil {
			return err
		}

		err = os.RemoveAll(tmpExportPath)
		if err != nil {
			return err
		}

		err = os.Chmod(filepath.Join(exportPath, fmt.Sprintf("%s-%s", binaryName, s.InternalName)), 0755)
		if err != nil {
			return err
		}

		return nil
	}

	err = os.MkdirAll(exportPath, 0755)
	if err != nil {
		return err
	}

	err = dbox.ContainerExportBin(s.InternalName, binary, exportPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubSystem) UnexportDesktopEntry(appName string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerUnexportDesktopEntry(s.InternalName, appName)
}

func (s *SubSystem) UnexportBin(binary string, exportPath string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerUnexportBin(s.InternalName, binary, exportPath)
}
