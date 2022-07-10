package tui

import (
	"regexp"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
)

var (
	GENERALLAYOUT       = []string{"PID", "COMM", "STATE", "CPU", "MEM", "R/s", "W/s"}
	GENERALDEFAULTORDER = "CPU"
	CPULAYOUT           = []string{"PID", "COMM", "CPU", "USERCPU", "SYSCPU", "PRI", "NICE", "PPID", "THR", "STARTTIME"}
	CPUDEFAULTORDER     = "CPU"
	MEMLAYOUT           = []string{"PID", "COMM", "MEM", "MINFLT", "MAJFLT", "VSIZE", "RSS"}
	MEMDEFAULTORDER     = "MEM"
	IOLAYOUT            = []string{"PID", "COMM", "DISK", "RCHAR", "WCHAR", "SYSCR", "SYSCW", "READ", "WRITE", "WCANCEL"}
	IODEFAULTORDER      = "DISK"
)

type ProcessView struct {
	Tui *tview.Application
	*tview.Flex
	Upper          *tview.Flex
	ProcessTable   *tview.Table
	SortView       *tview.List
	SortField      string
	descOrder      bool
	SortToogle     bool
	SearchView     *tview.InputField
	SearchToogle   bool
	SearchText     string
	VisibleColumns []string
	DefaultOrder   string
	Source         []model.Process
}

func NewProcessView(tui *tview.Application) *ProcessView {

	pv := &ProcessView{
		Tui:            tui,
		Flex:           tview.NewFlex(),
		Upper:          tview.NewFlex(),
		ProcessTable:   tview.NewTable(),
		SortView:       tview.NewList(),
		SortField:      "CPU",
		descOrder:      true,
		SortToogle:     false,
		SearchView:     tview.NewInputField(),
		SearchToogle:   false,
		SearchText:     "",
		VisibleColumns: GENERALLAYOUT,
		DefaultOrder:   GENERALDEFAULTORDER,
	}

	pv.ProcessTable.
		SetFixed(1, 1).
		SetBorder(true).
		SetTitle("Process").
		SetTitleAlign(tview.AlignLeft)

	pv.SortView.
		SetTitle("Sort by").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	pv.SortView.ShowSecondaryText(false)

	pv.SortView.SetSelectedFunc(func(i int, mainText, secondText string, r rune) {
		if pv.SortField == mainText {
			pv.descOrder = !pv.descOrder
		} else {
			pv.SortField = mainText
			pv.descOrder = true
		}

		pv.Update()
	})

	pv.SetSortContent(pv.VisibleColumns, pv.DefaultOrder)

	pv.SearchView.
		SetTitle("Search     ESC to close").
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft)
	pv.SearchView.
		SetLabel(">  ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorBlue).
		SetChangedFunc(func(text string) {
			pv.SearchText = text
			pv.Update()
		})

	pv.Flex.
		SetDirection(tview.FlexRow).
		AddItem(pv.Upper.
			AddItem(pv.SortView, 0, 0, false).
			AddItem(pv.ProcessTable, 0, 1, true), 0, 1, true).
		AddItem(pv.SearchView, 0, 0, false)

	pv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if pv.SearchView.HasFocus() {
			if event.Key() == tcell.KeyEsc {
				pv.SwitchSearchToogle()
				pv.UpdateLayerOut()
				return nil
			}
			return event
		}
		if event.Rune() == 's' {
			pv.SwitchSortToogle()
			pv.UpdateLayerOut()
			return nil
		}
		if event.Rune() == '/' {
			pv.SwitchSearchToogle()
			pv.UpdateLayerOut()
			return nil
		}

		if pv.ProcessTable.HasFocus() {
			if event.Rune() == 'g' {
				pv.SetVisibleColumns(GENERALLAYOUT, GENERALDEFAULTORDER)
				pv.Update()
				return nil
			} else if event.Rune() == 'c' {
				pv.SetVisibleColumns(CPULAYOUT, CPUDEFAULTORDER)
				pv.Update()
				return nil
			} else if event.Rune() == 'm' {
				pv.SetVisibleColumns(MEMLAYOUT, MEMDEFAULTORDER)
				pv.Update()
				return nil
			} else if event.Rune() == 'd' {
				pv.SetVisibleColumns(IOLAYOUT, IODEFAULTORDER)
				pv.Update()
				return nil
			}
			return event
		}

		return event
	})
	return pv
}

func (pv *ProcessView) SwitchSortToogle() {
	if pv.SortToogle == true {
		pv.SortToogle = false
	} else {
		pv.SortToogle = true
	}
}

func (pv *ProcessView) SwitchSearchToogle() {
	if pv.SearchToogle == true {
		pv.SearchToogle = false
	} else {
		pv.SearchToogle = true
	}
}

func (pv *ProcessView) UpdateLayerOut() {

	pv.Tui.SetFocus(pv.ProcessTable)
	sortWidth := 0
	if pv.SortToogle {
		sortWidth = 9
		pv.Tui.SetFocus(pv.SortView)
	}
	searchWidth := 0
	if pv.SearchToogle {
		searchWidth = 3
		pv.Tui.SetFocus(pv.SearchView)
	}

	pv.Upper.ResizeItem(pv.SortView, sortWidth, 0)
	pv.ResizeItem(pv.SearchView, searchWidth, 0)

}

func (pv *ProcessView) SetSortField(field string) {
	pv.SortField = field
}

func (pv *ProcessView) SetSource(s []model.Process) {
	pv.Source = s
	pv.Update()
}

func (pv *ProcessView) SetVisibleColumns(c []string, order string) {
	pv.VisibleColumns = c
	pv.SetSortContent(c, order)
}

func (pv *ProcessView) SetSortContent(visleCol []string, order string) {
	pv.SortView.Clear()
	pv.SortField = order
	pv.descOrder = true
	for i, f := range pv.VisibleColumns {
		pv.SortView.AddItem(f, "", 0, nil)
		if pv.SortField == f {
			pv.SortView.SetCurrentItem(i)
		}
	}
}

func (pv *ProcessView) Update() {
	pv.ProcessTable.Clear()
	pv.ProcessTable.SetOffset(0, 0)

	visbleData := []model.Process{}
	if pv.SearchText != "" {
		for _, s := range pv.Source {
			matched, _ := regexp.MatchString(pv.SearchText, s.Comm)
			if matched {
				visbleData = append(visbleData, s)
			}
		}
	} else {
		visbleData = append(visbleData, pv.Source...)
	}

	if pv.SortField != "" {
		sort.SliceStable(visbleData, func(i, j int) bool {
			return model.SortMap[pv.SortField](visbleData[i], visbleData[j])
		})
		if pv.descOrder == false {
			for i := 0; i < len(visbleData)/2; i++ {
				visbleData[i], visbleData[len(visbleData)-1-i] = visbleData[len(visbleData)-1-i], visbleData[i]
			}
		}
	}

	for i, col := range pv.VisibleColumns {
		orderFlag := ""
		if pv.SortField == col {
			if pv.descOrder {
				orderFlag = "▼"
			} else {
				orderFlag = "▲"
			}
		}
		pv.ProcessTable.SetCell(0, i, tview.NewTableCell(col+orderFlag).SetTextColor(tcell.ColorBlue))
	}
	for r := 0; r < len(visbleData); r++ {
		for i, col := range pv.VisibleColumns {
			pv.ProcessTable.SetCell(r+1,
				i,
				tview.NewTableCell(visbleData[r].
					GetRenderValue(col)).
					SetExpansion(1).
					SetAlign(tview.AlignLeft))
		}
	}
}
