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

func (header *Header) Update(sm *model.Model) {
	header.text.Clear()
	fmt.Fprintf(header.text, "%s    Elapsed: %ds    %s    Uptime: %-12s    Mode: %s %s\n",
		time.Unix(sm.Curr.TimeStamp, 0),
		sm.Curr.TimeStamp-sm.Prev.TimeStamp,
		sm.Curr.HostName,
		time.Duration(sm.Curr.TimeStamp-int64(sm.Curr.BootTime))*time.Second,
		sm.Mode,
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
