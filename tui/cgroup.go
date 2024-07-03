package tui

import (
	"fmt"
	"path/filepath"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

var (
	CGROUPGENERALLAYOUT = []string{"Name", "UsagePercent", "MemoryCurrent", "RbytePerSec", "WbytePerSec",
		"NrDescendants", "NrDyingDescendants", "Controllers"}
	CGROUPGENERALDEFAULTORDER = "Name"
	CGROUPCPULAYOUT           = []string{"Name", "UsagePercent", "UserPercent", "SystemPercent", "NrPeriodsPerSec", "NrThrottledPerSec", "ThrottledPercent", "NrBurstsPerSec", "BurstPercent"}
	CGROUPCPUDEFAULTORDER     = "Name"
	CGROUPMEMLAYOUT           = []string{"Name", "MemoryCurrent", "SwapCurrent", "Anon", "File", "KernelStack", "Slab", "Sock", "Shmem", "Zswap", "Zswapped", "FileMapped",
		"FileDirty", "FileWriteback", "AnonThp", "InactiveAnon", "ActiveAnon", "InactiveFile", "ActiveFile", "Unevictable",
		"SlabReclaimable", "SlabUnreclaimable", "PgfaultPerSec", "PgmajfaultPerSec", "WorkingsetRefaultPerSec", "WorkingsetActivatePerSec",
		"WorkingsetNodereclaimPerSec", "PgrefillPerSec", "PgscanPerSec", "PgstealPerSec", "PgactivatePerSec", "PgdeactivatePerSec", "PglazyfreePerSec", "PglazyfreedPerSec",
		"ZswpInPerSec", "ZswpOutPerSec", "ThpFaultAllocPerSec", "ThpCollapseAllocPerSec",
		"EventLow", "EventHigh", "EventMax", "EventOom", "EventOomKill"}
	CGROUPMEMDEFAULTORDER      = "Name"
	CGROUPIOLAYOUT             = []string{"Name", "RbytePerSec", "WbytePerSec", "RioPerSec", "WioPerSec", "DbytePerSec", "DioPerSec"}
	CGROUPIODEFAULTORDER       = "Name"
	CGROUPPRESSURELAYOUT       = []string{"Name", "CPUSomePressure", "CPUFullPressure", "MemorySomePressure", "MemoryFullPressure", "IOSomePressure", "IOFullPressure"}
	CGROUPPRESSUREDEFAULTORDER = "Name"
)

type Cgroup struct {
	*tview.Flex
	status             *tview.TextView
	statuxText         string
	noSelect           bool
	regions            []string
	currRegionIdx      int
	header             *tview.TextView
	cgroupView         *tview.Table
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
	visbleData         []*model.Cgroup
	source             *model.Model
}

func NewCgroup(status *tview.TextView) *Cgroup {

	cgroup := &Cgroup{
		Flex:           tview.NewFlex(),
		status:         status,
		header:         tview.NewTextView(),
		cgroupView:     tview.NewTable(),
		sortView:       tview.NewList(),
		sortField:      "Name",
		descOrder:      false,
		sortDisplay:    false,
		searchView:     tview.NewInputField(),
		searchDisplay:  false,
		searchText:     "",
		visibleColumns: CGROUPGENERALLAYOUT,
		defaultOrder:   CGROUPGENERALDEFAULTORDER,
	}

	cgroup.regions = []string{"g", "c", "m", "d", "p"}
	fmt.Fprintf(cgroup.header, `["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""]  ["%s"]%s[""] ["%s"]%s[""]`,
		"g", "General",
		"c", "CPU",
		"m", "Mem",
		"d", "I/O",
		"p", "Pressure")
	cgroup.header.SetRegions(true).Highlight("g")

	cgroup.cgroupView.
		SetFixed(1, 2).
		SetSelectable(true, false).
		SetSelectedFunc(func(row, column int) {
			if row > len(cgroup.visbleData) {
				return
			}
			selected := cgroup.visbleData[row-1]
			selected.IsExpand = !selected.IsExpand
			cgroup.noSelect = true
			cgroup.update()
		}).
		SetSelectionChangedFunc(func(row int, column int) {
			cgroup.status.Clear()
			idx := row - 1
			if 0 <= idx && idx < len(cgroup.visbleData) {
				c := cgroup.visbleData[idx]
				extra := filepath.Join(c.Path, c.Name)
				cgroup.statuxText = extra
				cgroup.status.SetText(extra)
			}
		})
	cgroup.SetBorder(true).
		SetTitle("Cgroup").
		SetTitleAlign(tview.AlignLeft)

	cgroup.sortView.
		ShowSecondaryText(false).
		SetSelectedFunc(func(i int, mainText, secondText string, r rune) {
			if cgroup.sortField == cgroup.visibleColumns[i] {
				cgroup.descOrder = !cgroup.descOrder
			} else {
				cgroup.sortField = cgroup.visibleColumns[i]
				cgroup.descOrder = true
				if cgroup.sortField == "Name" {
					cgroup.descOrder = false
				}
			}
			cgroup.update()
		}).
		SetTitle("Sort by").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	cgroup.setVisibleColumns(CGROUPGENERALLAYOUT, CGROUPGENERALDEFAULTORDER)

	cgroup.searchView.
		SetLabel(">  ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorTeal).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				input := cgroup.searchView.GetText()
				if input == "" {
					cgroup.searchText = ""
					cgroup.searchprogram = nil
					cgroup.update()
					return
				}
				program, err := expr.Compile(input, expr.Env(model.Cgroup{}), expr.AsBool())
				if err == nil {
					cgroup.searchText = input
					cgroup.searchprogram = program
					cgroup.update()
				} else {
					cgroup.status.Clear()
					fmt.Fprintf(cgroup.status, "%s", err)
				}
			}
		}).
		SetTitle("Search     ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	cgroup.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(cgroup.sortView, 0, 0, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(cgroup.header, 1, 0, false).
				AddItem(cgroup.cgroupView, 0, 1, true), 0, 1, true), 0, 1, true).
		AddItem(cgroup.searchView, 0, 0, false)

	return cgroup
}

