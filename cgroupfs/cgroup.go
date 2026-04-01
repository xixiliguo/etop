package cgroupfs

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

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
	cg := Cgroup{
		FullPath: fullPaht,
		Name:     name,
		bufPtr:   &buf,
	}
	return cg
}

func (c *Cgroup) Child(name string) Cgroup {
	child := Cgroup{
		FullPath: filepath.Join(c.FullPath, name),
		Name:     strings.Clone(name),
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

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR {
			return err
		}
	}
}

func (c *Cgroup) processFile(name string, fn func(i int, line string) error) error {

	var (
		fd  int
		err error
	)

	ignoringEINTR(func() error {
		fd, err = unix.Open(name, unix.O_RDONLY|unix.O_CLOEXEC, 0)
		return err
	})

	if err != nil {
		return nil // fmt.Errorf("cgroup open %s: %+w", name, err)
	}
	defer unix.Close(fd)

	*c.bufPtr = (*c.bufPtr)[:0]

	for {
		var (
			n int
			e error
		)

		ignoringEINTR(func() error {
			n, e = unix.Read(fd, (*c.bufPtr)[len(*c.bufPtr):cap(*c.bufPtr)])
			return e
		})

		if e != nil {
			if e == io.EOF {
				break
			}
			return nil // fmt.Errorf("cgroup read %s: %+w", name, e)
		}
		if n == 0 {
			break
		}
		*c.bufPtr = (*c.bufPtr)[:len(*c.bufPtr)+n]
		if len(*c.bufPtr) == cap(*c.bufPtr) {
			*c.bufPtr = append(*c.bufPtr, 0)[:len(*c.bufPtr)]
		}
	}

	i := 0
	for sub := range bytes.FieldsFuncSeq(*c.bufPtr, func(r rune) bool { return r == '\n' || r == '\r' }) {
		line := stringutil.ToString(sub)
		if err := fn(i, line); err != nil {
			return err
		}
		i++
	}
	return nil
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

	var re string
	err := c.processFile(fullPath, func(i int, line string) error {
		re = string(bytes.TrimSpace(stringutil.ToByte(line)))
		return nil
	})

	return re, err
}

type CgoupStat struct {
	NrDescendants      uint64
	NrDyingDescendants uint64
}

func (c *Cgroup) CgoupStat() (CgoupStat, error) {
	fullPath := c.path("cgroup.stat")
	stat := CgoupStat{
		NrDescendants:      math.MaxUint64,
		NrDyingDescendants: math.MaxUint64,
	}

	err := c.processFile(fullPath, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in cgroup.stat: '%s'", fullPath, line)
		}
		var err error
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

	err := c.processFile(fullPath, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in cpu.stat: '%s'", fullPath, line)
		}
		var err error
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
	PageTables             uint64
	SecPageTables          uint64
	PerCPU                 uint64
	Sock                   uint64
	Vmalloc                uint64
	Shmem                  uint64
	Zswap                  uint64
	Zswapped               uint64
	FileMapped             uint64
	FileDirty              uint64
	FileWriteback          uint64
	SwapCached             uint64
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
	Slab                   uint64
	WorkingsetRefaultAnon  uint64
	WorkingsetRefaultFile  uint64
	WorkingsetActivateAnon uint64
	WorkingsetActivateFile uint64
	WorkingsetRestoreAnon  uint64
	WorkingsetRestoreFile  uint64
	WorkingsetNodereclaim  uint64
	Pgscan                 uint64
	Pgsteal                uint64
	PgscanKswapd           uint64
	PgscanDirect           uint64
	PgscanKhugepaged       uint64
	PgstealKswapd          uint64
	PgstealDirect          uint64
	PgstealKhugepaged      uint64
	Pgfault                uint64
	Pgmajfault             uint64
	Pgrefill               uint64
	Pgactivate             uint64
	Pgdeactivate           uint64
	Pglazyfree             uint64
	Pglazyfreed            uint64
	ZswpIn                 uint64
	ZswpOut                uint64
	ZswpWb                 uint64
	ThpFaultAlloc          uint64
	ThpCollapseAlloc       uint64
}

