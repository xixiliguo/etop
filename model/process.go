package model

import (
	"os"
	"time"

	"github.com/xixiliguo/etop/store"
)

const (
	userHZ = 100
)

var (
	enableBootTimeTick = true
)

var DefaultProcessFields = []string{"Pid", "Comm", "State", "CPU", "Mem", "ReadBytePerSec", "WriteBytePerSec"}
var AllProcessFields = []string{"Pid", "Comm", "State", "Ppid", "NumThreads", "StartTime", "OnCPU", "CmdLine",
	"User", "System", "Priority", "Nice", "CPU",
	"MinFlt", "MajFlt", "VSize", "RSS", "Mem",
	"ReadCharPerSec", "WriteCharPerSec",
	"SyscRPerSec", "SyscWPerSec",
	"ReadBytePerSec", "WriteBytePerSec", "CancelledWriteBytePerSec", "Disk"}

type Process struct {
	Pid        int
	Comm       string
	State      string
	Ppid       int
	NumThreads int
	StartTime  uint64
	EndTime    uint64
	ExitCode   uint32
	OnCPU      uint
	CmdLine    string
	PCPU
	PMEM
	PIO
}

type PCPU struct {
	User     float64
	System   float64
	Priority int
	Nice     int
	CPU      float64
}

func (c *PCPU) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "User":
		cfg = Field{"User", Raw, 1, "%", 10, false}
	case "System":
		cfg = Field{"System", Raw, 1, "%", 10, false}
	case "Priority":
		cfg = Field{"Priority", Raw, 0, "", 10, false}
	case "Nice":
		cfg = Field{"Nice", Raw, 0, "", 10, false}
	case "CPU":
		cfg = Field{"CPU", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (c *PCPU) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "User":
		cfg = OpenMetricField{"User", Gauge, "", "", []string{"Pid"}}
	case "System":
		cfg = OpenMetricField{"System", Gauge, "", "", []string{"Pid"}}
	case "Priority":
		cfg = OpenMetricField{"Priority", Gauge, "", "", []string{"Pid"}}
	case "Nice":
		cfg = OpenMetricField{"Nice", Gauge, "", "", []string{"Pid"}}
	case "CPU":
		cfg = OpenMetricField{"CPU", Gauge, "", "", []string{"Pid"}}
	}
	return cfg
}

func (c *PCPU) GetRenderValue(field string, opt FieldOpt) string {
	cfg := c.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "User":
		s = cfg.Render(c.User)
	case "System":
		s = cfg.Render(c.System)
	case "Priority":
		s = cfg.Render(c.Priority)
	case "Nice":
		s = cfg.Render(c.Nice)
	case "CPU":
		s = cfg.Render(c.CPU)
	default:
		s = "no " + field + " for process cpu stat"
	}
	return s
}

type PMEM struct {
	MinFlt uint
	MajFlt uint
	VSize  uint
	RSS    int
	Mem    float64
}

func (m *PMEM) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "MinFlt":
		cfg = Field{"MinFlt", Raw, 0, "", 10, false}
	case "MajFlt":
		cfg = Field{"MajFlt", Raw, 0, "", 10, false}
	case "VSize":
		cfg = Field{"VSize", HumanReadableSize, 0, "", 10, false}
	case "RSS":
		cfg = Field{"RSS", HumanReadableSize, 0, "", 10, false}
	case "Mem":
		cfg = Field{"Mem", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (m *PMEM) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "MinFlt":
		cfg = OpenMetricField{"MinFlt", Gauge, "", "", []string{"Pid"}}
	case "MajFlt":
		cfg = OpenMetricField{"MajFlt", Gauge, "", "", []string{"Pid"}}
	case "VSize":
		cfg = OpenMetricField{"VSize", Gauge, "", "", []string{"Pid"}}
	case "RSS":
		cfg = OpenMetricField{"RSS", Gauge, "", "", []string{"Pid"}}
	case "Mem":
		cfg = OpenMetricField{"Mem", Gauge, "", "", []string{"Pid"}}
	}
	return cfg
}

func (m *PMEM) GetRenderValue(field string, opt FieldOpt) string {
	cfg := m.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "MinFlt":
		s = cfg.Render(m.MinFlt)
	case "MajFlt":
		s = cfg.Render(m.MajFlt)
	case "VSize":
		s = cfg.Render(m.VSize)
	case "RSS":
		s = cfg.Render(m.RSS)
	case "Mem":
		s = cfg.Render(m.Mem)
	default:
		s = "no " + field + " for process mem stat"
	}
	return s
}

