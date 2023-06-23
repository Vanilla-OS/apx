package core

import (
	"fmt"
	"log"
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
		InternalName: fmt.Sprintf("apx-%s-%s", stack.Name, name),
		Name:         name,
		Stack:        stack,
	}, nil
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
			InternalName: container.Name,
			Name:         container.Labels["name"],
			Stack:        stack,
			Status:       container.Status,
		}

		subsystems = append(subsystems, subsystem)
	}

	return subsystems, nil
}

func (s *SubSystem) Exec(args ...string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerExec(s.InternalName, args...)
}

func (s *SubSystem) Enter() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerEnter(s.InternalName)
}
