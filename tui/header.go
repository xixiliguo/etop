package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/version"
)

type Header struct {
	*tview.TextView
}

func NewHeader() *Header {

	header := &Header{
		TextView: tview.NewTextView(),
	}
	header.SetBorder(true)
	return header
}

func (header *Header) Update(sm *model.Model) {
	header.Clear()
	fmt.Fprintf(header, "%s    Elapsed: %ds    %s    Uptime: %-12s    Mode: %s %s\n",
		time.Unix(sm.Curr.TimeStamp, 0),
		sm.Curr.TimeStamp-sm.Prev.TimeStamp,
		sm.Curr.HostName,
		time.Duration(sm.Curr.TimeStamp-int64(sm.Curr.BootTime))*time.Second,
		sm.Mode,
		version.Version)
}
