package model

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/mattn/go-runewidth"
	"github.com/xixiliguo/etop/store"
)

var DefaultCgroupFields = []string{"Name", "NrDescendants", "NrDyingDescendants", "UsagePercent", "Controllers"}
var AllCgroupFields = []string{"Path", "Name", "Level", "Inode", "Controllers",
	"NrDescendants", "NrDyingDescendants",
	"UsagePercent", "UserPercent", "SystemPercent", "NrPeriodsPerSec", "NrThrottledPerSec", "ThrottledPercent", "NrBurstsPerSec", "BurstPercent",
	"Anon", "File", "KernelStack", "Slab", "Sock", "Shmem", "Zswap", "Zswapped", "FileMapped",
	"FileDirty", "FileWriteback", "AnonThp", "InactiveAnon", "ActiveAnon", "InactiveFile", "ActiveFile", "Unevictable",
	"SlabReclaimable", "SlabUnreclaimable", "Pgfault", "Pgmajfault", "WorkingsetRefault", "WorkingsetActivate",
	"WorkingsetNodereclaim", "Pgrefill", "Pgscan", "Pgsteal", "Pgactivate", "Pgdeactivate", "Pglazyfree", "Pglazyfreed",
	"ZswpIn", "ZswpOut", "ThpFaultAlloc", "ThpCollapseAlloc",
	"CpuSetCpus", "CpuSetCpusEffective", "CpuSetMems", "CpuSetMemsEffective", "CpuWeight", "CpuMax",
	"MemoryCurrent", "MemoryLow", "MemoryHigh", "MemoryMin", "MemoryMax", "MemoryPeak", "SwapCurrent", "SwapMax",
	"ZswapCurrent", "ZswapMax", "EventLow", "EventHigh", "EventMax", "EventOom", "EventOomKill",
	"RbytePerSec", "WbytePerSec", "RioPerSec", "WioPerSec", "DbytePerSec", "DioPerSec",
	"CPUSomePressure", "CPUFullPressure", "MemorySomePressure", "MemoryFullPressure", "IOSomePressure", "IOFullPressure"}

var FiledToCgroupFile = map[string]store.CgroupFile{
	"Anon":                        store.MemoryStatFile,
	"File":                        store.MemoryStatFile,
	"KernelStack":                 store.MemoryStatFile,
	"Slab":                        store.MemoryStatFile,
	"Sock":                        store.MemoryStatFile,
	"Shmem":                       store.MemoryStatFile,
	"Zswap":                       store.MemoryStatFile,
	"Zswapped":                    store.MemoryStatFile,
	"FileMapped":                  store.MemoryStatFile,
	"FileDirty":                   store.MemoryStatFile,
	"FileWriteback":               store.MemoryStatFile,
	"AnonThp":                     store.MemoryStatFile,
	"InactiveAnon":                store.MemoryStatFile,
	"ActiveAnon":                  store.MemoryStatFile,
	"InactiveFile":                store.MemoryStatFile,
	"ActiveFile":                  store.MemoryStatFile,
	"Unevictable":                 store.MemoryStatFile,
	"SlabReclaimable":             store.MemoryStatFile,
	"SlabUnreclaimable":           store.MemoryStatFile,
	"PgfaultPerSec":               store.MemoryStatFile,
	"PgmajfaultPerSec":            store.MemoryStatFile,
	"WorkingsetRefaultPerSec":     store.MemoryStatFile,
	"WorkingsetActivatePerSec":    store.MemoryStatFile,
	"WorkingsetNodereclaimPerSec": store.MemoryStatFile,
	"PgrefillPerSec":              store.MemoryStatFile,
	"PgscanPerSec":                store.MemoryStatFile,
	"PgstealPerSec":               store.MemoryStatFile,
	"PgactivatePerSec":            store.MemoryStatFile,
	"PgdeactivatePerSec":          store.MemoryStatFile,
	"PglazyfreePerSec":            store.MemoryStatFile,
	"PglazyfreedPerSec":           store.MemoryStatFile,
	"ZswpInPerSec":                store.MemoryStatFile,
	"ZswpOutPerSec":               store.MemoryStatFile,
	"ThpFaultAllocPerSec":         store.MemoryStatFile,
	"ThpCollapseAllocPerSec":      store.MemoryStatFile,
	"CpuSetCpus":                  store.CpuSetCpusFile,
	"CpuSetCpusEffective":         store.CpuSetCpusEffectiveFile,
	"CpuSetMems":                  store.CpuSetMemsFileFile,
	"CpuSetMemsEffective":         store.CpuSetMemsEffectiveFile,
	"CpuWeight":                   store.CpuWeightFile,
	"CpuMax":                      store.CpuMaxFile,
	"MemoryCurrent":               store.MemoryCurrentFile,
	"MemoryLow":                   store.MemoryLowFile,
	"MemoryHigh":                  store.MemoryHighFile,
	"MemoryMin":                   store.MemoryMinFile,
	"MemoryMax":                   store.MemoryMaxFile,
	"MemoryPeak":                  store.MemoryPeakFile,
	"SwapCurrent":                 store.MemorySwapCurrentFile,
	"SwapMax":                     store.MemorySwapMaxFile,
	"ZswapCurrent":                store.MemoryZswapCurrentFile,
	"ZswapMax":                    store.MemoryZswapMaxFile,
	"EventLow":                    store.MemoryEventsFile,
	"EventHigh":                   store.MemoryEventsFile,
	"EventMax":                    store.MemoryEventsFile,
	"EventOom":                    store.MemoryEventsFile,
	"EventOomKill":                store.MemoryEventsFile,
	"RbytePerSec":                 store.IoStatFile,
	"WbytePerSec":                 store.IoStatFile,
	"RioPerSec":                   store.IoStatFile,
	"WioPerSec":                   store.IoStatFile,
	"DbytePerSec":                 store.IoStatFile,
	"DioPerSec":                   store.IoStatFile,
	"CPUSomePressure":             store.CpuPressureFile,
	"CPUFullPressure":             store.CpuPressureFile,
	"MemorySomePressure":          store.MemoryPressureFile,
	"MemoryFullPressure":          store.MemoryPressureFile,
	"IOSomePressure":              store.IoPressureFile,
	"IOFullPressure":              store.IoPressureFile,
}

