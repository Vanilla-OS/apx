package cmdr

import "github.com/pterm/pterm"

type BulletListItem = pterm.BulletListItem
type Style = pterm.Style
type Color = pterm.Color

// Foreground colors. basic foreground colors 30 - 37.
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	// FgDefault revert default FG.
	FgDefault Color = 39
)

// Extra foreground color 90 - 97.
const (
	FgDarkGray Color = iota + 90
	FgLightRed
	FgLightGreen
	FgLightYellow
	FgLightBlue
	FgLightMagenta
	FgLightCyan
	FgLightWhite
	// FgGray is an alias of FgDarkGray.
	FgGray Color = 90
)

// Background colors. basic background colors 40 - 47.
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow // BgBrown like yellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	// BgDefault reverts to the default background.
	BgDefault Color = 49
)

// Extra background color 100 - 107.
const (
	BgDarkGray Color = iota + 100
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgLightWhite
	// BgGray is an alias of BgDarkGray.
	BgGray Color = 100
)

// Option settings.
const (
	Reset Color = iota
	Bold
	Fuzzy
	Italic
	Underscore
	Blink
	FastBlink
	Reverse
	Concealed
	Strikethrough
)

var (
	// Red is an alias for FgRed.Sprint.
	Red = FgRed.Sprint
	// Cyan is an alias for FgCyan.Sprint.
	Cyan = FgCyan.Sprint
	// Gray is an alias for FgGray.Sprint.
	Gray = FgGray.Sprint
	// Blue is an alias for FgBlue.Sprint.
	Blue = FgBlue.Sprint
	// Black is an alias for FgBlack.Sprint.
	Black = FgBlack.Sprint
	// Green is an alias for FgGreen.Sprint.
	Green = FgGreen.Sprint
	// White is an alias for FgWhite.Sprint.
	White = FgWhite.Sprint
	// Yellow is an alias for FgYellow.Sprint.
	Yellow = FgYellow.Sprint
	// Magenta is an alias for FgMagenta.Sprint.
	Magenta = FgMagenta.Sprint

	// Normal is an alias for FgDefault.Sprint.
	Normal = FgDefault.Sprint

	// extra light.

	// LightRed is a shortcut for FgLightRed.Sprint.
	LightRed = FgLightRed.Sprint
	// LightCyan is a shortcut for FgLightCyan.Sprint.
	LightCyan = FgLightCyan.Sprint
	// LightBlue is a shortcut for FgLightBlue.Sprint.
	LightBlue = FgLightBlue.Sprint
	// LightGreen is a shortcut for FgLightGreen.Sprint.
	LightGreen = FgLightGreen.Sprint
	// LightWhite is a shortcut for FgLightWhite.Sprint.
	LightWhite = FgLightWhite.Sprint
	// LightYellow is a shortcut for FgLightYellow.Sprint.
	LightYellow = FgLightYellow.Sprint
	// LightMagenta is a shortcut for FgLightMagenta.Sprint.
	LightMagenta = FgLightMagenta.Sprint
)

var (
	Info, Warning, Success, Fatal, Debug, Description, Error pterm.PrefixPrinter
	Spinner                                                  pterm.SpinnerPrinter
	ProgressBar                                              pterm.ProgressbarPrinter
	BulletList                                               pterm.BulletListPrinter
	TerminalSize                                             func() (int, int, error)
	TerminalWidth, TerminalHeight                            func() int
	NewStyle                                                 func(...Color) *Style
	EnableColor, DisableColor                                func()
)

func init() {
	Error = pterm.Error
	Info = pterm.Info
	Warning = pterm.Warning
	Success = pterm.Success
	Fatal = pterm.Fatal
	Debug = pterm.Debug
	Description = pterm.Description
	Spinner = pterm.DefaultSpinner
	ProgressBar = pterm.DefaultProgressbar
	TerminalSize = pterm.GetTerminalSize
	TerminalWidth = pterm.GetTerminalWidth
	TerminalHeight = pterm.GetTerminalHeight
	BulletList = pterm.DefaultBulletList
	NewStyle = pterm.NewStyle
	EnableColor = pterm.EnableColor
	DisableColor = pterm.DisableColor
}
