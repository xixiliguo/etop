package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultSystemFields = []string{"Load1", "Load5", "Load15", "NumCPU",
	"Processes", "Threads",
	"ProcessesRunning", "ProcessesBlocked",
	"ClonePerSec", "ContextSwitchPerSec"}

type System struct {
	Load1               float64
	Load5               float64
	Load15              float64
	NumCPU              uint64
	Processes           uint64
	Threads             uint64
	ProcessesRunning    uint64
	ProcessesBlocked    uint64
	ClonePerSec         float64
	ContextSwitchPerSec float64
}

func (sys *System) DefaultConfig(field string) Field {

	cfg := Field{}
	switch field {
	case "Load1":
		cfg = Field{"Load1", Raw, 0, "", 10, false}
	case "Load5":
		cfg = Field{"Load5", Raw, 0, "", 10, false}
	case "Load15":
		cfg = Field{"Load15", Raw, 0, "", 10, false}
	case "NumCPU":
		cfg = Field{"NumCPU", Raw, 0, "", 10, false}
	case "Processes":
		cfg = Field{"Process", Raw, 0, "", 10, false}
	case "Threads":
		cfg = Field{"Thread", Raw, 0, "", 10, false}
	case "ProcessesRunning":
		cfg = Field{"Running", Raw, 0, "", 10, false}
	case "ProcessesBlocked":
		cfg = Field{"Blocked", Raw, 0, "", 10, false}
	case "ClonePerSec":
		cfg = Field{"Clone/s", Raw, 1, "/s", 10, false}
	case "ContextSwitchPerSec":
		cfg = Field{"CtxSw/s", Raw, 1, "/s", 10, false}
	}
	return cfg
}

func (sys *System) GetRenderValue(field string, opt FieldOpt) string {

	cfg := sys.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Load1":
		s = cfg.Render(sys.Load1)
	case "Load5":
		s = cfg.Render(sys.Load5)
	case "Load15":
		s = cfg.Render(sys.Load15)
	case "NumCPU":
		s = cfg.Render(sys.NumCPU)
	case "Processes":
		s = cfg.Render(sys.Processes)
	case "Threads":
		s = cfg.Render(sys.Threads)
	case "ProcessesRunning":
		s = cfg.Render(sys.ProcessesRunning)
	case "ProcessesBlocked":
		s = cfg.Render(sys.ProcessesBlocked)
	case "ClonePerSec":
		s = cfg.Render(sys.ClonePerSec)
	case "ContextSwitchPerSec":
		s = cfg.Render(sys.ContextSwitchPerSec)
	default:
		s = "no " + field + " for cpu stat"
	}
	return s
}

func (sys *System) Collect(prev, curr *store.Sample) {

	sys.Load1 = curr.Load1
	sys.Load5 = curr.Load5
	sys.Load15 = curr.Load15
	sys.NumCPU = uint64(len(curr.CPU))
	sys.ProcessesRunning = curr.ProcessesRunning
	sys.ProcessesBlocked = curr.ProcessesBlocked

	interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)
	sys.ClonePerSec = float64(curr.ProcessCreated-prev.ProcessCreated) / float64(interval)
	sys.ContextSwitchPerSec = float64(curr.ContextSwitches-prev.ContextSwitches) / float64(interval)

}
