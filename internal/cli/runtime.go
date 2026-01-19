package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Pietro di Caprio <pietro@fabricators.ltd>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/vanilla-os/apx/v3/core"
)

// RootCmd does not have a Run method to trigger help automatically

func (c *SubsystemEnterCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}
	err = subSystem.Enter()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.enteringContainer"), err)
	}
	return nil
}

func (c *SubsystemRunCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}
	_, err = subSystem.Exec(false, false, c.Args...)
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.executingCommand"), err)
	}
	return nil
}

func (c *SubsystemInstallCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	pkgManager, err := subSystem.Stack.GetPkgManager()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.cantAccessPkgManager"), err)
	}

	finalArgs := pkgManager.GenCmd(pkgManager.CmdInstall, c.Args...)
	_, err = subSystem.Exec(false, false, finalArgs...)
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.executingCommand"), err)
	}

	if !c.NoExport {
		exportedN, err := subSystem.ExportDesktopEntries(c.Args...)
		if err == nil {
			Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.exportedApps"), exportedN)
		}
	}
	return nil
}

func (c *SubsystemRemoveCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	pkgManager, err := subSystem.Stack.GetPkgManager()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.cantAccessPkgManager"), err)
	}

	exportedN, err := subSystem.UnexportDesktopEntries(c.Args...)
	if err == nil {
		Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.unexportedApps"), exportedN)
	}

	finalArgs := pkgManager.GenCmd(pkgManager.CmdRemove, c.Args...)
	_, err = subSystem.Exec(false, false, finalArgs...)
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.executingCommand"), err)
	}
	return nil
}

func (c *SubsystemUpdateCmd) Run() error {
	return genericPkgManagerCommand(c.Name, "update")
}

func (c *SubsystemUpgradeCmd) Run() error {
	return genericPkgManagerCommand(c.Name, "upgrade")
}

func (c *SubsystemListCmd) Run() error {
	return genericPkgManagerCommand(c.Name, "list")
}

func (c *SubsystemSearchCmd) Run() error {
	return genericPkgManagerArgsCommand(c.Name, "search", c.Args)
}

func (c *SubsystemShowCmd) Run() error {
	return genericPkgManagerArgsCommand(c.Name, "show", c.Args)
}

func (c *SubsystemAutoRemoveCmd) Run() error {
	return genericPkgManagerCommand(c.Name, "autoremove")
}

func (c *SubsystemCleanCmd) Run() error {
	return genericPkgManagerCommand(c.Name, "clean")
}

func (c *SubsystemPurgeCmd) Run() error {
	return genericPkgManagerArgsCommand(c.Name, "purge", c.Args)
}

func (c *SubsystemStartCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}
	Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.startingContainer"), subSystem.Name)
	err = subSystem.Start()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.startingContainer"), err)
	}
	Apx.Log.Info(Apx.LC.Get("runtimeCommand.info.startedContainer"))
	return nil
}

func (c *SubsystemStopCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}
	Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.stoppingContainer"), subSystem.Name)
	err = subSystem.Stop()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.stoppingContainer"), err)
	}
	Apx.Log.Info(Apx.LC.Get("runtimeCommand.info.stoppedContainer"))
	return nil
}

func (c *SubsystemExportCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	if c.AppName == "" && c.Bin == "" {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.noAppNameOrBin"))
	}

	if c.AppName != "" && c.Bin != "" {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.sameAppOrBin"))
	}

	if c.AppName != "" {
		err := subSystem.ExportDesktopEntry(c.AppName)
		if err != nil {
			return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.exportingApp"), err)
		}
		Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.exportedApp"), c.AppName)
	} else {
		err := subSystem.ExportBin(c.Bin, c.BinOutput)
		if err != nil {
			return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.exportingBin"), err)
		}
		Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.exportedBin"), c.Bin)
	}
	return nil
}

func (c *SubsystemUnexportCmd) Run() error {
	subSystem, err := core.LoadSubSystem(c.Name, false)
	if err != nil {
		return err
	}

	if c.AppName == "" && c.Bin == "" {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.noAppNameOrBin"))
	}

	if c.AppName != "" && c.Bin != "" {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.sameAppOrBin"))
	}

	if c.AppName != "" {
		err := subSystem.UnexportDesktopEntry(c.AppName)
		if err != nil {
			return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.unexportingApp"), err)
		}
		Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.unexportedApp"), c.AppName)
	} else {
		err := subSystem.UnexportBin(c.Bin, c.BinOutput)
		if err != nil {
			return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.unexportingBin"), err)
		}
		Apx.Log.Infof(Apx.LC.Get("runtimeCommand.info.unexportedBin"), c.Bin)
	}
	return nil
}

// Helpers

func genericPkgManagerCommand(subsystemName string, action string) error {
	subSystem, err := core.LoadSubSystem(subsystemName, false)
	if err != nil {
		return err
	}
	pkgManager, err := subSystem.Stack.GetPkgManager()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.cantAccessPkgManager"), err)
	}

	cmdStr, err := pkgManagerCommands(pkgManager, action)
	if err != nil {
		return err
	}

	finalArgs := pkgManager.GenCmd(cmdStr)
	_, err = subSystem.Exec(false, false, finalArgs...)
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.executingCommand"), err)
	}
	return nil
}

func genericPkgManagerArgsCommand(subsystemName string, action string, args []string) error {
	subSystem, err := core.LoadSubSystem(subsystemName, false)
	if err != nil {
		return err
	}
	pkgManager, err := subSystem.Stack.GetPkgManager()
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.cantAccessPkgManager"), err)
	}

	cmdStr, err := pkgManagerCommands(pkgManager, action)
	if err != nil {
		return err
	}

	finalArgs := pkgManager.GenCmd(cmdStr, args...)
	_, err = subSystem.Exec(false, false, finalArgs...)
	if err != nil {
		return fmt.Errorf(Apx.LC.Get("runtimeCommand.error.executingCommand"), err)
	}
	return nil
}

func pkgManagerCommands(pkgManager *core.PkgManager, command string) (string, error) {
	switch command {
	case "autoremove":
		return pkgManager.CmdAutoRemove, nil
	case "clean":
		return pkgManager.CmdClean, nil
	case "install":
		return pkgManager.CmdInstall, nil
	case "list":
		return pkgManager.CmdList, nil
	case "purge":
		return pkgManager.CmdPurge, nil
	case "remove":
		return pkgManager.CmdRemove, nil
	case "search":
		return pkgManager.CmdSearch, nil
	case "show":
		return pkgManager.CmdShow, nil
	case "update":
		return pkgManager.CmdUpdate, nil
	case "upgrade":
		return pkgManager.CmdUpgrade, nil
	default:
		return "", fmt.Errorf(Apx.LC.Get("apx.errors.unknownCommand"), command)
	}
}
