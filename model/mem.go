package model

import (
	"fmt"

	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/util"
)

type MEM struct {
	MemTotal          uint64
	MemFree           uint64
	MemAvailable      uint64
	Buffers           uint64
	Cached            uint64
	SwapCached        uint64
	Active            uint64
	Inactive          uint64
	ActiveAnon        uint64
	InactiveAnon      uint64
	ActiveFile        uint64
	InactiveFile      uint64
	Unevictable       uint64
	Mlocked           uint64
	SwapTotal         uint64
	SwapFree          uint64
	Dirty             uint64
	Writeback         uint64
	AnonPages         uint64
	Mapped            uint64
	Shmem             uint64
	Slab              uint64
	SReclaimable      uint64
	SUnreclaim        uint64
	KernelStack       uint64
	PageTables        uint64
	NFSUnstable       uint64
	Bounce            uint64
	WritebackTmp      uint64
	CommitLimit       uint64
	CommittedAS       uint64
	VmallocTotal      uint64
	VmallocUsed       uint64
	VmallocChunk      uint64
	HardwareCorrupted uint64
	AnonHugePages     uint64
	ShmemHugePages    uint64
	ShmemPmdMapped    uint64
	CmaTotal          uint64
	CmaFree           uint64
	HugePagesTotal    uint64
	HugePagesFree     uint64
	HugePagesRsvd     uint64
	HugePagesSurp     uint64
	Hugepagesize      uint64
	DirectMap4k       uint64
	DirectMap2M       uint64
	DirectMap1G       uint64
}

func (m *MEM) GetRenderValue(field string) string {
	switch field {
	case "Total":
		return fmt.Sprintf("%s", util.GetHumanSize(m.MemTotal*1024))
	case "Free":
		return fmt.Sprintf("%s", util.GetHumanSize(m.MemFree*1024))
	case "Avail":
		return fmt.Sprintf("%s", util.GetHumanSize(m.MemAvailable*1024))
	case "Buffer":
		return fmt.Sprintf("%s", util.GetHumanSize(m.Buffers*1024))
	case "Cache":
		return fmt.Sprintf("%s", util.GetHumanSize(m.Cached*1024))
	case "MemTotal":
		return fmt.Sprintf("%d KB", m.MemTotal)
	case "MemFree":
		return fmt.Sprintf("%d KB", m.MemFree)
	case "MemAvailable":
		return fmt.Sprintf("%d KB", m.MemAvailable)
	case "Buffers":
		return fmt.Sprintf("%d KB", m.Buffers)
	case "Cached":
		return fmt.Sprintf("%d KB", m.Cached)
	case "SwapCached":
		return fmt.Sprintf("%d KB", m.SwapCached)
	case "Active":
		return fmt.Sprintf("%d KB", m.Active)
	case "Inactive":
		return fmt.Sprintf("%d KB", m.Inactive)
	case "ActiveAnon":
		return fmt.Sprintf("%d KB", m.ActiveAnon)
	case "InactiveAnon":
		return fmt.Sprintf("%d KB", m.InactiveAnon)
	case "ActiveFile":
		return fmt.Sprintf("%d KB", m.ActiveFile)
	case "InactiveFile":
		return fmt.Sprintf("%d KB", m.InactiveFile)
	case "Unevictable":
		return fmt.Sprintf("%d KB", m.Unevictable)
	case "Mlocked":
		return fmt.Sprintf("%d KB", m.Mlocked)
	case "SwapTotal":
		return fmt.Sprintf("%d KB", m.SwapTotal)
	case "SwapFree":
		return fmt.Sprintf("%d KB", m.SwapFree)
	case "Dirty":
		return fmt.Sprintf("%d KB", m.Dirty)
	case "Writeback":
		return fmt.Sprintf("%d KB", m.Writeback)
	case "AnonPages":
		return fmt.Sprintf("%d KB", m.AnonPages)
	case "Mapped":
		return fmt.Sprintf("%d KB", m.Mapped)
	case "Shmem":
		return fmt.Sprintf("%d KB", m.Shmem)
	case "Slab":
		return fmt.Sprintf("%d KB", m.Slab)
	case "SReclaimable":
		return fmt.Sprintf("%d KB", m.SReclaimable)
	case "SUnreclaim":
		return fmt.Sprintf("%d KB", m.SUnreclaim)
	case "KernelStack":
		return fmt.Sprintf("%d KB", m.KernelStack)
	case "PageTables":
		return fmt.Sprintf("%d KB", m.PageTables)
	case "NFSUnstable":
		return fmt.Sprintf("%d KB", m.NFSUnstable)
	case "Bounce":
		return fmt.Sprintf("%d KB", m.Bounce)
	case "WritebackTmp":
		return fmt.Sprintf("%d KB", m.WritebackTmp)
	case "CommitLimit":
		return fmt.Sprintf("%d KB", m.CommitLimit)
	case "CommittedAS":
		return fmt.Sprintf("%d KB", m.CommittedAS)
	case "VmallocTotal":
		return fmt.Sprintf("%d KB", m.VmallocTotal)
	case "VmallocUsed":
		return fmt.Sprintf("%d KB", m.VmallocUsed)
	case "VmallocChunk":
		return fmt.Sprintf("%d KB", m.VmallocChunk)
	case "HardwareCorrupted":
		return fmt.Sprintf("%d KB", m.HardwareCorrupted)
	case "AnonHugePages":
		return fmt.Sprintf("%d KB", m.AnonHugePages)
	case "ShmemHugePages":
		return fmt.Sprintf("%d KB", m.ShmemHugePages)
	case "ShmemPmdMapped":
		return fmt.Sprintf("%d KB", m.ShmemPmdMapped)
	case "CmaTotal":
		return fmt.Sprintf("%d KB", m.CmaTotal)
	case "CmaFree":
		return fmt.Sprintf("%d KB", m.CmaFree)
	case "HugePagesTotal":
		return fmt.Sprintf("%d", m.HugePagesTotal)
	case "HugePagesFree":
		return fmt.Sprintf("%d", m.HugePagesFree)
	case "HugePagesRsvd":
		return fmt.Sprintf("%d", m.HugePagesRsvd)
	case "HugePagesSurp":
		return fmt.Sprintf("%d", m.HugePagesSurp)
	case "Hugepagesize":
		return fmt.Sprintf("%d KB", m.Hugepagesize)
	case "DirectMap4k":
		return fmt.Sprintf("%d KB", m.DirectMap4k)
	case "DirectMap2M":
		return fmt.Sprintf("%d KB", m.DirectMap2M)
	case "DirectMap1G":
		return fmt.Sprintf("%d KB", m.DirectMap1G)
	}
	return ""
}

