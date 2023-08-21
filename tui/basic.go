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
	load   *tview.TextView
	proc   *tview.TextView
	cpu    *tview.TextView
	mem    *tview.TextView
	disk   *tview.TextView
	net    *tview.TextView
}

func NewBasic() *Basic {

	basic := &Basic{
		Box:    tview.NewBox(),
		layout: tview.NewFlex(),
		load:   tview.NewTextView(),
		proc:   tview.NewTextView(),
		cpu:    tview.NewTextView(),
		mem:    tview.NewTextView(),
		disk:   tview.NewTextView(),
		net:    tview.NewTextView(),
	}
	basic.SetBorder(true)
	basic.layout.SetDirection(tview.FlexRow).
		AddItem(basic.load, 1, 0, false).
		AddItem(basic.proc, 1, 0, false).
		AddItem(basic.cpu, 1, 0, false).
		AddItem(basic.mem, 1, 0, false).
		AddItem(basic.disk, 1, 0, false).
		AddItem(basic.net, 1, 0, false)

	return basic
}

func (basic *Basic) Update(sm *model.Model) {

	basic.load.Clear()
	fmt.Fprintf(basic.load, "%-7sLoad1 %9s%5sLoad5 %9s%5sLoad15 %8s",
		"Load", sm.Sys.GetRenderValue(sm.Config["system"], "Load1"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "Load5"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "Load15"))

	basic.proc.Clear()
	fmt.Fprintf(basic.proc, "%-7sProcess %7s%5sThread %8s%5sRunning %7s%5sBlocked %7s%5sClone %9s%5sCtxSw %9s",
		"Proc", sm.Sys.GetRenderValue(sm.Config["system"], "Processes"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "Threads"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "ProcessesRunning"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "ProcessesBlocked"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "ClonePerSec"), "",
		sm.Sys.GetRenderValue(sm.Config["system"], "ContextSwitchPerSec"))

	c := model.CPU{}
	for i := 0; i < len(sm.CPUs); i++ {
		if sm.CPUs[i].Index == "total" {
			c = sm.CPUs[i]
		}
	}
	basic.cpu.Clear()
	fmt.Fprintf(basic.cpu, "%-7sUser %10s%5sSystem %8s%5sIowait %8s%5sIdle %10s%5sIRQ %11s%5sSoftIRQ %7s",
		"CPU", c.GetRenderValue(sm.Config["cpu"], "User"), "",
		c.GetRenderValue(sm.Config["cpu"], "System"), "",
		c.GetRenderValue(sm.Config["cpu"], "Iowait"), "",
		c.GetRenderValue(sm.Config["cpu"], "Idle"), "",
		c.GetRenderValue(sm.Config["cpu"], "IRQ"), "",
		c.GetRenderValue(sm.Config["cpu"], "SoftIRQ"))

	basic.mem.Clear()
	fmt.Fprintf(basic.mem, "%-7sTotal %9s%5sFree %10s%5sAvail %9s%5sSlab %10s%5sBuffer %8s%5sCache %9s",
		"Mem", sm.MEM.GetRenderValue(sm.Config["memory"], "Total"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Free"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Avail"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "HSlab"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Buffer"), "",
		sm.MEM.GetRenderValue(sm.Config["memory"], "Cache"))

	basic.disk.Clear()
	fmt.Fprintf(basic.disk, "%-7s", "Disk")
	for _, disk := range sm.Disks.GetKeys() {
		d := sm.Disks[disk]
		fmt.Fprintf(basic.disk, "%-5s%10s|%-10s ",
			d.GetRenderValue(sm.Config["disk"], "Disk"),
			d.GetRenderValue(sm.Config["disk"], "ReadByte/s"),
			d.GetRenderValue(sm.Config["disk"], "WriteByte/s"))
	}

	basic.net.Clear()
	fmt.Fprintf(basic.net, "%-7s", "Net")
	for _, netDev := range sm.Nets.GetKeys() {
		d := sm.Nets[netDev]
		fmt.Fprintf(basic.net, "%-5s%10s|%-10s ",
			d.GetRenderValue(sm.Config["netdev"], "Name"),
			d.GetRenderValue(sm.Config["netdev"], "RxByte/s"),
			d.GetRenderValue(sm.Config["netdev"], "TxByte/s"))
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
