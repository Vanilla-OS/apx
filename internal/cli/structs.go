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
	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/cli"
)

var Apx *app.App

type RootCmd struct {
	cli.Base
	Version string

	Stacks      StacksCmd      `cmd:"stacks" help:"pr:apx.cmd.stacks"`
	Subsystems  SubsystemsCmd  `cmd:"subsystems" help:"pr:apx.cmd.subsystems"`
	PkgManagers PkgManagersCmd `cmd:"pkgmanagers" help:"pr:apx.cmd.pkgmanagers"`

	DynamicSubsystems *map[string]*SubsystemCmd `cmd:"*" help:"apx.subsystem"`
}

type SubsystemCmd struct {
	cli.Base
	Name string `json:"-"`

	Enter      SubsystemEnterCmd      `cmd:"enter" help:"pr:apx.cmd.subsystem.enter"`
	Run        SubsystemRunCmd        `cmd:"run" help:"pr:apx.cmd.subsystem.run"`
	Install    SubsystemInstallCmd    `cmd:"install" help:"pr:apx.cmd.subsystem.install"`
	Remove     SubsystemRemoveCmd     `cmd:"remove" help:"pr:apx.cmd.subsystem.remove"`
	Update     SubsystemUpdateCmd     `cmd:"update" help:"pr:apx.cmd.subsystem.update"`
	Upgrade    SubsystemUpgradeCmd    `cmd:"upgrade" help:"pr:apx.cmd.subsystem.upgrade"`
	List       SubsystemListCmd       `cmd:"list" help:"pr:apx.cmd.subsystem.list"`
	Search     SubsystemSearchCmd     `cmd:"search" help:"pr:apx.cmd.subsystem.search"`
	Show       SubsystemShowCmd       `cmd:"show" help:"pr:apx.cmd.subsystem.show"`
	Export     SubsystemExportCmd     `cmd:"export" help:"pr:apx.cmd.subsystem.export"`
	Unexport   SubsystemUnexportCmd   `cmd:"unexport" help:"pr:apx.cmd.subsystem.unexport"`
	Start      SubsystemStartCmd      `cmd:"start" help:"pr:apx.cmd.subsystem.start"`
	Stop       SubsystemStopCmd       `cmd:"stop" help:"pr:apx.cmd.subsystem.stop"`
	AutoRemove SubsystemAutoRemoveCmd `cmd:"autoremove" help:"pr:apx.cmd.subsystem.autoremove"`
	Clean      SubsystemCleanCmd      `cmd:"clean" help:"pr:apx.cmd.subsystem.clean"`
	Purge      SubsystemPurgeCmd      `cmd:"purge" help:"pr:apx.cmd.subsystem.purge"`
}

type SubsystemEnterCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemRunCmd struct {
	cli.Base
	Name string   `json:"-"`
	Args []string `arg:"" optional:"" name:"command" help:"pr:apx.arg.command"`
}

type SubsystemInstallCmd struct {
	cli.Base
	Name     string   `json:"-"`
	NoExport bool     `flag:"short:n, long:no-export, name:pr:apx.cmd.subsystem.install.options.noExport"`
	Args     []string `arg:"" optional:"" name:"packages" help:"pr:apx.arg.packages"`
}

type SubsystemRemoveCmd struct {
	cli.Base
	Name string   `json:"-"`
	Args []string `arg:"" optional:"" name:"packages" help:"pr:apx.arg.packages"`
}

type SubsystemUpdateCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemUpgradeCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemListCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemSearchCmd struct {
	cli.Base
	Name string   `json:"-"`
	Args []string `arg:"" optional:"" name:"query" help:"pr:apx.arg.query"`
}

type SubsystemShowCmd struct {
	cli.Base
	Name string   `json:"-"`
	Args []string `arg:"" optional:"" name:"package" help:"pr:apx.arg.package"`
}

