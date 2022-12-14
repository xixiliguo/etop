package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

type System struct {
	*tview.Box
	layout  *tview.Flex
	header  *tview.TextView
	content *tview.Pages
	cpu     *tview.Table
	mem     *tview.Table
	disk    *tview.Table
	net     *tview.Table
	source  *model.System
}

func NewSystem() *System {

	system := &System{
		Box:     tview.NewBox(),
		layout:  tview.NewFlex(),
		header:  tview.NewTextView(),
		content: tview.NewPages(),
		cpu:     tview.NewTable().SetFixed(1, 1),
		mem:     tview.NewTable().SetFixed(1, 1),
		disk:    tview.NewTable().SetFixed(1, 1),
		net:     tview.NewTable().SetFixed(1, 1),
	}

	system.SetTitle("System").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	system.content.
		AddPage("CPU", system.cpu, true, true).
		AddPage("MEM", system.mem, true, false).
		AddPage("DISK", system.disk, true, false).
		AddPage("NET", system.net, true, false)

	system.layout.SetDirection(tview.FlexRow).
		AddItem(system.header, 1, 0, false).
		AddItem(system.content, 0, 1, true)

	fmt.Fprintf(system.header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]`,
		"c", "CPU",
		"m", "MEM",
		"d", "DISK",
		"n", "NET")
	system.header.SetRegions(true).Highlight("c")

	system.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'c' {
			system.header.Highlight("c")
			system.content.SwitchToPage("CPU")
			return nil
		}

		if event.Rune() == 'm' {
			system.header.Highlight("m")
			system.content.SwitchToPage("MEM")
			return nil
		}

		if event.Rune() == 'd' {
			system.header.Highlight("d")
			system.content.SwitchToPage("DISK")
			return nil
		}

		if event.Rune() == 'n' {
			system.header.Highlight("n")
			system.content.SwitchToPage("NET")
			return nil
		}

		return event
	})

	return system
}

func (system *System) SetSource(source *model.System) {
	system.source = source
	system.DrawCPUInfo()
	system.DrawMEMInfo()
	system.DrawDiskInfo()
	system.DrawNetInfo()
}

func (system *System) DrawCPUInfo() {
	system.cpu.Clear()
	system.cpu.SetOffset(0, 0)

	visbleCols := []string{"Index", "User", "Nice",
		"System", "Idle", "Iowait", "IRQ",
		"SoftIRQ", "Steal", "Guest", "GuestNice"}

	for i, col := range visbleCols {
		if col == "Index" {
			col = ""
		}
		system.cpu.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(system.source.CPUs); r++ {
		c := system.source.CPUs[r]
		for i, col := range visbleCols {
			system.cpu.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(system.source.Config["cpu"], col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}
}

func (system *System) DrawMEMInfo() {
	system.mem.Clear()
	system.mem.SetOffset(0, 0)

	items := []string{"MemTotal", "MemFree", "MemAvailable",
		"Buffers", "Cached", "SwapCached", "Active",
		"Inactive", "ActiveAnon", "InactiveAnon", "Unevictable",
		"Mlocked", "SwapTotal", "SwapFree", "Dirty",
		"Writeback", "AnonPages", "Mapped", "Shmem",
		"Slab", "SReclaimable", "SUnreclaim", "KernelStack",
		"PageTables", "NFSUnstable", "Bounce", "WritebackTmp",
		"CommitLimit", "CommittedAS", "VmallocTotal", "VmallocUsed",
		"VmallocChunk", "HardwareCorrupted", "AnonHugePages", "ShmemHugePages",
		"ShmemPmdMapped", "CmaTotal", "CmaFree", "HugePagesTotal",
		"HugePagesFree", "HugePagesRsvd", "HugePagesSurp", "Hugepagesize",
		"DirectMap4k", "DirectMap2M", "DirectMap1G",
	}
	for i, v := range []string{"Field", "Value"} {
		system.mem.SetCell(0, i, tview.NewTableCell(v).SetTextColor(tcell.ColorBlue))
	}

	for i, item := range items {
		system.mem.SetCell(i+1,
			0,
			tview.NewTableCell(item).
				SetExpansion(0).
				SetAlign(tview.AlignLeft))
		system.mem.SetCell(i+1,
			1,
			tview.NewTableCell(system.source.MEM.GetRenderValue(system.source.Config["memory"], item)).
				SetExpansion(0).
				SetAlign(tview.AlignRight))
	}

}

func (system *System) DrawDiskInfo() {
	system.disk.Clear()
	system.disk.SetOffset(0, 0)

	visbleCols := []string{
		"Disk", "Busy",
		"Read", "R/s",
		"Write", "W/s",
		"Await", "Avio",
	}

	for i, col := range visbleCols {

		system.disk.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(system.source.Disks); r++ {
		c := system.source.Disks[r]
		for i, col := range visbleCols {
			system.disk.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}

}

func (system *System) DrawNetInfo() {
	system.net.Clear()
	system.net.SetOffset(0, 0)

	visbleCols := []string{
		"Name",
		"R/s", "Rp/s",
		"T/s", "Tp/s",
		"RxBytes", "RxPackets", "RxErrors", "RxDropped", "RxFIFO", "RxFrame", "RxCompressed", "RxMulticast",
		"TxBytes", "TxPackets", "TxErrors", "TxDropped", "TxFIFO", "TxCollisions", "TxCarrier", "TxCompressed",
	}

	for i, col := range visbleCols {

		system.net.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(system.source.Nets); r++ {
		c := system.source.Nets[r]
		for i, col := range visbleCols {
			system.net.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}

}

func (system *System) HasFocus() bool {
	return system.layout.HasFocus()
}

func (system *System) Focus(delegate func(p tview.Primitive)) {
	delegate(system.layout)
}

func (system *System) Draw(screen tcell.Screen) {
	system.Box.DrawForSubclass(screen, system)
	x, y, width, height := system.Box.GetInnerRect()

	system.layout.SetRect(x, y, width, height)
	system.layout.Draw(screen)
}
