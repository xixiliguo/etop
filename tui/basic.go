package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

const (
	cpuFmtStr      = "%-7sUser %10s%5sSystem %8s%5sIowait %8s%5sIdle %10s%5sIRQ %11s%5sSoftIRQ %7s"
	cpuBusyFmtStr  = "%-7sUser %10s%5sSystem %8s%5sIowait %8s%5sIdle [red]%10s[white]%5sIRQ %11s%5sSoftIRQ %7s"
	memFmtStr      = "%-7sTotal %9s%5sFree %10s%5sAvail %9s%5sSlab %10s%5sBuffer %8s%5sCache %9s"
	memBusyFmtStr  = "%-7sTotal %9s%5sFree %10s%5sAvail [red]%9s[white]%5sSlab %10s%5sBuffer %8s%5sCache %9s"
	diskFmtStr     = "%-5s%10s|%-10s "
	diskBusyFmtStr = "%-5s[red]%10s|%-10s[white] "
)

type Basic struct {
	*tview.Flex
	load *tview.TextView
	proc *tview.TextView
	cpu  *tview.TextView
	mem  *tview.TextView
	disk *tview.TextView
	net  *tview.TextView
}

func NewBasic() *Basic {

	basic := &Basic{
		Flex: tview.NewFlex(),
		load: tview.NewTextView(),
		proc: tview.NewTextView(),
		cpu:  tview.NewTextView().SetDynamicColors(true),
		mem:  tview.NewTextView().SetDynamicColors(true),
		disk: tview.NewTextView().SetDynamicColors(true),
		net:  tview.NewTextView(),
	}
	basic.SetBorder(true)
	basic.SetDirection(tview.FlexRow).
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
		"Load", sm.Sys.GetRenderValue("Load1", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("Load5", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("Load15", model.FieldOpt{}))

	basic.proc.Clear()
	fmt.Fprintf(basic.proc, "%-7sProcess %7s%5sThread %8s%5sRunning %7s%5sBlocked %7s%5sClone %9s%5sCtxSw %9s",
		"Proc", sm.Sys.GetRenderValue("Processes", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("Threads", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("ProcessesRunning", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("ProcessesBlocked", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("ClonePerSec", model.FieldOpt{}), "",
		sm.Sys.GetRenderValue("ContextSwitchPerSec", model.FieldOpt{}))
	var fmtStr string
	c := model.CPU{}
	for i := 0; i < len(sm.CPUs); i++ {
		if sm.CPUs[i].Index == "total" {
			c = sm.CPUs[i]
		}
	}
	basic.cpu.Clear()
	fmtStr = cpuFmtStr
	if c.Idle <= (100 - CPUBusy) {
		fmtStr = cpuBusyFmtStr
	}
	fmt.Fprintf(basic.cpu, fmtStr,
		"CPU", c.GetRenderValue("User", model.FieldOpt{}), "",
		c.GetRenderValue("System", model.FieldOpt{}), "",
		c.GetRenderValue("Iowait", model.FieldOpt{}), "",
		c.GetRenderValue("Idle", model.FieldOpt{}), "",
		c.GetRenderValue("IRQ", model.FieldOpt{}), "",
		c.GetRenderValue("SoftIRQ", model.FieldOpt{}))

	basic.mem.Clear()
	fmtStr = memFmtStr
	if float64(sm.MEM.MemAvailable*100/sm.MEM.MemTotal) <= (100 - MemBusy) {
		fmtStr = memBusyFmtStr
	}
	fmt.Fprintf(basic.mem, fmtStr,
		"Mem", sm.MEM.GetRenderValue("Total", model.FieldOpt{}), "",
		sm.MEM.GetRenderValue("Free", model.FieldOpt{}), "",
		sm.MEM.GetRenderValue("Avail", model.FieldOpt{}), "",
		sm.MEM.GetRenderValue("HSlab", model.FieldOpt{}), "",
		sm.MEM.GetRenderValue("Buffer", model.FieldOpt{}), "",
		sm.MEM.GetRenderValue("Cache", model.FieldOpt{}))

	basic.disk.Clear()
	fmt.Fprintf(basic.disk, "%-7s", "Disk")
	for _, disk := range sm.Disks.GetKeys() {
		d := sm.Disks[disk]
		fmtStr = diskFmtStr
		if d.Util >= DiskBusy {
			fmtStr = diskBusyFmtStr
		}
		fmt.Fprintf(basic.disk, fmtStr,
			d.GetRenderValue("Disk", model.FieldOpt{}),
			d.GetRenderValue("ReadByte/s", model.FieldOpt{}),
			d.GetRenderValue("WriteByte/s", model.FieldOpt{}))
	}

	basic.net.Clear()
	fmt.Fprintf(basic.net, "%-7s", "Net")
	for _, netDev := range sm.Nets.GetKeys() {
		d := sm.Nets[netDev]
		fmt.Fprintf(basic.net, "%-5s%10s|%-10s ",
			d.GetRenderValue("Name", model.FieldOpt{}),
			d.GetRenderValue("RxByte/s", model.FieldOpt{}),
			d.GetRenderValue("TxByte/s", model.FieldOpt{}))
	}

}