type SubsystemExportCmd struct {
	cli.Base
	Name      string   `json:"-"`
	AppName   string   `flag:"short:a, long:app-name, name:pr:apx.cmd.subsystem.export.options.appName"`
	Bin       string   `flag:"short:b, long:bin, name:pr:apx.cmd.subsystem.export.options.bin"`
	BinOutput string   `flag:"short:o, long:bin-output, name:pr:apx.cmd.subsystem.export.options.binOutput"`
	Args      []string `arg:"" optional:"" name:"applications" help:"pr:apx.arg.applications"`
}

type SubsystemUnexportCmd struct {
	cli.Base
	Name      string   `json:"-"`
	AppName   string   `flag:"short:a, long:app-name, name:pr:apx.cmd.subsystem.unexport.options.appName"`
	Bin       string   `flag:"short:b, long:bin, name:pr:apx.cmd.subsystem.unexport.options.bin"`
	BinOutput string   `flag:"short:o, long:bin-output, name:pr:apx.cmd.subsystem.unexport.options.binOutput"`
	Args      []string `arg:"" optional:"" name:"applications" help:"pr:apx.arg.applications"`
}

type SubsystemStartCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemStopCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemAutoRemoveCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemCleanCmd struct {
	cli.Base
	Name string `json:"-"`
}

type SubsystemPurgeCmd struct {
	cli.Base
	Name string   `json:"-"`
	Args []string `arg:"" optional:"" name:"packages" help:"pr:apx.arg.packages"`
}

// Stacks

type StacksCmd struct {
	cli.Base
	List   StacksListCmd   `cmd:"list" help:"pr:apx.cmd.stacks.list"`
	Show   StacksShowCmd   `cmd:"show" help:"pr:apx.cmd.stacks.show"`
	New    StacksNewCmd    `cmd:"new" help:"pr:apx.cmd.stacks.new"`
	Update StacksUpdateCmd `cmd:"update" help:"pr:apx.cmd.stacks.update"`
	Rm     StacksRmCmd     `cmd:"rm" help:"pr:apx.cmd.stacks.rm"`
	Export StacksExportCmd `cmd:"export" help:"pr:apx.cmd.stacks.export"`
	Import StacksImportCmd `cmd:"import" help:"pr:apx.cmd.stacks.import"`
}

type StacksListCmd struct {
	cli.Base
	Json bool `flag:"short:j, long:json, name:pr:apx.cmd.stacks.list.options.json"`
}

type StacksShowCmd struct {
	cli.Base
	Args []string `arg:"" optional:"" name:"stack" help:"pr:apx.arg.stack"`
}

type StacksNewCmd struct {
	cli.Base
	NoPrompt   bool   `flag:"short:y, long:no-prompt, name:pr:apx.cmd.stacks.new.options.noPrompt"`
	Name       string `flag:"short:n, long:name, name:pr:apx.cmd.stacks.new.options.name"`
	BaseImage  string `flag:"short:b, long:base, name:pr:apx.cmd.stacks.new.options.base"`
	Packages   string `flag:"short:p, long:packages, name:pr:apx.cmd.stacks.new.options.packages"`
	PkgManager string `flag:"short:k, long:pkg-manager, name:pr:apx.cmd.stacks.new.options.pkgManager"`
}

type StacksUpdateCmd struct {
	cli.Base
	NoPrompt   bool     `flag:"short:y, long:no-prompt, name:pr:apx.cmd.stacks.update.options.noPrompt"`
	Name       string   `flag:"short:n, long:name, name:pr:apx.cmd.stacks.update.options.name"`
	BaseImage  string   `flag:"short:b, long:base, name:pr:apx.cmd.stacks.update.options.base"`
	Packages   string   `flag:"short:p, long:packages, name:pr:apx.cmd.stacks.update.options.packages"`
	PkgManager string   `flag:"short:k, long:pkg-manager, name:pr:apx.cmd.stacks.update.options.pkgManager"`
	Args       []string `arg:"" optional:"" name:"stack" help:"pr:apx.arg.stack"`
}

type PkgManagersShowCmd struct {
	cli.Base
	Args []string `arg:"" optional:"" name:"pkgmanager" help:"pr:apx.arg.pkgmanager"`
}

