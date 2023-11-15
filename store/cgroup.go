package store

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	CgroupV2MountPoint  = "/sys/fs/cgroup"
	ErrInvalidFormat    = errors.New("cgroups: parsing file with invalid format failed")
	ErrInvalidGroupPath = errors.New("cgroups: invalid group path")
)

type CgroupFile int

const (
	CpuStatFile CgroupFile = iota
	MemoryStatFile
	CpuSetCpusFile
	CpuSetCpusEffectiveFile
	CpuSetMemsFileFile
	CpuSetMemsEffectiveFile
	CpuWeightFile
	CpuMaxFile
	MemoryCurrentFile
	MemoryLowFile
	MemoryHighFile
	MemoryMinFile
	MemoryMaxFile
	MemoryPeakFile
	MemorySwapCurrentFile
	MemorySwapMaxFile
	MemoryZswapCurrentFile
	MemoryZswapMaxFile
	MemoryEventsFile
	IoStatFile
	CpuPressureFile
	MemoryPressureFile
	IoPressureFile
)

var isCgroup2 = sync.OnceValue[bool](func() bool {
	var st unix.Statfs_t
	if err := unix.Statfs(CgroupV2MountPoint, &st); err != nil {
		return false
	}
	if st.Type == unix.CGROUP2_SUPER_MAGIC {
		return true
	}
	return false
})

type CgroupSample struct {
	Path        string
	Name        string
	Level       int
	Inode       uint64
	Child       map[string]*CgroupSample
	IsNotExist  map[CgroupFile]struct{}
	Controllers string
	CgoupStat
	CPUStat
	MemoryStat
	CpuSetCpus          string
	CpuSetCpusEffective string
	CpuSetMems          string
	CpuSetMemsEffective string
	CpuWeight           uint64
	CpuMax              string
	MemoryCurrent       uint64
	MemoryLow           uint64
	MemoryHigh          uint64
	MemoryMin           uint64
	MemoryMax           uint64
	MemoryPeak          uint64
	SwapCurrent         uint64
	SwapMax             uint64
	ZswapCurrent        uint64
	ZswapMax            uint64
	MemoryEvents
	IOStats        []IOStat
	CpuPressure    PSIStats
	MemoryPressure PSIStats
	IOPressure     PSIStats
}

