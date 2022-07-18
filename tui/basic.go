package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/util"
)

type BasicView struct {
	*tview.Flex
	PRC  *tview.TextView
	CPL  *tview.TextView
	CPU  *tview.TextView
	MEM  *tview.TextView
	DISK *tview.TextView
	NET  *tview.TextView
}

func NewBasicView() *BasicView {

	bv := &BasicView{
		Flex: tview.NewFlex(),
		PRC:  tview.NewTextView(),
		CPL:  tview.NewTextView(),
		CPU:  tview.NewTextView(),
		MEM:  tview.NewTextView(),
		DISK: tview.NewTextView(),
		NET:  tview.NewTextView(),
	}
	bv.SetBorder(true)
	bv.SetDirection(tview.FlexRow).
		AddItem(bv.PRC, 1, 0, false).
		AddItem(bv.CPL, 1, 0, false).
		AddItem(bv.CPU, 1, 0, false).
		AddItem(bv.MEM, 1, 0, false).
		AddItem(bv.DISK, 1, 0, false).
		AddItem(bv.NET, 1, 0, false)

	return bv
}

func (bv *BasicView) Update(sm *model.System) {

	bv.PRC.Clear()
	fmt.Fprintf(bv.PRC, "%-10sProcess %4d%5sThread %5d%5sClone %4d/s",
		"PRC", sm.Prcesses, "", sm.Threads, "", sm.Clones)

	bv.CPL.Clear()
	fmt.Fprintf(bv.CPL, "%-10sAvg1 %7.2f%5sAvg5 %7.2f%5sAvg15 %6.2f%5sTrun %7d%5sTslpu %6d%5sCtxsw %4d/s",
		"CPL", sm.Curr.Load1, "", sm.Curr.Load5, "", sm.Curr.Load15, "", sm.Curr.ProcessesRunning, "",
		sm.Curr.ProcessesBlocked, "",
		sm.ContextSwitch)

	c := model.CPU{}
	for i := 0; i < len(sm.CPUs); i++ {
		if sm.CPUs[i].Index == "total" {
			c = sm.CPUs[i]
		}
	}
	bv.CPU.Clear()
	fmt.Fprintf(bv.CPU, "%-10sUser %6.1f%%%5sSystem %4.1f%%%5sIowait %4.1f%%%5sIdle %6.1f%%%5sIRQ %7.1f%%%5sSoftIRQ %3.1f%%",
		"CPU", c.User, "", c.System, "", c.Iowait, "", c.Idle, "", c.IRQ, "", c.SoftIRQ)
	bv.MEM.Clear()
	fmt.Fprintf(bv.MEM, "%-10sTotal %6s%5sFree %7s%5sAvail %6s%5sSlab %7s%5sBuffer %5s%5sCache %6s",
		"MEM", sm.MEM.GetRenderValue("Total"), "",
		sm.MEM.GetRenderValue("Free"), "",
		sm.MEM.GetRenderValue("Avail"), "",
		util.GetHumanSize(sm.MEM.Slab*1024), "",
		sm.MEM.GetRenderValue("Buffer"), "",
		sm.MEM.GetRenderValue("Cache"))

	bv.DISK.Clear()
	fmt.Fprintf(bv.DISK, "%-10s", "DSK")
	for _, d := range sm.Disks {
		fmt.Fprintf(bv.DISK, "%-7s%7s|%-7s ", d.GetRenderValue("Disk"), d.GetRenderValue("R/s"), d.GetRenderValue("W/s"))
	}

	bv.NET.Clear()
	fmt.Fprintf(bv.NET, "%-10s", "NET")
	for _, d := range sm.Nets {
		fmt.Fprintf(bv.NET, "%-7s%7s|%-7s ", d.GetRenderValue("Name"), d.GetRenderValue("R/s"), d.GetRenderValue("T/s"))
	}

}
