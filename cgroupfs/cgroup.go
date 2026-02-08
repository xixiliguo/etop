package cgroupfs

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
	"golang.org/x/sys/unix"
)

var (
	CgroupV2MountPoint = "/sys/fs/cgroup"
)

var IsCgroup2 = sync.OnceValue(func() bool {
	var st unix.Statfs_t
	if err := unix.Statfs(CgroupV2MountPoint, &st); err != nil {
		return false
	}
	if st.Type == unix.CGROUP2_SUPER_MAGIC {
		return true
	}

	return false
})

type Cgroup struct {
	FullPath string
	Name     string
	bufPtr   *[]byte
}

func NewCgroup(fullPaht string, name string) Cgroup {
	buf := make([]byte, 1024)
	child := Cgroup{
		FullPath: fullPaht,
		Name:     name,
		bufPtr:   &buf,
	}
	return child
}

func (c *Cgroup) Child(name string) Cgroup {
	child := Cgroup{
		FullPath: filepath.Join(c.FullPath, name),
		Name:     name,
		bufPtr:   c.bufPtr,
	}
	return child
}

func (c *Cgroup) path(file string) string {
	*c.bufPtr = (*c.bufPtr)[:0]
	*c.bufPtr = append(*c.bufPtr, CgroupV2MountPoint...)

	*c.bufPtr = append(*c.bufPtr, c.FullPath...)

	*c.bufPtr = append(*c.bufPtr, "/"...)
	*c.bufPtr = append(*c.bufPtr, file...)

	return stringutil.ToString(*c.bufPtr)
}

func (c *Cgroup) Inode() (uint64, error) {

	info, err := os.Stat(c.path(""))
	if err != nil {
		return 0, err
	}
	t := info.Sys().(*syscall.Stat_t)
	return t.Ino, err
}

func (c *Cgroup) Controllers() (string, error) {

	fullPath := c.path("cgroup.controllers")
	b, err := os.ReadFile(fullPath)
	return stringutil.ToString(bytes.TrimSpace(b)), err
}

type CgoupStat struct {
	NrDescendants      uint64
	NrDyingDescendants uint64
}

func (c *Cgroup) CgoupStat() (CgoupStat, error) {
	fullPath := c.path("cgroup.stat")

	f, err := os.Open(fullPath)
	if err != nil {
		return CgoupStat{}, err
	}
	defer f.Close()

	stat := CgoupStat{
		NrDescendants:      math.MaxUint64,
		NrDyingDescendants: math.MaxUint64,
	}
	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in cgroup.stat: '%s'", fullPath, line)
		}
		switch fields[0] {
		case "nr_descendants":
			stat.NrDescendants, err = strconv.ParseUint(fields[1], 10, 64)
		case "nr_dying_descendants":
			stat.NrDyingDescendants, err = strconv.ParseUint(fields[1], 10, 64)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return stat, err
}

type CPUStat struct {
	UsageUsec     uint64
	UserUsec      uint64
	SystemUsec    uint64
	NrPeriods     uint64
	NrThrottled   uint64
	ThrottledUsec uint64
	NrBursts      uint64
	BurstUsec     uint64
}