type CgoupStat struct {
	NrDescendants      uint64
	NrDyingDescendants uint64
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

type MemoryStat struct {
	Anon                  uint64
	File                  uint64
	KernelStack           uint64
	Slab                  uint64
	Sock                  uint64
	Shmem                 uint64
	Zswap                 uint64
	Zswapped              uint64
	FileMapped            uint64
	FileDirty             uint64
	FileWriteback         uint64
	AnonThp               uint64
	InactiveAnon          uint64
	ActiveAnon            uint64
	InactiveFile          uint64
	ActiveFile            uint64
	Unevictable           uint64
	SlabReclaimable       uint64
	SlabUnreclaimable     uint64
	Pgfault               uint64
	Pgmajfault            uint64
	WorkingsetRefault     uint64
	WorkingsetActivate    uint64
	WorkingsetNodereclaim uint64
	Pgrefill              uint64
	Pgscan                uint64
	Pgsteal               uint64
	Pgactivate            uint64
	Pgdeactivate          uint64
	Pglazyfree            uint64
	Pglazyfreed           uint64
	ZswpIn                uint64
	ZswpOut               uint64
	ThpFaultAlloc         uint64
	ThpCollapseAlloc      uint64
}

type MemoryEvents struct {
	Low     uint64
	High    uint64
	Max     uint64
	Oom     uint64
	OomKill uint64
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

func walkCgroupNode(cSample *CgroupSample, log *slog.Logger) error {
	fullPath := filepath.Join(CgroupV2MountPoint, cSample.Path, cSample.Name)
	es, err := os.ReadDir(fullPath)
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		log.Warn(msg)
		return err
	}
	for _, e := range es {
		if e.IsDir() {
			if info, err := e.Info(); err == nil {
				t := info.Sys().(*syscall.Stat_t)
				new := &CgroupSample{
					Path:       filepath.Join(cSample.Path, cSample.Name),
					Name:       e.Name(),
					Level:      cSample.Level + 1,
					Inode:      t.Ino,
					Child:      make(map[string]*CgroupSample),
					IsNotExist: make(map[CgroupFile]struct{}),
				}
				if err := new.collectStat(); err != nil {
					msg := fmt.Sprintf("%s", err)
					log.Warn(msg)
					continue
				}
				cSample.Child[new.Name] = new
				if err := walkCgroupNode(new, log); err != nil {
					msg := fmt.Sprintf("%s", err)
					log.Warn(msg)
					continue
				}
			}
		}
	}
	return nil
}

func (node *CgroupSample) collectStat() error {
	fullPath := filepath.Join(CgroupV2MountPoint, node.Path, node.Name)

	b, err := os.ReadFile(filepath.Join(fullPath, "cgroup.controllers"))
	if err != nil {
		return err
	}
	node.Controllers = string(b)

	out := make(map[string]uint64, 50)
	if err := readKVStatsFile(fullPath, "cgroup.stat", out); err == nil {
		node.CgoupStat = CgoupStat{
			NrDescendants:      out["nr_descendants"],
			NrDyingDescendants: out["nr_dying_descendants"],
		}
	}
	if err := readKVStatsFile(fullPath, "cpu.stat", out); err == nil {
		node.CPUStat = CPUStat{
			UsageUsec:     out["usage_usec"],
			UserUsec:      out["user_usec"],
			SystemUsec:    out["system_usec"],
			NrPeriods:     out["nr_periods"],
			NrThrottled:   out["nr_throttled"],
			ThrottledUsec: out["throttled_usec"],
			NrBursts:      out["nr_bursts"],
			BurstUsec:     out["burst_usec"],
		}
	}

	if err := readKVStatsFile(fullPath, "memory.stat", out); err == nil {
		node.MemoryStat = MemoryStat{
			Anon:                  out["anon"],
			File:                  out["file"],
			KernelStack:           out["kernel_stack"],
			Slab:                  out["slab"],
			Sock:                  out["sock"],
			Shmem:                 out["shmem"],
			Zswap:                 out["zswap"],
			Zswapped:              out["zswapped"],
			FileMapped:            out["file_mapped"],
			FileDirty:             out["file_dirty"],
			FileWriteback:         out["file_writeback"],
			AnonThp:               out["anon_thp"],
			InactiveAnon:          out["inactive_anon"],
			ActiveAnon:            out["active_anon"],
			InactiveFile:          out["inactive_file"],
			ActiveFile:            out["active_file"],
			Unevictable:           out["unevictable"],
			SlabReclaimable:       out["slab_reclaimable"],
			SlabUnreclaimable:     out["slab_unreclaimable"],
			Pgfault:               out["pgfault"],
			Pgmajfault:            out["pgmajfault"],
			WorkingsetRefault:     out["workingset_refault"],
			WorkingsetActivate:    out["workingset_activate"],
			WorkingsetNodereclaim: out["workingset_nodereclaim"],
			Pgrefill:              out["pgrefill"],
			Pgscan:                out["pgscan"],
			Pgsteal:               out["pgsteal"],
			Pgactivate:            out["pgactivate"],
			Pgdeactivate:          out["pgdeactivate"],
			Pglazyfree:            out["pglazyfree"],
			Pglazyfreed:           out["pglazyfreed"],
			ZswpIn:                out["zswpin"],
			ZswpOut:               out["zswpout"],
			ThpFaultAlloc:         out["thp_fault_alloc"],
			ThpCollapseAlloc:      out["thp_collapse_alloc"],
		}
	} else {
		node.IsNotExist[MemoryStatFile] = struct{}{}
	}

	if node.CpuSetCpus, err = getStatFileContentString(filepath.Join(fullPath, "cpuset.cpus")); err != nil {
		node.IsNotExist[CpuSetCpusFile] = struct{}{}
	}
	if node.CpuSetCpusEffective, err = getStatFileContentString(filepath.Join(fullPath, "cpuset.cpus.effective")); err != nil {
		node.IsNotExist[CpuSetCpusEffectiveFile] = struct{}{}
	}
	if node.CpuSetMems, err = getStatFileContentString(filepath.Join(fullPath, "cpuset.mems")); err != nil {
		node.IsNotExist[CpuSetMemsFileFile] = struct{}{}
	}
	if node.CpuSetMemsEffective, err = getStatFileContentString(filepath.Join(fullPath, "cpuset.mems.effective")); err != nil {
		node.IsNotExist[CpuSetMemsEffectiveFile] = struct{}{}
	}

	if node.CpuWeight, err = getStatFileContentUint64(filepath.Join(fullPath, "cpu.weight")); err != nil {
		node.IsNotExist[CpuWeightFile] = struct{}{}
	}
	if node.CpuMax, err = getStatFileContentString(filepath.Join(fullPath, "cpu.max")); err != nil {
		node.IsNotExist[CpuMaxFile] = struct{}{}
	}

	if node.MemoryCurrent, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.current")); err != nil {
		node.IsNotExist[MemoryCurrentFile] = struct{}{}
	}
	if node.MemoryLow, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.low")); err != nil {
		node.IsNotExist[MemoryLowFile] = struct{}{}
	}
	if node.MemoryHigh, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.high")); err != nil {
		node.IsNotExist[MemoryHighFile] = struct{}{}
	}
	if node.MemoryMin, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.min")); err != nil {
		node.IsNotExist[MemoryMinFile] = struct{}{}
	}
	if node.MemoryMax, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.max")); err != nil {
		node.IsNotExist[MemoryMaxFile] = struct{}{}
	}
	if node.MemoryPeak, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.peak")); err != nil {
		node.IsNotExist[MemoryPeakFile] = struct{}{}
	}
	if node.SwapCurrent, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.swap.current")); err != nil {
		node.IsNotExist[MemorySwapCurrentFile] = struct{}{}
	}
	if node.SwapMax, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.swap.max")); err != nil {
		node.IsNotExist[MemorySwapMaxFile] = struct{}{}
	}
	if node.ZswapCurrent, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.zswap.current")); err != nil {
		node.IsNotExist[MemoryZswapCurrentFile] = struct{}{}
	}
	if node.ZswapCurrent, err = getStatFileContentUint64(filepath.Join(fullPath, "memory.zswap.max")); err != nil {
		node.IsNotExist[MemoryZswapMaxFile] = struct{}{}
	}

	memoryEvents := make(map[string]uint64)
	if err := readKVStatsFile(fullPath, "memory.events", memoryEvents); err == nil {
		node.MemoryEvents = MemoryEvents{
			Low:     memoryEvents["low"],
			High:    memoryEvents["high"],
			Max:     memoryEvents["max"],
			Oom:     memoryEvents["oom"],
			OomKill: memoryEvents["oom_kill"],
		}
	} else {
		node.IsNotExist[MemoryEventsFile] = struct{}{}
	}
	if node.IOStats, err = readIoStats(fullPath); err != nil {
		node.IsNotExist[IoStatFile] = struct{}{}
	}

	if node.CpuPressure, err = getStatPSIFromFile(filepath.Join(fullPath, "cpu.pressure")); err != nil {
		node.IsNotExist[CpuPressureFile] = struct{}{}
	}
	if node.MemoryPressure, err = getStatPSIFromFile(filepath.Join(fullPath, "memory.pressure")); err != nil {
		node.IsNotExist[MemoryPressureFile] = struct{}{}
	}
	if node.IOPressure, err = getStatPSIFromFile(filepath.Join(fullPath, "io.pressure")); err != nil {
		node.IsNotExist[IoPressureFile] = struct{}{}
	}

	return nil
}