type Cgroup struct {
	Path        string
	Name        string
	Level       int
	Inode       uint64
	IsExpand    bool
	Child       map[string]*Cgroup
	IsNotExist  map[store.CgroupFile]struct{}
	Controllers string
	store.CgoupStat
	UsagePercent      float64
	UserPercent       float64
	SystemPercent     float64
	NrPeriodsPerSec   float64
	NrThrottledPerSec float64
	ThrottledPercent  float64
	NrBurstsPerSec    float64
	BurstPercent      float64
	// from memory.stat
	Anon                        uint64
	File                        uint64
	KernelStack                 uint64
	Slab                        uint64
	Sock                        uint64
	Shmem                       uint64
	Zswap                       uint64
	Zswapped                    uint64
	FileMapped                  uint64
	FileDirty                   uint64
	FileWriteback               uint64
	AnonThp                     uint64
	InactiveAnon                uint64
	ActiveAnon                  uint64
	InactiveFile                uint64
	ActiveFile                  uint64
	Unevictable                 uint64
	SlabReclaimable             uint64
	SlabUnreclaimable           uint64
	PgfaultPerSec               float64
	PgmajfaultPerSec            float64
	WorkingsetRefaultPerSec     float64
	WorkingsetActivatePerSec    float64
	WorkingsetNodereclaimPerSec float64
	PgrefillPerSec              float64
	PgscanPerSec                float64
	PgstealPerSec               float64
	PgactivatePerSec            float64
	PgdeactivatePerSec          float64
	PglazyfreePerSec            float64
	PglazyfreedPerSec           float64
	ZswpInPerSec                float64
	ZswpOutPerSec               float64
	ThpFaultAllocPerSec         float64
	ThpCollapseAllocPerSec      float64
	CpuSetCpus                  string
	CpuSetCpusEffective         string
	CpuSetMems                  string
	CpuSetMemsEffective         string
	CpuWeight                   uint64
	CpuMax                      string
	MemoryCurrent               uint64
	MemoryLow                   uint64
	MemoryHigh                  uint64
	MemoryMin                   uint64
	MemoryMax                   uint64
	MemoryPeak                  uint64
	SwapCurrent                 uint64
	SwapMax                     uint64
	ZswapCurrent                uint64
	ZswapMax                    uint64
	// memory event
	EventLow     uint64
	EventHigh    uint64
	EventMax     uint64
	EventOom     uint64
	EventOomKill uint64
	// io.stat
	RbytePerSec float64
	WbytePerSec float64
	RioPerSec   float64
	WioPerSec   float64
	DbytePerSec float64
	DioPerSec   float64
	// pressure file
	CPUSomePressure    float64
	CPUFullPressure    float64
	MemorySomePressure float64
	MemoryFullPressure float64
	IOSomePressure     float64
	IOFullPressure     float64
}

