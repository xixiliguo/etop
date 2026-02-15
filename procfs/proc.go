package procfs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

//go:generate go tool stringer -linecomment -output=proc_string.go -type=ProcState,ProcPolicy

// Proc provides information about a running process.
type Proc struct {
	PID int
	fs  *FS
}

func (p Proc) path(file string) string {
	p.fs.bufName = p.fs.bufName[:0]
	p.fs.bufName = append(p.fs.bufName, p.fs.mountPoint...)
	p.fs.bufName = append(p.fs.bufName, "/"...)

	p.fs.bufName = strconv.AppendInt(p.fs.bufName, int64(p.PID), 10)

	p.fs.bufName = append(p.fs.bufName, "/"...)
	p.fs.bufName = append(p.fs.bufName, file...)
	return stringutil.ToString(p.fs.bufName)
}

type ProcState byte

const (
	Running         ProcState = 'R'
	Sleeping        ProcState = 'S'
	Uninterruptible ProcState = 'D'
	Idle            ProcState = 'I'
	Zombie          ProcState = 'Z'
	Stopped         ProcState = 'T'
	TracingStopped  ProcState = 't'
	Dead            ProcState = 'X'
	Deadx           ProcState = 'x'
	Wakekill        ProcState = 'K'
	Waking          ProcState = 'W'
	Parked          ProcState = 'P'
)

type ProcPolicy uint64

const (
	NORMAL ProcPolicy = iota
	FIFO
	RR
	BATCH
	_
	IDLE
	DEADLINE
)

// ProcStat provides status information about the process,
// read from /proc/[pid]/stat.
type ProcStat struct {
	// The process ID.
	PID int
	// The filename of the executable.
	Comm string
	// The process state.
	State ProcState
	// The PID of the parent of this process.
	PPID int
	// The process group ID of the process.
	PGRP int
	// The session ID of the process.
	Session int
	// The controlling terminal of the process.
	TTY int
	// The ID of the foreground process group of the controlling terminal of
	// the process.
	TPGID int
	// The kernel flags word of the process.
	Flags uint64
	// The number of minor faults the process has made which have not required
	// loading a memory page from disk.
	MinFlt uint64
	// The number of minor faults that the process's waited-for children have
	// made.
	CMinFlt uint64
	// The number of major faults the process has made which have required
	// loading a memory page from disk.
	MajFlt uint64
	// The number of major faults that the process's waited-for children have
	// made.
	CMajFlt uint64
	// Amount of time that this process has been scheduled in user mode,
	// measured in clock ticks.
	UTime uint64
	// Amount of time that this process has been scheduled in kernel mode,
	// measured in clock ticks.
	STime uint64
	// Amount of time that this process's waited-for children have been
	// scheduled in user mode, measured in clock ticks.
	CUTime uint64
	// Amount of time that this process's waited-for children have been
	// scheduled in kernel mode, measured in clock ticks.
	CSTime uint64
	// For processes running a real-time scheduling policy, this is the negated
	// scheduling priority, minus one.
	Priority int
	// The nice value, a value in the range 19 (low priority) to -20 (high
	// priority).
	Nice int
	// Number of threads in this process.
	NumThreads int
	// The time the process started after system boot, the value is expressed
	// in clock ticks.
	Starttime uint64
	// Virtual memory size in bytes.
	VSize uint64
	// Resident set size in pages.
	RSS uint64
	// Soft limit in bytes on the rss of the process.
	RSSLimit uint64
	// CPU number last executed on.
	Processor int
	// Real-time scheduling priority, a number in the range 1 to 99 for processes
	// scheduled under a real-time policy, or 0, for non-real-time processes.
	RTPriority uint64
	// Scheduling policy.
	Policy ProcPolicy
	// Aggregated block I/O delays, measured in clock ticks (centiseconds).
	DelayAcctBlkIOTicks uint64
	// Guest time of the process (time spent running a virtual CPU for a guest
	// operating system), measured in clock ticks.
	GuestTime uint64
	// Guest time of the process's children, measured in clock ticks.
	CGuestTime uint64
}

