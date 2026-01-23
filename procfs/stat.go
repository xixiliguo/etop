package procfs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

// CPUStat shows how much time the cpu spend in various stages.
type CPUStat struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

// SoftIRQStat represent the softirq statistics as exported in the procfs stat file.
// A nice introduction can be found at https://0xax.gitbooks.io/linux-insides/content/interrupts/interrupts-9.html
// It is possible to get per-cpu stats by reading `/proc/softirqs`.
type SoftIRQStat struct {
	Hi          uint64
	Timer       uint64
	NetTx       uint64
	NetRx       uint64
	Block       uint64
	BlockIoPoll uint64
	Tasklet     uint64
	Sched       uint64
	Hrtimer     uint64
	Rcu         uint64
}

// Stat represents kernel/system statistics.
type Stat struct {
	// Boot time in seconds since the Epoch.
	BootTime uint64
	// Summed up cpu statistics.
	CPUTotal CPUStat
	// Per-CPU statistics.
	CPU map[int]CPUStat
	// Number of times interrupts were handled, which contains numbered and unnumbered IRQs.
	IRQTotal uint64
	// Number of times a numbered IRQ was triggered.
	IRQ []uint64
	// Number of times a context switch happened.
	ContextSwitches uint64
	// Number of times a process was created.
	ProcessCreated uint64
	// Number of processes currently running.
	ProcessesRunning uint64
	// Number of processes currently blocked (waiting for IO).
	ProcessesBlocked uint64
	// Number of times a softirq was scheduled.
	SoftIRQTotal uint64
	// Detailed softirq statistics.
	SoftIRQ SoftIRQStat
}

func parseCPUStat(fields []string) (CPUStat, error) {
	cpuStat := CPUStat{}
	var err error

	if cpuStat.User, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.Nice, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.System, err = strconv.ParseFloat(fields[3], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.Idle, err = strconv.ParseFloat(fields[4], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.Iowait, err = strconv.ParseFloat(fields[5], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.IRQ, err = strconv.ParseFloat(fields[6], 64); err != nil {
		return cpuStat, err
	}
	if cpuStat.SoftIRQ, err = strconv.ParseFloat(fields[7], 64); err != nil {
		return cpuStat, err
	}

	if cpuStat.Steal, err = strconv.ParseFloat(fields[8], 64); err != nil {
		return cpuStat, err
	}

	if cpuStat.Guest, err = strconv.ParseFloat(fields[9], 64); err != nil {
		return cpuStat, err
	}

	if cpuStat.GuestNice, err = strconv.ParseFloat(fields[10], 64); err != nil {
		return cpuStat, err
	}

	cpuStat.User /= userHZ
	cpuStat.Nice /= userHZ
	cpuStat.System /= userHZ
	cpuStat.Idle /= userHZ
	cpuStat.Iowait /= userHZ
	cpuStat.IRQ /= userHZ
	cpuStat.SoftIRQ /= userHZ
	cpuStat.Steal /= userHZ
	cpuStat.Guest /= userHZ
	cpuStat.GuestNice /= userHZ

	return cpuStat, nil
}

func parseSoftIRQStat(fields []string) (SoftIRQStat, uint64, error) {
	softIRQStat := SoftIRQStat{}
	var total uint64
	var err error

	if total, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Hi, err = strconv.ParseUint(fields[2], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Timer, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.NetTx, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.NetRx, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Block, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.BlockIoPoll, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Tasklet, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Sched, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Hrtimer, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
		return softIRQStat, total, err
	}
	if softIRQStat.Rcu, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
		return softIRQStat, total, err
	}

	return softIRQStat, total, nil
}

func (fs FS) Stat() (Stat, error) {

	stat := Stat{
		CPU: make(map[int]CPUStat),
	}

	path := fs.path("stat")
	f, err := os.Open(path)
	if err != nil {
		return stat, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var err error

		var fields [12]string
		nFields := stringutil.FieldsN(line, fields[:])

		if nFields < 2 {
			return fmt.Errorf("unexpected line in stat: '%s'", line)
		}

		field0 := fields[0]
		switch {
		case strings.HasPrefix(field0, "cpu"):
			if nFields < 11 {
				return fmt.Errorf("unexpected line in stat: '%s'", line)
			}
			if cpuStat, err := parseCPUStat(fields[:]); err != nil {
				return err
			} else {
				if field0 == "cpu" {
					stat.CPUTotal = cpuStat
				} else {
					if n, err := strconv.Atoi(field0[3:]); err != nil {
						return err
					} else {
						stat.CPU[n] = cpuStat
					}
				}
			}
		case field0 == "ctx":
			if stat.ContextSwitches, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			}
		case field0 == "btime":
			if stat.BootTime, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			}
		case field0 == "processes":
			if stat.ProcessCreated, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			}
		case field0 == "procs_running":
			if stat.ProcessesRunning, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			}
		case field0 == "procs_blocked":
			if stat.ProcessesBlocked, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
				return err
			}
		case field0 == "softirq":
			if stat.SoftIRQ, stat.SoftIRQTotal, err = parseSoftIRQStat(fields[:]); err != nil {
				return err
			}
		}
		return nil
	})

	return stat, err
}
