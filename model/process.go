package model

import (
	"fmt"
	"time"

	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/util"
)

const (
	userHZ = 100
)

type Process struct {
	PID        int
	Comm       string
	State      string
	PPID       int
	NumThreads int
	Starttime  uint64
	PCPU
	PMEM
	PIO
}

type PCPU struct {
	UTime    float64
	STime    float64
	Priority int
	Nice     int
	CPUUsage int
}

func (c *PCPU) GetRenderValue(field string) string {

	switch field {
	case "USERCPU":
		return fmt.Sprintf("%.2fs", c.UTime)
	case "SYSCPU":
		return fmt.Sprintf("%.2fs", c.STime)
	case "PRI":
		return fmt.Sprintf("%d", c.Priority)
	case "NICE":
		return fmt.Sprintf("%d", c.Nice)
	case "CPU":
		return fmt.Sprintf("%d%%", c.CPUUsage)
	}
	return ""
}

type PMEM struct {
	MinFlt   uint
	MajFlt   uint
	VSize    uint
	RSS      int
	MEMUsage int
}

func (m *PMEM) GetRenderValue(field string) string {
	switch field {
	case "MINFLT":
		return fmt.Sprintf("%d", m.MinFlt)
	case "MAJFLT":
		return fmt.Sprintf("%d", m.MajFlt)
	case "VSIZE":
		return util.GetHumanSize(int(m.VSize))
	case "RSS":
		return util.GetHumanSize((m.RSS))
	case "MEM":
		return fmt.Sprintf("%d%%", m.MEMUsage)

	}
	return ""
}

type PIO struct {
	RChar               uint64
	WChar               uint64
	SyscR               uint64
	SyscW               uint64
	ReadBytes           uint64
	WriteBytes          uint64
	ReadBytesPerSec     uint64
	WriteBytesPerSec    uint64
	CancelledWriteBytes int64
	DiskUage            int
}

func (i *PIO) GetRenderValue(field string) string {
	switch field {
	case "RCHAR":
		return util.GetHumanSize(int(i.RChar))
	case "WCHAR":
		return util.GetHumanSize(int(i.WChar))
	case "SYSCR":
		return fmt.Sprintf("%d", i.SyscR)
	case "SYSCW":
		return fmt.Sprintf("%d", i.SyscW)
	case "READ":
		return util.GetHumanSize(int(i.ReadBytes))
	case "WRITE":
		return util.GetHumanSize(int(i.WriteBytes))
	case "WCANCEL":
		return util.GetHumanSize(int(i.CancelledWriteBytes))
	case "R/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(i.ReadBytesPerSec))
	case "W/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(i.WriteBytesPerSec))
	case "DISK":
		return fmt.Sprintf("%d%%", i.DiskUage)
	}
	return ""
}

func (p *Process) GetRenderValue(field string) string {
	switch field {
	case "PID":
		return fmt.Sprintf("%d", p.PID)
	case "COMM":
		return fmt.Sprintf("%-20.20s", p.Comm)
	case "STATE":
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
		return fmt.Sprintf("%s", stodesc[p.State])
	case "PPID":
		return fmt.Sprintf("%d", p.PPID)
	case "THR":
		return fmt.Sprintf("%d", p.NumThreads)
	case "STARTTIME":
		return fmt.Sprintf("%s", time.Unix(int64(p.Starttime), 0))
	case "USERCPU", "SYSCPU", "PRI", "NICE", "CPU":
		return p.PCPU.GetRenderValue(field)
	case "MINFLT", "MAJFLT", "VSIZE", "RSS", "MEM":
		return p.PMEM.GetRenderValue(field)
	case "RCHAR", "WCHAR", "SYSCR", "SYSCW", "READ", "WRITE", "WCANCEL", "R/s", "W/s", "DISK":
		return p.PIO.GetRenderValue(field)
	}
	return ""
}

type sortFunc func(i, j Process) bool

//go:generate go run sort.go

type ProcessMap map[int]Process

func (processMap ProcessMap) Collect(prev, curr *store.Sample) (processes, threads int) {

	for k := range processMap {
		delete(processMap, k)
	}

	interval := curr.TimeStamp - prev.TimeStamp

	totalIO := uint64(0)

	for pid := range curr.ProcSamples {

		bootTime := curr.SystemSample.BootTime
		sample := curr.ProcSamples[pid]
		p := Process{
			PID:        sample.PID,
			Comm:       sample.Comm,
			State:      sample.State,
			PPID:       sample.PPID,
			NumThreads: sample.NumThreads,
			Starttime:  bootTime + sample.Starttime/userHZ,
		}

		// get cpu info
		p.UTime = float64((sample.UTime - prev.ProcSamples[pid].UTime)) / userHZ
		p.STime = float64((sample.STime - prev.ProcSamples[pid].STime)) / userHZ
		p.Priority = sample.Priority
		p.Nice = sample.Nice
		p.CPUUsage = int((p.UTime + p.STime) * 100 / float64(interval))

		p.MinFlt = sample.MinFlt
		p.MajFlt = sample.MajFlt
		p.VSize = sample.VSize
		p.RSS = sample.RSS * curr.PageSize
		p.MEMUsage = p.RSS * 100 / 1024 / int(*curr.MemTotal)

		p.RChar = sample.RChar - prev.ProcSamples[pid].RChar
		p.WChar = sample.WChar - prev.ProcSamples[pid].WChar
		p.SyscR = sample.SyscR - prev.ProcSamples[pid].SyscR
		p.SyscW = sample.SyscW - prev.ProcSamples[pid].SyscW
		p.ReadBytes = sample.ReadBytes - prev.ProcSamples[pid].ReadBytes
		p.WriteBytes = sample.WriteBytes - prev.ProcSamples[pid].WriteBytes
		p.CancelledWriteBytes = sample.CancelledWriteBytes - prev.ProcSamples[pid].CancelledWriteBytes
		p.ReadBytesPerSec = p.ReadBytes / uint64(interval)
		p.WriteBytesPerSec = p.WriteBytes / uint64(interval)

		processMap[pid] = p

		totalIO += p.ReadBytes + p.WriteBytes
		processes += 1
		threads += p.NumThreads

	}
	if totalIO != 0 {
		for pid, proc := range processMap {
			proc.DiskUage = int((proc.ReadBytes + proc.WriteBytes) * 100 / totalIO)
			processMap[pid] = proc
		}
	}

	return processes, threads
}