type PkgManagersUpdateCmd struct {
	cli.Base
	NoPrompt   bool     `flag:"short:y, long:no-prompt, name:pr:apx.cmd.pkgmanagers.new.options.noPrompt"`
	Name       string   `flag:"short:n, long:name, name:pr:apx.cmd.pkgmanagers.new.options.name"`
	NeedSudo   bool     `flag:"short:S, long:need-sudo, name:pr:apx.cmd.pkgmanagers.new.options.needSudo"`
	AutoRemove string   `flag:"short:a, long:autoremove, name:pr:apx.cmd.pkgmanagers.new.options.autoremove"`
	Clean      string   `flag:"short:c, long:clean, name:pr:apx.cmd.pkgmanagers.new.options.clean"`
	Install    string   `flag:"short:i, long:install, name:pr:apx.cmd.pkgmanagers.new.options.install"`
	List       string   `flag:"short:l, long:list, name:pr:apx.cmd.pkgmanagers.new.options.list"`
	Purge      string   `flag:"short:p, long:purge, name:pr:apx.cmd.pkgmanagers.new.options.purge"`
	Remove     string   `flag:"short:r, long:remove, name:pr:apx.cmd.pkgmanagers.new.options.remove"`
	Search     string   `flag:"short:s, long:search, name:pr:apx.cmd.pkgmanagers.new.options.search"`
	Show       string   `flag:"short:w, long:show, name:pr:apx.cmd.pkgmanagers.new.options.show"`
	Update     string   `flag:"short:u, long:update, name:pr:apx.cmd.pkgmanagers.new.options.update"`
	Upgrade    string   `flag:"short:U, long:upgrade, name:pr:apx.cmd.pkgmanagers.new.options.upgrade"`
	Args       []string `arg:"" optional:"" name:"pkgmanager" help:"pr:apx.arg.pkgmanager"`
}

type StacksRmCmd struct {
	cli.Base
	Name  string `flag:"short:n, long:name, name:pr:apx.cmd.stacks.rm.options.name"`
	Force bool   `flag:"short:f, long:force, name:pr:apx.cmd.stacks.rm.options.force"`
}

type StacksExportCmd struct {
	cli.Base
	Name   string `flag:"short:n, long:name, name:pr:apx.cmd.stacks.export.options.name"`
	Output string `flag:"short:o, long:output, name:pr:apx.cmd.stacks.export.options.output"`
}

type StacksImportCmd struct {
	cli.Base
	Input string `flag:"short:i, long:input, name:pr:apx.cmd.stacks.import.options.input"`
}

// Subsystems

type SubsystemsCmd struct {
	cli.Base
	List  SubsystemsListCmd  `cmd:"list" help:"pr:apx.cmd.subsystems.list"`
	New   SubsystemsNewCmd   `cmd:"new" help:"pr:apx.cmd.subsystems.new"`
	Rm    SubsystemsRmCmd    `cmd:"rm" help:"pr:apx.cmd.subsystems.rm"`
	Reset SubsystemsResetCmd `cmd:"reset" help:"pr:apx.cmd.subsystems.reset"`
}

type SubsystemsListCmd struct {
	cli.Base
	Json bool `flag:"short:j, long:json, name:pr:apx.cmd.subsystem.list.options.json"`
}

type SubsystemsNewCmd struct {
	cli.Base
	Stack string `flag:"short:s, long:stack, name:pr:apx.cmd.subsystem.new.options.stack"`
	Name  string `flag:"short:n, long:name, name:pr:apx.cmd.subsystem.new.options.name"`
	Home  string `flag:"short:H, long:home, name:pr:apx.cmd.subsystem.new.options.home"`
	Init  bool   `flag:"short:i, long:init, name:pr:apx.cmd.subsystem.new.options.init"`
}

type SubsystemsRmCmd struct {
	cli.Base
	Name  string `flag:"short:n, long:name, name:pr:apx.cmd.subsystem.rm.options.name"`
	Force bool   `flag:"short:f, long:force, name:pr:apx.cmd.subsystem.rm.options.force"`
}