func (c *Cgroup) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Cgroup":
		cfg = Field{"Cgroup", Raw, 0, "", 50, false}
	case "Path":
		cfg = Field{"Path", Raw, 0, "", 50, false}
	case "Name":
		cfg = Field{"Name", Raw, 0, "", 50, false}
	case "Level":
		cfg = Field{"Level", Raw, 0, "", 10, false}
	case "Inode":
		cfg = Field{"Inode", Raw, 0, "", 10, false}
	case "Controllers":
		cfg = Field{"Controllers", Raw, 0, "", 10, false}
	case "NrDescendants":
		cfg = Field{"NrDescendants", Raw, 0, "", 10, false}
	case "NrDyingDescendants":
		cfg = Field{"NrDyingDescendants", Raw, 0, "", 10, false}
	case "UsagePercent":
		cfg = Field{"CPU", Raw, 1, "%", 10, false}
	case "UserPercent":
		cfg = Field{"UserCPU", Raw, 1, "%", 10, false}
	case "SystemPercent":
		cfg = Field{"SystemCPU", Raw, 1, "%", 10, false}
	case "NrPeriodsPerSec":
		cfg = Field{"NrPeriods/s", Raw, 1, "/s", 10, false}
	case "NrThrottledPerSec":
		cfg = Field{"NrThrottled/s", Raw, 1, "/s", 10, false}
	case "ThrottledPercent":
		cfg = Field{"Throttled", Raw, 1, "%", 10, false}
	case "NrBurstsPerSec":
		cfg = Field{"NrBursts/s", Raw, 1, "/s", 10, false}
	case "BurstPercent":
		cfg = Field{"Burst", Raw, 1, "%", 10, false}
	case "Anon":
		cfg = Field{"Anon", HumanReadableSize, 1, "", 10, false}
	case "File":
		cfg = Field{"File", HumanReadableSize, 1, "", 10, false}
	case "KernelStack":
		cfg = Field{"KernelStack", HumanReadableSize, 1, "", 10, false}
	case "Slab":
		cfg = Field{"Slab", HumanReadableSize, 1, "", 10, false}
	case "Sock":
		cfg = Field{"Sock", HumanReadableSize, 1, "", 10, false}
	case "Shmem":
		cfg = Field{"Shmem", HumanReadableSize, 1, "", 10, false}
	case "Zswap":
		cfg = Field{"Zswap", HumanReadableSize, 1, "", 10, false}
	case "Zswapped":
		cfg = Field{"Zswapped", HumanReadableSize, 1, "", 10, false}
	case "FileMapped":
		cfg = Field{"FileMapped", HumanReadableSize, 1, "", 10, false}
	case "FileDirty":
		cfg = Field{"FileDirty", HumanReadableSize, 1, "", 10, false}
	case "FileWriteback":
		cfg = Field{"FileWriteback", HumanReadableSize, 1, "", 10, false}
	case "AnonThp":
		cfg = Field{"AnonThp", HumanReadableSize, 1, "", 10, false}
	case "InactiveAnon":
		cfg = Field{"InactiveAnon", HumanReadableSize, 1, "", 10, false}
	case "ActiveAnon":
		cfg = Field{"ActiveAnon", HumanReadableSize, 1, "", 10, false}
	case "InactiveFile":
		cfg = Field{"InactiveFile", HumanReadableSize, 1, "", 10, false}
	case "ActiveFile":
		cfg = Field{"ActiveFile", HumanReadableSize, 1, "", 10, false}
	case "Unevictable":
		cfg = Field{"Unevictable", HumanReadableSize, 1, "", 10, false}
	case "SlabReclaimable":
		cfg = Field{"SlabReclaimable", HumanReadableSize, 1, "", 10, false}
	case "SlabUnreclaimable":
		cfg = Field{"SlabUnreclaimable", HumanReadableSize, 1, "", 10, false}
	case "PgfaultPerSec":
		cfg = Field{"Pgfault/s", Raw, 0, "/s", 10, false}
	case "PgmajfaultPerSec":
		cfg = Field{"Pgmajfault/s", Raw, 0, "/s", 10, false}
	case "WorkingsetRefaultPerSec":
		cfg = Field{"WorkingsetRefault/s", Raw, 0, "/s", 10, false}
	case "WorkingsetActivatePerSec":
		cfg = Field{"WorkingsetActivate/s", Raw, 0, "/s", 10, false}
	case "WorkingsetNodereclaimPerSec":
		cfg = Field{"WorkingsetNodereclaim/s", Raw, 0, "/s", 10, false}
	case "PgrefillPerSec":
		cfg = Field{"Pgrefill/s", Raw, 0, "/s", 10, false}
	case "PgscanPerSec":
		cfg = Field{"Pgscan/s", Raw, 0, "/s", 10, false}
	case "PgstealPerSec":
		cfg = Field{"Pgsteal/s", Raw, 0, "/s", 10, false}
	case "PgactivatePerSec":
		cfg = Field{"Pgactivate/s", Raw, 0, "/s", 10, false}
	case "PgdeactivatePerSec":
		cfg = Field{"Pgdeactivate/s", Raw, 0, "/s", 10, false}
	case "PglazyfreePerSec":
		cfg = Field{"Pglazyfree/s", Raw, 0, "/s", 10, false}
	case "PglazyfreedPerSec":
		cfg = Field{"Pglazyfreed/s", Raw, 0, "/s", 10, false}
	case "ZswpInPerSec":
		cfg = Field{"ZswpIn/s", Raw, 0, "/s", 10, false}
	case "ZswpOutPerSec":
		cfg = Field{"ZswpOut/s", Raw, 0, "/s", 10, false}
	case "ThpFaultAllocPerSec":
		cfg = Field{"ThpFaultAlloc/s", Raw, 0, "/s", 10, false}
	case "ThpCollapseAllocPerSec":
		cfg = Field{"ThpCollapseAlloc/s", Raw, 0, "/s", 10, false}
	case "CpuSetCpus":
		cfg = Field{"CpuSetCpus", Raw, 0, "", 10, false}
	case "CpuSetCpusEffective":
		cfg = Field{"CpuSetCpusEffective", Raw, 0, "", 10, false}
	case "CpuSetMems":
		cfg = Field{"CpuSetMems", Raw, 0, "", 10, false}
	case "CpuSetMemsEffective":
		cfg = Field{"CpuSetMemsEffective", Raw, 0, "", 10, false}
	case "CpuWeight":
		cfg = Field{"CpuWeight", Raw, 0, "", 10, false}
	case "CpuMax":
		cfg = Field{"CpuMax", Raw, 0, "", 10, false}
	case "MemoryCurrent":
		cfg = Field{"Memory", HumanReadableSize, 1, "", 10, false}
	case "MemoryLow":
		cfg = Field{"MemoryLow", HumanReadableSize, 1, "", 10, false}
	case "MemoryHigh":
		cfg = Field{"MemoryHigh", HumanReadableSize, 1, "", 10, false}
	case "MemoryMin":
		cfg = Field{"MemoryMin", HumanReadableSize, 1, "", 10, false}
	case "MemoryMax":
		cfg = Field{"MemoryMax", HumanReadableSize, 1, "", 10, false}
	case "MemoryPeak":
		cfg = Field{"MemoryPeak", HumanReadableSize, 1, "", 10, false}
	case "SwapCurrent":
		cfg = Field{"Swap", HumanReadableSize, 1, "", 10, false}
	case "SwapMax":
		cfg = Field{"SwapMax", HumanReadableSize, 1, "", 10, false}
	case "ZswapCurrent":
		cfg = Field{"Zswap", HumanReadableSize, 1, "", 10, false}
	case "ZswapMax":
		cfg = Field{"ZswapMax", HumanReadableSize, 1, "", 10, false}
	case "EventLow":
		cfg = Field{"EventLow", Raw, 0, "", 10, false}
	case "EventHigh":
		cfg = Field{"EventHigh", Raw, 0, "", 10, false}
	case "EventMax":
		cfg = Field{"EventMax", Raw, 0, "", 10, false}
	case "EventOom":
		cfg = Field{"EventOom", Raw, 0, "", 10, false}
	case "EventOomKill":
		cfg = Field{"EventOomKill", Raw, 0, "", 10, false}
	case "RbytePerSec":
		cfg = Field{"Rbyte/s", HumanReadableSize, 1, "/s", 10, false}
	case "WbytePerSec":
		cfg = Field{"Wbyte/s", HumanReadableSize, 1, "/s", 10, false}
	case "RioPerSec":
		cfg = Field{"Rio/s", Raw, 1, "/s", 10, false}
	case "WioPerSec":
		cfg = Field{"Wio/s", Raw, 1, "/s", 10, false}
	case "DbytePerSec":
		cfg = Field{"Dbyte/s", HumanReadableSize, 1, "/s", 10, false}
	case "DioPerSec":
		cfg = Field{"Dio/s", Raw, 1, "/s", 10, false}
	case "CPUSomePressure":
		cfg = Field{"CPUSomePressure", Raw, 0, "%", 10, false}
	case "CPUFullPressure":
		cfg = Field{"CPUFullPressure", Raw, 0, "%", 10, false}
	case "MemorySomePressure":
		cfg = Field{"MemorySomePressure", Raw, 0, "%", 10, false}
	case "MemoryFullPressure":
		cfg = Field{"MemoryFullPressure", Raw, 0, "%", 10, false}
	case "IOSomePressure":
		cfg = Field{"IOSomePressure", Raw, 0, "%", 10, false}
	case "IOFullPressure":
		cfg = Field{"IOFullPressure", Raw, 0, "%", 10, false}
	}
	return cfg
}