func readKVStatsFile(path string, file string, out map[string]uint64) error {
	f, err := os.Open(filepath.Join(path, file))
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		name, value, err := parseKV(s.Text())
		if err != nil {
			return fmt.Errorf("error while parsing %s (line=%q): %w", filepath.Join(path, file), s.Text(), err)
		}
		out[name] = value
	}
	return s.Err()
}
func parseKV(raw string) (string, uint64, error) {
	parts := strings.Fields(raw)
	if len(parts) != 2 {
		return "", 0, ErrInvalidFormat
	}
	v, err := parseUint(parts[1], 10, 64)
	return parts[0], v, err
}

func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

func getStatFileContentString(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err == nil {
		b := bytes.TrimSpace(b)
		return string(b), nil
	}
	return "", err
}

func getStatFileContentUint64(filePath string) (uint64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// We expect an unsigned 64 bit integer, or a "max" string
	// in some cases.
	buf := make([]byte, 32)
	n, err := f.Read(buf)
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(string(buf[:n]))
	if trimmed == "max" {
		return math.MaxUint64, nil
	}

	res, err := parseUint(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func readIoStats(path string) ([]IOStat, error) {
	// more details on the io.stat file format: https://www.kernel.org/doc/Documentation/cgroup-v2.txt
	var usage []IOStat
	fpath := filepath.Join(path, "io.stat")
	currentData, err := os.ReadFile(fpath)
	if err != nil {
		return usage, err
	}
	entries := strings.Split(string(currentData), "\n")

	for _, entry := range entries {
		parts := strings.Split(entry, " ")
		if len(parts) < 2 {
			continue
		}
		majmin := strings.Split(parts[0], ":")
		if len(majmin) != 2 {
			continue
		}
		major, err := strconv.ParseUint(majmin[0], 10, 0)
		if err != nil {
			return usage, err
		}
		minor, err := strconv.ParseUint(majmin[1], 10, 0)
		if err != nil {
			return usage, err
		}
		parts = parts[1:]
		ioStat := IOStat{
			Major: major,
			Minor: minor,
		}
		for _, s := range parts {
			keyPairValue := strings.Split(s, "=")
			if len(keyPairValue) != 2 {
				continue
			}
			v, err := strconv.ParseUint(keyPairValue[1], 10, 0)
			if err != nil {
				continue
			}
			switch keyPairValue[0] {
			case "rbytes":
				ioStat.Rbytes = v
			case "wbytes":
				ioStat.Wbytes = v
			case "rios":
				ioStat.Rios = v
			case "wios":
				ioStat.Wios = v
			case "dbytes":
				ioStat.Dbytes = v
			case "dios":
				ioStat.Dios = v
			}
		}
		usage = append(usage, ioStat)
	}
	return usage, nil
}

func getStatPSIFromFile(path string) (PSIStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return PSIStats{}, err
	}
	defer f.Close()

	psiStats := PSIStats{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Fields(sc.Text())
		var pv *PSIData
		switch parts[0] {
		case "some":
			psiStats.Some = PSIData{}
			pv = &psiStats.Some
		case "full":
			psiStats.Full = PSIData{}
			pv = &psiStats.Full
		}
		if pv != nil {
			err = parsePSIData(parts[1:], pv)
			if err != nil {
				return PSIStats{}, err
			}
		}
	}

	if err := sc.Err(); err != nil {
		return PSIStats{}, err
	}
	return psiStats, nil
}

func parsePSIData(psi []string, data *PSIData) error {
	for _, f := range psi {
		kv := strings.SplitN(f, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid PSI data: %q", f)
		}
		var pv *float64
		switch kv[0] {
		case "avg10":
			pv = &data.Avg10
		case "avg60":
			pv = &data.Avg60
		case "avg300":
			pv = &data.Avg300
		case "total":
			v, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid %s PSI value: %w", kv[0], err)
			}
			data.Total = v
		}
		if pv != nil {
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return fmt.Errorf("invalid %s PSI value: %w", kv[0], err)
			}
			*pv = v
		}
	}
	return nil
}