func (c *Cgroup) CPUStat() (CPUStat, error) {
	fullPath := c.path("cpu.stat")

	f, err := os.Open(fullPath)
	if err != nil {
		return CPUStat{}, err
	}
	defer f.Close()

	cpuStat := CPUStat{
		UsageUsec:     math.MaxUint64,
		UserUsec:      math.MaxUint64,
		SystemUsec:    math.MaxUint64,
		NrPeriods:     math.MaxUint64,
		NrThrottled:   math.MaxUint64,
		ThrottledUsec: math.MaxUint64,
		NrBursts:      math.MaxUint64,
		BurstUsec:     math.MaxUint64,
	}
	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in cpu.stat: '%s'", fullPath, line)
		}
		switch fields[0] {
		case "usage_usec":
			cpuStat.UsageUsec, err = strconv.ParseUint(fields[1], 10, 64)
		case "user_usec":
			cpuStat.UserUsec, err = strconv.ParseUint(fields[1], 10, 64)
		case "system_usec":
			cpuStat.SystemUsec, err = strconv.ParseUint(fields[1], 10, 64)
		case "nr_periods":
			cpuStat.NrPeriods, err = strconv.ParseUint(fields[1], 10, 64)
		case "nr_throttled":
			cpuStat.NrThrottled, err = strconv.ParseUint(fields[1], 10, 64)
		case "throttled_usec":
			cpuStat.ThrottledUsec, err = strconv.ParseUint(fields[1], 10, 64)
		case "nr_bursts":
			cpuStat.NrBursts, err = strconv.ParseUint(fields[1], 10, 64)
		case "burst_usec":
			cpuStat.BurstUsec, err = strconv.ParseUint(fields[1], 10, 64)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return cpuStat, err
}

type MemoryStat struct {
	Anon                   uint64
	File                   uint64
	Kernel                 uint64
	KernelStack            uint64
	Slab                   uint64
	Sock                   uint64
	Shmem                  uint64
	Zswap                  uint64
	Zswapped               uint64
	FileMapped             uint64
	FileDirty              uint64
	FileWriteback          uint64
	AnonThp                uint64
	FileThp                uint64
	ShmemThp               uint64
	InactiveAnon           uint64
	ActiveAnon             uint64
	InactiveFile           uint64
	ActiveFile             uint64
	Unevictable            uint64
	SlabReclaimable        uint64
	SlabUnreclaimable      uint64
	Pgfault                uint64
	Pgmajfault             uint64
	WorkingsetRefaultAnon  uint64
	WorkingsetRefaultFile  uint64
	WorkingsetActivateAnon uint64
	WorkingsetActivateFile uint64
	WorkingsetRestoreAnon  uint64
	WorkingsetRestoreFile  uint64
	WorkingsetNodereclaim  uint64
	Pgrefill               uint64
	Pgscan                 uint64
	Pgsteal                uint64
	Pgactivate             uint64
	Pgdeactivate           uint64
	Pglazyfree             uint64
	Pglazyfreed            uint64
	ZswpIn                 uint64
	ZswpOut                uint64
	ThpFaultAlloc          uint64
	ThpCollapseAlloc       uint64
}

