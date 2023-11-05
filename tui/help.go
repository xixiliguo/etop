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
	'b'             - open Dialog to search specific sample

process view:
	's'             - sort by specific field
	'/'             - filter
	'g'             - show general info
	'c'             - show cpu info
	'm'             - show memory info
	'd'             - show disk info

system view:
	'c'             - show cpu info
	'm'             - show memory info
	'v'             - show vm info
	'd'             - show disk info
	'n'             - show networ info

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

	fmt.Fprintf(help, content)

	return help
}