type SubsystemsResetCmd struct {
	cli.Base
	Name  string `flag:"short:n, long:name, name:pr:apx.cmd.subsystem.reset.options.name"`
	Force bool   `flag:"short:f, long:force, name:pr:apx.cmd.subsystem.reset.options.force"`
}

// PkgManagers

type PkgManagersCmd struct {
	cli.Base
	List   PkgManagersListCmd   `cmd:"list" help:"pr:apx.cmd.pkgmanagers.list"`
	Show   PkgManagersShowCmd   `cmd:"show" help:"pr:apx.cmd.pkgmanagers.show"`
	New    PkgManagersNewCmd    `cmd:"new" help:"pr:apx.cmd.pkgmanagers.new"`
	Rm     PkgManagersRmCmd     `cmd:"rm" help:"pr:apx.cmd.pkgmanagers.rm"`
	Export PkgManagersExportCmd `cmd:"export" help:"pr:apx.cmd.pkgmanagers.export"`
	Import PkgManagersImportCmd `cmd:"import" help:"pr:apx.cmd.pkgmanagers.import"`
	Update PkgManagersUpdateCmd `cmd:"update" help:"pr:apx.cmd.pkgmanagers.update"`
}

type PkgManagersListCmd struct {
	cli.Base
	Json bool `flag:"short:j, long:json, name:pr:apx.cmd.pkgmanagers.list.options.json"`
}

type PkgManagersNewCmd struct {
	cli.Base
	NoPrompt   bool   `flag:"short:y, long:no-prompt, name:pr:apx.cmd.pkgmanagers.new.options.noPrompt"`
	Name       string `flag:"short:n, long:name, name:pr:apx.cmd.pkgmanagers.new.options.name"`
	NeedSudo   bool   `flag:"short:S, long:need-sudo, name:pr:apx.cmd.pkgmanagers.new.options.needSudo"`
	AutoRemove string `flag:"short:a, long:autoremove, name:pr:apx.cmd.pkgmanagers.new.options.autoremove"`
	Clean      string `flag:"short:c, long:clean, name:pr:apx.cmd.pkgmanagers.new.options.clean"`
	Install    string `flag:"short:i, long:install, name:pr:apx.cmd.pkgmanagers.new.options.install"`
	List       string `flag:"short:l, long:list, name:pr:apx.cmd.pkgmanagers.new.options.list"`
	Purge      string `flag:"short:p, long:purge, name:pr:apx.cmd.pkgmanagers.new.options.purge"`
	Remove     string `flag:"short:r, long:remove, name:pr:apx.cmd.pkgmanagers.new.options.remove"`
	Search     string `flag:"short:s, long:search, name:pr:apx.cmd.pkgmanagers.new.options.search"`
	Show       string `flag:"short:w, long:show, name:pr:apx.cmd.pkgmanagers.new.options.show"`
	Update     string `flag:"short:u, long:update, name:pr:apx.cmd.pkgmanagers.new.options.update"`
	Upgrade    string `flag:"short:U, long:upgrade, name:pr:apx.cmd.pkgmanagers.new.options.upgrade"`
}

type PkgManagersRmCmd struct {
	cli.Base
	Name  string `flag:"short:n, long:name, name:pr:apx.cmd.pkgmanagers.rm.options.name"`
	Force bool   `flag:"short:f, long:force, name:pr:apx.cmd.pkgmanagers.rm.options.force"`
}

type PkgManagersExportCmd struct {
	cli.Base
	Name   string `flag:"short:n, long:name, name:pr:apx.cmd.pkgmanagers.export.options.name"`
	Output string `flag:"short:o, long:output, name:pr:apx.cmd.pkgmanagers.export.options.output"`
}

type PkgManagersImportCmd struct {
	cli.Base
	Input string `flag:"short:i, long:input, name:pr:apx.cmd.pkgmanagers.import.options.input"`
}
