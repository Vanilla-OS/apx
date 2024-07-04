package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Apx is a wrapper around multiple package managers to install Packages and run commands inside a managed container.
*/

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Stack represents a stack in Apx, a set of instructions to build a container.
type Stack struct {
	Name       string
	Base       string
	Packages   []string
	PkgManager string
	BuiltIn    bool // If true, the stack is built-in (stored in /usr/share/apx/stacks) and cannot be removed by the user
}

// NewStack creates a new Stack instance.
func NewStack(name, base string, packages []string, pkgManager string, builtIn bool) *Stack {
	return &Stack{
		Name:       name,
		Base:       base,
		Packages:   packages,
		PkgManager: pkgManager,
		BuiltIn:    builtIn,
	}
}

// LoadStack loads a stack from the specified path.
func LoadStack(name string) (*Stack, error) {
	usrStackFile := SelectYamlFile(apx.Cnf.UserStacksPath, name)
	stack, err := LoadStackFromPath(usrStackFile)
	if err != nil {
		stackFile := SelectYamlFile(apx.Cnf.StacksPath, name)
		stack, err = LoadStackFromPath(stackFile)
	}
	return stack, err
}

// LoadStackFromPath loads a stack from the specified path.
func LoadStackFromPath(path string) (*Stack, error) {
	stack := &Stack{}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("stack not found")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, stack)
	if err != nil {
		return nil, err
	}

	if stack.Name == "" || stack.Base == "" || stack.PkgManager == "" {
		return nil, errors.New("invalid stack file")
	}

	return stack, nil
}

// Save saves the stack to a YAML file.
func (stack *Stack) Save() error {
	data, err := yaml.Marshal(stack)
	if err != nil {
		return err
	}

	filePath := SelectYamlFile(apx.Cnf.UserStacksPath, stack.Name)
	err = os.WriteFile(filePath, data, 0644)
	return err
}

// GetPkgManager returns the package manager of the stack.
func (stack *Stack) GetPkgManager() (*PkgManager, error) {
	pkgManager, err := LoadPkgManager(stack.PkgManager)
	if err != nil {
		return nil, err
	}

	return pkgManager, nil
}

// Remove removes the stack file.
func (stack *Stack) Remove() error {
	if stack.BuiltIn {
		return errors.New("cannot remove built-in stack")
	}

	filePath := SelectYamlFile(apx.Cnf.UserStacksPath, stack.Name)
	err := os.Remove(filePath)
	return err
}

// Export exports the stack YAML to the specified path.
func (stack *Stack) Export(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	filePath := SelectYamlFile(path, stack.Name)
	data, err := yaml.Marshal(stack)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ListStacks returns a list of all stacks.
func ListStacks() []*Stack {
	stacks := make([]*Stack, 0)

	stacksFromEtc := listStacksFromPath(apx.Cnf.UserStacksPath)
	stacks = append(stacks, stacksFromEtc...)

	if apx.Cnf.UserStacksPath == apx.Cnf.StacksPath {
		// user install
		return stacks
	}

	stacksFromShare := listStacksFromPath(apx.Cnf.StacksPath)
	stacks = append(stacks, stacksFromShare...)

	return stacks
}

// ListStackForPkgManager returns a list of stacks for the specified package manager.
func ListStackForPkgManager(pkgManager string) []*Stack {
	stacks := make([]*Stack, 0)

	stacksFromEtc := listStacksFromPath(apx.Cnf.UserStacksPath)
	for _, stack := range stacksFromEtc {
		if stack.PkgManager == pkgManager {
			stacks = append(stacks, stack)
		}
	}

	stacksFromShare := listStacksFromPath(apx.Cnf.StacksPath)
	for _, stack := range stacksFromShare {
		if stack.PkgManager == pkgManager {
			stacks = append(stacks, stack)
		}
	}

	return stacks
}

// listStacksFromPath returns a list of stacks from the specified path.
// this func does not return an error, since Apx is meant to be portable and
// the main directory can be missing, while the user directory is always created.
func listStacksFromPath(path string) []*Stack {
	stacks := make([]*Stack, 0)

	files, err := os.ReadDir(path)
	if err != nil {
		return stacks
	}

	for _, file := range files {
		extension := filepath.Ext(file.Name())

		if !file.IsDir() && (extension == ".yaml" || extension == ".yml") {
			// Remove the ".yaml" or ".yml" extension
			stackName := file.Name()[:(len(file.Name()) - len(extension))]
			stack, err := LoadStack(stackName)
			if err == nil {
				stacks = append(stacks, stack)
			}
		}
	}

	return stacks
}

// StackExists checks if a stack exists.
func StackExists(name string) bool {
	s, _ := LoadStack(name)
	return s != nil
}