func (c *Cgroup) MemoryStat() (MemoryStat, error) {
	memStat := MemoryStat{
		Anon:                   math.MaxUint64,
		File:                   math.MaxUint64,
		KernelStack:            math.MaxUint64,
		Slab:                   math.MaxUint64,
		Sock:                   math.MaxUint64,
		Shmem:                  math.MaxUint64,
		Zswap:                  math.MaxUint64,
		Zswapped:               math.MaxUint64,
		FileMapped:             math.MaxUint64,
		FileDirty:              math.MaxUint64,
		FileWriteback:          math.MaxUint64,
		AnonThp:                math.MaxUint64,
		InactiveAnon:           math.MaxUint64,
		ActiveAnon:             math.MaxUint64,
		InactiveFile:           math.MaxUint64,
		ActiveFile:             math.MaxUint64,
		Unevictable:            math.MaxUint64,
		SlabReclaimable:        math.MaxUint64,
		SlabUnreclaimable:      math.MaxUint64,
		Pgfault:                math.MaxUint64,
		Pgmajfault:             math.MaxUint64,
		WorkingsetRefaultAnon:  math.MaxUint64,
		WorkingsetRefaultFile:  math.MaxUint64,
		WorkingsetActivateAnon: math.MaxUint64,
		WorkingsetActivateFile: math.MaxUint64,
		WorkingsetRestoreAnon:  math.MaxUint64,
		WorkingsetRestoreFile:  math.MaxUint64,
		WorkingsetNodereclaim:  math.MaxUint64,
		Pgrefill:               math.MaxUint64,
		Pgscan:                 math.MaxUint64,
		Pgsteal:                math.MaxUint64,
		Pgactivate:             math.MaxUint64,
		Pgdeactivate:           math.MaxUint64,
		Pglazyfree:             math.MaxUint64,
		Pglazyfreed:            math.MaxUint64,
		ZswpIn:                 math.MaxUint64,
		ZswpOut:                math.MaxUint64,
		ThpFaultAlloc:          math.MaxUint64,
		ThpCollapseAlloc:       math.MaxUint64,
	}

	fullPath := c.path("memory.stat")

	f, err := os.Open(fullPath)
	if err != nil {
		return MemoryStat{}, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in memory.stat: '%s'", fullPath, line)
		}
		switch fields[0] {
		case "anon":
			memStat.Anon, err = strconv.ParseUint(fields[1], 10, 64)
		case "file":
			memStat.File, err = strconv.ParseUint(fields[1], 10, 64)
		case "kernel":
			memStat.Kernel, err = strconv.ParseUint(fields[1], 10, 64)
		case "kernel_stack":
			memStat.KernelStack, err = strconv.ParseUint(fields[1], 10, 64)
		case "slab":
			memStat.Slab, err = strconv.ParseUint(fields[1], 10, 64)
		case "sock":
			memStat.Sock, err = strconv.ParseUint(fields[1], 10, 64)
		case "shmem":
			memStat.Shmem, err = strconv.ParseUint(fields[1], 10, 64)
		case "zswap":
			memStat.Zswap, err = strconv.ParseUint(fields[1], 10, 64)
		case "zswapped":
			memStat.Zswapped, err = strconv.ParseUint(fields[1], 10, 64)
		case "file_mapped":
			memStat.FileMapped, err = strconv.ParseUint(fields[1], 10, 64)
		case "file_dirty":
			memStat.FileDirty, err = strconv.ParseUint(fields[1], 10, 64)
		case "file_writeback":
			memStat.FileWriteback, err = strconv.ParseUint(fields[1], 10, 64)
		case "anon_thp":
			memStat.AnonThp, err = strconv.ParseUint(fields[1], 10, 64)
		case "file_thp":
			memStat.FileThp, err = strconv.ParseUint(fields[1], 10, 64)
		case "shmem_thp":
			memStat.ShmemThp, err = strconv.ParseUint(fields[1], 10, 64)
		case "inactive_anon":
			memStat.InactiveAnon, err = strconv.ParseUint(fields[1], 10, 64)
		case "active_anon":
			memStat.ActiveAnon, err = strconv.ParseUint(fields[1], 10, 64)
		case "inactive_file":
			memStat.InactiveFile, err = strconv.ParseUint(fields[1], 10, 64)
		case "active_file":
			memStat.ActiveFile, err = strconv.ParseUint(fields[1], 10, 64)
		case "unevictable":
			memStat.Unevictable, err = strconv.ParseUint(fields[1], 10, 64)
		case "slab_reclaimable":
			memStat.SlabReclaimable, err = strconv.ParseUint(fields[1], 10, 64)
		case "slab_unreclaimable":
			memStat.SlabUnreclaimable, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgfault":
			memStat.Pgfault, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgmajfault":
			memStat.Pgmajfault, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_refault_anon":
			memStat.WorkingsetRefaultAnon, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_refault_file":
			memStat.WorkingsetRefaultFile, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_activate_anon":
			memStat.WorkingsetActivateAnon, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_activate_file":
			memStat.WorkingsetActivateFile, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_restore_anon":
			memStat.WorkingsetRestoreAnon, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_restore_file":
			memStat.WorkingsetRestoreFile, err = strconv.ParseUint(fields[1], 10, 64)
		case "workingset_nodereclaim":
			memStat.WorkingsetNodereclaim, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgrefill":
			memStat.Pgrefill, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgscan":
			memStat.Pgscan, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgsteal":
			memStat.Pgsteal, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgactivate":
			memStat.Pgactivate, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgdeactivate":
			memStat.Pgdeactivate, err = strconv.ParseUint(fields[1], 10, 64)
		case "pglazyfree":
			memStat.Pglazyfree, err = strconv.ParseUint(fields[1], 10, 64)
		case "pglazyfreed":
			memStat.Pglazyfreed, err = strconv.ParseUint(fields[1], 10, 64)
		case "zswpin":
			memStat.ZswpIn, err = strconv.ParseUint(fields[1], 10, 64)
		case "zswpout":
			memStat.ZswpOut, err = strconv.ParseUint(fields[1], 10, 64)
		case "thp_fault_alloc":
			memStat.ThpFaultAlloc, err = strconv.ParseUint(fields[1], 10, 64)
		case "thp_collapse_alloc":
			memStat.ThpCollapseAlloc, err = strconv.ParseUint(fields[1], 10, 64)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return memStat, err
}

type MemoryEvents struct {
	Low          uint64
	High         uint64
	Max          uint64
	Oom          uint64
	OomKill      uint64
	oomGroupKill uint64
}

func (c *Cgroup) MemoryEvents() (MemoryEvents, error) {
	fullPath := c.path("memory.events")

	event := MemoryEvents{
		Low:     math.MaxUint64,
		High:    math.MaxUint64,
		Max:     math.MaxUint64,
		Oom:     math.MaxUint64,
		OomKill: math.MaxUint64,
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return event, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in memory.events: '%s'", fullPath, line)
		}
		switch fields[0] {
		case "low":
			event.Low, err = strconv.ParseUint(fields[1], 10, 64)
		case "high":
			event.High, err = strconv.ParseUint(fields[1], 10, 64)
		case "max":
			event.Max, err = strconv.ParseUint(fields[1], 10, 64)
		case "oom":
			event.Oom, err = strconv.ParseUint(fields[1], 10, 64)
		case "oom_kill":
			event.OomKill, err = strconv.ParseUint(fields[1], 10, 64)
		case "oom_group_kill":
			event.oomGroupKill, err = strconv.ParseUint(fields[1], 10, 64)

		}
		if err != nil {
			return err
		}
		return nil
	})
	return event, err
}