type PIO struct {
	RChar                    uint64
	WChar                    uint64
	ReadCharPerSec           float64
	WriteCharPerSec          float64
	SyscR                    uint64
	SyscW                    uint64
	SyscRPerSec              float64
	SyscWPerSec              float64
	ReadBytes                uint64
	WriteBytes               uint64
	CancelledWriteBytes      int64
	ReadBytePerSec           float64
	WriteBytePerSec          float64
	CancelledWriteBytePerSec float64
	Disk                     float64
}

func (i *PIO) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "ReadCharPerSec":
		cfg = Field{"ReadChar/s", HumanReadableSize, 1, "/s", 10, false}
	case "WriteCharPerSec":
		cfg = Field{"WriteChar/s", HumanReadableSize, 1, "/s", 10, false}
	case "SyscR":
		cfg = Field{"SyscR", Raw, 0, "", 10, false}
	case "SyscW":
		cfg = Field{"SyscW", Raw, 0, "", 10, false}
	case "SyscRPerSec":
		cfg = Field{"SyscR/s", Raw, 1, "/s", 10, false}
	case "SyscWPerSec":
		cfg = Field{"SyscW/s", Raw, 1, "/s", 10, false}
	case "ReadBytes":
		cfg = Field{"ReadBytes", HumanReadableSize, 0, "", 10, false}
	case "WriteBytes":
		cfg = Field{"WriteBytes", HumanReadableSize, 0, "", 10, false}
	case "CancelledWriteBytes":
		cfg = Field{"CancelledWriteBytes", HumanReadableSize, 0, "", 10, false}
	case "ReadBytePerSec":
		cfg = Field{"ReadByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "WriteBytePerSec":
		cfg = Field{"WriteByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "CancelledWriteBytePerSec":
		cfg = Field{"CancelledWriteByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "Disk":
		cfg = Field{"Disk", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (i *PIO) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "ReadCharPerSec":
		cfg = OpenMetricField{"ReadCharPerSec", Gauge, "", "", []string{"Pid"}}
	case "WriteCharPerSec":
		cfg = OpenMetricField{"WriteCharPerSec", Gauge, "", "", []string{"Pid"}}
	case "SyscR":
		cfg = OpenMetricField{"SyscR", Gauge, "", "", []string{"Pid"}}
	case "SyscW":
		cfg = OpenMetricField{"SyscW", Gauge, "", "", []string{"Pid"}}
	case "SyscRPerSec":
		cfg = OpenMetricField{"SyscRPerSec", Gauge, "", "", []string{"Pid"}}
	case "SyscWPerSec":
		cfg = OpenMetricField{"SyscWPerSec", Gauge, "", "", []string{"Pid"}}
	case "ReadBytes":
		cfg = OpenMetricField{"ReadBytes", Gauge, "", "", []string{"Pid"}}
	case "WriteBytes":
		cfg = OpenMetricField{"WriteBytes", Gauge, "", "", []string{"Pid"}}
	case "CancelledWriteBytes":
		cfg = OpenMetricField{"CancelledWriteBytes", Gauge, "", "", []string{"Pid"}}
	case "ReadBytePerSec":
		cfg = OpenMetricField{"ReadBytePerSec", Gauge, "", "", []string{"Pid"}}
	case "WriteBytePerSec":
		cfg = OpenMetricField{"WriteBytePerSec", Gauge, "", "", []string{"Pid"}}
	case "CancelledWriteBytePerSec":
		cfg = OpenMetricField{"CancelledWriteBytePerSec", Gauge, "", "", []string{"Pid"}}
	case "Disk":
		cfg = OpenMetricField{"Disk", Gauge, "", "", []string{"Pid"}}
	}
	return cfg
}

func (i *PIO) GetRenderValue(field string, opt FieldOpt) string {
	cfg := i.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Rchar":
		s = cfg.Render(i.RChar)
	case "Wchar":
		s = cfg.Render(i.WChar)
	case "ReadCharPerSec":
		s = cfg.Render(i.ReadCharPerSec)
	case "WriteCharPerSec":
		s = cfg.Render(i.WriteCharPerSec)
	case "SyscR":
		s = cfg.Render(i.SyscR)
	case "SyscW":
		s = cfg.Render(i.SyscW)
	case "SyscRPerSec":
		s = cfg.Render(i.SyscRPerSec)
	case "SyscWPerSec":
		s = cfg.Render(i.SyscWPerSec)
	case "ReadBytes":
		s = cfg.Render(i.ReadBytes)
	case "WriteBytes":
		s = cfg.Render(i.WriteBytes)
	case "CancelledWriteBytes":
		s = cfg.Render(i.CancelledWriteBytes)
	case "ReadBytePerSec":
		s = cfg.Render(i.ReadBytePerSec)
	case "WriteBytePerSec":
		s = cfg.Render(i.WriteBytePerSec)
	case "CancelledWriteBytePerSec":
		s = cfg.Render(i.CancelledWriteBytePerSec)
	case "Disk":
		s = cfg.Render(i.Disk)
	default:
		s = "no " + field + " for process io stat"
	}
	return s
}

func (p *Process) DefaultConfig(field string) Field {

	cfg := Field{}
	switch field {
	case "Pid":
		cfg = Field{"Pid", Raw, 0, "", 10, false}
	case "Comm":
		cfg = Field{"Comm", Raw, 0, "", 16, false}
	case "State":
		cfg = Field{"State", Raw, 0, "", 10, false}
	case "Ppid":
		cfg = Field{"Ppid", Raw, 0, "", 10, false}
	case "NumThreads":
		cfg = Field{"NumThreads", Raw, 0, "", 10, false}
	case "StartTime":
		cfg = Field{"StartTime", Raw, 0, "", 10, false}
	case "OnCPU":
		cfg = Field{"OnCPU", Raw, 0, "", 10, false}
	case "CmdLine":
		cfg = Field{"CmdLine", Raw, 0, "", 10, false}
	case "User", "System", "Priority", "Nice", "CPU":
		return p.PCPU.DefaultConfig(field)
	case "MinFlt", "MajFlt", "VSize", "RSS", "Mem":
		return p.PMEM.DefaultConfig(field)
	case "Rchar", "Wchar", "ReadCharPerSec", "WriteCharPerSec",
		"SyscR", "SyscW", "SyscRPerSec", "SyscWPerSec",
		"ReadBytes", "WriteBytes", "CancelledWriteBytes", "ReadBytePerSec", "WriteBytePerSec", "CancelledWriteBytePerSec", "Disk":
		return p.PIO.DefaultConfig(field)
	}
	return cfg
}

func (p *Process) DefaultOMConfig(field string) OpenMetricField {

	cfg := OpenMetricField{}
	switch field {
	case "Pid":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"Pid"}}
	case "Comm":
		cfg = OpenMetricField{"Comm", Gauge, "", "", []string{"Pid"}}
	case "State":
		cfg = OpenMetricField{"State", Gauge, "", "", []string{"Pid"}}
	case "Ppid":
		cfg = OpenMetricField{"Ppid", Gauge, "", "", []string{"Pid"}}
	case "NumThreads":
		cfg = OpenMetricField{"NumThreads", Gauge, "", "", []string{"Pid"}}
	case "StartTime":
		cfg = OpenMetricField{"StartTime", Gauge, "", "", []string{"Pid"}}
	case "OnCPU":
		cfg = OpenMetricField{"OnCPU", Gauge, "", "", []string{"Pid"}}
	case "CmdLine":
		cfg = OpenMetricField{"CmdLine", Gauge, "", "", []string{"Pid"}}
	case "User", "System", "Priority", "Nice", "CPU":
		return p.PCPU.DefaultOMConfig(field)
	case "MinFlt", "MajFlt", "VSize", "RSS", "Mem":
		return p.PMEM.DefaultOMConfig(field)
	case "Rchar", "Wchar", "ReadCharPerSec", "WriteCharPerSec",
		"SyscR", "SyscW", "SyscRPerSec", "SyscWPerSec",
		"ReadBytes", "WriteBytes", "CancelledWriteBytes", "ReadBytePerSec", "WriteBytePerSec", "CancelledWriteBytePerSec", "Disk":
		return p.PIO.DefaultOMConfig(field)
	}
	return cfg
}

