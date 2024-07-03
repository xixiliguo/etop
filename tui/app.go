package tui

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/util"
)

const (
	LIVE = iota
	REPORT
)

type TUI struct {
	*tview.Application
	pages   *tview.Pages
	base    *tview.Flex
	header  *Header
	basic   *Basic
	process *Process
	system  *System
	cgroup  *Cgroup
	detail  *tview.Pages
	status  *tview.TextView
	search  *InputDialog
	help    *Help
	log     *slog.Logger
	mode    int
	sm      *model.Model
}

func NewTUI() *TUI {
	tui := &TUI{
		Application: tview.NewApplication(),
		pages:       tview.NewPages(),
		base:        tview.NewFlex(),
		header:      NewHeader(),
		basic:       NewBasic(),
		system:      NewSystem(),
		detail:      tview.NewPages(),
		status:      tview.NewTextView(),
		search:      NewInputDialog(),
		help:        NewHelp(),
	}

	tui.log = util.CreateLogger(tui.status, true)

	tui.process = NewProcess(tui.status)
	tui.cgroup = NewCgroup(tui.status)

	tui.status.SetBorder(true)

	tui.initSearch()
	tui.initDetails()
	tui.initBase()
	tui.initPages()

	return tui
}

func (tui *TUI) initSearch() {

	tui.search.form.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			timeStamp, err := util.ConvertToUnixTime(tui.search.form.GetText())
			if err != nil {
				msg := fmt.Sprintf("user-input time: %s", err)
				tui.status.Clear()
				tui.log.Error(msg)
				return
			}
			if err := tui.sm.CollectSampleByTime(timeStamp); err != nil {
				msg := fmt.Sprintf("search sample by %s: %s", time.Unix(timeStamp, 0).Format(time.RFC3339), err)
				tui.status.Clear()
				tui.log.Error(msg)
				return
			}
			tui.header.Update(tui.sm)
			tui.basic.Update(tui.sm)
			tui.process.SetSource(tui.sm)
			tui.system.SetSource(tui.sm)
			tui.cgroup.SetSource(tui.sm)
			tui.search.form.SetText("")
			tui.pages.HidePage("search")
			return
		}
	})
}

func (tui *TUI) initDetails() {
	tui.detail.AddPage("Process", tui.process, true, true)
	tui.detail.AddPage("System", tui.system, true, false)
	tui.detail.AddPage("Cgroup", tui.cgroup, true, false)

	tui.detail.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 || event.Name() == "Alt+Rune[1]" {
			tui.detail.SwitchToPage("Process")
			return nil
		} else if event.Key() == tcell.KeyF2 || event.Name() == "Alt+Rune[2]" {
			tui.detail.SwitchToPage("System")
			return nil
		} else if event.Key() == tcell.KeyF3 || event.Name() == "Alt+Rune[3]" {
			if len(tui.cgroup.visbleData) > 1 {
				tui.detail.SwitchToPage("Cgroup")
			}
			return nil
		}
		name, _ := tui.detail.GetFrontPage()
		if name == "Cgroup" && event.Rune() == 'z' {
			fullPath := tui.cgroup.SelectedCgroupPath()
			filterRule := fmt.Sprintf("Cgroup startsWith %q", fullPath)
			if err := tui.process.SetFilterRule(filterRule); err != nil {
				tui.status.Clear()
				fmt.Fprintf(tui.status, "%s %s", fullPath, err)
			} else {
				tui.detail.SwitchToPage("Process")
			}
			return nil
		}
		return event
	})
}
func (tui *TUI) initBase() {
	tui.base.SetDirection(tview.FlexRow).
		AddItem(tui.header, 3, 1, false).
		AddItem(tui.basic, 8, 1, false).
		AddItem(tui.detail, 0, 1, true).
		AddItem(tui.status, 3, 0, false)
}