// Stat returns the current status information of the process.
func (p Proc) Stat() (ProcStat, error) {

	path := p.path("stat")

	f, err := os.Open(path)
	if err != nil {
		return ProcStat{}, err
	}

	s := ProcStat{PID: p.PID}

	err = fileutil.ProcessFile(f, func(line string) error {
		var (
			l = strings.Index(line, "(")
			r = strings.LastIndex(line, ")")
		)

		if l < 0 || r < 0 {
			return fmt.Errorf("unexpected format, couldn't extract comm %q", line)
		}

		s.Comm = strings.Clone(line[l+1 : r])

		// Check the following resources for the details about the particular stat
		// fields and their data types:
		// * https://man7.org/linux/man-pages/man5/proc.5.html
		// * https://man7.org/linux/man-pages/man3/scanf.3.html

		var fields [52]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 52 {
			return fmt.Errorf("pid %d: unexpected line in stat: '%s'", p.PID, line)
		}

		for fieldIdx, field := range fields {
			switch fieldIdx {
			case 2: // State
				s.State = ProcState(field[0])
			case 3: // PPID
				s.PPID, err = strconv.Atoi(field)
			case 4: // PGRP
				s.PGRP, err = strconv.Atoi(field)
			case 5: // Session
				s.Session, err = strconv.Atoi(field)
			case 6: // TTY
				s.TTY, err = strconv.Atoi(field)
			case 7: // TPGID
				s.TPGID, err = strconv.Atoi(field)
			case 8: // Flags
				s.Flags, err = strconv.ParseUint(field, 10, 64)
			case 9: // MinFlt
				s.MinFlt, err = strconv.ParseUint(field, 10, 64)
			case 10: // CMinFlt
				s.CMinFlt, err = strconv.ParseUint(field, 10, 64)
			case 11: // MajFlt
				s.MajFlt, err = strconv.ParseUint(field, 10, 64)
			case 12: // CMajFlt
				s.CMajFlt, err = strconv.ParseUint(field, 10, 64)
			case 13: // UTime
				s.UTime, err = strconv.ParseUint(field, 10, 64)
			case 14: // STime
				s.STime, err = strconv.ParseUint(field, 10, 64)
			case 15: // CUTime
				s.CUTime, err = strconv.ParseUint(field, 10, 64)
			case 16: // CSTime
				s.CSTime, err = strconv.ParseUint(field, 10, 64)
			case 17: // Priority
				s.Priority, _ = strconv.Atoi(field)
			case 18: // Nice
				s.Nice, err = strconv.Atoi(field)
			case 19: // NumThreads
				s.NumThreads, err = strconv.Atoi(field)
			case 21: // Starttime
				s.Starttime, err = strconv.ParseUint(field, 10, 64)
			case 22: // VSize
				s.VSize, err = strconv.ParseUint(field, 10, 64)
			case 23: // RSS
				s.RSS, err = strconv.ParseUint(field, 10, 64)
			case 24: // RSSLimit
				s.RSSLimit, err = strconv.ParseUint(field, 10, 64)
			case 38: // Processor
				s.Processor, err = strconv.Atoi(field)
			case 39: // RTPriority
				s.RTPriority, err = strconv.ParseUint(field, 10, 64)
			case 40: // Policy
				var val uint64
				val, err = strconv.ParseUint(field, 10, 64)
				s.Policy = ProcPolicy(val)
			case 41: // DelayAcctBlkIOTicks
				s.DelayAcctBlkIOTicks, err = strconv.ParseUint(field, 10, 64)
			case 42: // GuestTime
				s.GuestTime, err = strconv.ParseUint(field, 10, 64)
			case 43: // CGuestTime
				s.CGuestTime, err = strconv.ParseUint(field, 10, 64)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})

	return s, err
}

// ProcIO models the content of /proc/<pid>/io.
type ProcIO struct {
	// Chars read.
	RChar uint64
	// Chars written.
	WChar uint64
	// Read syscalls.
	SyscR uint64
	// Write syscalls.
	SyscW uint64
	// Bytes read.
	ReadBytes uint64
	// Bytes written.
	WriteBytes uint64
	// Bytes written, but taking into account truncation. See
	// Documentation/filesystems/proc.txt in the kernel sources for
	// detailed explanation.
	CancelledWriteBytes uint64
}

// IO creates a new ProcIO instance from a given Proc instance.
func (p Proc) IO() (ProcIO, error) {

	path := p.path("io")
	f, err := os.Open(path)
	if err != nil {
		return ProcIO{}, err
	}
	defer f.Close()

	pio := ProcIO{}

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("pid %d: unexpected line in stat: '%s'", p.PID, line)
		}
		switch fields[0] {
		case "rchar:":
			pio.RChar, err = strconv.ParseUint(fields[1], 10, 64)
		case "wchar:":
			pio.WChar, err = strconv.ParseUint(fields[1], 10, 64)
		case "syscr:":
			pio.SyscR, err = strconv.ParseUint(fields[1], 10, 64)
		case "syscw:":
			pio.SyscW, err = strconv.ParseUint(fields[1], 10, 64)
		case "read_bytes:":
			pio.ReadBytes, err = strconv.ParseUint(fields[1], 10, 64)
		case "write_bytes:":
			pio.WriteBytes, err = strconv.ParseUint(fields[1], 10, 64)
		case "cancelled_write_bytes:":
			pio.CancelledWriteBytes, err = strconv.ParseUint(fields[1], 10, 64)
		}
		if err != nil {
			return err
		}
		return nil
	})

	return pio, err
}