func (cgroup *Cgroup) Focus(delegate func(p tview.Primitive)) {
	if cgroup.searchDisplay {
		delegate(cgroup.searchView)
		return
	}
	if cgroup.sortDisplay {
		delegate(cgroup.sortView)
		return
	}
	delegate(cgroup.cgroupView)
}

func (cgroup *Cgroup) setRegionAndSwitchView(region string) {
	for i, r := range cgroup.regions {
		if r == region {
			cgroup.currRegionIdx = i
		}
	}
	cgroup.header.Highlight(region)
	cgroup.prevVisibleColumns = cgroup.visibleColumns
	switch region {
	case "g":
		cgroup.setVisibleColumns(CGROUPGENERALLAYOUT, CGROUPGENERALDEFAULTORDER)
	case "c":
		cgroup.setVisibleColumns(CGROUPCPULAYOUT, CGROUPCPUDEFAULTORDER)
	case "m":
		cgroup.setVisibleColumns(CGROUPMEMLAYOUT, CGROUPMEMDEFAULTORDER)
	case "d":
		cgroup.setVisibleColumns(CGROUPIOLAYOUT, CGROUPIODEFAULTORDER)
	case "p":
		cgroup.setVisibleColumns(CGROUPPRESSURELAYOUT, CGROUPPRESSUREDEFAULTORDER)
	}
	cgroup.update()
}

