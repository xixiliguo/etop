package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultMEMFields = []string{
	"MemTotal", "MemFree", "MemAvailable",
	"Buffers", "Cached", "SwapCached", "Active",
	"Inactive", "ActiveAnon", "InactiveAnon", "Unevictable",
	"Mlocked", "SwapTotal", "SwapFree", "Dirty",
	"Writeback", "AnonPages", "Mapped", "Shmem",
	"Slab", "SReclaimable", "SUnreclaim", "KernelStack",
	"PageTables", "NFSUnstable", "Bounce", "WritebackTmp",
	"CommitLimit", "CommittedAS", "VmallocTotal", "VmallocUsed",
	"VmallocChunk", "HardwareCorrupted", "AnonHugePages", "ShmemHugePages",
	"ShmemPmdMapped", "CmaTotal", "CmaFree", "HugePagesTotal",
	"HugePagesFree", "HugePagesRsvd", "HugePagesSurp", "Hugepagesize",
	"DirectMap4k", "DirectMap2M", "DirectMap1G",
}

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

func (m *MEM) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Total":
		cfg = Field{"Total", HumanReadableSize, 0, "", 10, false}
	case "Free":
		cfg = Field{"Free", HumanReadableSize, 0, "", 10, false}
	case "Avail":
		cfg = Field{"Avail", HumanReadableSize, 0, "", 10, false}
	case "HSlab":
		cfg = Field{"Slab", HumanReadableSize, 0, "", 10, false}
	case "Buffer":
		cfg = Field{"Buffer", HumanReadableSize, 0, "", 10, false}
	case "Cache":
		cfg = Field{"Cache", HumanReadableSize, 0, "", 10, false}
	case "MemTotal":
		cfg = Field{"MemTotal", Raw, 0, " KB", 10, false}
	case "MemFree":
		cfg = Field{"MemFree", Raw, 0, " KB", 10, false}
	case "MemAvailable":
		cfg = Field{"MemAvailable", Raw, 0, " KB", 10, false}
	case "Buffers":
		cfg = Field{"Buffers", Raw, 0, " KB", 10, false}
	case "Cached":
		cfg = Field{"Cached", Raw, 0, " KB", 10, false}
	case "SwapCached":
		cfg = Field{"SwapCached", Raw, 0, " KB", 10, false}
	case "Active":
		cfg = Field{"Active", Raw, 0, " KB", 10, false}
	case "Inactive":
		cfg = Field{"Inactive", Raw, 0, " KB", 10, false}
	case "ActiveAnon":
		cfg = Field{"ActiveAnon", Raw, 0, " KB", 10, false}
	case "InactiveAnon":
		cfg = Field{"InactiveAnon", Raw, 0, " KB", 10, false}
	case "ActiveFile":
		cfg = Field{"ActiveFile", Raw, 0, " KB", 10, false}
	case "InactiveFile":
		cfg = Field{"InactiveFile", Raw, 0, " KB", 10, false}
	case "Unevictable":
		cfg = Field{"Unevictable", Raw, 0, " KB", 10, false}
	case "Mlocked":
		cfg = Field{"Mlocked", Raw, 0, " KB", 10, false}
	case "SwapTotal":
		cfg = Field{"SwapTotal", Raw, 0, " KB", 10, false}
	case "SwapFree":
		cfg = Field{"SwapFree", Raw, 0, " KB", 10, false}
	case "Dirty":
		cfg = Field{"Dirty", Raw, 0, " KB", 10, false}
	case "Writeback":
		cfg = Field{"Writeback", Raw, 0, " KB", 10, false}
	case "AnonPages":
		cfg = Field{"AnonPages", Raw, 0, " KB", 10, false}
	case "Mapped":
		cfg = Field{"Mapped", Raw, 0, " KB", 10, false}
	case "Shmem":
		cfg = Field{"Shmem", Raw, 0, " KB", 10, false}
	case "Slab":
		cfg = Field{"Slab", Raw, 0, " KB", 10, false}
	case "SReclaimable":
		cfg = Field{"SReclaimable", Raw, 0, " KB", 10, false}
	case "SUnreclaim":
		cfg = Field{"SUnreclaim", Raw, 0, " KB", 10, false}
	case "KernelStack":
		cfg = Field{"KernelStack", Raw, 0, " KB", 10, false}
	case "PageTables":
		cfg = Field{"PageTables", Raw, 0, " KB", 10, false}
	case "NFSUnstable":
		cfg = Field{"NFSUnstable", Raw, 0, " KB", 10, false}
	case "Bounce":
		cfg = Field{"Bounce", Raw, 0, " KB", 10, false}
	case "WritebackTmp":
		cfg = Field{"WritebackTmp", Raw, 0, " KB", 10, false}
	case "CommitLimit":
		cfg = Field{"CommitLimit", Raw, 0, " KB", 10, false}
	case "CommittedAS":
		cfg = Field{"CommittedAS", Raw, 0, " KB", 10, false}
	case "VmallocTotal":
		cfg = Field{"VmallocTotal", Raw, 0, " KB", 10, false}
	case "VmallocUsed":
		cfg = Field{"VmallocUsed", Raw, 0, " KB", 10, false}
	case "VmallocChunk":
		cfg = Field{"VmallocChunk", Raw, 0, " KB", 10, false}
	case "HardwareCorrupted":
		cfg = Field{"HardwareCorrupted", Raw, 0, " KB", 10, false}
	case "AnonHugePages":
		cfg = Field{"AnonHugePages", Raw, 0, " KB", 10, false}
	case "ShmemHugePages":
		cfg = Field{"ShmemHugePages", Raw, 0, " KB", 10, false}
	case "ShmemPmdMapped":
		cfg = Field{"ShmemPmdMapped", Raw, 0, " KB", 10, false}
	case "CmaTotal":
		cfg = Field{"CmaTotal", Raw, 0, " KB", 10, false}
	case "CmaFree":
		cfg = Field{"CmaFree", Raw, 0, " KB", 10, false}
	case "HugePagesTotal":
		cfg = Field{"HugePagesTotal", Raw, 0, " KB", 10, false}
	case "HugePagesFree":
		cfg = Field{"HugePagesFree", Raw, 0, " KB", 10, false}
	case "HugePagesRsvd":
		cfg = Field{"HugePagesRsvd", Raw, 0, " KB", 10, false}
	case "HugePagesSurp":
		cfg = Field{"HugePagesSurp", Raw, 0, " KB", 10, false}
	case "Hugepagesize":
		cfg = Field{"Hugepagesize", Raw, 0, " KB", 10, false}
	case "DirectMap4k":
		cfg = Field{"DirectMap4k", Raw, 0, " KB", 10, false}
	case "DirectMap2M":
		cfg = Field{"DirectMap2M", Raw, 0, " KB", 10, false}
	case "DirectMap1G":
		cfg = Field{"DirectMap1G", Raw, 0, " KB", 10, false}
	}
	return cfg
}

