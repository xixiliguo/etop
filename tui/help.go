package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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

type Help struct {
	*tview.Box
	text *tview.TextView
}

func NewHelp() *Help {

	help := &Help{
		Box:  tview.NewBox(),
		text: tview.NewTextView(),
	}

	help.SetTitle("Help").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	fmt.Fprintf(help.text, content)

	return help
}

func (help *Help) HasFocus() bool {
	return help.text.HasFocus()
}

func (help *Help) Focus(delegate func(p tview.Primitive)) {
	delegate(help.text)
}

func (help *Help) Draw(screen tcell.Screen) {
	help.Box.DrawForSubclass(screen, help)
	x, y, width, height := help.Box.GetInnerRect()

	help.text.SetRect(x, y, width, height)
	help.text.Draw(screen)
}
