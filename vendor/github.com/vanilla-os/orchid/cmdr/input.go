package cmdr

import "github.com/pterm/pterm"

var (
	Confirm pterm.InteractiveConfirmPrinter
)

func init() {
	Confirm = pterm.DefaultInteractiveConfirm
}
