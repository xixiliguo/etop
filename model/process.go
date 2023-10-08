package model

import (
	"time"

	"github.com/xixiliguo/etop/store"
)

const (
	userHZ = 100
)

var DefaultProcessFields = []string{"Pid", "Comm", "State", "CPU", "Mem", "R/s", "W/s"}
var AllProcessFields = []string{"Pid", "Comm", "State", "Ppid", "Thr", "StartTime", "OnCPU", "CmdLine",
	"UserCPU", "SysCPU", "Pri", "Nice", "CPU",
	"Minflt", "Majflt", "Vsize", "RSS", "Mem",
	"Rchar/s", "Wchar/s",
	"Syscr/s", "Syscw/s",
	"R/s", "W/s", "CW/s", "Disk"}

type Process struct {
	Pid        int
	Comm       string
	State      string
	Ppid       int
	NumThreads int
	StartTime  uint64
	EndTime    uint64
	ExitCode   uint32
	Processor  uint
	CmdLine    string
	PCPU
	PMEM
	PIO
}

type PCPU struct {
	UTime    float64
	STime    float64
	Priority int
	Nice     int
	CPUUsage float64
}

func (c *PCPU) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "UserCPU":
		cfg = Field{"UserCPU", Raw, 1, "%", 10, false}
	case "SysCPU":
		cfg = Field{"SysCPU", Raw, 1, "%", 10, false}
	case "Pri":
		cfg = Field{"Pri", Raw, 0, "", 10, false}
	case "Nice":
		cfg = Field{"Nice", Raw, 0, "", 10, false}
	case "CPU":
		cfg = Field{"CPU", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (c *PCPU) GetRenderValue(field string, opt FieldOpt) string {
	cfg := c.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "UserCPU":
		s = cfg.Render(c.UTime)
	case "SysCPU":
		s = cfg.Render(c.STime)
	case "Pri":
		s = cfg.Render(c.Priority)
	case "Nice":
		s = cfg.Render(c.Nice)
	case "CPU":
		s = cfg.Render(c.CPUUsage)
	default:
		s = "no " + field + " for process cpu stat"
	}
	return s
}

type PMEM struct {
	MinFlt   uint
	MajFlt   uint
	VSize    uint
	RSS      int
	MemUsage float64
}

func (m *PMEM) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Minflt":
		cfg = Field{"Minflt", Raw, 0, "", 10, false}
	case "Majflt":
		cfg = Field{"Majflt", Raw, 0, "", 10, false}
	case "Vsize":
		cfg = Field{"Vsize", HumanReadableSize, 0, "", 10, false}
	case "RSS":
		cfg = Field{"RSS", HumanReadableSize, 0, "", 10, false}
	case "Mem":
		cfg = Field{"Mem", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (m *PMEM) GetRenderValue(field string, opt FieldOpt) string {
	cfg := m.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Minflt":
		s = cfg.Render(m.MinFlt)
	case "Majflt":
		s = cfg.Render(m.MajFlt)
	case "Vsize":
		s = cfg.Render(m.VSize)
	case "RSS":
		s = cfg.Render(m.RSS)
	case "Mem":
		s = cfg.Render(m.MemUsage)
	default:
		s = "no " + field + " for process mem stat"
	}
	return s
}

type PIO struct {
	RChar                     uint64
	WChar                     uint64
	RCharPerSec               float64
	WCharPerSec               float64
	SyscR                     uint64
	SyscW                     uint64
	SyscRPerSec               float64
	SyscWPerSec               float64
	ReadBytes                 uint64
	WriteBytes                uint64
	CancelledWriteBytes       int64
	ReadBytesPerSec           float64
	WriteBytesPerSec          float64
	CancelledWriteBytesPerSec float64
	DiskUage                  float64
}

