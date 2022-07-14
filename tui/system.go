package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

type SystemView struct {
	Tui *tview.Application
	*tview.Flex
	Header    *tview.TextView
	Content   *tview.Pages
	CPUTable  *tview.Table
	MEMTable  *tview.Table
	DiskTable *tview.Table
	NetTable  *tview.Table
	Source    *model.System
}

func NewSystemView(tui *tview.Application) *SystemView {

	sysv := &SystemView{
		Tui:       tui,
		Flex:      tview.NewFlex(),
		Header:    tview.NewTextView(),
		Content:   tview.NewPages(),
		CPUTable:  tview.NewTable().SetFixed(1, 1),
		MEMTable:  tview.NewTable().SetFixed(1, 1),
		DiskTable: tview.NewTable().SetFixed(1, 1),
		NetTable:  tview.NewTable().SetFixed(1, 1),
	}

	sysv.SetTitle("System").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	sysv.Content.
		AddPage("CPU", sysv.CPUTable, true, true).
		AddPage("MEM", sysv.MEMTable, true, false).
		AddPage("DISK", sysv.DiskTable, true, false).
		AddPage("NET", sysv.NetTable, true, false)

	sysv.SetDirection(tview.FlexRow).
		AddItem(sysv.Header, 1, 0, false).
		AddItem(sysv.Content, 0, 1, true)

	fmt.Fprintf(sysv.Header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]`,
		"c", "CPU",
		"m", "MEM",
		"d", "DISK",
		"n", "NET")
	sysv.Header.SetRegions(true).Highlight("c")

	sysv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Rune() == 'c' {
			sysv.Header.Highlight("c")
			sysv.Content.SwitchToPage("CPU")
			return nil
		}

		if event.Rune() == 'm' {
			sysv.Header.Highlight("m")
			sysv.Content.SwitchToPage("MEM")
			return nil
		}

		if event.Rune() == 'd' {
			sysv.Header.Highlight("d")
			sysv.Content.SwitchToPage("DISK")
			return nil
		}

		if event.Rune() == 'n' {
			sysv.Header.Highlight("n")
			sysv.Content.SwitchToPage("NET")
			return nil
		}

		return event
	})

	return sysv
}

func (sysv *SystemView) SetSource(source *model.System) {
	sysv.Source = source
	sysv.DrawCPUInfo()
	sysv.DrawMEMInfo()
	sysv.DrawDiskInfo()
	sysv.DrawNetInfo()
}

func (sysv *SystemView) DrawCPUInfo() {
	sysv.CPUTable.Clear()
	sysv.CPUTable.SetOffset(0, 0)

	visbleCols := []string{"Index", "User", "Nice",
		"System", "Idle", "Iowait", "IRQ",
		"SoftIRQ", "Steal", "Guest", "GuestNice"}

	for i, col := range visbleCols {
		if col == "Index" {
			col = ""
		}
		sysv.CPUTable.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(sysv.Source.CPUs); r++ {
		c := sysv.Source.CPUs[r]
		for i, col := range visbleCols {
			sysv.CPUTable.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}
}

func (sysv *SystemView) DrawMEMInfo() {
	sysv.MEMTable.Clear()
	sysv.MEMTable.SetOffset(0, 0)

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
		sysv.MEMTable.SetCell(0, i, tview.NewTableCell(v).SetTextColor(tcell.ColorBlue))
	}

	for i, item := range items {
		sysv.MEMTable.SetCell(i+1,
			0,
			tview.NewTableCell(item).
				SetExpansion(0).
				SetAlign(tview.AlignLeft))
		sysv.MEMTable.SetCell(i+1,
			1,
			tview.NewTableCell(sysv.Source.MEM.GetRenderValue(item)).
				SetExpansion(0).
				SetAlign(tview.AlignRight))
	}

}

func (sysv *SystemView) DrawDiskInfo() {
	sysv.DiskTable.Clear()
	sysv.DiskTable.SetOffset(0, 0)

	visbleCols := []string{
		"Disk", "Busy",
		"Read", "R/s",
		"Write", "W/s",
		"Await", "Avio",
	}

	for i, col := range visbleCols {

		sysv.DiskTable.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(sysv.Source.Disks); r++ {
		c := sysv.Source.Disks[r]
		for i, col := range visbleCols {
			sysv.DiskTable.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}

}

func (sysv *SystemView) DrawNetInfo() {
	sysv.NetTable.Clear()
	sysv.NetTable.SetOffset(0, 0)

	visbleCols := []string{
		"Name",
		"R/s", "Rp/s",
		"T/s", "Tp/s",
		"RxBytes", "RxPackets", "RxErrors", "RxDropped", "RxFIFO", "RxFrame", "RxCompressed", "RxMulticast",
		"TxBytes", "TxPackets", "TxErrors", "TxDropped", "TxFIFO", "TxCollisions", "TxCarrier", "TxCompressed",
	}

	for i, col := range visbleCols {

		sysv.NetTable.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(sysv.Source.Nets); r++ {
		c := sysv.Source.Nets[r]
		for i, col := range visbleCols {
			sysv.NetTable.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}

}