func (c *Cgroup) GetRenderValue(field string, opt FieldOpt) string {

	cfg := c.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""

	if file, ok := FiledToCgroupFile[field]; ok {
		if _, ok := c.IsNotExist[file]; ok {
			if cfg.FixWidth == true {
				return fmt.Sprintf("%-[1]*s", cfg.Width, "-")
			}
			return "-"
		}
	}

	switch field {
	case "Cgroup":
		s = cfg.Render(filepath.Join(c.Path, c.Name))
	case "Path":
		s = cfg.Render(c.Path)
	case "Name":
		indents := ""
		if c.Level > 1 {
			indents = strings.Repeat("   ", c.Level-1)
		}

		s = indents + "└── " + c.Name
		if c.IsExpand == false && len(c.Child) != 0 {
			s = indents + "└─+ " + c.Name
		}
		if c.Level == 0 {
			s = "/"
		}
		if cfg.FixWidth == true {
			pad := cfg.Width - runewidth.StringWidth(s)
			if pad > 0 {
				s = s + strings.Repeat(" ", pad)
			}
		}
	case "Level":
		s = cfg.Render(c.Level)
	case "Inode":
		s = cfg.Render(c.Inode)
	case "Controllers":
		s = cfg.Render(c.Controllers)
	case "NrDescendants":
		s = cfg.Render(c.NrDescendants)
	case "NrDyingDescendants":
		s = cfg.Render(c.NrDyingDescendants)
	case "UsagePercent":
		s = cfg.Render(c.UsagePercent)
	case "UserPercent":
		s = cfg.Render(c.UserPercent)
	case "SystemPercent":
		s = cfg.Render(c.SystemPercent)
	case "NrPeriodsPerSec":
		s = cfg.Render(c.NrPeriodsPerSec)
	case "NrThrottledPerSec":
		s = cfg.Render(c.NrThrottledPerSec)
	case "ThrottledPercent":
		s = cfg.Render(c.ThrottledPercent)
	case "NrBurstsPerSec":
		s = cfg.Render(c.NrBurstsPerSec)
	case "BurstPercent":
		s = cfg.Render(c.BurstPercent)
	case "Anon":
		s = cfg.Render(c.Anon)
	case "File":
		s = cfg.Render(c.File)
	case "KernelStack":
		s = cfg.Render(c.KernelStack)
	case "Slab":
		s = cfg.Render(c.Slab)
	case "Sock":
		s = cfg.Render(c.Sock)
	case "Shmem":
		s = cfg.Render(c.Shmem)
	case "Zswap":
		s = cfg.Render(c.Zswap)
	case "Zswapped":
		s = cfg.Render(c.Zswapped)
	case "FileMapped":
		s = cfg.Render(c.FileMapped)
	case "FileDirty":
		s = cfg.Render(c.FileDirty)
	case "FileWriteback":
		s = cfg.Render(c.FileWriteback)
	case "AnonThp":
		s = cfg.Render(c.AnonThp)
	case "InactiveAnon":
		s = cfg.Render(c.InactiveAnon)
	case "ActiveAnon":
		s = cfg.Render(c.ActiveAnon)
	case "InactiveFile":
		s = cfg.Render(c.InactiveFile)
	case "ActiveFile":
		s = cfg.Render(c.ActiveFile)
	case "Unevictable":
		s = cfg.Render(c.Unevictable)
	case "SlabReclaimable":
		s = cfg.Render(c.SlabReclaimable)
	case "SlabUnreclaimable":
		s = cfg.Render(c.SlabUnreclaimable)
	case "PgfaultPerSec":
		s = cfg.Render(c.PgfaultPerSec)
	case "PgmajfaultPerSec":
		s = cfg.Render(c.PgmajfaultPerSec)
	case "WorkingsetRefaultPerSec":
		s = cfg.Render(c.WorkingsetRefaultPerSec)
	case "WorkingsetActivatePerSec":
		s = cfg.Render(c.WorkingsetActivatePerSec)
	case "WorkingsetNodereclaimPerSec":
		s = cfg.Render(c.WorkingsetNodereclaimPerSec)
	case "PgrefillPerSec":
		s = cfg.Render(c.PgrefillPerSec)
	case "PgscanPerSec":
		s = cfg.Render(c.PgscanPerSec)
	case "PgstealPerSec":
		s = cfg.Render(c.PgstealPerSec)
	case "PgactivatePerSec":
		s = cfg.Render(c.PgactivatePerSec)
	case "PgdeactivatePerSec":
		s = cfg.Render(c.PgdeactivatePerSec)
	case "PglazyfreePerSec":
		s = cfg.Render(c.PglazyfreePerSec)
	case "PglazyfreedPerSec":
		s = cfg.Render(c.PglazyfreedPerSec)
	case "ZswpInPerSec":
		s = cfg.Render(c.ZswpInPerSec)
	case "ZswpOutPerSec":
		s = cfg.Render(c.ZswpOutPerSec)
	case "ThpFaultAllocPerSec":
		s = cfg.Render(c.ThpFaultAllocPerSec)
	case "ThpCollapseAllocPerSec":
		s = cfg.Render(c.ThpCollapseAllocPerSec)
	case "CpuSetCpus":
		s = cfg.Render(c.CpuSetCpus)
	case "CpuSetCpusEffective":
		s = cfg.Render(c.CpuSetCpusEffective)
	case "CpuSetMems":
		s = cfg.Render(c.CpuSetMems)
	case "CpuSetMemsEffective":
		s = cfg.Render(c.CpuSetMemsEffective)
	case "CpuWeight":
		s = cfg.Render(c.CpuWeight)
	case "CpuMax":
		s = cfg.Render(c.CpuMax)
	case "MemoryCurrent":
		s = cfg.Render(c.MemoryCurrent)
	case "MemoryLow":
		s = cfg.Render(c.MemoryLow)
	case "MemoryHigh":
		s = cfg.Render(c.MemoryHigh)
	case "MemoryMin":
		s = cfg.Render(c.MemoryMin)
	case "MemoryMax":
		s = cfg.Render(c.MemoryMax)
	case "MemoryPeak":
		s = cfg.Render(c.MemoryPeak)
	case "SwapCurrent":
		s = cfg.Render(c.SwapCurrent)
	case "SwapMax":
		s = cfg.Render(c.SwapMax)
	case "ZswapCurrent":
		s = cfg.Render(c.ZswapCurrent)
	case "ZswapMax":
		s = cfg.Render(c.ZswapMax)
	case "EventLow":
		s = cfg.Render(c.EventLow)
	case "EventHigh":
		s = cfg.Render(c.EventHigh)
	case "EventMax":
		s = cfg.Render(c.EventMax)
	case "EventOom":
		s = cfg.Render(c.EventOom)
	case "EventOomKill":
		s = cfg.Render(c.EventOomKill)
	case "RbytePerSec":
		s = cfg.Render(c.RbytePerSec)
	case "WbytePerSec":
		s = cfg.Render(c.WbytePerSec)
	case "RioPerSec":
		s = cfg.Render(c.RioPerSec)
	case "WioPerSec":
		s = cfg.Render(c.WioPerSec)
	case "DbytePerSec":
		s = cfg.Render(c.DbytePerSec)
	case "DioPerSec":
		s = cfg.Render(c.DioPerSec)
	case "CPUSomePressure":
		s = cfg.Render(c.CPUSomePressure)
	case "CPUFullPressure":
		s = cfg.Render(c.CPUFullPressure)
	case "MemorySomePressure":
		s = cfg.Render(c.MemorySomePressure)
	case "MemoryFullPressure":
		s = cfg.Render(c.MemoryFullPressure)
	case "IOSomePressure":
		s = cfg.Render(c.IOSomePressure)
	case "IOFullPressure":
		s = cfg.Render(c.IOFullPressure)
	default:
		s = "no " + field + " for cgroup stat"
	}
	return s
}

