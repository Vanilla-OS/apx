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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// PkgManager represents a package manager in Apx, a set of instructions to handle a package manager.
type PkgManager struct {
	// Model values:
	// 1: name will be used as the main command;
	// 2: each command is the whole command
	// Default: 2
	// DEPRECATION WARNING: Model 1 will be removed in the future, please
	// update your configuration files to use model 2.
	Model         int
	Name          string
	NeedSudo      bool
	CmdAutoRemove string
	CmdClean      string
	CmdInstall    string
	CmdList       string
	CmdPurge      string
	CmdRemove     string
	CmdSearch     string
	CmdShow       string
	CmdUpdate     string
	CmdUpgrade    string

	// BuiltIn:
	// If true, the package manager is built-in (stored in
	// /usr/share/apx/pkg-managers) and cannot be removed by the user
	BuiltIn bool
}

// NewPkgManager creates a new PkgManager instance.
func NewPkgManager(name string, needSudo bool, autoRemove, clean, install, list, purge, remove, search, show, update, upgrade string, builtIn bool) *PkgManager {
	return &PkgManager{
		Name:          name,
		NeedSudo:      needSudo,
		CmdAutoRemove: autoRemove,
		CmdClean:      clean,
		CmdInstall:    install,
		CmdList:       list,
		CmdPurge:      purge,
		CmdRemove:     remove,
		CmdSearch:     search,
		CmdShow:       show,
		CmdUpdate:     update,
		CmdUpgrade:    upgrade,
		BuiltIn:       builtIn,
		Model:         2,
	}
}

// LoadPkgManager loads a package manager from the specified path.
func LoadPkgManager(name string) (*PkgManager, error) {
	userPkgFile := ChooseYamlFile(apx.Cnf.UserPkgManagersPath, name)
	pkgManager, err := loadPkgManagerFromPath(userPkgFile)
	if err != nil {
		pkgFile := ChooseYamlFile(apx.Cnf.PkgManagersPath, name)
		pkgManager, err = loadPkgManagerFromPath(pkgFile)
	}
	return pkgManager, err
}

// loadPkgManagerFromPath loads a package manager from the specified path.
func loadPkgManagerFromPath(path string) (*PkgManager, error) {
	pkgManager := &PkgManager{}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("package manager not found")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, pkgManager)
	if err != nil {
		return nil, err
	}

	return pkgManager, nil
}

// Save saves the package manager to the specified path.
func (pkgManager *PkgManager) Save() error {
	data, err := yaml.Marshal(pkgManager)
	if err != nil {
		return err
	}

	filePath := ChooseYamlFile(apx.Cnf.UserPkgManagersPath, pkgManager.Name)
	err = os.WriteFile(filePath, data, 0644)
	return err
}

// Remove removes the package manager from the specified path.
func (pkgManager *PkgManager) Remove() error {
	if pkgManager.BuiltIn {
		return errors.New("cannot remove built-in package manager")
	}

	filePath := ChooseYamlFile(apx.Cnf.UserPkgManagersPath, pkgManager.Name)
	err := os.Remove(filePath)
	return err
}

// GenCmd generates the command to run inside the container.
func (pkgManager *PkgManager) GenCmd(cmd string, args ...string) []string {
	finalArgs := make([]string, 0)

	if pkgManager.NeedSudo {
		finalArgs = append(finalArgs, "sudo")
	}

	if pkgManager.Model == 0 || pkgManager.Model == 1 {
		// no-translate (DEPRECATION WARNING)
		fmt.Println("!!! DEPRECATION WARNING: Model 1 will be removed in the future, please update your Apx package manager to use model 2.")
		finalArgs = append(finalArgs, pkgManager.Name)
		finalArgs = append(finalArgs, cmd)
		finalArgs = append(finalArgs, args...)
	} else {
		cmdItems := strings.Fields(cmd)
		finalArgs = append(finalArgs, cmdItems...)
		finalArgs = append(finalArgs, args...)
	}

	return finalArgs
}

// ListPkgManagers lists all the package managers.
func ListPkgManagers() []*PkgManager {
	pkgManagers := make([]*PkgManager, 0)

	pkgManagersFromEtc := listPkgManagersFromPath(apx.Cnf.UserPkgManagersPath)
	pkgManagers = append(pkgManagers, pkgManagersFromEtc...)

	if apx.Cnf.PkgManagersPath == apx.Cnf.UserPkgManagersPath {
		// user install
		return pkgManagers
	}

	pkgManagersFromShare := listPkgManagersFromPath(apx.Cnf.PkgManagersPath)
	pkgManagers = append(pkgManagers, pkgManagersFromShare...)

	return pkgManagers
}

// listPkgManagersFromPath lists all the package managers from the specified path.
func listPkgManagersFromPath(path string) []*PkgManager {
	pkgManagers := make([]*PkgManager, 0)

	files, err := os.ReadDir(path)
	if err != nil {
		return pkgManagers
	}

	for _, file := range files {
		extension := filepath.Ext(file.Name())

		if !file.IsDir() && (extension == ".yaml" || extension == ".yml") {
			// Remove the ".yaml" or ".yml" extension
			pkgManagerName := file.Name()[:(len(file.Name()) - len(extension))]
			pkgManager, err := LoadPkgManager(pkgManagerName)
			if err == nil {
				pkgManagers = append(pkgManagers, pkgManager)
			}
		}
	}

	return pkgManagers
}

// PkgManagerExists checks if the package manager exists.
func PkgManagerExists(name string) bool {
	_, err := LoadPkgManager(name)
	return err == nil
}

// LoadPkgManagerFromPath loads a package manager from the specified path.
func LoadPkgManagerFromPath(path string) (*PkgManager, error) {
	pkgManager := &PkgManager{}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("package manager not found")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, pkgManager)
	if err != nil {
		return nil, err
	}

	if pkgManager.Model == 0 {
		pkgManager.Model = 1 // assuming old model if not specified
	}

	return pkgManager, nil
}

// Export exports the package manager to the specified path.
func (pkgManager *PkgManager) Export(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	filePath := ChooseYamlFile(path, pkgManager.Name)
	data, err := yaml.Marshal(pkgManager)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