type IOStat struct {
	Major  uint64
	Minor  uint64
	Rbytes uint64
	Wbytes uint64
	Rios   uint64
	Wios   uint64
	Dbytes uint64
	Dios   uint64
}

func (c *Cgroup) IOStats() ([]IOStat, error) {
	fullPath := c.path("io.stat")

	ioStats := []IOStat{}

	f, err := os.Open(fullPath)
	if err != nil {
		return ioStats, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [7]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 7 {
			return fmt.Errorf("%s: unexpected line in io.stat: '%s'", fullPath, line)
		}
		stat := IOStat{}

		idx := strings.Index(fields[0], ":")
		stat.Major, err = strconv.ParseUint(fields[0][:idx], 10, 64)
		if err != nil {
			return err
		}
		stat.Minor, err = strconv.ParseUint(fields[0][idx+1:], 10, 64)
		if err != nil {
			return err
		}

		for _, field := range fields[1:] {
			idx := strings.Index(field, "=")
			switch field[:idx] {
			case "rbytes":
				stat.Rbytes, err = strconv.ParseUint(field[idx+1:], 10, 64)
			case "wbytes":
				stat.Wbytes, err = strconv.ParseUint(field[idx+1:], 10, 64)
			case "rios":
				stat.Rios, err = strconv.ParseUint(field[idx+1:], 10, 64)
			case "wios":
				stat.Wios, err = strconv.ParseUint(field[idx+1:], 10, 64)
			case "dbytes":
				stat.Dbytes, err = strconv.ParseUint(field[idx+1:], 10, 64)
			case "dios":
				stat.Dios, err = strconv.ParseUint(field[idx+1:], 10, 64)
			}
			if err != nil {
				return err
			}
		}
		ioStats = append(ioStats, stat)
		return nil
	})
	return ioStats, err
}

