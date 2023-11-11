package tui

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

var (
	GENERALLAYOUT       = []string{"Pid", "Comm", "State", "CPU", "Mem", "ReadBytePerSec", "WriteBytePerSec"}
	GENERALDEFAULTORDER = "CPU"
	CPULAYOUT           = []string{"Pid", "Comm", "CPU", "User", "System", "Priority", "Nice", "Ppid", "NumThreads", "OnCPU", "StartTime"}
	CPUDEFAULTORDER     = "CPU"
	MEMLAYOUT           = []string{"Pid", "Comm", "Mem", "MinFlt", "MajFlt", "VSize", "RSS"}
	MEMDEFAULTORDER     = "Mem"
	IOLAYOUT            = []string{"Pid", "Comm", "Disk", "ReadCharPerSec", "WriteCharPerSec", "SyscRPerSec", "SyscWPerSec", "ReadBytePerSec", "WriteBytePerSec", "CancelledWriteBytePerSec"}
	IODEFAULTORDER      = "Disk"
)

type Process struct {
	*tview.Flex
	status             *tview.TextView
	regions            []string
	currRegionIdx      int
	header             *tview.TextView
	processView        *tview.Table
	processDisplay     bool
	sortView           *tview.List
	sortField          string
	descOrder          bool
	sortDisplay        bool
	searchView         *tview.InputField
	searchDisplay      bool
	searchText         string
	searchprogram      *vm.Program
	prevVisibleColumns []string
	visibleColumns     []string
	visibleColumnsText []string
	defaultOrder       string
	visbleData         []model.Process
	source             *model.Model
}