func (m *MEM) Collect(prev, curr *store.Sample) {

	*m = MEM{
		MemTotal:          getValueOrDefault(curr.MemTotal),
		MemFree:           getValueOrDefault(curr.MemFree),
		MemAvailable:      getValueOrDefault(curr.MemAvailable),
		Buffers:           getValueOrDefault(curr.Buffers),
		Cached:            getValueOrDefault(curr.Cached),
		SwapCached:        getValueOrDefault(curr.SwapCached),
		Active:            getValueOrDefault(curr.Active),
		Inactive:          getValueOrDefault(curr.Inactive),
		ActiveAnon:        getValueOrDefault(curr.ActiveAnon),
		InactiveAnon:      getValueOrDefault(curr.InactiveAnon),
		ActiveFile:        getValueOrDefault(curr.ActiveFile),
		InactiveFile:      getValueOrDefault(curr.InactiveFile),
		Unevictable:       getValueOrDefault(curr.Unevictable),
		Mlocked:           getValueOrDefault(curr.Mlocked),
		SwapTotal:         getValueOrDefault(curr.SwapTotal),
		SwapFree:          getValueOrDefault(curr.SwapFree),
		Dirty:             getValueOrDefault(curr.Dirty),
		Writeback:         getValueOrDefault(curr.Writeback),
		AnonPages:         getValueOrDefault(curr.AnonPages),
		Mapped:            getValueOrDefault(curr.Mapped),
		Shmem:             getValueOrDefault(curr.Shmem),
		Slab:              getValueOrDefault(curr.Slab),
		SReclaimable:      getValueOrDefault(curr.SReclaimable),
		SUnreclaim:        getValueOrDefault(curr.SUnreclaim),
		KernelStack:       getValueOrDefault(curr.KernelStack),
		PageTables:        getValueOrDefault(curr.PageTables),
		NFSUnstable:       getValueOrDefault(curr.NFSUnstable),
		Bounce:            getValueOrDefault(curr.Bounce),
		WritebackTmp:      getValueOrDefault(curr.WritebackTmp),
		CommitLimit:       getValueOrDefault(curr.CommitLimit),
		CommittedAS:       getValueOrDefault(curr.CommittedAS),
		VmallocTotal:      getValueOrDefault(curr.VmallocTotal),
		VmallocUsed:       getValueOrDefault(curr.VmallocUsed),
		VmallocChunk:      getValueOrDefault(curr.VmallocChunk),
		HardwareCorrupted: getValueOrDefault(curr.HardwareCorrupted),
		AnonHugePages:     getValueOrDefault(curr.AnonHugePages),
		ShmemHugePages:    getValueOrDefault(curr.ShmemHugePages),
		ShmemPmdMapped:    getValueOrDefault(curr.ShmemPmdMapped),
		CmaTotal:          getValueOrDefault(curr.CmaTotal),
		CmaFree:           getValueOrDefault(curr.CmaFree),
		HugePagesTotal:    getValueOrDefault(curr.HugePagesTotal),
		HugePagesFree:     getValueOrDefault(curr.HugePagesFree),
		HugePagesRsvd:     getValueOrDefault(curr.HugePagesRsvd),
		HugePagesSurp:     getValueOrDefault(curr.HugePagesSurp),
		Hugepagesize:      getValueOrDefault(curr.Hugepagesize),
		DirectMap4k:       getValueOrDefault(curr.DirectMap4k),
		DirectMap2M:       getValueOrDefault(curr.DirectMap2M),
	}

}

func getValueOrDefault(m *uint64) uint64 {
	if m == nil {
		return 0
	}
	return *m
}