func (p *Process) GetRenderValue(field string, opt FieldOpt) string {

	cfg := p.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Pid":
		s = cfg.Render(p.Pid)
	case "Comm":
		s = cfg.Render(p.Comm)
	case "State":
		stodesc := map[string]string{
			"R": "Running",
			"S": "Sleeping",
			"D": "Uninterruptible",
			"I": "Idle",
			"Z": "Zombie",
			"T": "Stopped",
			"t": "Tracing stop",
			"X": "Dead",
			"x": "Dead",
			"K": "Wakekill",
			"W": "Waking",
			"P": "Parked",
		}
		s = cfg.Render(stodesc[p.State])
	case "Ppid":
		s = cfg.Render(p.Ppid)
	case "NumThreads":
		s = cfg.Render(p.NumThreads)
	case "StartTime":
		startTime := time.Unix(int64(p.StartTime), 0).Format(time.RFC3339)
		s = cfg.Render(startTime)
	case "OnCPU":
		s = cfg.Render(p.OnCPU)
	case "CmdLine":
		s = cfg.Render(p.CmdLine)
	case "User", "System", "Priority", "Nice", "CPU":
		return p.PCPU.GetRenderValue(field, opt)
	case "MinFlt", "MajFlt", "VSize", "RSS", "Mem":
		return p.PMEM.GetRenderValue(field, opt)
	case "Rchar", "Wchar", "ReadCharPerSec", "WriteCharPerSec",
		"SyscR", "SyscW", "SyscRPerSec", "SyscWPerSec",
		"ReadBytes", "WriteBytes", "CancelledWriteBytes", "ReadBytePerSec", "WriteBytePerSec", "CancelledWriteBytePerSec", "Disk":
		return p.PIO.GetRenderValue(field, opt)
	}
	return s
}