func (m *MEM) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "Total":
		cfg = OpenMetricField{"Total", Gauge, "", "", []string{}}
	case "Free":
		cfg = OpenMetricField{"Free", Gauge, "", "", []string{}}
	case "Avail":
		cfg = OpenMetricField{"Avail", Gauge, "", "", []string{}}
	case "HSlab":
		cfg = OpenMetricField{"Slab", Gauge, "", "", []string{}}
	case "Buffer":
		cfg = OpenMetricField{"Buffer", Gauge, "", "", []string{}}
	case "Cache":
		cfg = OpenMetricField{"Cache", Gauge, "", "", []string{}}
	case "MemTotal":
		cfg = OpenMetricField{"MemTotal", Gauge, "", "", []string{}}
	case "MemFree":
		cfg = OpenMetricField{"MemFree", Gauge, "", "", []string{}}
	case "MemAvailable":
		cfg = OpenMetricField{"MemAvailable", Gauge, "", "", []string{}}
	case "Buffers":
		cfg = OpenMetricField{"Buffers", Gauge, "", "", []string{}}
	case "Cached":
		cfg = OpenMetricField{"Cached", Gauge, "", "", []string{}}
	case "SwapCached":
		cfg = OpenMetricField{"SwapCached", Gauge, "", "", []string{}}
	case "Active":
		cfg = OpenMetricField{"Active", Gauge, "", "", []string{}}
	case "Inactive":
		cfg = OpenMetricField{"Inactive", Gauge, "", "", []string{}}
	case "ActiveAnon":
		cfg = OpenMetricField{"ActiveAnon", Gauge, "", "", []string{}}
	case "InactiveAnon":
		cfg = OpenMetricField{"InactiveAnon", Gauge, "", "", []string{}}
	case "ActiveFile":
		cfg = OpenMetricField{"ActiveFile", Gauge, "", "", []string{}}
	case "InactiveFile":
		cfg = OpenMetricField{"InactiveFile", Gauge, "", "", []string{}}
	case "Unevictable":
		cfg = OpenMetricField{"Unevictable", Gauge, "", "", []string{}}
	case "Mlocked":
		cfg = OpenMetricField{"Mlocked", Gauge, "", "", []string{}}
	case "SwapTotal":
		cfg = OpenMetricField{"SwapTotal", Gauge, "", "", []string{}}
	case "SwapFree":
		cfg = OpenMetricField{"SwapFree", Gauge, "", "", []string{}}
	case "Dirty":
		cfg = OpenMetricField{"Dirty", Gauge, "", "", []string{}}
	case "Writeback":
		cfg = OpenMetricField{"Writeback", Gauge, "", "", []string{}}
	case "AnonPages":
		cfg = OpenMetricField{"AnonPages", Gauge, "", "", []string{}}
	case "Mapped":
		cfg = OpenMetricField{"Mapped", Gauge, "", "", []string{}}
	case "Shmem":
		cfg = OpenMetricField{"Shmem", Gauge, "", "", []string{}}
	case "Slab":
		cfg = OpenMetricField{"Slab", Gauge, "", "", []string{}}
	case "SReclaimable":
		cfg = OpenMetricField{"SReclaimable", Gauge, "", "", []string{}}
	case "SUnreclaim":
		cfg = OpenMetricField{"SUnreclaim", Gauge, "", "", []string{}}
	case "KernelStack":
		cfg = OpenMetricField{"KernelStack", Gauge, "", "", []string{}}
	case "PageTables":
		cfg = OpenMetricField{"PageTables", Gauge, "", "", []string{}}
	case "NFSUnstable":
		cfg = OpenMetricField{"NFSUnstable", Gauge, "", "", []string{}}
	case "Bounce":
		cfg = OpenMetricField{"Bounce", Gauge, "", "", []string{}}
	case "WritebackTmp":
		cfg = OpenMetricField{"WritebackTmp", Gauge, "", "", []string{}}
	case "CommitLimit":
		cfg = OpenMetricField{"CommitLimit", Gauge, "", "", []string{}}
	case "CommittedAS":
		cfg = OpenMetricField{"CommittedAS", Gauge, "", "", []string{}}
	case "VmallocTotal":
		cfg = OpenMetricField{"VmallocTotal", Gauge, "", "", []string{}}
	case "VmallocUsed":
		cfg = OpenMetricField{"VmallocUsed", Gauge, "", "", []string{}}
	case "VmallocChunk":
		cfg = OpenMetricField{"VmallocChunk", Gauge, "", "", []string{}}
	case "HardwareCorrupted":
		cfg = OpenMetricField{"HardwareCorrupted", Gauge, "", "", []string{}}
	case "AnonHugePages":
		cfg = OpenMetricField{"AnonHugePages", Gauge, "", "", []string{}}
	case "ShmemHugePages":
		cfg = OpenMetricField{"ShmemHugePages", Gauge, "", "", []string{}}
	case "ShmemPmdMapped":
		cfg = OpenMetricField{"ShmemPmdMapped", Gauge, "", "", []string{}}
	case "CmaTotal":
		cfg = OpenMetricField{"CmaTotal", Gauge, "", "", []string{}}
	case "CmaFree":
		cfg = OpenMetricField{"CmaFree", Gauge, "", "", []string{}}
	case "HugePagesTotal":
		cfg = OpenMetricField{"HugePagesTotal", Gauge, "", "", []string{}}
	case "HugePagesFree":
		cfg = OpenMetricField{"HugePagesFree", Gauge, "", "", []string{}}
	case "HugePagesRsvd":
		cfg = OpenMetricField{"HugePagesRsvd", Gauge, "", "", []string{}}
	case "HugePagesSurp":
		cfg = OpenMetricField{"HugePagesSurp", Gauge, "", "", []string{}}
	case "Hugepagesize":
		cfg = OpenMetricField{"Hugepagesize", Gauge, "", "", []string{}}
	case "DirectMap4k":
		cfg = OpenMetricField{"DirectMap4k", Gauge, "", "", []string{}}
	case "DirectMap2M":
		cfg = OpenMetricField{"DirectMap2M", Gauge, "", "", []string{}}
	case "DirectMap1G":
		cfg = OpenMetricField{"DirectMap1G", Gauge, "", "", []string{}}
	}
	return cfg
}

