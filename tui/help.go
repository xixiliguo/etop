package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

var content = `
	t               - show next sample
	<Shift>+t       - show preview sample
	<F1>, <Alt>+1   - switch to process view
	<F2>, <Alt>+2   - switch to system view
	'b'             - open dialog to search specific sample

process view:
	's'             - show/hide sort view
	'/'             - show/hide filter view
	'g'             - show process-level general info
	'c'             - show process-level cpu info
	'm'             - show process-level memory info
	'd'             - show process-level disk info

system view:
	'c'             - show system-level cpu info
	'm'             - show system-level memory info
	'v'             - show system-level vm info
	'd'             - show system-level disk info
	'n'             - show system-level network info

	Type 'ESC' to close
`

type Help struct {
	*tview.TextView
}

func NewHelp() *Help {

	help := &Help{
		TextView: tview.NewTextView(),
	}

	help.SetTitle("Help").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	fmt.Fprint(help, content)

	return help
}
