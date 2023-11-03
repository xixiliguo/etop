package tui

import (
	"fmt"

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
		AddPage("Mem", system.mem, true, false).
		AddPage("Vm", system.vm, true, false).
		AddPage("Disk", system.disk, true, false).
		AddPage("Net", system.net, true, false)

	system.layout.SetDirection(tview.FlexRow).
		AddItem(system.header, 1, 0, false).
		AddItem(system.content, 0, 1, true)

	system.regions = []string{"c", "m", "v", "d", "n"}
	system.regionToPage = map[string]string{
		"c": "CPU",
		"m": "Mem",
		"v": "Vm",
		"d": "Disk",
		"n": "Net",
	}
	fmt.Fprintf(system.header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]`,
		"c", "CPU",
		"m", "Mem",
		"v", "Vm",
		"d", "Disk",
		"n", "Net")
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

	visbleCols := model.DefaultCPUFields

	for i, col := range visbleCols {
		if col == "Index" {
			col = ""
		}
		system.cpu.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorTeal))
	}
	for r := 0; r < len(system.source.CPUs); r++ {
		c := system.source.CPUs[r]
		for i, col := range visbleCols {
			system.cpu.SetCell(r+1,
				i,
				tview.NewTableCell(c.GetRenderValue(col, model.FieldOpt{})).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}
}

func (system *System) DrawMEMInfo() {
	system.mem.Clear()
	system.mem.SetOffset(0, 0)

	items := model.DefaultMEMFields
	for i, v := range []string{"Field", "Value"} {
		system.mem.SetCell(0, i, tview.NewTableCell(v).SetTextColor(tcell.ColorTeal))
	}

	for i, item := range items {
		system.mem.SetCell(i+1,
			0,
			tview.NewTableCell(item).
				SetExpansion(0).
				SetAlign(tview.AlignLeft))
		system.mem.SetCell(i+1,
			1,
			tview.NewTableCell(system.source.MEM.GetRenderValue(item, model.FieldOpt{})).
				SetExpansion(0).
				SetAlign(tview.AlignRight))
	}

}

func (system *System) DrawVMInfo() {
	system.vm.Clear()
	system.vm.SetOffset(0, 0)

	items := model.DefaultVmFields
	for i, v := range []string{"Field", "Value"} {
		system.vm.SetCell(0, i, tview.NewTableCell(v).SetTextColor(tcell.ColorTeal))
	}

	for i, item := range items {
		system.vm.SetCell(i+1,
			0,
			tview.NewTableCell(item).
				SetExpansion(0).
				SetAlign(tview.AlignLeft))
		system.vm.SetCell(i+1,
			1,
			tview.NewTableCell(system.source.Vm.GetRenderValue(item, model.FieldOpt{})).
				SetExpansion(0).
				SetAlign(tview.AlignRight))
	}

}

func (system *System) DrawDiskInfo() {
	system.disk.Clear()
	system.disk.SetOffset(0, 0)

	visbleCols := model.DefaultDiskFields

	for i, col := range visbleCols {

		system.disk.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorTeal))
	}

	r := 0
	for _, n := range system.source.Disks.GetKeys() {
		disk := system.source.Disks[n]
		for i, col := range visbleCols {
			system.disk.SetCell(r+1,
				i,
				tview.NewTableCell(disk.GetRenderValue(col, model.FieldOpt{})).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
		r++
	}

}

func (system *System) DrawNetInfo() {
	system.net.Clear()
	system.net.SetOffset(0, 0)

	visbleCols := model.DefaultNetDevFields

	for i, col := range visbleCols {

		system.net.SetCell(0, i, tview.NewTableCell(col).SetTextColor(tcell.ColorTeal))
	}

	r := 0
	for _, n := range system.source.Nets.GetKeys() {
		net := system.source.Nets[n]
		for i, col := range visbleCols {
			system.net.SetCell(r+1,
				i,
				tview.NewTableCell(net.GetRenderValue(col, model.FieldOpt{})).
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
