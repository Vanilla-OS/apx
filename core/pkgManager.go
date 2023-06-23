package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install Packages and run commands inside a managed container.
*/

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/vanilla-os/apx/settings"
)

// PkgManager represents a package manager in Apx, a set of instructions to handle a package manager.
type PkgManager struct {
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
	BuiltIn       bool // If true, the package manager is built-in (stored in /usr/share/apx/pkg-managers) and cannot be removed by the user
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
	}
}

// LoadPkgManager loads a package manager from the specified path.
func LoadPkgManager(name string) (*PkgManager, error) {
	pkgManager, err := loadPkgManagerFromPath(filepath.Join(settings.Cnf.UserPkgManagersPath, name+".yaml"))
	if err != nil {
		pkgManager, err = loadPkgManagerFromPath(filepath.Join(settings.Cnf.PkgManagersPath, name+".yaml"))
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

	data, err := ioutil.ReadAll(f)
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

	filePath := filepath.Join(settings.Cnf.UserPkgManagersPath, pkgManager.Name+".yaml")
	err = ioutil.WriteFile(filePath, data, 0644)
	return err
}

// Remove removes the package manager from the specified path.
func (pkgManager *PkgManager) Remove() error {
	if pkgManager.BuiltIn {
		return errors.New("cannot remove built-in package manager")
	}

	filePath := filepath.Join(settings.Cnf.UserPkgManagersPath, pkgManager.Name+".yaml")
	err := os.Remove(filePath)
	return err
}

// GenCmd generates the command to run inside the container.
func (pkgManager *PkgManager) GenCmd(cmd string, args ...string) []string {
	finalArgs := make([]string, 0)

	if pkgManager.NeedSudo {
		finalArgs = append(finalArgs, "sudo")
	}

	finalArgs = append(finalArgs, pkgManager.Name)
	finalArgs = append(finalArgs, cmd)
	finalArgs = append(finalArgs, args...)

	return finalArgs
}

// ListPkgManagers lists all the package managers.
func ListPkgManagers() []*PkgManager {
	pkgManagers := make([]*PkgManager, 0)

	pkgManagersFromEtc := listPkgManagersFromPath(settings.Cnf.UserPkgManagersPath)
	pkgManagers = append(pkgManagers, pkgManagersFromEtc...)

	pkgManagersFromShare := listPkgManagersFromPath(settings.Cnf.PkgManagersPath)
	pkgManagers = append(pkgManagers, pkgManagersFromShare...)

	return pkgManagers
}

// listPkgManagersFromPath lists all the package managers from the specified path.
func listPkgManagersFromPath(path string) []*PkgManager {
	pkgManagers := make([]*PkgManager, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return pkgManagers
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			pkgManagerName := file.Name()[:len(file.Name())-5] // Remove the ".yaml" extension
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
