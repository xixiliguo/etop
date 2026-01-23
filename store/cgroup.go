package store

import (
	"errors"
	"os"
	"sync"

	"github.com/xixiliguo/etop/cgroupfs"
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
	FullPath    string
	Name        string
	Level       int
	Inode       uint64
	Child       map[string]CgroupSample
	Controllers string
	cgroupfs.CgoupStat
	cgroupfs.CPUStat
	cgroupfs.MemoryStat
	cgroupfs.Property
	cgroupfs.MemoryEvents
	IOStats        []cgroupfs.IOStat
	CpuPressure    cgroupfs.PSIStats
	MemoryPressure cgroupfs.PSIStats
	IOPressure     cgroupfs.PSIStats
}

func walkCgroupNode(level int, cg cgroupfs.Cgroup) (CgroupSample, error) {
	root := CgroupSample{
		FullPath: cg.FullPath,
		Name:     cg.Name,
		Level:    level,
		Child:    map[string]CgroupSample{},
	}

	root.Inode, _ = cg.Inode()
	root.Controllers, _ = cg.Controllers()
	root.CgoupStat, _ = cg.CgoupStat()
	root.CPUStat, _ = cg.CPUStat()
	root.MemoryStat, _ = cg.MemoryStat()
	root.Property, _ = cg.Properties()
	root.MemoryEvents, _ = cg.MemoryEvents()
	root.IOStats, _ = cg.IOStats()
	root.CpuPressure, _ = cg.PSIStats("cpu.pressure")
	root.MemoryPressure, _ = cg.PSIStats("memory.pressure")
	root.IOPressure, _ = cg.PSIStats("io.pressure")

	es, err := os.ReadDir(CgroupV2MountPoint + cg.FullPath)
	if err != nil {
		return root, err
	}

	for _, e := range es {
		if e.IsDir() {
			child := cg.Child(e.Name())
			childSample, err := walkCgroupNode(level+1, child)
			if err != nil {
				continue
			}
			root.Child[child.Name] = childSample

		}
	}
	return root, nil
}
