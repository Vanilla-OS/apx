package cmdr

import "github.com/pterm/pterm"

var (
	Confirm pterm.InteractiveConfirmPrinter
	Prompt  pterm.InteractiveTextInputPrinter
)

func init() {
	Confirm = pterm.DefaultInteractiveConfirm
	Prompt = pterm.DefaultInteractiveTextInput
}
