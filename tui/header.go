package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/version"
)

type HeaderView struct {
	*tview.TextView
}

func NewHeaderView() *HeaderView {

	bv := &HeaderView{
		TextView: tview.NewTextView(),
	}
	bv.SetBorder(true)
	return bv
}

func (hv *HeaderView) Update(sm *model.System) {
	hv.Clear()
	fmt.Fprintf(hv.TextView, "%s    Elapsed: %4ds    %s    Uptime: %s    %s\n",
		time.Unix(sm.Curr.CurrTime, 0),
		sm.Curr.CurrTime-sm.Prev.CurrTime,
		sm.Curr.HostName,
		time.Duration(sm.Curr.CurrTime-int64(sm.Curr.BootTime))*time.Second,
		version.Version)
}