func (i *PIO) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Rchar/s":
		cfg = Field{"Rchar/s", HumanReadableSize, 1, "/s", 10, false}
	case "Wchar/s":
		cfg = Field{"Wchar/s", HumanReadableSize, 1, "/s", 10, false}
	case "Syscr":
		cfg = Field{"Syscr", Raw, 0, "", 10, false}
	case "Syscw":
		cfg = Field{"Syscw", Raw, 0, "", 10, false}
	case "Syscr/s":
		cfg = Field{"Syscr/s", Raw, 1, "/s", 10, false}
	case "Syscw/s":
		cfg = Field{"Syscw/s", Raw, 1, "/s", 10, false}
	case "Read":
		cfg = Field{"Read", HumanReadableSize, 0, "", 10, false}
	case "Write":
		cfg = Field{"Write", HumanReadableSize, 0, "", 10, false}
	case "Wcancel":
		cfg = Field{"Wcancel", HumanReadableSize, 0, "", 10, false}
	case "R/s":
		cfg = Field{"R/s", HumanReadableSize, 1, "/s", 10, false}
	case "W/s":
		cfg = Field{"W/s", HumanReadableSize, 1, "/s", 10, false}
	case "CW/s":
		cfg = Field{"CW/s", HumanReadableSize, 1, "/s", 10, false}
	case "Disk":
		cfg = Field{"Disk", Raw, 1, "%", 10, false}
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
	case "Rchar/s":
		s = cfg.Render(i.RCharPerSec)
	case "Wchar/s":
		s = cfg.Render(i.WCharPerSec)
	case "Syscr":
		s = cfg.Render(i.SyscR)
	case "Syscw":
		s = cfg.Render(i.SyscW)
	case "Syscr/s":
		s = cfg.Render(i.SyscRPerSec)
	case "Syscw/s":
		s = cfg.Render(i.SyscWPerSec)
	case "Read":
		s = cfg.Render(i.ReadBytes)
	case "Write":
		s = cfg.Render(i.WriteBytes)
	case "Wcancel":
		s = cfg.Render(i.CancelledWriteBytes)
	case "R/s":
		s = cfg.Render(i.ReadBytesPerSec)
	case "W/s":
		s = cfg.Render(i.WriteBytesPerSec)
	case "CW/s":
		s = cfg.Render(i.CancelledWriteBytesPerSec)
	case "Disk":
		s = cfg.Render(i.DiskUage)
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
	case "Thr":
		cfg = Field{"Thr", Raw, 0, "", 10, false}
	case "StartTime":
		cfg = Field{"StartTime", Raw, 0, "", 10, false}
	case "OnCPU":
		cfg = Field{"OnCPU", Raw, 0, "", 10, false}
	case "CmdLine":
		cfg = Field{"CmdLine", Raw, 0, "", 10, false}
	case "UserCPU", "SysCPU", "Pri", "Nice", "CPU":
		return p.PCPU.DefaultConfig(field)
	case "Minflt", "Majflt", "Vsize", "RSS", "Mem":
		return p.PMEM.DefaultConfig(field)
	case "Rchar", "Wchar", "Rchar/s", "Wchar/s",
		"Syscr", "Syscw", "Syscr/s", "Syscw/s",
		"Read", "Write", "Wcancel", "R/s", "W/s", "CW/s", "Disk":
		return p.PIO.DefaultConfig(field)
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
	case "Thr":
		s = cfg.Render(p.NumThreads)
	case "StartTime":
		startTime := time.Unix(int64(p.StartTime), 0).Format(time.RFC3339)
		s = cfg.Render(startTime)
	case "OnCPU":
		s = cfg.Render(p.Processor)
	case "CmdLine":
		s = cfg.Render(p.CmdLine)
	case "UserCPU", "SysCPU", "Pri", "Nice", "CPU":
		return p.PCPU.GetRenderValue(field, opt)
	case "Minflt", "Majflt", "Vsize", "RSS", "Mem":
		return p.PMEM.GetRenderValue(field, opt)
	case "Rchar", "Wchar", "Rchar/s", "Wchar/s",
		"Syscr", "Syscw", "Syscr/s", "Syscw/s",
		"Read", "Write", "Wcancel", "R/s", "W/s", "CW/s", "Disk":
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

		bootTime := curr.SystemSample.BootTime
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
			StartTime:  bootTime + new.Starttime/userHZ,
			Processor:  new.Processor,
			CmdLine:    new.CmdLine,
		}
		if new.EndTime != 0 {
			// exited process from ebpf have not cmdline info
			// use old one
			p.CmdLine = old.CmdLine
			p.EndTime = bootTime + new.EndTime/userHZ
			p.ExitCode = new.ExitCode
		}

		// get cpu info
		p.UTime = SubWithInterval(float64(new.UTime), float64(old.UTime), float64(interval))
		p.STime = SubWithInterval(float64(new.STime), float64(old.STime), float64(interval))
		p.Priority = new.Priority
		p.Nice = new.Nice
		p.CPUUsage = p.UTime + p.STime

		p.MinFlt = new.MinFlt - old.MinFlt
		p.MajFlt = new.MajFlt - old.MajFlt
		p.VSize = new.VSize
		p.RSS = new.RSS * curr.PageSize
		p.MemUsage = float64(p.RSS) * 100 / 1024 / float64(*curr.MemTotal)

		p.RChar = new.RChar - old.RChar
		p.WChar = new.WChar - old.WChar
		p.RCharPerSec = SubWithInterval(float64(new.RChar), float64(old.RChar), float64(interval))
		p.WCharPerSec = SubWithInterval(float64(new.WChar), float64(old.WChar), float64(interval))
		p.SyscR = new.SyscR - old.SyscR
		p.SyscW = new.SyscW - old.SyscW
		p.SyscRPerSec = SubWithInterval(float64(new.SyscR), float64(old.SyscR), float64(interval))
		p.SyscWPerSec = SubWithInterval(float64(new.SyscW), float64(old.SyscW), float64(interval))
		p.ReadBytes = new.ReadBytes - old.ReadBytes
		p.WriteBytes = new.WriteBytes - old.WriteBytes
		p.CancelledWriteBytes = new.CancelledWriteBytes - old.CancelledWriteBytes
		p.ReadBytesPerSec = float64(p.ReadBytes) / float64(interval)
		p.WriteBytesPerSec = float64(p.WriteBytes) / float64(interval)
		p.CancelledWriteBytesPerSec = float64(p.CancelledWriteBytes) / float64(interval)
		processMap[pid] = p

		totalIO += p.ReadBytes + p.WriteBytes
		processes += 1
		threads += uint64(p.NumThreads)

	}
	if totalIO != 0 {
		for pid, proc := range processMap {
			proc.DiskUage = float64(proc.ReadBytes+proc.WriteBytes) * 100 / float64(totalIO)
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
