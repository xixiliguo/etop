package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

type Basic struct {
	*tview.Box
	layout *tview.Flex
	prc    *tview.TextView
	cpl    *tview.TextView
	cpu    *tview.TextView
	mem    *tview.TextView
	disk   *tview.TextView
	net    *tview.TextView
}

func NewBasic() *Basic {

	basic := &Basic{
		Box:    tview.NewBox(),
		layout: tview.NewFlex(),
		prc:    tview.NewTextView(),
		cpl:    tview.NewTextView(),
		cpu:    tview.NewTextView(),
		mem:    tview.NewTextView(),
		disk:   tview.NewTextView(),
		net:    tview.NewTextView(),
	}
	basic.SetBorder(true)
	basic.layout.SetDirection(tview.FlexRow).
		AddItem(basic.prc, 1, 0, false).
		AddItem(basic.cpl, 1, 0, false).
		AddItem(basic.cpu, 1, 0, false).
		AddItem(basic.mem, 1, 0, false).
		AddItem(basic.disk, 1, 0, false).
		AddItem(basic.net, 1, 0, false)

	return basic
}

func (basic *Basic) Update(sm *model.System) {

	basic.prc.Clear()
	fmt.Fprintf(basic.prc, "%-10sProcess %4d%5sThread %5d%5sClone %4d/s",
		"PRC", sm.Prcesses, "", sm.Threads, "", sm.Clones)

	basic.cpl.Clear()
	fmt.Fprintf(basic.cpl, "%-10sAvg1 %7.2f%5sAvg5 %7.2f%5sAvg15 %6.2f%5sTrun %7d%5sTslpu %6d%5sCtxsw %4d/s",
		"CPL", sm.Curr.Load1, "", sm.Curr.Load5, "", sm.Curr.Load15, "", sm.Curr.ProcessesRunning, "",
		sm.Curr.ProcessesBlocked, "",
		sm.ContextSwitch)

	c := model.CPU{}
	for i := 0; i < len(sm.CPUs); i++ {
		if sm.CPUs[i].Index == "total" {
			c = sm.CPUs[i]
		}
	}
	basic.cpu.Clear()
	fmt.Fprintf(basic.cpu, "%-10sUser %6.1f%%%5sSystem %4.1f%%%5sIowait %4.1f%%%5sIdle %6.1f%%%5sIRQ %7.1f%%%5sSoftIRQ %3.1f%%",
		"CPU", c.User, "", c.System, "", c.Iowait, "", c.Idle, "", c.IRQ, "", c.SoftIRQ)
	basic.mem.Clear()
	fmt.Fprintf(basic.mem, "%-10sTotal %6s%5sFree %7s%5sAvail %6s%5sSlab %7s%5sBuffer %5s%5sCache %6s",
		"MEM", sm.MEM.GetRenderValue(sm.Config["memory"], "Total"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Free"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Avail"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "HSlab"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Buffer"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Cache"))

	basic.disk.Clear()
	fmt.Fprintf(basic.disk, "%-10s", "DSK")
	for _, d := range sm.Disks {
		fmt.Fprintf(basic.disk, "%-7s%7s|%-7s ", d.GetRenderValue("Disk"), d.GetRenderValue("R/s"), d.GetRenderValue("W/s"))
	}

	basic.net.Clear()
	fmt.Fprintf(basic.net, "%-10s", "NET")
	for _, d := range sm.Nets {
		fmt.Fprintf(basic.net, "%-7s%7s|%-7s ", d.GetRenderValue("Name"), d.GetRenderValue("R/s"), d.GetRenderValue("T/s"))
	}

}

func (basic *Basic) HasFocus() bool {
	return basic.layout.HasFocus()
}

func (basic *Basic) Focus(delegate func(p tview.Primitive)) {
	delegate(basic.layout)
}

func (basic *Basic) Draw(screen tcell.Screen) {
	basic.Box.DrawForSubclass(screen, basic)
	x, y, width, height := basic.Box.GetInnerRect()

	basic.layout.SetRect(x, y, width, height)
	basic.layout.Draw(screen)
}