func (tui *TUI) initPages() {
	tui.pages.AddPage("base", tui.base, true, true).
		AddPage("search", tui.search, true, false).
		AddPage("help", tui.help, true, false)

	tui.pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		name, _ := tui.pages.GetFrontPage()
		if name == "search" || name == "help" {
			if event.Key() == tcell.KeyEsc {
				tui.pages.HidePage(name)
				return nil
			}
			return event
		}
		if cur := tui.GetFocus(); cur == tui.process.searchView || cur == tui.cgroup.searchView {
			return event
		}
		if event.Rune() == 't' {
			if tui.mode == REPORT {
				if err := tui.sm.CollectNext(); err != nil {
					msg := fmt.Sprintf("next sample: %s", err)
					tui.status.Clear()
					tui.log.Error(msg)
					return nil
				}
				tui.header.Update(tui.sm)
				tui.basic.Update(tui.sm)
				tui.process.SetSource(tui.sm)
				tui.system.SetSource(tui.sm)
				tui.cgroup.SetSource(tui.sm)
			}
			return nil
		} else if event.Rune() == 'T' {
			if tui.mode == REPORT {
				if err := tui.sm.CollectPrev(); err != nil {
					msg := fmt.Sprintf("previous sample: %s", err)
					tui.status.Clear()
					tui.log.Error(msg)
					return nil
				}
				tui.header.Update(tui.sm)
				tui.basic.Update(tui.sm)
				tui.process.SetSource(tui.sm)
				tui.system.SetSource(tui.sm)
				tui.cgroup.SetSource(tui.sm)
			}
			return nil
		} else if event.Rune() == 'b' {
			if tui.mode == REPORT {
				tui.pages.ShowPage("search")
			}
			return nil
		} else if event.Rune() == 'h' {
			tui.pages.ShowPage("help")
			return nil
		} else if event.Rune() == 'q' {
			tui.Stop()
			return nil
		}
		return event
	})
}

func (tui *TUI) Run(path string, beginTime string) error {
	tui.mode = REPORT

	local, err := store.NewLocalStore(
		store.WithPathAndLogger(path, tui.log),
	)
	if err != nil {
		return err
	}
	sm, err := model.NewSysModel(local, tui.log)
	if err != nil {
		return err
	}

	begin := int64(0)
	if begin, err = util.ConvertToUnixTime(beginTime); err != nil {
		return err
	}

	if err := sm.CollectSampleByTime(begin); err != nil {
		return err
	}

	tui.sm = sm
	tui.header.Update(sm)
	tui.basic.Update(sm)
	tui.process.SetSource(sm)
	tui.system.SetSource(sm)
	tui.cgroup.SetSource(sm)

	if err := tui.Application.SetRoot(tui.pages, true).SetFocus(tui.pages).Run(); err != nil {
		return err
	}
	return nil
}

func (tui *TUI) RunWithLive(interval time.Duration) error {

	exit := store.NewExitProcess(tui.log)
	go exit.Collect()

	tui.mode = LIVE
	sm, err := model.NewSysModel(nil, tui.log)
	if err != nil {
		return err
	}
	sm.Mode = "live"

	tui.sm = sm
	go func() {
		for {

			start := time.Now()
			if err := sm.CollectLiveSample(exit); err != nil {
				return
			}

			tui.QueueUpdateDraw(func() {
				tui.header.Update(sm)
				tui.basic.Update(sm)
				tui.process.SetSource(sm)
				tui.system.SetSource(sm)
				tui.cgroup.SetSource(sm)
			})

			collectDuration := time.Since(start)
			sleepDuration := time.Duration(1 * time.Second)
			if interval-collectDuration > 1*time.Second {
				sleepDuration = interval - collectDuration
			}
			time.Sleep(sleepDuration)
		}
	}()

	if err := tui.Application.SetRoot(tui.pages, true).SetFocus(tui.pages).Run(); err != nil {
		return err
	}
	return nil
}