func (m *MEM) GetRenderValue(field string, opt FieldOpt) string {
	cfg := m.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Total":
		s = cfg.Render(m.MemTotal * 1024)
	case "Free":
		s = cfg.Render(m.MemFree * 1024)
	case "Avail":
		s = cfg.Render(m.MemAvailable * 1024)
	case "HSlab":
		s = cfg.Render(m.Slab * 1024)
	case "Buffer":
		s = cfg.Render(m.Buffers * 1024)
	case "Cache":
		s = cfg.Render(m.Cached * 1024)
	case "MemTotal":
		s = cfg.Render(m.MemTotal)
	case "MemFree":
		s = cfg.Render(m.MemFree)
	case "MemAvailable":
		s = cfg.Render(m.MemAvailable)
	case "Buffers":
		s = cfg.Render(m.Buffers)
	case "Cached":
		s = cfg.Render(m.Cached)
	case "SwapCached":
		s = cfg.Render(m.SwapCached)
	case "Active":
		s = cfg.Render(m.Active)
	case "Inactive":
		s = cfg.Render(m.Inactive)
	case "ActiveAnon":
		s = cfg.Render(m.ActiveAnon)
	case "InactiveAnon":
		s = cfg.Render(m.InactiveAnon)
	case "ActiveFile":
		s = cfg.Render(m.ActiveFile)
	case "InactiveFile":
		s = cfg.Render(m.InactiveFile)
	case "Unevictable":
		s = cfg.Render(m.Unevictable)
	case "Mlocked":
		s = cfg.Render(m.Mlocked)
	case "SwapTotal":
		s = cfg.Render(m.SwapTotal)
	case "SwapFree":
		s = cfg.Render(m.SwapFree)
	case "Dirty":
		s = cfg.Render(m.Dirty)
	case "Writeback":
		s = cfg.Render(m.Writeback)
	case "AnonPages":
		s = cfg.Render(m.AnonPages)
	case "Mapped":
		s = cfg.Render(m.Mapped)
	case "Shmem":
		s = cfg.Render(m.Shmem)
	case "Slab":
		s = cfg.Render(m.Slab)
	case "SReclaimable":
		s = cfg.Render(m.SReclaimable)
	case "SUnreclaim":
		s = cfg.Render(m.SUnreclaim)
	case "KernelStack":
		s = cfg.Render(m.KernelStack)
	case "PageTables":
		s = cfg.Render(m.PageTables)
	case "NFSUnstable":
		s = cfg.Render(m.NFSUnstable)
	case "Bounce":
		s = cfg.Render(m.Bounce)
	case "WritebackTmp":
		s = cfg.Render(m.WritebackTmp)
	case "CommitLimit":
		s = cfg.Render(m.CommitLimit)
	case "CommittedAS":
		s = cfg.Render(m.CommittedAS)
	case "VmallocTotal":
		s = cfg.Render(m.VmallocTotal)
	case "VmallocUsed":
		s = cfg.Render(m.VmallocUsed)
	case "VmallocChunk":
		s = cfg.Render(m.VmallocChunk)
	case "HardwareCorrupted":
		s = cfg.Render(m.HardwareCorrupted)
	case "AnonHugePages":
		s = cfg.Render(m.AnonHugePages)
	case "ShmemHugePages":
		s = cfg.Render(m.ShmemHugePages)
	case "ShmemPmdMapped":
		s = cfg.Render(m.ShmemPmdMapped)
	case "CmaTotal":
		s = cfg.Render(m.CmaTotal)
	case "CmaFree":
		s = cfg.Render(m.CmaFree)
	case "HugePagesTotal":
		s = cfg.Render(m.HugePagesTotal)
	case "HugePagesFree":
		s = cfg.Render(m.HugePagesFree)
	case "HugePagesRsvd":
		s = cfg.Render(m.HugePagesRsvd)
	case "HugePagesSurp":
		s = cfg.Render(m.HugePagesSurp)
	case "Hugepagesize":
		s = cfg.Render(m.Hugepagesize)
	case "DirectMap4k":
		s = cfg.Render(m.DirectMap4k)
	case "DirectMap2M":
		s = cfg.Render(m.DirectMap2M)
	case "DirectMap1G":
		s = cfg.Render(m.DirectMap1G)
	default:
		s = "no " + field + " for mem stat"
	}
	return s
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

func getValueOrDefault[T uint64 | float64](m *T) T {
	if m == nil {
		return 0
	}
	return *m
}