func (c *Cgroup) MemoryStat() (MemoryStat, error) {
	memStat := MemoryStat{
		Anon:                   math.MaxUint64,
		File:                   math.MaxUint64,
		Kernel:                 math.MaxUint64,
		KernelStack:            math.MaxUint64,
		PageTables:             math.MaxUint64,
		SecPageTables:          math.MaxUint64,
		PerCPU:                 math.MaxUint64,
		Sock:                   math.MaxUint64,
		Vmalloc:                math.MaxUint64,
		Shmem:                  math.MaxUint64,
		Zswap:                  math.MaxUint64,
		Zswapped:               math.MaxUint64,
		FileMapped:             math.MaxUint64,
		FileDirty:              math.MaxUint64,
		FileWriteback:          math.MaxUint64,
		SwapCached:             math.MaxUint64,
		AnonThp:                math.MaxUint64,
		FileThp:                math.MaxUint64,
		ShmemThp:               math.MaxUint64,
		InactiveAnon:           math.MaxUint64,
		ActiveAnon:             math.MaxUint64,
		InactiveFile:           math.MaxUint64,
		ActiveFile:             math.MaxUint64,
		Unevictable:            math.MaxUint64,
		SlabReclaimable:        math.MaxUint64,
		SlabUnreclaimable:      math.MaxUint64,
		Slab:                   math.MaxUint64,
		WorkingsetRefaultAnon:  math.MaxUint64,
		WorkingsetRefaultFile:  math.MaxUint64,
		WorkingsetActivateAnon: math.MaxUint64,
		WorkingsetActivateFile: math.MaxUint64,
		WorkingsetRestoreAnon:  math.MaxUint64,
		WorkingsetRestoreFile:  math.MaxUint64,
		WorkingsetNodereclaim:  math.MaxUint64,
		Pgscan:                 math.MaxUint64,
		Pgsteal:                math.MaxUint64,
		PgscanKswapd:           math.MaxUint64,
		PgscanDirect:           math.MaxUint64,
		PgscanKhugepaged:       math.MaxUint64,
		PgstealKswapd:          math.MaxUint64,
		PgstealDirect:          math.MaxUint64,
		PgstealKhugepaged:      math.MaxUint64,
		Pgfault:                math.MaxUint64,
		Pgmajfault:             math.MaxUint64,
		Pgrefill:               math.MaxUint64,
		Pgactivate:             math.MaxUint64,
		Pgdeactivate:           math.MaxUint64,
		Pglazyfree:             math.MaxUint64,
		Pglazyfreed:            math.MaxUint64,
		ZswpIn:                 math.MaxUint64,
		ZswpOut:                math.MaxUint64,
		ZswpWb:                 math.MaxUint64,
		ThpFaultAlloc:          math.MaxUint64,
		ThpCollapseAlloc:       math.MaxUint64,
	}

	fullPath := c.path("memory.stat")

	err := c.processFile(fullPath, func(i int, line string) error {
		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in memory.stat: '%s'", fullPath, line)
		}

		var err error
		switch fields[0] {
		case "anon":
			memStat.Anon, err = strconv.ParseUint(fields[1], 10, 64)
		case "file":
			memStat.File, err = strconv.ParseUint(fields[1], 10, 64)
		case "kernel":
			memStat.Kernel, err = strconv.ParseUint(fields[1], 10, 64)
		case "kernel_stack":
			memStat.KernelStack, err = strconv.ParseUint(fields[1], 10, 64)
		case "pagetables":
			memStat.PageTables, err = strconv.ParseUint(fields[1], 10, 64)
		case "sec_pagetables":
			memStat.SecPageTables, err = strconv.ParseUint(fields[1], 10, 64)
		case "percpu":
			memStat.PerCPU, err = strconv.ParseUint(fields[1], 10, 64)
		case "sock":
			memStat.Sock, err = strconv.ParseUint(fields[1], 10, 64)
		case "vmalloc":
			memStat.Vmalloc, err = strconv.ParseUint(fields[1], 10, 64)
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
		case "swapcached":
			memStat.SwapCached, err = strconv.ParseUint(fields[1], 10, 64)
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
		case "slab":
			memStat.Slab, err = strconv.ParseUint(fields[1], 10, 64)
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
		case "pgscan":
			memStat.Pgscan, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgsteal":
			memStat.Pgsteal, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgscan_kswapd":
			memStat.PgscanKswapd, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgscan_direct":
			memStat.PgscanDirect, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgscan_khugepaged":
			memStat.PgscanKhugepaged, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgsteal_kswapd":
			memStat.PgstealKswapd, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgsteal_direct":
			memStat.PgstealDirect, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgsteal_khugepaged":
			memStat.PgstealKhugepaged, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgfault":
			memStat.Pgfault, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgmajfault":
			memStat.Pgmajfault, err = strconv.ParseUint(fields[1], 10, 64)
		case "pgrefill":
			memStat.Pgrefill, err = strconv.ParseUint(fields[1], 10, 64)
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
		case "zswpwb":
			memStat.ZswpWb, err = strconv.ParseUint(fields[1], 10, 64)
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
	Low           uint64
	High          uint64
	Max           uint64
	Oom           uint64
	OomKill       uint64
	OomGroupKill  uint64
	SockThrottled uint64
}

func (c *Cgroup) MemoryEvents() (MemoryEvents, error) {
	fullPath := c.path("memory.events")

	event := MemoryEvents{
		Low:           math.MaxUint64,
		High:          math.MaxUint64,
		Max:           math.MaxUint64,
		Oom:           math.MaxUint64,
		OomKill:       math.MaxUint64,
		OomGroupKill:  math.MaxUint64,
		SockThrottled: math.MaxUint64,
	}

	c.processFile(fullPath, func(i int, line string) error {

		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 2 {
			return fmt.Errorf("%s: unexpected line in memory.events: '%s'", fullPath, line)
		}

		var err error
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
			event.OomGroupKill, err = strconv.ParseUint(fields[1], 10, 64)
		case "sock_throttled":
			event.SockThrottled, err = strconv.ParseUint(fields[1], 10, 64)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return event, nil
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

	err := c.processFile(fullPath, func(i int, line string) error {

		var fields [7]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 7 {
			return fmt.Errorf("%s: unexpected line in io.stat: '%s'", fullPath, line)
		}
		stat := IOStat{
			Rbytes: math.MaxUint64,
			Wbytes: math.MaxUint64,
			Rios:   math.MaxUint64,
			Wios:   math.MaxUint64,
			Dbytes: math.MaxUint64,
			Dios:   math.MaxUint64,
		}

		var err error
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
			Avg10:  math.MaxFloat64,
			Avg60:  math.MaxFloat64,
			Avg300: math.MaxFloat64,
			Total:  math.MaxUint64,
		},
		Full: PSIData{
			Avg10:  math.MaxFloat64,
			Avg60:  math.MaxFloat64,
			Avg300: math.MaxFloat64,
			Total:  math.MaxUint64,
		},
	}

	err := c.processFile(fullPath, func(i int, line string) error {

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

		var err error
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
	MemoryCurrent                uint64
	MemoryLow                    uint64
	MemoryHigh                   uint64
	MemoryMin                    uint64
	MemoryMax                    uint64
	MemoryOOMGroup               uint64
	MemorySwapCurrent            uint64
	MemorySwapMax                uint64
	MemoryZSwapCurrent           uint64
	MemoryZSwapMax               uint64
	CpuWeight                    uint64
	CpuMax                       string
	CpuSetCpus                   string
	CpuSetCpusEffective          string
	CpuSetCpusExclusive          string
	CpuSetCpusExclusiveEffective string
	TidsCurrent                  uint64
	TidsMax                      uint64
}

const MaxCgroupPropertyUintValue = math.MaxUint64 - 1
const NoExistCgroupPropertyStrValue = "no-exist"

func (c *Cgroup) getUintFromPropertyFile(file string) (uint64, error) {

	var re uint64
	err := c.processFile(file, func(i int, line string) error {
		s := stringutil.ToString(bytes.TrimSpace(stringutil.ToByte(line)))
		if s == "max" {
			re = MaxCgroupPropertyUintValue
			return nil
		}
		re, _ = strconv.ParseUint(line, 10, 64)
		return nil
	})
	if err != nil {
		return math.MaxUint64, nil
	}
	return re, nil
}

func (c *Cgroup) getStrFromPropertyFile(file string) (string, error) {

	re := ""
	err := c.processFile(file, func(i int, line string) error {
		re = string(bytes.TrimSpace(stringutil.ToByte(line)))
		return nil
	})
	if err != nil {
		return NoExistCgroupPropertyStrValue, nil
	}

	return re, nil
}

func (c *Cgroup) Properties() (Property, error) {

	p := Property{}
	var err error
	if p.MemoryCurrent, err = c.getUintFromPropertyFile(c.path("memory.current")); err != nil {
		return p, err
	}
	if p.MemoryLow, err = c.getUintFromPropertyFile(c.path("memory.low")); err != nil {
		return p, err
	}
	if p.MemoryHigh, err = c.getUintFromPropertyFile(c.path("memory.high")); err != nil {
		return p, err
	}
	if p.MemoryMin, err = c.getUintFromPropertyFile(c.path("memory.min")); err != nil {
		return p, err
	}
	if p.MemoryMax, err = c.getUintFromPropertyFile(c.path("memory.max")); err != nil {
		return p, err
	}

	if p.MemoryOOMGroup, err = c.getUintFromPropertyFile(c.path("memory.oom.group")); err != nil {
		return p, err
	}

	if p.MemorySwapCurrent, err = c.getUintFromPropertyFile(c.path("memory.swap.current")); err != nil {
		return p, err
	}

	if p.MemorySwapMax, err = c.getUintFromPropertyFile(c.path("memory.swap.max")); err != nil {
		return p, err
	}

	if p.MemoryZSwapCurrent, err = c.getUintFromPropertyFile(c.path("memory.zswap.current")); err != nil {
		return p, err
	}

	if p.MemoryZSwapMax, err = c.getUintFromPropertyFile(c.path("memory.zswap.max")); err != nil {
		return p, err
	}

	if p.CpuMax, err = c.getStrFromPropertyFile(c.path("cpu.max")); err != nil {
		return p, err
	}
	if p.CpuWeight, err = c.getUintFromPropertyFile(c.path("cpu.weight")); err != nil {
		return p, err
	}

	if p.CpuSetCpus, err = c.getStrFromPropertyFile(c.path("cpuset.cpus")); err != nil {
		return p, err
	}
	if p.CpuSetCpusEffective, err = c.getStrFromPropertyFile(c.path("cpuset.cpus.effective")); err != nil {
		return p, err
	}

	if p.CpuSetCpusExclusive, err = c.getStrFromPropertyFile(c.path("cpuset.cpus.exclusive")); err != nil {
		return p, err
	}
	if p.CpuSetCpusExclusiveEffective, err = c.getStrFromPropertyFile(c.path("cpuset.cpus.exclusive.effective")); err != nil {
		return p, err
	}

	if p.TidsCurrent, err = c.getUintFromPropertyFile(c.path("pids.current")); err != nil {
		return p, err
	}
	if p.TidsMax, err = c.getUintFromPropertyFile(c.path("pids.max")); err != nil {
		return p, err
	}

	return p, nil
}