func NewProcess(status *tview.TextView) *Process {

	process := &Process{
		Flex:           tview.NewFlex(),
		status:         status,
		header:         tview.NewTextView(),
		processView:    tview.NewTable(),
		sortView:       tview.NewList(),
		sortField:      "CPU",
		descOrder:      true,
		sortDisplay:    false,
		searchView:     tview.NewInputField(),
		searchDisplay:  false,
		searchText:     "",
		visibleColumns: GENERALLAYOUT,
		defaultOrder:   GENERALDEFAULTORDER,
	}

	process.regions = []string{"g", "c", "m", "d"}
	fmt.Fprintf(process.header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]`,
		"g", "General",
		"c", "CPU",
		"m", "Mem",
		"d", "I/O")
	process.header.SetRegions(true).Highlight("g")

	process.processView.
		SetFixed(1, 2).
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(row int, column int) {
			process.status.Clear()
			idx := row - 1
			if 0 <= idx && idx < len(process.visbleData) {
				p := process.visbleData[idx]
				extra := p.CmdLine
				if extra == "" {
					extra = p.Comm
				}
				fmt.Fprintf(process.status, "%s", extra)
				if p.State == "x" || p.State == "X" {
					fmt.Fprintf(process.status, " exit code %d at %s",
						p.ExitCode,
						time.Unix(int64(p.EndTime), 0).Format(time.RFC3339))
				}

			}
		})
	process.SetBorder(true).
		SetTitle("Process").
		SetTitleAlign(tview.AlignLeft)

	process.sortView.
		ShowSecondaryText(false).
		SetSelectedFunc(func(i int, mainText, secondText string, r rune) {
			if process.sortField == process.visibleColumns[i] {
				process.descOrder = !process.descOrder
			} else {
				process.sortField = process.visibleColumns[i]
				process.descOrder = true
			}
			process.update()
		}).
		SetTitle("Sort by").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	process.setVisibleColumns(GENERALLAYOUT, GENERALDEFAULTORDER)

	process.searchView.
		SetLabel(">  ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorTeal).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				input := process.searchView.GetText()
				if input == "" {
					process.searchText = ""
					process.searchprogram = nil
					process.update()
					return
				}
				program, err := expr.Compile(input, expr.Env(model.Process{}), expr.AsBool())
				if err == nil {
					process.searchText = input
					process.searchprogram = program
					process.update()
				} else {
					process.status.Clear()
					fmt.Fprintf(process.status, "%s", err)
				}
			}
		}).
		SetTitle("Search     ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	process.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(process.sortView, 0, 0, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(process.header, 1, 0, false).
				AddItem(process.processView, 0, 1, true), 0, 1, true), 0, 1, true).
		AddItem(process.searchView, 0, 0, false)

	return process
}

func (process *Process) Focus(delegate func(p tview.Primitive)) {
	if process.searchDisplay {
		delegate(process.searchView)
		return
	}
	if process.sortDisplay {
		delegate(process.sortView)
		return
	}
	delegate(process.processView)
	return
}

func (process *Process) setRegionAndSwitchView(region string) {
	for i, r := range process.regions {
		if r == region {
			process.currRegionIdx = i
		}
	}
	process.header.Highlight(region)
	process.prevVisibleColumns = process.visibleColumns
	switch region {
	case "g":
		process.setVisibleColumns(GENERALLAYOUT, GENERALDEFAULTORDER)
	case "c":
		process.setVisibleColumns(CPULAYOUT, CPUDEFAULTORDER)
	case "m":
		process.setVisibleColumns(MEMLAYOUT, MEMDEFAULTORDER)
	case "d":
		process.setVisibleColumns(IOLAYOUT, IODEFAULTORDER)
	}
	process.update()
	return
}

func (process *Process) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return process.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if process.searchDisplay {
			if event.Key() == tcell.KeyEsc {
				process.searchDisplay = false
				process.ResizeItem(process.searchView, 0, 0)
				process.Focus(setFocus)
				return
			}
			if handler := process.searchView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if event.Rune() == 's' {
			upper := process.GetItem(0).(*tview.Flex)
			sortWidth := 0
			if process.sortDisplay == true {
				process.sortDisplay = false
				sortWidth = 0
			} else {
				process.sortDisplay = true
				sortWidth = 9
			}
			upper.ResizeItem(process.sortView, sortWidth, 0)
			process.Focus(setFocus)
			return
		}
		if event.Rune() == '/' {
			searchWidth := 0
			if process.searchDisplay == true {
				process.searchDisplay = false
				searchWidth = 0
			} else {
				process.searchDisplay = true
				searchWidth = 3
			}
			process.ResizeItem(process.searchView, searchWidth, 0)
			process.Focus(setFocus)
			return
		}

		if process.sortDisplay {
			if handler := process.sortView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if process.processView.HasFocus() {

			if event.Key() == tcell.KeyTab {
				nextId := (process.currRegionIdx + 1) % len(process.regions)
				region := process.regions[nextId]
				process.setRegionAndSwitchView(region)
				return
			} else if event.Key() == tcell.KeyBacktab {
				nextId := (process.currRegionIdx - 1 + len(process.regions)) % len(process.regions)
				region := process.regions[nextId]
				process.setRegionAndSwitchView(region)
				return
			} else if event.Rune() == 'g' {
				process.setRegionAndSwitchView("g")
				return
			} else if event.Rune() == 'c' {
				process.setRegionAndSwitchView("c")
				return
			} else if event.Rune() == 'm' {
				process.setRegionAndSwitchView("m")
				return
			} else if event.Rune() == 'd' {
				process.setRegionAndSwitchView("d")
				return
			}
			if handler := process.processView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
			return
		}
		return
	})
}

func (process *Process) SetSource(s *model.Model) {
	process.source = s
	process.update()
}

func (process *Process) setVisibleColumns(cols []string, order string) {

	process.sortView.Clear()
	process.sortField = order
	process.descOrder = true

	process.visibleColumns = cols
	process.visibleColumnsText = make([]string, len(process.visibleColumns))
	p := model.Process{}
	for i, col := range cols {
		text := p.DefaultConfig(col).Name
		process.visibleColumnsText[i] = text
		process.sortView.AddItem(text, "", 0, nil)
		if process.sortField == col {
			process.sortView.SetCurrentItem(i)
		}
	}
}

func (process *Process) update() {
	process.processView.Clear()
	process.processView.SetOffset(0, 0)
	title := "Process"
	if process.searchText != "" {
		title += " Filter: " + process.searchText
	}
	process.SetTitle(title)

	process.visbleData = process.visbleData[:0]
	if process.searchText != "" {
		for _, s := range process.source.Processes {
			output, _ := expr.Run(process.searchprogram, s)
			if output.(bool) {
				process.visbleData = append(process.visbleData, s)
			}
		}
	} else {
		for _, p := range process.source.Processes {
			process.visbleData = append(process.visbleData, p)
		}
	}

	if process.sortField != "" {
		sort.SliceStable(process.visbleData, func(i, j int) bool {
			return model.SortMap[process.sortField](process.visbleData[i], process.visbleData[j])
		})
		if process.descOrder == false {
			for i := 0; i < len(process.visbleData)/2; i++ {
				process.visbleData[i],
					process.visbleData[len(process.visbleData)-1-i] =
					process.visbleData[len(process.visbleData)-1-i],
					process.visbleData[i]
			}
		}
	}

	for i, col := range process.visibleColumns {
		text := process.visibleColumnsText[i]
		orderFlag := ""
		if process.sortField == col {
			if process.descOrder {
				orderFlag = "▼"
			} else {
				orderFlag = "▲"
			}
		}
		process.processView.SetCell(0, i, tview.NewTableCell(text+orderFlag).SetTextColor(tcell.ColorTeal).SetSelectable(false))
	}
	for r := 0; r < len(process.visbleData); r++ {
		for i, col := range process.visibleColumns {
			width := 0
			if col == "Comm" {
				width = 16
			}
			process.processView.SetCell(r+1,
				i,
				tview.NewTableCell(process.visbleData[r].
					GetRenderValue(col, model.FieldOpt{FixWidth: true})).
					SetExpansion(1).
					SetAlign(tview.AlignLeft).
					SetMaxWidth(width))
		}
	}

	if slices.Compare(process.prevVisibleColumns, process.visibleColumns) != 0 {
		process.processView.Select(1, 0)
	}
}