func (c *Cgroup) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "Path":
		cfg = OpenMetricField{"Path", Gauge, "", "", []string{"Cgroup"}}
	case "Name":
		cfg = OpenMetricField{"Name", Gauge, "", "", []string{"Cgroup"}}
	case "Level":
		cfg = OpenMetricField{"Level", Gauge, "", "", []string{"Cgroup"}}
	case "Inode":
		cfg = OpenMetricField{"Inode", Gauge, "", "", []string{"Cgroup"}}
	case "Controllers":
		cfg = OpenMetricField{"Controllers", Gauge, "", "", []string{"Cgroup"}}
	case "NrDescendants":
		cfg = OpenMetricField{"NrDescendants", Gauge, "", "", []string{"Cgroup"}}
	case "NrDyingDescendants":
		cfg = OpenMetricField{"NrDyingDescendants", Gauge, "", "", []string{"Cgroup"}}
	case "UsagePercent":
		cfg = OpenMetricField{"CPU", Gauge, "", "", []string{"Cgroup"}}
	case "UserPercent":
		cfg = OpenMetricField{"UserCPU", Gauge, "", "", []string{"Cgroup"}}
	case "SystemPercent":
		cfg = OpenMetricField{"SystemCPU", Gauge, "", "", []string{"Cgroup"}}
	case "NrPeriodsPerSec":
		cfg = OpenMetricField{"NrPeriodsPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "NrThrottledPerSec":
		cfg = OpenMetricField{"NrThrottledPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "ThrottledPercent":
		cfg = OpenMetricField{"Throttled", Gauge, "", "", []string{"Cgroup"}}
	case "NrBurstsPerSec":
		cfg = OpenMetricField{"NrBurstsPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "BurstPercent":
		cfg = OpenMetricField{"Burst", Gauge, "", "", []string{"Cgroup"}}
	case "Anon":
		cfg = OpenMetricField{"Anon", Gauge, "", "", []string{"Cgroup"}}
	case "File":
		cfg = OpenMetricField{"File", Gauge, "", "", []string{"Cgroup"}}
	case "KernelStack":
		cfg = OpenMetricField{"KernelStack", Gauge, "", "", []string{"Cgroup"}}
	case "Slab":
		cfg = OpenMetricField{"Slab", Gauge, "", "", []string{"Cgroup"}}
	case "Sock":
		cfg = OpenMetricField{"Sock", Gauge, "", "", []string{"Cgroup"}}
	case "Shmem":
		cfg = OpenMetricField{"Shmem", Gauge, "", "", []string{"Cgroup"}}
	case "Zswap":
		cfg = OpenMetricField{"Zswap", Gauge, "", "", []string{"Cgroup"}}
	case "Zswapped":
		cfg = OpenMetricField{"Zswapped", Gauge, "", "", []string{"Cgroup"}}
	case "FileMapped":
		cfg = OpenMetricField{"FileMapped", Gauge, "", "", []string{"Cgroup"}}
	case "FileDirty":
		cfg = OpenMetricField{"FileDirty", Gauge, "", "", []string{"Cgroup"}}
	case "FileWriteback":
		cfg = OpenMetricField{"FileWriteback", Gauge, "", "", []string{"Cgroup"}}
	case "AnonThp":
		cfg = OpenMetricField{"AnonThp", Gauge, "", "", []string{"Cgroup"}}
	case "InactiveAnon":
		cfg = OpenMetricField{"InactiveAnon", Gauge, "", "", []string{"Cgroup"}}
	case "ActiveAnon":
		cfg = OpenMetricField{"ActiveAnon", Gauge, "", "", []string{"Cgroup"}}
	case "InactiveFile":
		cfg = OpenMetricField{"InactiveFile", Gauge, "", "", []string{"Cgroup"}}
	case "ActiveFile":
		cfg = OpenMetricField{"ActiveFile", Gauge, "", "", []string{"Cgroup"}}
	case "Unevictable":
		cfg = OpenMetricField{"Unevictable", Gauge, "", "", []string{"Cgroup"}}
	case "SlabReclaimable":
		cfg = OpenMetricField{"SlabReclaimable", Gauge, "", "", []string{"Cgroup"}}
	case "SlabUnreclaimable":
		cfg = OpenMetricField{"SlabUnreclaimable", Gauge, "", "", []string{"Cgroup"}}
	case "PgfaultPerSec":
		cfg = OpenMetricField{"PgfaultPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgmajfaultPerSec":
		cfg = OpenMetricField{"PgmajfaultPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "WorkingsetRefaultPerSec":
		cfg = OpenMetricField{"WorkingsetRefaultPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "WorkingsetActivatePerSec":
		cfg = OpenMetricField{"WorkingsetActivatePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "WorkingsetNodereclaimPerSec":
		cfg = OpenMetricField{"WorkingsetNodereclaimPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgrefillPerSec":
		cfg = OpenMetricField{"PgrefillPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgscanPerSec":
		cfg = OpenMetricField{"PgscanPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgstealPerSec":
		cfg = OpenMetricField{"PgstealPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgactivatePerSec":
		cfg = OpenMetricField{"PgactivatePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PgdeactivatePerSec":
		cfg = OpenMetricField{"PgdeactivatePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PglazyfreePerSec":
		cfg = OpenMetricField{"PglazyfreePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "PglazyfreedPerSec":
		cfg = OpenMetricField{"PglazyfreedPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "ZswpInPerSec":
		cfg = OpenMetricField{"ZswpInPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "ZswpOutPerSec":
		cfg = OpenMetricField{"ZswpOutPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "ThpFaultAllocPerSec":
		cfg = OpenMetricField{"ThpFaultAllocPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "ThpCollapseAllocPerSec":
		cfg = OpenMetricField{"ThpCollapseAllocPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "CpuSetCpus":
		cfg = OpenMetricField{"CpuSetCpus", Gauge, "", "", []string{"Cgroup"}}
	case "CpuSetCpusEffective":
		cfg = OpenMetricField{"CpuSetCpusEffective", Gauge, "", "", []string{"Cgroup"}}
	case "CpuSetMems":
		cfg = OpenMetricField{"CpuSetMems", Gauge, "", "", []string{"Cgroup"}}
	case "CpuSetMemsEffective":
		cfg = OpenMetricField{"CpuSetMemsEffective", Gauge, "", "", []string{"Cgroup"}}
	case "CpuWeight":
		cfg = OpenMetricField{"CpuWeight", Gauge, "", "", []string{"Cgroup"}}
	case "CpuMax":
		cfg = OpenMetricField{"CpuMax", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryCurrent":
		cfg = OpenMetricField{"Memory", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryLow":
		cfg = OpenMetricField{"MemoryLow", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryHigh":
		cfg = OpenMetricField{"MemoryHigh", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryMin":
		cfg = OpenMetricField{"MemoryMin", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryMax":
		cfg = OpenMetricField{"MemoryMax", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryPeak":
		cfg = OpenMetricField{"MemoryPeak", Gauge, "", "", []string{"Cgroup"}}
	case "SwapCurrent":
		cfg = OpenMetricField{"Swap", Gauge, "", "", []string{"Cgroup"}}
	case "SwapMax":
		cfg = OpenMetricField{"SwapMax", Gauge, "", "", []string{"Cgroup"}}
	case "ZswapCurrent":
		cfg = OpenMetricField{"Zswap", Gauge, "", "", []string{"Cgroup"}}
	case "ZswapMax":
		cfg = OpenMetricField{"ZswapMax", Gauge, "", "", []string{"Cgroup"}}
	case "EventLow":
		cfg = OpenMetricField{"EventLow", Gauge, "", "", []string{"Cgroup"}}
	case "EventHigh":
		cfg = OpenMetricField{"EventHigh", Gauge, "", "", []string{"Cgroup"}}
	case "EventMax":
		cfg = OpenMetricField{"EventMax", Gauge, "", "", []string{"Cgroup"}}
	case "EventOom":
		cfg = OpenMetricField{"EventOom", Gauge, "", "", []string{"Cgroup"}}
	case "EventOomKill":
		cfg = OpenMetricField{"EventOomKill", Gauge, "", "", []string{"Cgroup"}}
	case "RbytePerSec":
		cfg = OpenMetricField{"RbytePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "WbytePerSec":
		cfg = OpenMetricField{"WbytePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "RioPerSec":
		cfg = OpenMetricField{"RioPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "WioPerSec":
		cfg = OpenMetricField{"WioPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "DbytePerSec":
		cfg = OpenMetricField{"DbytePerSec", Gauge, "", "", []string{"Cgroup"}}
	case "DioPerSec":
		cfg = OpenMetricField{"DioPerSec", Gauge, "", "", []string{"Cgroup"}}
	case "CPUSomePressure":
		cfg = OpenMetricField{"CPUSomePressure", Gauge, "", "", []string{"Cgroup"}}
	case "CPUFullPressure":
		cfg = OpenMetricField{"CPUFullPressure", Gauge, "", "", []string{"Cgroup"}}
	case "MemorySomePressure":
		cfg = OpenMetricField{"MemorySomePressure", Gauge, "", "", []string{"Cgroup"}}
	case "MemoryFullPressure":
		cfg = OpenMetricField{"MemoryFullPressure", Gauge, "", "", []string{"Cgroup"}}
	case "IOSomePressure":
		cfg = OpenMetricField{"IOSomePressure", Gauge, "", "", []string{"Cgroup"}}
	case "IOFullPressure":
		cfg = OpenMetricField{"IOFullPressure", Gauge, "", "", []string{"Cgroup"}}
	}
	return cfg
}

func (c *Cgroup) Collect(prev, curr *store.CgroupSample, interval int64) {
	if curr == nil {
		*c = Cgroup{
			Child: make(map[string]*Cgroup),
		}
		return
	}

	if prev == nil || curr.Name != prev.Name || curr.Inode != prev.Inode {
		prev = &store.CgroupSample{
			Child: make(map[string]*store.CgroupSample),
		}
	}

	*c = Cgroup{
		Path:                        curr.Path,
		Name:                        curr.Name,
		Level:                       curr.Level,
		Inode:                       curr.Inode,
		IsExpand:                    true,
		Child:                       make(map[string]*Cgroup),
		IsNotExist:                  curr.IsNotExist,
		Controllers:                 curr.Controllers,
		CgoupStat:                   curr.CgoupStat,
		UsagePercent:                PercentWithInterval(curr.UsageUsec, prev.UsageUsec, interval*1000000),
		UserPercent:                 PercentWithInterval(curr.UserUsec, prev.UserUsec, interval*1000000),
		SystemPercent:               PercentWithInterval(curr.SystemUsec, prev.SystemUsec, interval*1000000),
		NrPeriodsPerSec:             SubWithInterval(curr.NrPeriods, prev.NrPeriods, uint64(interval)),
		NrThrottledPerSec:           SubWithInterval(curr.NrThrottled, prev.NrThrottled, uint64(interval)),
		ThrottledPercent:            PercentWithInterval(curr.ThrottledUsec, prev.ThrottledUsec, interval*1000000),
		NrBurstsPerSec:              SubWithInterval(curr.NrBursts, prev.NrBursts, uint64(interval)),
		BurstPercent:                PercentWithInterval(curr.BurstUsec, prev.BurstUsec, interval*1000000),
		Anon:                        curr.Anon,
		File:                        curr.File,
		KernelStack:                 curr.KernelStack,
		Slab:                        curr.Slab,
		Sock:                        curr.Sock,
		Shmem:                       curr.Shmem,
		Zswap:                       curr.Zswap,
		Zswapped:                    curr.Zswapped,
		FileMapped:                  curr.FileMapped,
		FileDirty:                   curr.FileDirty,
		FileWriteback:               curr.FileWriteback,
		AnonThp:                     curr.AnonThp,
		InactiveAnon:                curr.InactiveAnon,
		ActiveAnon:                  curr.ActiveAnon,
		InactiveFile:                curr.InactiveFile,
		ActiveFile:                  curr.ActiveFile,
		Unevictable:                 curr.Unevictable,
		SlabReclaimable:             curr.SlabReclaimable,
		SlabUnreclaimable:           curr.SlabUnreclaimable,
		PgfaultPerSec:               SubWithInterval(curr.Pgfault, prev.Pgfault, uint64(interval)),
		PgmajfaultPerSec:            SubWithInterval(curr.Pgmajfault, prev.Pgmajfault, uint64(interval)),
		WorkingsetRefaultPerSec:     SubWithInterval(curr.WorkingsetRefault, prev.WorkingsetRefault, uint64(interval)),
		WorkingsetActivatePerSec:    SubWithInterval(curr.WorkingsetActivate, prev.WorkingsetActivate, uint64(interval)),
		WorkingsetNodereclaimPerSec: SubWithInterval(curr.WorkingsetNodereclaim, prev.WorkingsetNodereclaim, uint64(interval)),
		PgrefillPerSec:              SubWithInterval(curr.Pgrefill, prev.Pgrefill, uint64(interval)),
		PgscanPerSec:                SubWithInterval(curr.Pgscan, prev.Pgscan, uint64(interval)),
		PgstealPerSec:               SubWithInterval(curr.Pgsteal, prev.Pgsteal, uint64(interval)),
		PgactivatePerSec:            SubWithInterval(curr.Pgactivate, prev.Pgactivate, uint64(interval)),
		PgdeactivatePerSec:          SubWithInterval(curr.Pgdeactivate, prev.Pgdeactivate, uint64(interval)),
		PglazyfreePerSec:            SubWithInterval(curr.Pglazyfree, prev.Pglazyfree, uint64(interval)),
		PglazyfreedPerSec:           SubWithInterval(curr.Pglazyfreed, prev.Pglazyfreed, uint64(interval)),
		ZswpInPerSec:                SubWithInterval(curr.ZswpIn, prev.ZswpIn, uint64(interval)),
		ZswpOutPerSec:               SubWithInterval(curr.ZswpOut, prev.ZswpOut, uint64(interval)),
		ThpFaultAllocPerSec:         SubWithInterval(curr.ThpFaultAlloc, prev.ThpFaultAlloc, uint64(interval)),
		ThpCollapseAllocPerSec:      SubWithInterval(curr.ThpCollapseAlloc, prev.ThpCollapseAlloc, uint64(interval)),
		CpuSetCpus:                  curr.CpuSetCpus,
		CpuSetCpusEffective:         curr.CpuSetCpusEffective,
		CpuSetMems:                  curr.CpuSetMems,
		CpuSetMemsEffective:         curr.CpuSetMemsEffective,
		CpuWeight:                   curr.CpuWeight,
		CpuMax:                      curr.CpuMax,
		MemoryCurrent:               curr.MemoryCurrent,
		MemoryLow:                   curr.MemoryLow,
		MemoryHigh:                  curr.MemoryHigh,
		MemoryMin:                   curr.MemoryMin,
		MemoryMax:                   curr.MemoryMax,
		MemoryPeak:                  curr.MemoryPeak,
		SwapCurrent:                 curr.SwapCurrent,
		SwapMax:                     curr.SwapMax,
		ZswapCurrent:                curr.ZswapCurrent,
		ZswapMax:                    curr.ZswapMax,
		EventLow:                    curr.MemoryEvents.Low - prev.MemoryEvents.Low,
		EventHigh:                   curr.MemoryEvents.High - prev.MemoryEvents.High,
		EventMax:                    curr.MemoryEvents.Max - prev.MemoryEvents.Max,
		EventOom:                    curr.MemoryEvents.Oom - prev.MemoryEvents.Oom,
		EventOomKill:                curr.MemoryEvents.OomKill - prev.MemoryEvents.OomKill,
		CPUSomePressure:             curr.CpuPressure.Some.Avg60,
		CPUFullPressure:             curr.CpuPressure.Full.Avg60,
		MemorySomePressure:          curr.MemoryPressure.Some.Avg60,
		MemoryFullPressure:          curr.MemoryPressure.Full.Avg60,
		IOSomePressure:              curr.IOPressure.Some.Avg60,
		IOFullPressure:              curr.IOPressure.Full.Avg60,
	}

	var currRbyte, currWbyte, currRio, currWio, currDbyte, currDio uint64
	for _, line := range curr.IOStats {
		currRbyte += line.Rbytes
		currWbyte += line.Wbytes
		currRio += line.Rios
		currWio += line.Wios
		currDbyte += line.Dbytes
		currDio += line.Dios
	}

	for _, line := range prev.IOStats {
		currRbyte -= line.Rbytes
		currWbyte -= line.Wbytes
		currRio -= line.Rios
		currWio -= line.Wios
		currDbyte -= line.Dbytes
		currDio -= line.Dios
	}

	c.RbytePerSec = float64(currRbyte) / float64(interval)
	c.WbytePerSec = float64(currWbyte) / float64(interval)
	c.RioPerSec = float64(currRio) / float64(interval)
	c.WioPerSec = float64(currWio) / float64(interval)
	c.DbytePerSec = float64(currDbyte) / float64(interval)
	c.DioPerSec = float64(currDio) / float64(interval)

	for _, currChild := range curr.Child {
		prevChild := prev.Child[currChild.Name]

		child := &Cgroup{
			Child: make(map[string]*Cgroup),
		}
		child.Collect(prevChild, currChild, interval)
		c.Child[child.Name] = child
	}
	return
}

func (c *Cgroup) GetChildCgroupByNames(names []string) *Cgroup {
	if len(names) == 0 {
		return c
	}

	name := names[0]
	for _, child := range c.Child {
		if child.Name == name {
			return child.GetChildCgroupByNames(names[1:])
		}
	}
	return nil
}

func (c *Cgroup) Iterate(searchprogram *vm.Program, sortField string, descOrder bool) []*Cgroup {

	isMatch := true

	if searchprogram != nil {
		output, _ := expr.Run(searchprogram, c)
		if output.(bool) == false {
			isMatch = false
		}
	}

	childs := c.sortChild(sortField, descOrder)
	childReSult := []*Cgroup{}
	for _, c := range childs {
		if isMatch {
			childReSult = append(childReSult, c.Iterate(nil, sortField, descOrder)...)
		} else {
			childReSult = append(childReSult, c.Iterate(searchprogram, sortField, descOrder)...)
		}
	}

	if isMatch == false && len(childReSult) == 0 {
		return []*Cgroup{}
	}
	cgs := []*Cgroup{c}
	if c.IsExpand == false {
		return cgs
	}
	cgs = append(cgs, childReSult...)
	return cgs
}

func (c *Cgroup) sortChild(sortField string, descOrder bool) []*Cgroup {
	childs := []*Cgroup{}
	for _, child := range c.Child {
		childs = append(childs, child)
	}

	sort.SliceStable(childs, func(i, j int) bool {
		return childs[i].Name < childs[j].Name
	})

	sort.SliceStable(childs, func(i, j int) bool {
		switch sortField {
		case "Name":
			return childs[i].Name > childs[j].Name
		case "Level":
			return childs[i].Level > childs[j].Level
		case "Inode":
			return childs[i].Inode > childs[j].Inode
		case "Controllers":
			return childs[i].Controllers > childs[j].Controllers
		case "NrDescendants":
			return childs[i].NrDescendants > childs[j].NrDescendants
		case "NrDyingDescendants":
			return childs[i].NrDyingDescendants > childs[j].NrDyingDescendants
		case "UsagePercent":
			return childs[i].UsagePercent > childs[j].UsagePercent
		case "UserPercent":
			return childs[i].UserPercent > childs[j].UserPercent
		case "SystemPercent":
			return childs[i].SystemPercent > childs[j].SystemPercent
		case "NrPeriodsPerSec":
			return childs[i].NrPeriodsPerSec > childs[j].NrPeriodsPerSec
		case "NrThrottledPerSec":
			return childs[i].NrThrottledPerSec > childs[j].NrThrottledPerSec
		case "ThrottledPercent":
			return childs[i].ThrottledPercent > childs[j].ThrottledPercent
		case "NrBurstsPerSec":
			return childs[i].NrBurstsPerSec > childs[j].NrBurstsPerSec
		case "BurstPercent":
			return childs[i].BurstPercent > childs[j].BurstPercent
		case "Anon":
			return childs[i].Anon > childs[j].Anon
		case "File":
			return childs[i].File > childs[j].File
		case "KernelStack":
			return childs[i].KernelStack > childs[j].KernelStack
		case "Slab":
			return childs[i].Slab > childs[j].Slab
		case "Sock":
			return childs[i].Sock > childs[j].Sock
		case "Shmem":
			return childs[i].Shmem > childs[j].Shmem
		case "Zswap":
			return childs[i].Zswap > childs[j].Zswap
		case "Zswapped":
			return childs[i].Zswapped > childs[j].Zswapped
		case "FileMapped":
			return childs[i].FileMapped > childs[j].FileMapped
		case "FileDirty":
			return childs[i].FileDirty > childs[j].FileDirty
		case "FileWriteback":
			return childs[i].FileWriteback > childs[j].FileWriteback
		case "AnonThp":
			return childs[i].AnonThp > childs[j].AnonThp
		case "InactiveAnon":
			return childs[i].InactiveAnon > childs[j].InactiveAnon
		case "ActiveAnon":
			return childs[i].ActiveAnon > childs[j].ActiveAnon
		case "InactiveFile":
			return childs[i].InactiveFile > childs[j].InactiveFile
		case "ActiveFile":
			return childs[i].ActiveFile > childs[j].ActiveFile
		case "Unevictable":
			return childs[i].Unevictable > childs[j].Unevictable
		case "SlabReclaimable":
			return childs[i].SlabReclaimable > childs[j].SlabReclaimable
		case "SlabUnreclaimable":
			return childs[i].SlabUnreclaimable > childs[j].SlabUnreclaimable
		case "PgfaultPerSec":
			return childs[i].PgfaultPerSec > childs[j].PgfaultPerSec
		case "PgmajfaultPerSec":
			return childs[i].PgmajfaultPerSec > childs[j].PgmajfaultPerSec
		case "WorkingsetRefaultPerSec":
			return childs[i].WorkingsetRefaultPerSec > childs[j].WorkingsetRefaultPerSec
		case "WorkingsetActivatePerSec":
			return childs[i].WorkingsetActivatePerSec > childs[j].WorkingsetActivatePerSec
		case "WorkingsetNodereclaimPerSec":
			return childs[i].WorkingsetNodereclaimPerSec > childs[j].WorkingsetNodereclaimPerSec
		case "PgrefillPerSec":
			return childs[i].PgrefillPerSec > childs[j].PgrefillPerSec
		case "PgscanPerSec":
			return childs[i].PgscanPerSec > childs[j].PgscanPerSec
		case "PgstealPerSec":
			return childs[i].PgstealPerSec > childs[j].PgstealPerSec
		case "PgactivatePerSec":
			return childs[i].PgactivatePerSec > childs[j].PgactivatePerSec
		case "PgdeactivatePerSec":
			return childs[i].PgdeactivatePerSec > childs[j].PgdeactivatePerSec
		case "PglazyfreePerSec":
			return childs[i].PglazyfreePerSec > childs[j].PglazyfreePerSec
		case "PglazyfreedPerSec":
			return childs[i].PglazyfreedPerSec > childs[j].PglazyfreedPerSec
		case "ZswpInPerSec":
			return childs[i].ZswpInPerSec > childs[j].ZswpInPerSec
		case "ZswpOutPerSec":
			return childs[i].ZswpOutPerSec > childs[j].ZswpOutPerSec
		case "ThpFaultAllocPerSec":
			return childs[i].ThpFaultAllocPerSec > childs[j].ThpFaultAllocPerSec
		case "ThpCollapseAllocPerSec":
			return childs[i].ThpCollapseAllocPerSec > childs[j].ThpCollapseAllocPerSec
		case "CpuSetCpus":
			return childs[i].CpuSetCpus > childs[j].CpuSetCpus
		case "CpuSetCpusEffective":
			return childs[i].CpuSetCpusEffective > childs[j].CpuSetCpusEffective
		case "CpuSetMems":
			return childs[i].CpuSetMems > childs[j].CpuSetMems
		case "CpuSetMemsEffective":
			return childs[i].CpuSetMemsEffective > childs[j].CpuSetMemsEffective
		case "CpuWeight":
			return childs[i].CpuWeight > childs[j].CpuWeight
		case "CpuMax":
			return childs[i].CpuMax > childs[j].CpuMax
		case "MemoryCurrent":
			return childs[i].MemoryCurrent > childs[j].MemoryCurrent
		case "MemoryLow":
			return childs[i].MemoryLow > childs[j].MemoryLow
		case "MemoryHigh":
			return childs[i].MemoryHigh > childs[j].MemoryHigh
		case "MemoryMin":
			return childs[i].MemoryMin > childs[j].MemoryMin
		case "MemoryMax":
			return childs[i].MemoryMax > childs[j].MemoryMax
		case "MemoryPeak":
			return childs[i].MemoryPeak > childs[j].MemoryPeak
		case "SwapCurrent":
			return childs[i].SwapCurrent > childs[j].SwapCurrent
		case "SwapMax":
			return childs[i].SwapMax > childs[j].SwapMax
		case "ZswapCurrent":
			return childs[i].ZswapCurrent > childs[j].ZswapCurrent
		case "ZswapMax":
			return childs[i].ZswapMax > childs[j].ZswapMax
		case "EventLow":
			return childs[i].EventLow > childs[j].EventLow
		case "EventHigh":
			return childs[i].EventHigh > childs[j].EventHigh
		case "EventMax":
			return childs[i].EventMax > childs[j].EventMax
		case "EventOom":
			return childs[i].EventOom > childs[j].EventOom
		case "EventOomKill":
			return childs[i].EventOomKill > childs[j].EventOomKill
		case "RbytePerSec":
			return childs[i].RbytePerSec > childs[j].RbytePerSec
		case "WbytePerSec":
			return childs[i].WbytePerSec > childs[j].WbytePerSec
		case "RioPerSec":
			return childs[i].RioPerSec > childs[j].RioPerSec
		case "WioPerSec":
			return childs[i].WioPerSec > childs[j].WioPerSec
		case "DbytePerSec":
			return childs[i].DbytePerSec > childs[j].DbytePerSec
		case "DioPerSec":
			return childs[i].DioPerSec > childs[j].DioPerSec
		case "CPUSomePressure":
			return childs[i].CPUSomePressure > childs[j].CPUSomePressure
		case "CPUFullPressure":
			return childs[i].CPUFullPressure > childs[j].CPUFullPressure
		case "MemorySomePressure":
			return childs[i].MemorySomePressure > childs[j].MemorySomePressure
		case "MemoryFullPressure":
			return childs[i].MemoryFullPressure > childs[j].MemoryFullPressure
		case "IOSomePressure":
			return childs[i].IOSomePressure > childs[j].IOSomePressure
		case "IOFullPressure":
			return childs[i].IOFullPressure > childs[j].IOFullPressure
		}
		return false
	})

	if descOrder == false {
		for i := 0; i < len(childs)/2; i++ {
			childs[i], childs[len(childs)-1-i] = childs[len(childs)-1-i], childs[i]
		}
	}
	return childs
}
