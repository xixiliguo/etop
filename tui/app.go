package tui

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/store"
)

const (
	LIVE = iota
	REPORT
)

type TUI struct {
	*tview.Application
	pages   *tview.Pages
	flex    *tview.Flex
	header  *HeaderView
	basic   *BasicView
	process *ProcessView
	system  *SystemView
	detail  *tview.Pages
	search  *InputDialog
	log     *log.Logger
	mode    int
	sm      *model.System
}

func NewTUI(log *log.Logger) *TUI {
	t := &TUI{
		Application: tview.NewApplication(),
		pages:       tview.NewPages(),
		flex:        tview.NewFlex(),
		log:         log,
	}
	t.header = NewHeaderView()
	t.basic = NewBasicView()
	t.process = NewProcessView(t.Application)
	t.system = NewSystemView(t.Application)

	t.detail = tview.NewPages()
	t.detail.AddPage("Process", t.process, true, true)
	t.detail.AddPage("System", t.system, true, false)

	t.detail.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			t.detail.SwitchToPage("Process")
			t.process.UpdateLayerOut()
			return nil
		} else if event.Key() == tcell.KeyF2 {
			t.detail.SwitchToPage("System")
			return nil
		}
		return event
	})

	t.search = NewInputDialog()
	t.search.form.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := t.sm.CollectSampleByTime(t.search.form.GetText()); err != nil {
				t.log.Printf("%s", err)
				return
			}
			t.header.Update(t.sm)
			t.basic.Update(t.sm)
			t.process.SetSource(t.sm.ProcessList)
			t.system.SetSource(t.sm)
			t.search.form.SetText("")
			t.pages.HidePage("search")
			return
		}
	})

	t.flex.SetDirection(tview.FlexRow).
		AddItem(t.header, 3, 1, false).
		AddItem(t.basic, 8, 1, false).
		AddItem(t.detail, 0, 1, true)

	t.pages.AddPage("base", t.flex, true, true).
		AddPage("search", t.search, true, false).
		AddPage("help", NewHelpView(), true, false)

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		name, _ := t.pages.GetFrontPage()
		if name == "search" || name == "help" {
			if event.Key() == tcell.KeyEsc {
				t.pages.HidePage(name)
				return nil
			}
			return event
		}
		if t.GetFocus() == t.process.SearchView {
			return event
		}
		if event.Rune() == 't' {
			if t.mode == REPORT {
				if err := t.sm.CollectNext(); err != nil {
					t.log.Printf("%s", err)
					return nil
				}
				t.header.Update(t.sm)
				t.basic.Update(t.sm)
				t.process.SetSource(t.sm.ProcessList)
				t.system.SetSource(t.sm)
			}
			return nil
		} else if event.Rune() == 'T' {
			if t.mode == REPORT {
				if err := t.sm.CollectPrev(); err != nil {
					t.log.Printf("%s", err)
					return nil
				}
				t.header.Update(t.sm)
				t.basic.Update(t.sm)
				t.process.SetSource(t.sm.ProcessList)
				t.system.SetSource(t.sm)
			}
			return nil
		} else if event.Rune() == 'b' {
			if t.mode == REPORT {
				t.pages.ShowPage("search")
				t.SetFocus(t.search)
			}
			return nil
		} else if event.Rune() == 'h' {
			t.pages.ShowPage("help")
			return nil
		} else if event.Rune() == 'q' {
			t.Stop()
			return nil
		}
		return event
	})
	return t
}

func (t *TUI) Run(inputFileName string) error {
	t.mode = REPORT
	local, err := store.NewLocalStoreWithReadOnly(inputFileName, t.log)
	if err != nil {
		return err
	}

	sm, err := model.NewSysModel(local, t.log)
	if err != nil {
		return err
	}

	t.sm = sm
	t.header.Update(sm)
	t.basic.Update(sm)
	t.process.SetSource(sm.ProcessList)
	t.system.SetSource(sm)

	if err := t.Application.SetRoot(t.pages, true).SetFocus(t.pages).Run(); err != nil {
		return err
	}
	return nil
}

func (t *TUI) RunWithLive(interval time.Duration) error {
	t.mode = LIVE
	sm, err := model.NewSysModelWithLive(t.log)
	if err != nil {
		return err
	}

	t.sm = sm
	go func() {
		for {

			start := time.Now()
			if err := sm.CollectLiveSample(); err != nil {
				return
			}

			t.QueueUpdateDraw(func() {
				t.header.Update(sm)
				t.basic.Update(sm)
				t.process.SetSource(sm.ProcessList)
				t.system.SetSource(sm)
			})

			collectDuration := time.Now().Sub(start)
			sleepDuration := time.Duration(1 * time.Second)
			if interval-collectDuration > 1*time.Second {
				sleepDuration = interval - collectDuration
			}
			time.Sleep(sleepDuration)
		}
	}()

	if err := t.Application.SetRoot(t.pages, true).SetFocus(t.pages).Run(); err != nil {
		return err
	}
	return nil
}
