package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputDialog struct {
	*tview.Flex
	form *tview.InputField
}

func NewInputDialog() *InputDialog {

	input := &InputDialog{
		Flex: tview.NewFlex(),
		form: tview.NewInputField(),
	}

	input.form.
		SetTitle("Search sample (e.g [yyyy-mm-dd ]hh:mm)      ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)
	input.form.
		SetLabel(">  ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorTeal)

	input.AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(input.form, 3, 1, true).
			AddItem(nil, 0, 1, false), 70, 1, true).
		AddItem(nil, 0, 1, false)

	return input
}
