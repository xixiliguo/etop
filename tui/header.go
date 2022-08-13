package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/version"
)

type Header struct {
	*tview.Box
	text *tview.TextView
}

func NewHeader() *Header {

	header := &Header{
		Box:  tview.NewBox(),
		text: tview.NewTextView(),
	}
	header.SetBorder(true)
	return header
}

func (header *Header) Update(sm *model.System) {
	header.text.Clear()
	fmt.Fprintf(header.text, "%s    Elapsed: %4ds    %s    Uptime: %s    %s\n",
		time.Unix(sm.Curr.CurrTime, 0),
		sm.Curr.CurrTime-sm.Prev.CurrTime,
		sm.Curr.HostName,
		time.Duration(sm.Curr.CurrTime-int64(sm.Curr.BootTime))*time.Second,
		version.Version)
}

func (header *Header) HasFocus() bool {
	return header.text.HasFocus()
}

func (header *Header) Focus(delegate func(p tview.Primitive)) {
	delegate(header.text)
}

func (header *Header) Draw(screen tcell.Screen) {
	header.Box.DrawForSubclass(screen, header)
	x, y, width, height := header.Box.GetInnerRect()

	header.text.SetRect(x, y, width, height)
	header.text.Draw(screen)
}