type sortFunc func(i, j Process) bool

//go:generate go run sort.go

type ProcessMap map[int]Process

func (processMap ProcessMap) Collect(prev, curr *store.Sample) (processes, threads uint64) {

	for k := range processMap {
		delete(processMap, k)
	}

	interval := curr.TimeStamp - prev.TimeStamp

	totalIO := uint64(0)

	for pid := range curr.ProcSamples {

		bootTime := curr.BootTimeTick
		if bootTime == 0 || enableBootTimeTick == false {
			bootTime = curr.BootTime * 100
		}
		new := curr.ProcSamples[pid]
		old := prev.ProcSamples[pid]

		if old.Starttime != new.Starttime {
			// new created process during samples
			old = store.ProcSample{}
		}

		p := Process{
			Pid:        new.PID,
			Comm:       new.Comm,
			State:      new.State,
			Ppid:       new.PPID,
			NumThreads: new.NumThreads,
			StartTime:  (bootTime + new.Starttime) / userHZ,
			OnCPU:      new.Processor,
			CmdLine:    new.CmdLine,
		}
		if new.EndTime != 0 {
			// exited process from ebpf have not cmdline info
			// use old one
			p.CmdLine = old.CmdLine
			p.EndTime = (bootTime + new.EndTime) / userHZ
			p.ExitCode = new.ExitCode
		}

		// get cpu info
		p.User = SubWithInterval(float64(new.UTime), float64(old.UTime), float64(interval))
		p.System = SubWithInterval(float64(new.STime), float64(old.STime), float64(interval))
		p.Priority = new.Priority
		p.Nice = new.Nice
		p.CPU = p.User + p.System

		p.MinFlt = new.MinFlt - old.MinFlt
		p.MajFlt = new.MajFlt - old.MajFlt
		p.VSize = new.VSize
		p.RSS = new.RSS * curr.PageSize
		p.Mem = float64(p.RSS) * 100 / 1024 / float64(*curr.MemTotal)

		p.RChar = new.RChar - old.RChar
		p.WChar = new.WChar - old.WChar
		p.ReadCharPerSec = SubWithInterval(float64(new.RChar), float64(old.RChar), float64(interval))
		p.WriteCharPerSec = SubWithInterval(float64(new.WChar), float64(old.WChar), float64(interval))
		p.SyscR = new.SyscR - old.SyscR
		p.SyscW = new.SyscW - old.SyscW
		p.SyscRPerSec = SubWithInterval(float64(new.SyscR), float64(old.SyscR), float64(interval))
		p.SyscWPerSec = SubWithInterval(float64(new.SyscW), float64(old.SyscW), float64(interval))
		p.ReadBytes = new.ReadBytes - old.ReadBytes
		p.WriteBytes = new.WriteBytes - old.WriteBytes
		p.CancelledWriteBytes = new.CancelledWriteBytes - old.CancelledWriteBytes
		p.ReadBytePerSec = float64(p.ReadBytes) / float64(interval)
		p.WriteBytePerSec = float64(p.WriteBytes) / float64(interval)
		p.CancelledWriteBytePerSec = float64(p.CancelledWriteBytes) / float64(interval)
		processMap[pid] = p

		totalIO += p.ReadBytes + p.WriteBytes
		processes += 1
		threads += uint64(p.NumThreads)

	}
	if totalIO != 0 {
		for pid, proc := range processMap {
			proc.Disk = float64(proc.ReadBytes+proc.WriteBytes) * 100 / float64(totalIO)
			processMap[pid] = proc
		}
	}

	return processes, threads
}

func SubWithInterval[T int | int64 | float64](curr, prev, interval T) T {
	if interval == 0 {
		return 0
	}
	if curr < prev {
		return 0
	}
	return (curr - prev) / interval
}

func init() {
	if os.Getenv("boottimetick") == "off" {
		enableBootTimeTick = false
	}
}
