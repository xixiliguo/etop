package tui

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

type System struct {
	*tview.Box
	layout           *tview.Flex
	header           *tview.TextView
	regions          []string
	currentRegionIdx int
	regionToPage     map[string]string
	content          *tview.Pages
	cpu              *tview.Table
	mem              *tview.Table
	vm               *tview.Table
	disk             *tview.Table
	net              *tview.Table
	source           *model.Model
}

func NewSystem() *System {

	system := &System{
		Box:     tview.NewBox(),
		layout:  tview.NewFlex(),
		header:  tview.NewTextView(),
		content: tview.NewPages(),
		cpu:     tview.NewTable().SetFixed(1, 1),
		mem:     tview.NewTable().SetFixed(1, 1),
		vm:      tview.NewTable().SetFixed(1, 1),
		disk:    tview.NewTable().SetFixed(1, 1),
		net:     tview.NewTable().SetFixed(1, 1),
	}

	system.SetTitle("System").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	system.content.
		AddPage("CPU", system.cpu, true, true).
		AddPage("MEM", system.mem, true, false).
		AddPage("VM", system.vm, true, false).
		AddPage("DISK", system.disk, true, false).
		AddPage("NET", system.net, true, false)

	system.layout.SetDirection(tview.FlexRow).
		AddItem(system.header, 1, 0, false).
		AddItem(system.content, 0, 1, true)

	system.regions = []string{"c", "m", "v", "d", "n"}
	system.regionToPage = map[string]string{
		"c": "CPU",
		"m": "MEM",
		"v": "VM",
		"d": "DISK",
		"n": "NET",
	}
	fmt.Fprintf(system.header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]`,
		"c", "CPU",
		"m", "MEM",
		"v", "VM",
		"d", "DISK",
		"n", "NET")
	system.header.SetRegions(true).Highlight("c")

	return system
}

func (system *System) SetSource(source *model.Model) {
	system.source = source
	system.DrawCPUInfo()
	system.DrawMEMInfo()
	system.DrawVMInfo()
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

func (system *System) DrawVMInfo() {
	system.vm.Clear()
	system.vm.SetOffset(0, 0)

	items := []string{"PageIn", "PageOut",
		"SwapIn", "SwapOut",
		"PageScanKswapd", "PageScanDirect",
		"PageStealKswapd", "PageStealDirect", "OOMKill"}
	for i, v := range []string{"Field", "Value"} {
		system.vm.SetCell(0, i, tview.NewTableCell(v).SetTextColor(tcell.ColorBlue))
	}

	for i, item := range items {
		system.vm.SetCell(i+1,
			0,
			tview.NewTableCell(item).
				SetExpansion(0).
				SetAlign(tview.AlignLeft))
		system.vm.SetCell(i+1,
			1,
			tview.NewTableCell(system.source.Vm.GetRenderValue(system.source.Config["vm"], item)).
				SetExpansion(0).
				SetAlign(tview.AlignRight))
	}

}

func (system *System) DrawDiskInfo() {
	system.disk.Clear()
	system.disk.SetOffset(0, 0)

	visbleCols := []string{
		"Disk", "Util",
		"Read/s", "ReadByte/s",
		"Write/s", "WriteByte/s",
		"AvgQueueLength", "AvgWait", "AvgIOTime",
	}

	for i, col := range visbleCols {

		system.disk.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}

	names := make([]string, 0, len(system.source.Disks))
	for n := range system.source.Disks {
		names = append(names, n)
	}
	sort.Strings(names)

	r := 0
	for _, n := range names {
		disk := system.source.Disks[n]
		for i, col := range visbleCols {
			system.disk.SetCell(r+1,
				i,
				tview.NewTableCell(disk.GetRenderValue(system.source.Config["disk"], col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
		r++
	}

}

func (system *System) DrawNetInfo() {
	system.net.Clear()
	system.net.SetOffset(0, 0)

	visbleCols := []string{
		"Name",
		"RxByte/s", "RxPacket/s",
		"TxByte/s", "TxPacket/s",
		"RxErrors", "RxDropped", "RxFIFO", "RxFrame",
		"TxErrors", "TxDropped", "TxFIFO", "TxCollisions",
	}

	for i, col := range visbleCols {

		system.net.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorBlue))
	}

	names := make([]string, 0, len(system.source.Nets))
	for n := range system.source.Nets {
		names = append(names, n)
	}
	sort.Strings(names)

	r := 0
	for _, n := range names {
		net := system.source.Nets[n]
		for i, col := range visbleCols {
			system.net.SetCell(r+1,
				i,
				tview.NewTableCell(net.GetRenderValue(system.source.Config["netdev"], col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
		r++
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

func (system *System) setRegionAndSwitchPage(region string) {
	for i, r := range system.regions {
		if r == region {
			system.currentRegionIdx = i
		}
	}
	system.header.Highlight(region)
	system.content.SwitchToPage(system.regionToPage[region])
	return
}

func (system *System) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return system.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

		if k := event.Rune(); k == 'c' || k == 'm' || k == 'v' || k == 'd' || k == 'n' {
			s := string(k)
			system.setRegionAndSwitchPage(s)
			return
		}

		if event.Key() == tcell.KeyTab {
			nextId := (system.currentRegionIdx + 1) % len(system.regions)
			s := system.regions[nextId]
			system.setRegionAndSwitchPage(s)
			return
		}
		if event.Key() == tcell.KeyBacktab {
			nextId := (system.currentRegionIdx - 1 + len(system.regions)) % len(system.regions)
			s := system.regions[nextId]
			system.setRegionAndSwitchPage(s)
			return
		}
		if system.content.HasFocus() {
			if Handler := system.content.InputHandler(); Handler != nil {
				Handler(event, setFocus)
				return
			}
		}
	})
}
