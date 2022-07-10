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

	inputDialog := &InputDialog{
		Flex: tview.NewFlex(),
		form: tview.NewInputField(),
	}

	inputDialog.form.
		SetTitle("Search sample by specific time (e.g hh:mm)      ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)
	inputDialog.form.
		SetLabel(">  ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorBlue)

	inputDialog.Flex.
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(inputDialog.form, 3, 1, true).
			AddItem(nil, 0, 1, false), 70, 1, true).
		AddItem(nil, 0, 1, false)

	return inputDialog
}
