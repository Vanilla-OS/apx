module github.com/vanilla-os/apx/v2

go 1.24.4

require (
	github.com/google/uuid v1.6.0
	github.com/vanilla-os/sdk v0.0.0-20251230185722-751386443204
	gopkg.in/yaml.v2 v2.4.0
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/containerd/console v1.0.5 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gookit/color v1.6.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mirkobrombin/go-cli-builder/v2 v2.0.4 // indirect
	github.com/mirkobrombin/go-struct-flags/v2 v2.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/phuslu/log v1.0.88 // indirect
	github.com/pterm/pterm v0.12.81 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/vorlif/spreak v0.6.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

require (
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0
)

replace github.com/vanilla-os/sdk => ../sdk

replace github.com/mirkobrombin/go-cli-builder/v2 => ../go-tools/go-cli-builder
