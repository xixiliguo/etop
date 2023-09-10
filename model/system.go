package model

import (
	"fmt"

	"github.com/xixiliguo/etop/store"
)

var DefaultSystemFields = []string{"Load1", "Load5", "Load15",
	"Processes", "Threads",
	"ProcessesRunning", "ProcessesBlocked",
	"ClonePerSec", "ContextSwitchPerSec"}

type System struct {
	Load1               float64
	Load5               float64
	Load15              float64
	Processes           uint64
	Threads             uint64
	ProcessesRunning    uint64
	ProcessesBlocked    uint64
	ClonePerSec         float64
	ContextSwitchPerSec float64
}

func (sys *System) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for system stat", field)
	switch field {
	case "Load1":
		s = config[field].Render(sys.Load1)
	case "Load5":
		s = config[field].Render(sys.Load5)
	case "Load15":
		s = config[field].Render(sys.Load15)
	case "Processes":
		s = config[field].Render(sys.Processes)
	case "Threads":
		s = config[field].Render(sys.Threads)
	case "ProcessesRunning":
		s = config[field].Render(sys.ProcessesRunning)
	case "ProcessesBlocked":
		s = config[field].Render(sys.ProcessesBlocked)
	case "ClonePerSec":
		s = config[field].Render(sys.ClonePerSec)
	case "ContextSwitchPerSec":
		s = config[field].Render(sys.ContextSwitchPerSec)
	}
	return s
}

func (sys *System) Collect(prev, curr *store.Sample) {

	sys.Load1 = curr.Load1
	sys.Load5 = curr.Load5
	sys.Load15 = curr.Load15
	sys.ProcessesRunning = curr.ProcessesRunning
	sys.ProcessesBlocked = curr.ProcessesBlocked

	interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)
	sys.ClonePerSec = float64(curr.ProcessCreated-prev.ProcessCreated) / float64(interval)
	sys.ContextSwitchPerSec = float64(curr.ContextSwitches-prev.ContextSwitches) / float64(interval)

	return
}
