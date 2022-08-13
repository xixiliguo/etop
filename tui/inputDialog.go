package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputDialog struct {
	*tview.Box
	layout *tview.Flex
	form   *tview.InputField
}

func NewInputDialog() *InputDialog {

	input := &InputDialog{
		Box:    tview.NewBox(),
		layout: tview.NewFlex(),
		form:   tview.NewInputField(),
	}

	input.form.
		SetTitle("Search sample (e.g [yyyy-mm-dd ]hh:mm)      ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)
	input.form.
		SetLabel(">  ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorBlue)

	input.layout.
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(input.form, 3, 1, true).
			AddItem(nil, 0, 1, false), 70, 1, true).
		AddItem(nil, 0, 1, false)

	return input
}

func (input *InputDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return input.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if handler := input.layout.InputHandler(); handler != nil {
			handler(event, setFocus)
			return
		}
	})
}

func (input *InputDialog) HasFocus() bool {
	return input.layout.HasFocus()
}

func (input *InputDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(input.layout)
}

func (input *InputDialog) Draw(screen tcell.Screen) {

	x, y, width, height := input.Box.GetInnerRect()

	input.layout.SetRect(x, y, width, height)
	input.layout.Draw(screen)
}
