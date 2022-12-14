package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/xixiliguo/etop/store"
)

var DefaultMEMFields = []string{"Total", "Free", "Avail"}

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

func (m *MEM) GetRenderValue(config RenderConfig, field string) string {
	s := config[field].Render(m.MemTotal)
	switch field {
	case "Total":
		s = config[field].Render(m.MemTotal * 1024)
	case "Free":
		s = config[field].Render(m.MemFree * 1024)
	case "Avail":
		s = config[field].Render(m.MemAvailable * 1024)
	case "Buffer":
		s = config[field].Render(m.Buffers * 1024)
	case "Cache":
		s = config[field].Render(m.Cached * 1024)
	case "MemTotal":
		s = config[field].Render(m.MemTotal)
	case "MemFree":
		s = config[field].Render(m.MemFree)
	case "MemAvailable":
		s = config[field].Render(m.MemAvailable)
	case "Buffers":
		s = config[field].Render(m.Buffers)
	case "Cached":
		s = config[field].Render(m.Cached)
	case "SwapCached":
		s = config[field].Render(m.SwapCached)
	case "Active":
		s = config[field].Render(m.Active)
	case "Inactive":
		s = config[field].Render(m.Inactive)
	case "ActiveAnon":
		s = config[field].Render(m.ActiveAnon)
	case "InactiveAnon":
		s = config[field].Render(m.InactiveAnon)
	case "ActiveFile":
		s = config[field].Render(m.ActiveFile)
	case "InactiveFile":
		s = config[field].Render(m.InactiveFile)
	case "Unevictable":
		s = config[field].Render(m.Unevictable)
	case "Mlocked":
		s = config[field].Render(m.Mlocked)
	case "SwapTotal":
		s = config[field].Render(m.SwapTotal)
	case "SwapFree":
		s = config[field].Render(m.SwapFree)
	case "Dirty":
		s = config[field].Render(m.Dirty)
	case "Writeback":
		s = config[field].Render(m.Writeback)
	case "AnonPages":
		s = config[field].Render(m.AnonPages)
	case "Mapped":
		s = config[field].Render(m.Mapped)
	case "Shmem":
		s = config[field].Render(m.Shmem)
	case "Slab":
		s = config[field].Render(m.Slab)
	case "SReclaimable":
		s = config[field].Render(m.SReclaimable)
	case "SUnreclaim":
		s = config[field].Render(m.SUnreclaim)
	case "KernelStack":
		s = config[field].Render(m.KernelStack)
	case "PageTables":
		s = config[field].Render(m.PageTables)
	case "NFSUnstable":
		s = config[field].Render(m.NFSUnstable)
	case "Bounce":
		s = config[field].Render(m.Bounce)
	case "WritebackTmp":
		s = config[field].Render(m.WritebackTmp)
	case "CommitLimit":
		s = config[field].Render(m.CommitLimit)
	case "CommittedAS":
		s = config[field].Render(m.CommittedAS)
	case "VmallocTotal":
		s = config[field].Render(m.VmallocTotal)
	case "VmallocUsed":
		s = config[field].Render(m.VmallocUsed)
	case "VmallocChunk":
		s = config[field].Render(m.VmallocChunk)
	case "HardwareCorrupted":
		s = config[field].Render(m.HardwareCorrupted)
	case "AnonHugePages":
		s = config[field].Render(m.AnonHugePages)
	case "ShmemHugePages":
		s = config[field].Render(m.ShmemHugePages)
	case "ShmemPmdMapped":
		s = config[field].Render(m.ShmemPmdMapped)
	case "CmaTotal":
		s = config[field].Render(m.CmaTotal)
	case "CmaFree":
		s = config[field].Render(m.CmaFree)
	case "HugePagesTotal":
		s = config[field].Render(m.HugePagesTotal)
	case "HugePagesFree":
		s = config[field].Render(m.HugePagesFree)
	case "HugePagesRsvd":
		s = config[field].Render(m.HugePagesRsvd)
	case "HugePagesSurp":
		s = config[field].Render(m.HugePagesSurp)
	case "Hugepagesize":
		s = config[field].Render(m.Hugepagesize)
	case "DirectMap4k":
		s = config[field].Render(m.DirectMap4k)
	case "DirectMap2M":
		s = config[field].Render(m.DirectMap2M)
	case "DirectMap1G":
		s = config[field].Render(m.DirectMap1G)
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

func getValueOrDefault(m *uint64) uint64 {
	if m == nil {
		return 0
	}
	return *m
}

func (m *MEM) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
		row := strings.Builder{}
		row.WriteString(dateTime)
		for _, f := range opt.Fields {
			renderValue := m.GetRenderValue(config, f)
			if f == opt.SelectField && opt.Filter != nil {
				if opt.Filter.MatchString(renderValue) == false {
					continue
				}
			}
			row.WriteString(" ")
			row.WriteString(renderValue)
		}
		row.WriteString("\n")
		opt.Output.WriteString(row.String())

	case "json":
		t := []any{}

		row := make(map[string]string)
		row["Timestamp"] = dateTime
		for _, f := range opt.Fields {
			renderValue := m.GetRenderValue(config, f)
			if f == opt.SelectField && opt.Filter != nil {
				if opt.Filter.MatchString(renderValue) == false {
					continue
				}
			}
			row[config[f].Name] = renderValue
		}
		t = append(t, row)

		b, _ := json.Marshal(t)
		opt.Output.Write(b)
	}

}