type PSIStats struct {
	Some PSIData
	Full PSIData
}

type PSIData struct {
	Avg10  float64
	Avg60  float64
	Avg300 float64
	Total  uint64
}

func (c *Cgroup) PSIStats(file string) (PSIStats, error) {
	fullPath := c.path(file)

	psi := PSIStats{
		Some: PSIData{
			Avg10:  math.NaN(),
			Avg60:  math.NaN(),
			Avg300: math.NaN(),
			Total:  math.MaxUint64,
		},
		Full: PSIData{
			Avg10:  math.NaN(),
			Avg60:  math.NaN(),
			Avg300: math.NaN(),
			Total:  math.MaxUint64,
		},
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return psi, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [5]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 5 {
			return fmt.Errorf("%s: unexpected line in %s: '%s'", fullPath, file, line)
		}
		var psiData *PSIData
		switch fields[0] {
		case "some":
			psiData = &psi.Some
		case "full":
			psiData = &psi.Full
		}
		if psiData == nil {
			return fmt.Errorf("%s: no some/full in %s: '%s'", fullPath, file, line)
		}

		for _, field := range fields[1:] {
			idx := strings.Index(field, "=")
			switch field[:idx] {
			case "avg10":
				psiData.Avg10, err = strconv.ParseFloat(field[idx+1:], 64)
			case "avg60":
				psiData.Avg60, err = strconv.ParseFloat(field[idx+1:], 64)
			case "avg300":
				psiData.Avg300, err = strconv.ParseFloat(field[idx+1:], 64)
			case "total":
				psiData.Total, err = strconv.ParseUint(field[idx+1:], 10, 64)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	return psi, err
}

type Property struct {
	MemoryLow           uint64
	MemoryHigh          uint64
	MemoryMin           uint64
	MemoryMax           uint64
	MemorySwapMax       uint64
	MemoryZSwapMax      uint64
	CPUWeight           uint64
	CPUMax              string
	CpuSetCpus          string
	CpuSetCpusEffective string
	PidsCurrent         uint64
	PidsMax             uint64
}

const MaxCgroupPropertyUintValue = math.MaxUint64 - 1
const MaxCgroupPropertyStrValue = "no-exist"

func getUintFromPropertyFile(file string) uint64 {
	f, err := os.Open(file)
	if err != nil {
		return math.MaxUint64
	}
	defer f.Close()
	var buf [32]byte
	n, err := f.Read(buf[:])
	if err != nil {
		return math.MaxUint64
	}
	line := stringutil.ToString(buf[:n])
	if line == "max" {
		return MaxCgroupPropertyUintValue
	}
	count, _ := strconv.ParseUint(line, 10, 64)
	return count
}

func getStrFromPropertyFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return MaxCgroupPropertyStrValue
	}
	defer f.Close()
	var buf [32]byte
	n, err := f.Read(buf[:])
	if err != nil {
		return MaxCgroupPropertyStrValue
	}
	return string(buf[:n])
}

func (c *Cgroup) Properties() (Property, error) {

	p := Property{}
	p.MemoryLow = getUintFromPropertyFile("memory.low")
	p.MemoryHigh = getUintFromPropertyFile("memory.high")
	p.MemoryMin = getUintFromPropertyFile("memory.min")
	p.MemoryMax = getUintFromPropertyFile("memory.max")

	p.MemorySwapMax = getUintFromPropertyFile("memory.swap.max")
	p.MemoryZSwapMax = getUintFromPropertyFile("memory.zswap.max")

	p.CPUMax = getStrFromPropertyFile("cpu.max")
	p.CPUWeight = getUintFromPropertyFile("cpu.weight")

	p.CpuSetCpus = getStrFromPropertyFile("cpuset.cpus")
	p.CpuSetCpusEffective = getStrFromPropertyFile("cpuset.cpus.effective")

	p.PidsCurrent = getUintFromPropertyFile("pids.current")
	p.PidsMax = getUintFromPropertyFile("pids.max")

	return p, nil
}
