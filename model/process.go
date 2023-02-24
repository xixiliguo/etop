package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/xixiliguo/etop/store"
)

const (
	userHZ = 100
)

var DefaultProcessFields = []string{"Pid", "Comm", "State", "CPU", "Mem", "R/s", "W/s"}

type Process struct {
	Pid        int
	Comm       string
	State      string
	Ppid       int
	NumThreads int
	StartTime  uint64
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

func (c *PCPU) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for process cpu stat", field)
	switch field {
	case "UserCPU":
		s = config[field].Render(c.UTime)
	case "SysCPU":
		s = config[field].Render(c.STime)
	case "Pri":
		s = config[field].Render(c.Priority)
	case "Nice":
		s = config[field].Render(c.Nice)
	case "CPU":
		s = config[field].Render(c.CPUUsage)
	}
	return s
}

type PMEM struct {
	MinFlt   uint
	MajFlt   uint
	VSize    uint
	RSS      int
	MemUsage int
}

func (m *PMEM) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for process mem stat", field)
	switch field {
	case "Minflt":
		s = config[field].Render(m.MinFlt)
	case "Majflt":
		s = config[field].Render(m.MajFlt)
	case "Vsize":
		s = config[field].Render(m.VSize)
	case "RSS":
		s = config[field].Render(m.RSS)
	case "Mem":
		s = config[field].Render(m.MemUsage)
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

func (i *PIO) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for process io stat", field)

	switch field {
	case "Rchar":
		s = config[field].Render(i.RChar)
	case "Wchar":
		s = config[field].Render(i.WChar)
	case "Rchar/s":
		s = config[field].Render(i.RCharPerSec)
	case "Wchar/s":
		s = config[field].Render(i.WCharPerSec)
	case "Syscr":
		s = config[field].Render(i.SyscR)
	case "Syscw":
		s = config[field].Render(i.SyscW)
	case "Syscr/s":
		s = config[field].Render(i.SyscRPerSec)
	case "Syscw/s":
		s = config[field].Render(i.SyscWPerSec)
	case "Read":
		s = config[field].Render(i.ReadBytes)
	case "Write":
		s = config[field].Render(i.WriteBytes)
	case "Wcancel":
		s = config[field].Render(i.CancelledWriteBytes)
	case "R/s":
		s = config[field].Render(i.ReadBytesPerSec)
	case "W/s":
		s = config[field].Render(i.WriteBytesPerSec)
	case "CW/s":
		s = config[field].Render(i.CancelledWriteBytesPerSec)
	case "Disk":
		s = config[field].Render(i.DiskUage)
	}
	return s
}

func (p *Process) GetRenderValue(config RenderConfig, field string) string {

	s := fmt.Sprintf("no %s for process stat", field)
	switch field {
	case "Pid":
		s = config[field].Render(p.Pid)
	case "Comm":
		s = config[field].Render(p.Comm)
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
		s = config[field].Render(stodesc[p.State])
	case "Ppid":
		s = config[field].Render(p.Ppid)
	case "Thr":
		s = config[field].Render(p.NumThreads)
	case "StartTime":
		startTime := time.Unix(int64(p.StartTime), 0).Format(time.RFC3339)
		s = config[field].Render(startTime)
	case "UserCPU", "SysCPU", "Pri", "Nice", "CPU":
		return p.PCPU.GetRenderValue(config, field)
	case "Minflt", "Majflt", "Vsize", "RSS", "Mem":
		return p.PMEM.GetRenderValue(config, field)
	case "Rchar", "Wchar", "Rchar/s", "Wchar/s",
		"Syscr", "Syscw", "Syscr/s", "Syscw/s",
		"Read", "Write", "Wcancel", "R/s", "W/s", "CW/s", "Disk":
		return p.PIO.GetRenderValue(config, field)
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
		p.MemUsage = p.RSS * 100 / 1024 / int(*curr.MemTotal)

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

func (processMap ProcessMap) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)

	processList := []Process{}
	for _, p := range processMap {
		processList = append(processList, p)
	}

	sort.SliceStable(processList, func(i, j int) bool {
		return SortMap[opt.SortField](processList[i], processList[j])
	})
	if opt.AscendingOrder == true {
		for i := 0; i < len(processList)/2; i++ {
			processList[i], processList[len(processList)-1-i] = processList[len(processList)-1-i], processList[i]
		}
	}

	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
		cnt := 0
	looptext:
		for _, p := range processList {
			row := strings.Builder{}
			row.WriteString(dateTime)
			for _, f := range opt.Fields {
				renderValue := p.GetRenderValue(config, f)
				if f == opt.SelectField && opt.Filter != nil {
					if opt.Filter.MatchString(renderValue) == false {
						continue looptext
					}
				}
				row.WriteString(" ")
				row.WriteString(renderValue)
			}
			row.WriteString("\n")

			opt.Output.WriteString(row.String())
			cnt++
			if opt.Top > 0 && opt.Top == cnt {
				break
			}
		}
	case "json":
		t := []any{}
		cnt := 0
	loopjson:
		for _, p := range processList {
			row := make(map[string]string)
			row["Timestamp"] = dateTime
			for _, f := range opt.Fields {
				renderValue := p.GetRenderValue(config, f)
				if f == opt.SelectField && opt.Filter != nil {
					if opt.Filter.MatchString(renderValue) == false {
						continue loopjson
					}
				}
				row[config[f].Name] = renderValue
			}
			t = append(t, row)
			cnt++
			if opt.Top > 0 && opt.Top == cnt {
				break
			}
		}
		b, _ := json.Marshal(t)
		opt.Output.Write(b)
	}
}

func SubWithInterval[T int | int64 | float64](curr, prev, interval T) T {
	if interval == 0 {
		return 0
	}
	return (curr - prev) / interval
}
