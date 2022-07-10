package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

var content = `
    <F1>           - switch to process view
    <F2>           - switch to system view
    'g'            - show gernal info (process view only)
    'c'            - show cpu-relate info
    'm'            - show memory-relate info
    'd'            - show disk-relate info
    'n'            - show network-relate info (process view only)

	Type 'ESC' to close
`

type HelpView struct {
	*tview.TextView
}

func NewHelpView() *HelpView {

	help := &HelpView{
		TextView: tview.NewTextView(),
	}

	help.SetTitle("Help").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	fmt.Fprintf(help, content)

	return help
}
