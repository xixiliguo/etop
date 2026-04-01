package store

import (
	"errors"
	"math"
	"sync"

	"github.com/xixiliguo/etop/cgroupfs"
	"github.com/xixiliguo/etop/internal/fileutil"
	"golang.org/x/sys/unix"
)

var (
	CgroupV2MountPoint  = "/sys/fs/cgroup"
	ErrInvalidFormat    = errors.New("cgroups: parsing file with invalid format failed")
	ErrInvalidGroupPath = errors.New("cgroups: invalid group path")
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
	RxPacket       uint64
	RxByte         uint64
	TxPacket       uint64
	TxByte         uint64
}

func walkCgroupNode(level int, cg cgroupfs.Cgroup, c *CgroupNetStat) (CgroupSample, error) {
	root := CgroupSample{
		FullPath: cg.FullPath,
		Name:     cg.Name,
		Level:    level,
		Child:    map[string]CgroupSample{},
	}

	var err error

	if root.Inode, err = cg.Inode(); err != nil {
		return root, err
	}
	if root.Controllers, err = cg.Controllers(); err != nil {
		return root, err
	}
	if root.CgoupStat, err = cg.CgoupStat(); err != nil {
		return root, err
	}
	if root.CPUStat, err = cg.CPUStat(); err != nil {
		return root, err
	}
	if root.MemoryStat, err = cg.MemoryStat(); err != nil {
		return root, err
	}
	if root.Property, err = cg.Properties(); err != nil {
		return root, err
	}
	if root.MemoryEvents, err = cg.MemoryEvents(); err != nil {
		return root, err
	}
	if root.IOStats, err = cg.IOStats(); err != nil {
		return root, err
	}
	if root.CpuPressure, err = cg.PSIStats("cpu.pressure"); err != nil {
		return root, err
	}
	if root.MemoryPressure, err = cg.PSIStats("memory.pressure"); err != nil {
		return root, err
	}
	if root.IOPressure, err = cg.PSIStats("io.pressure"); err != nil {
		return root, err
	}

	root.RxPacket = math.MaxUint64
	root.RxByte = math.MaxUint64
	root.TxPacket = math.MaxUint64
	root.TxByte = math.MaxUint64

	if c != nil && c.Stats != nil {
		if s, err := c.NetStat(root.Inode); err == nil {
			root.RxPacket = s.RxPacket
			root.RxByte = s.RxByte
			root.TxPacket = s.TxPacket
			root.TxByte = s.TxByte

		}
	}

	fileutil.SubDirWalk(CgroupV2MountPoint+cg.FullPath, func(subDir string) error {
		child := cg.Child(subDir)
		childSample, err := walkCgroupNode(level+1, child, c)
		if err != nil {
			return nil
		}
		root.Child[child.Name] = childSample
		return nil
	})

	return root, nil
}