func (cgroup *Cgroup) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cgroup.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if cgroup.searchDisplay {
			if event.Key() == tcell.KeyEsc {
				cgroup.searchDisplay = false
				cgroup.ResizeItem(cgroup.searchView, 0, 0)
				cgroup.Focus(setFocus)
				return
			}
			if handler := cgroup.searchView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if event.Rune() == 's' {
			upper := cgroup.GetItem(0).(*tview.Flex)
			sortWidth := 0
			if cgroup.sortDisplay {
				cgroup.sortDisplay = false
				sortWidth = 0
			} else {
				cgroup.sortDisplay = true
				sortWidth = 13
			}
			upper.ResizeItem(cgroup.sortView, sortWidth, 0)
			cgroup.Focus(setFocus)
			return
		}
		if event.Rune() == '/' {
			searchWidth := 0
			if cgroup.searchDisplay {
				cgroup.searchDisplay = false
				searchWidth = 0
			} else {
				cgroup.searchDisplay = true
				searchWidth = 3
			}
			cgroup.ResizeItem(cgroup.searchView, searchWidth, 0)
			cgroup.Focus(setFocus)
			return
		}

		if cgroup.sortDisplay {
			if handler := cgroup.sortView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if cgroup.cgroupView.HasFocus() {

			if event.Key() == tcell.KeyTab {
				nextId := (cgroup.currRegionIdx + 1) % len(cgroup.regions)
				region := cgroup.regions[nextId]
				cgroup.setRegionAndSwitchView(region)
				return
			} else if event.Key() == tcell.KeyBacktab {
				nextId := (cgroup.currRegionIdx - 1 + len(cgroup.regions)) % len(cgroup.regions)
				region := cgroup.regions[nextId]
				cgroup.setRegionAndSwitchView(region)
				return
			} else if event.Rune() == 'g' {
				cgroup.setRegionAndSwitchView("g")
				return
			} else if event.Rune() == 'c' {
				cgroup.setRegionAndSwitchView("c")
				return
			} else if event.Rune() == 'm' {
				cgroup.setRegionAndSwitchView("m")
				return
			} else if event.Rune() == 'd' {
				cgroup.setRegionAndSwitchView("d")
				return
			} else if event.Rune() == 'p' {
				cgroup.setRegionAndSwitchView("p")
				return
			}
			if handler := cgroup.cgroupView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
			return
		}
	})
}

func (cgroup *Cgroup) SelectedCgroupPath() string {
	row, _ := cgroup.cgroupView.GetSelection()
	c := cgroup.visbleData[row-1]
	return filepath.Join(c.Path, c.Name)
}

func (cgroup *Cgroup) SetSource(s *model.Model) {
	cgroup.source = s
	cgroup.update()
}

func (cgroup *Cgroup) setVisibleColumns(cols []string, order string) {

	cgroup.sortView.Clear()
	cgroup.sortField = order
	cgroup.descOrder = false

	cgroup.visibleColumns = cols
	cgroup.visibleColumnsText = make([]string, len(cgroup.visibleColumns))
	p := model.Cgroup{}
	for i, col := range cols {
		text := p.DefaultConfig(col).Name
		cgroup.visibleColumnsText[i] = text
		cgroup.sortView.AddItem(text, "", 0, nil)
		if cgroup.sortField == col {
			cgroup.sortView.SetCurrentItem(i)
		}
	}
}

func (cgroup *Cgroup) update() {
	row, column := cgroup.cgroupView.GetOffset()
	cgroup.cgroupView.Clear()
	cgroup.cgroupView.SetOffset(row, column)
	title := "Cgroup"
	if cgroup.searchText != "" {
		title += " Filter: " + cgroup.searchText
	}
	cgroup.SetTitle(title)

	cgroup.visbleData = cgroup.source.Cgroup.Iterate(cgroup.searchprogram, cgroup.sortField, cgroup.descOrder)

	for i, col := range cgroup.visibleColumns {
		text := cgroup.visibleColumnsText[i]
		orderFlag := ""
		if cgroup.sortField == col {
			if cgroup.descOrder {
				orderFlag = "▼"
			} else {
				orderFlag = "▲"
			}
		}
		cgroup.cgroupView.SetCell(0, i, tview.NewTableCell(text+orderFlag).SetTextColor(tcell.ColorTeal).SetSelectable(false))
	}
	for r := 0; r < len(cgroup.visbleData); r++ {
		for i, col := range cgroup.visibleColumns {
			width := 0
			if col == "Name" {
				width = 50
			}
			cgroup.cgroupView.SetCell(r+1,
				i,
				tview.NewTableCell(cgroup.visbleData[r].
					GetRenderValue(col, model.FieldOpt{FixWidth: true})).
					SetExpansion(1).
					SetAlign(tview.AlignLeft).
					SetMaxWidth(width))
		}
	}

	if cgroup.noSelect {
		cgroup.noSelect = false
		return
	}
	cgroup.cgroupView.Select(1, 0)
}