// ProcSchedstat contains the values from `/proc/<pid>/schedstat`.
type ProcSchedstat struct {
	RunningNanoseconds uint64
	WaitingNanoseconds uint64
	RunTimeslices      uint64
}

// Schedstat returns task scheduling information for the process.
func (p Proc) Schedstat() (ProcSchedstat, error) {

	path := p.path("schedstat")
	f, err := os.Open(path)
	if err != nil {
		return ProcSchedstat{}, err
	}

	s := ProcSchedstat{}

	err = fileutil.ProcessFile(f, func(line string) error {

		var fields [3]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 3 {
			return fmt.Errorf("pid %d: unexpected line in schedstat: '%s'", p.PID, line)
		}

		for fieldIdx, field := range fields {
			switch fieldIdx {
			case 0: // RunningNanoseconds
				s.RunningNanoseconds, err = strconv.ParseUint(field, 10, 64)
				if err != nil {
					return err
				}
			case 1: // WaitingNanoseconds
				s.WaitingNanoseconds, err = strconv.ParseUint(field, 10, 64)
				if err != nil {
					return err
				}
			case 2: // RunTimeslices
				s.RunTimeslices, err = strconv.ParseUint(field, 10, 64)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})

	return s, err
}

func (p Proc) Cgroup() (string, error) {

	path := p.path("cgroup")
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	cgroup := ""
	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		var fields [3]string
		nFields := stringutil.SplitN(line, ":", fields[:])
		if nFields < 3 {
			return fmt.Errorf("pid %d: unexpected line in cgroup: '%s'", p.PID, line)
		}
		if fields[0] == "0" && fields[1] == "" {
			end := len(fields[2])
			for ; end > 0 && fields[2][end-1] == '\n'; end-- {
			}
			cgroup = strings.Clone(fields[2][:end])
		}
		return nil
	})
	return cgroup, nil
}

func (p Proc) CmdLine() (string, error) {

	path := p.path("cmdline")

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	var c string
	err = fileutil.ProcessFile(f, func(line string) error {
		c = strings.ReplaceAll(strings.Clone(line), "\x00", " ")
		return nil
	})

	return c, nil
}
