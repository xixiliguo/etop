package model

import (
	"math"
	"sort"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/mattn/go-runewidth"
	"github.com/xixiliguo/etop/cgroupfs"
	"github.com/xixiliguo/etop/store"
)

var DefaultCgroupFields = []string{"Name", "NrDescendants", "NrDyingDescendants", "UsagePercent", "Controllers"}
var AllCgroupFields = []string{"Path", "Name", "Level", "Inode", "Controllers",
	"NrDescendants", "NrDyingDescendants",
	"UsagePercent", "UserPercent", "SystemPercent", "NrPeriodsPerSec", "NrThrottledPerSec", "ThrottledPercent", "NrBurstsPerSec", "BurstPercent",
	"Anon", "File", "Kernel", "KernelStack", "PageTables", "SecPageTables", "PerCPU", "Sock", "Vmalloc", "Shmem", "Zswap", "Zswapped", "FileMapped",
	"FileDirty", "FileWriteback", "SwapCached", "AnonThp", "FileThp", "ShmemThp", "InactiveAnon", "ActiveAnon", "InactiveFile", "ActiveFile", "Unevictable",
	"SlabReclaimable", "SlabUnreclaimable", "Slab", "WorkingsetRefaultPerSec", "WorkingsetActivatePerSec", "WorkingsetRestorePerSec", "WorkingsetNodereclaimPerSec",
	"PgscanPerSec", "PgstealPerSec", "PgscanKswapdPerSec", "PgscanDirectPerSec", "PgscanKhugepagedPerSec", "PgstealKswapdPerSec", "PgstealDirectPerSec", "PgstealKhugepagedPerSec",
	"PgfaultPerSec", "PgmajfaultPerSec", "PgrefillPerSec", "PgactivatePerSec", "PgdeactivatePerSec", "PglazyfreePerSec", "PglazyfreedPerSec",
	"ZswpInPerSec", "ZswpOutPerSec", "ZswpWbPerSec", "ThpFaultAllocPerSec", "ThpCollapseAllocPerSec",
	"EventLowPerSec", "EventHighPerSec", "EventMaxPerSec", "EventOomPerSec", "EventOomKillPerSec", "EventOomGroupKillPerSec", "EventSockThrottledPerSec",
	"RbytePerSec", "WbytePerSec", "RioPerSec", "WioPerSec", "DbytePerSec", "DioPerSec",
	"CPUSomePressure", "CPUFullPressure", "MemorySomePressure", "MemoryFullPressure", "IOSomePressure", "IOFullPressure",
	"MemoryCurrent", "MemoryLow", "MemoryHigh", "MemoryMin", "MemoryMax", "MemoryOOMGroup", "MemorySwapCurrent", "MemorySwapMax", "MemoryZSwapCurrent", "MemoryZSwapMax",
	"CpuWeight", "CpuMax", "CpuSetCpus", "CpuSetCpusEffective", "CpuSetMems", "CpuSetMemsEffective",
	"TidsCurrent", "TidsMax",
	"RxPacketPerSec", "RxBytePerSec", "TxPacketPerSec", "TxBytePerSec"}

type Cgroup struct {
	FullPath    string
	Name        string
	Level       int
	Inode       uint64
	IsExpand    bool
	Child       map[string]*Cgroup
	Controllers string
	cgroupfs.CgoupStat
	UsagePercent                 float64 // from cpu.stat
	UserPercent                  float64
	SystemPercent                float64
	NrPeriodsPerSec              float64
	NrThrottledPerSec            float64
	ThrottledPercent             float64
	NrBurstsPerSec               float64
	BurstPercent                 float64
	Anon                         uint64 // from memory.stat
	File                         uint64
	Kernel                       uint64
	KernelStack                  uint64
	PageTables                   uint64
	SecPageTables                uint64
	PerCPU                       uint64
	Sock                         uint64
	Vmalloc                      uint64
	Shmem                        uint64
	Zswap                        uint64
	Zswapped                     uint64
	FileMapped                   uint64
	FileDirty                    uint64
	FileWriteback                uint64
	SwapCached                   uint64
	AnonThp                      uint64
	FileThp                      uint64
	ShmemThp                     uint64
	InactiveAnon                 uint64
	ActiveAnon                   uint64
	InactiveFile                 uint64
	ActiveFile                   uint64
	Unevictable                  uint64
	SlabReclaimable              uint64
	SlabUnreclaimable            uint64
	Slab                         uint64
	WorkingsetRefaultPerSec      float64
	WorkingsetActivatePerSec     float64
	WorkingsetRestorePerSec      float64
	WorkingsetNodereclaimPerSec  float64
	PgscanPerSec                 float64
	PgstealPerSec                float64
	PgscanKswapdPerSec           float64
	PgscanDirectPerSec           float64
	PgscanKhugepagedPerSec       float64
	PgstealKswapdPerSec          float64
	PgstealDirectPerSec          float64
	PgstealKhugepagedPerSec      float64
	PgfaultPerSec                float64
	PgmajfaultPerSec             float64
	PgrefillPerSec               float64
	PgactivatePerSec             float64
	PgdeactivatePerSec           float64
	PglazyfreePerSec             float64
	PglazyfreedPerSec            float64
	ZswpInPerSec                 float64
	ZswpOutPerSec                float64
	ZswpWbPerSec                 float64
	ThpFaultAllocPerSec          float64
	ThpCollapseAllocPerSec       float64
	EventLowPerSec               float64 // memory event
	EventHighPerSec              float64
	EventMaxPerSec               float64
	EventOomPerSec               float64
	EventOomKillPerSec           float64
	EventOomGroupKillPerSec      float64
	EventSockThrottledPerSec     float64
	RbytePerSec                  float64 // io.stat
	WbytePerSec                  float64
	RioPerSec                    float64
	WioPerSec                    float64
	DbytePerSec                  float64
	DioPerSec                    float64
	CPUSomePressure              float64 // pressure file
	CPUFullPressure              float64
	MemorySomePressure           float64
	MemoryFullPressure           float64
	IOSomePressure               float64
	IOFullPressure               float64
	MemoryCurrent                uint64 //Property
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
	RxPacketPerSec               float64
	RxBytePerSec                 float64
	TxPacketPerSec               float64
	TxBytePerSec                 float64
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
	case "Kernel":
		cfg = Field{"Kernel", HumanReadableSize, 1, "", 10, false}
	case "KernelStack":
		cfg = Field{"KernelStack", HumanReadableSize, 1, "", 10, false}
	case "PageTables":
		cfg = Field{"PageTables", HumanReadableSize, 1, "", 10, false}
	case "SecPageTables":
		cfg = Field{"SecPageTables", HumanReadableSize, 1, "", 10, false}
	case "PerCPU":
		cfg = Field{"PerCPU", HumanReadableSize, 1, "", 10, false}
	case "Sock":
		cfg = Field{"Sock", HumanReadableSize, 1, "", 10, false}
	case "Vmalloc":
		cfg = Field{"Vmalloc", HumanReadableSize, 1, "", 10, false}
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
	case "FileThp":
		cfg = Field{"FileThp", HumanReadableSize, 1, "", 10, false}
	case "ShmemThp":
		cfg = Field{"ShmemThp", HumanReadableSize, 1, "", 10, false}
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
	case "Slab":
		cfg = Field{"Slab", HumanReadableSize, 1, "", 10, false}
	case "WorkingsetRefaultPerSec":
		cfg = Field{"WorkingsetRefault/s", Raw, 0, "/s", 10, false}
	case "WorkingsetActivatePerSec":
		cfg = Field{"WorkingsetActivate/s", Raw, 0, "/s", 10, false}
	case "WorkingsetRestorePerSec":
		cfg = Field{"WorkingsetRestore/s", Raw, 0, "/s", 10, false}
	case "WorkingsetNodereclaimPerSec":
		cfg = Field{"WorkingsetNodereclaim/s", Raw, 0, "/s", 10, false}
	case "PgscanPerSec":
		cfg = Field{"Pgscan/s", Raw, 0, "/s", 10, false}
	case "PgstealPerSec":
		cfg = Field{"Pgsteal/s", Raw, 0, "/s", 10, false}
	case "PgscanKswapdPerSec":
		cfg = Field{"PgscanKswapd/s", Raw, 0, "/s", 10, false}
	case "PgscanDirectPerSec":
		cfg = Field{"PgscanDirect/s", Raw, 0, "/s", 10, false}
	case "PgscanKhugepagedPerSec":
		cfg = Field{"PgscanKhugepaged/s", Raw, 0, "/s", 10, false}
	case "PgstealKswapdPerSec":
		cfg = Field{"PgstealKswapd/s", Raw, 0, "/s", 10, false}
	case "PgstealDirectPerSec":
		cfg = Field{"PgstealDirect/s", Raw, 0, "/s", 10, false}
	case "PgstealKhugepagedPerSec":
		cfg = Field{"PgstealKhugepaged/s", Raw, 0, "/s", 10, false}
	case "PgfaultPerSec":
		cfg = Field{"Pgfault/s", Raw, 0, "/s", 10, false}
	case "PgmajfaultPerSec":
		cfg = Field{"Pgmajfault/s", Raw, 0, "/s", 10, false}
	case "PgrefillPerSec":
		cfg = Field{"Pgrefill/s", Raw, 0, "/s", 10, false}
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
	case "ZswpWbPerSec":
		cfg = Field{"ZswpWb/s", Raw, 0, "/s", 10, false}
	case "ThpFaultAllocPerSec":
		cfg = Field{"ThpFaultAlloc/s", Raw, 0, "/s", 10, false}
	case "ThpCollapseAllocPerSec":
		cfg = Field{"ThpCollapseAlloc/s", Raw, 0, "/s", 10, false}
	case "EventLowPerSec":
		cfg = Field{"EventLow/s", Raw, 0, "/s", 10, false}
	case "EventHighPerSec":
		cfg = Field{"EventHigh/s", Raw, 0, "/s", 10, false}
	case "EventMaxPerSec":
		cfg = Field{"EventMax/s", Raw, 0, "/s", 10, false}
	case "EventOomPerSec":
		cfg = Field{"EventOom/s", Raw, 0, "/s", 10, false}
	case "EventOomKillPerSec":
		cfg = Field{"EventOomKill/s", Raw, 0, "/s", 10, false}
	case "EventOomGroupKillPerSec":
		cfg = Field{"EventOomGroupKill/s", Raw, 0, "/s", 10, false}
	case "EventSockThrottledPerSec":
		cfg = Field{"EventSockThrottled/s", Raw, 0, "/s", 10, false}
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
	case "MemoryOOMGroup":
		cfg = Field{"MemoryOOMGroup", Raw, 1, "", 10, false}
	case "MemorySwapCurrent":
		cfg = Field{"Swap", HumanReadableSize, 1, "", 10, false}
	case "MemorySwapMax":
		cfg = Field{"SwapMax", HumanReadableSize, 1, "", 10, false}
	case "MemoryZSwapCurrent":
		cfg = Field{"Zswap", HumanReadableSize, 1, "", 10, false}
	case "MemoryZSwapMax":
		cfg = Field{"ZswapMax", HumanReadableSize, 1, "", 10, false}
	case "CpuWeight":
		cfg = Field{"CpuWeight", Raw, 0, "", 10, false}
	case "CpuMax":
		cfg = Field{"CpuMax", Raw, 0, "", 10, false}
	case "CpuSetCpus":
		cfg = Field{"CpuSetCpus", Raw, 0, "", 10, false}
	case "CpuSetCpusEffective":
		cfg = Field{"CpuSetCpusEffective", Raw, 0, "", 10, false}
	case "CpuSetCpusExclusive":
		cfg = Field{"CpuSetCpusExclusive", Raw, 0, "", 10, false}
	case "CpuSetCpusExclusiveEffective":
		cfg = Field{"CpuSetCpusExclusiveEffective", Raw, 0, "", 10, false}
	case "TidsCurrent":
		cfg = Field{"Tids", Raw, 0, "", 10, false}
	case "TidsMax":
		cfg = Field{"TidsMax", Raw, 0, "", 10, false}
	case "RxPacketPerSec":
		cfg = Field{"Rpkt/s", Raw, 1, "/s", 10, false}
	case "RxBytePerSec":
		cfg = Field{"Rbyte/s", HumanReadableSize, 1, "/s", 10, false}
	case "TxPacketPerSec":
		cfg = Field{"Tpkt/s", Raw, 1, "/s", 10, false}
	case "TxBytePerSec":
		cfg = Field{"Tbyte/s", HumanReadableSize, 1, "/s", 10, false}
	}
	return cfg
}

func (c *Cgroup) GetRenderValue(field string, opt FieldOpt) string {

	cfg := c.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""

	switch field {
	case "Cgroup", "Path":
		s = cfg.Render(c.FullPath)
	case "Name":
		indents := ""
		if c.Level > 1 {
			indents = strings.Repeat("   ", c.Level-1)
		}

		s = indents + "└─ " + c.Name
		if !c.IsExpand && len(c.Child) != 0 {
			s = indents + "└+ " + c.Name
		}
		if c.Level == 0 {
			s = "/"
		}
		if cfg.FixWidth {
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
	case "Kernel":
		s = cfg.Render(c.Kernel)
	case "KernelStack":
		s = cfg.Render(c.KernelStack)
	case "PageTables":
		s = cfg.Render(c.PageTables)
	case "SecPageTables":
		s = cfg.Render(c.SecPageTables)
	case "PerCPU":
		s = cfg.Render(c.PerCPU)
	case "Sock":
		s = cfg.Render(c.Sock)
	case "Vmalloc":
		s = cfg.Render(c.Vmalloc)
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
	case "SwapCached":
		s = cfg.Render(c.SwapCached)
	case "AnonThp":
		s = cfg.Render(c.AnonThp)
	case "FileThp":
		s = cfg.Render(c.FileThp)
	case "ShmemThp":
		s = cfg.Render(c.ShmemThp)
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
	case "Slab":
		s = cfg.Render(c.Slab)
	case "WorkingsetRefaultPerSec":
		s = cfg.Render(c.WorkingsetRefaultPerSec)
	case "WorkingsetActivatePerSec":
		s = cfg.Render(c.WorkingsetActivatePerSec)
	case "WorkingsetRestorePerSec":
		s = cfg.Render(c.WorkingsetRestorePerSec)
	case "WorkingsetNodereclaimPerSec":
		s = cfg.Render(c.WorkingsetNodereclaimPerSec)
	case "PgscanPerSec":
		s = cfg.Render(c.PgscanPerSec)
	case "PgstealPerSec":
		s = cfg.Render(c.PgstealPerSec)
	case "PgscanKswapdPerSec":
		s = cfg.Render(c.PgscanKswapdPerSec)
	case "PgscanDirectPerSec":
		s = cfg.Render(c.PgscanDirectPerSec)
	case "PgscanKhugepagedPerSec":
		s = cfg.Render(c.PgscanKhugepagedPerSec)
	case "PgstealKswapdPerSec":
		s = cfg.Render(c.PgstealKswapdPerSec)
	case "PgstealDirectPerSec":
		s = cfg.Render(c.PgstealDirectPerSec)
	case "PgstealKhugepagedPerSec":
		s = cfg.Render(c.PgstealKhugepagedPerSec)
	case "PgfaultPerSec":
		s = cfg.Render(c.PgfaultPerSec)
	case "PgmajfaultPerSec":
		s = cfg.Render(c.PgmajfaultPerSec)
	case "PgrefillPerSec":
		s = cfg.Render(c.PgrefillPerSec)
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
	case "ZswpWbPerSec":
		s = cfg.Render(c.ZswpWbPerSec)
	case "ThpFaultAllocPerSec":
		s = cfg.Render(c.ThpFaultAllocPerSec)
	case "ThpCollapseAllocPerSec":
		s = cfg.Render(c.ThpCollapseAllocPerSec)
	case "EventLowPerSec":
		s = cfg.Render(c.EventLowPerSec)
	case "EventHighPerSec":
		s = cfg.Render(c.EventHighPerSec)
	case "EventMaxPerSec":
		s = cfg.Render(c.EventMaxPerSec)
	case "EventOomPerSec":
		s = cfg.Render(c.EventOomPerSec)
	case "EventOomKillPerSec":
		s = cfg.Render(c.EventOomKillPerSec)
	case "EventOomGroupKillPerSec":
		s = cfg.Render(c.EventOomGroupKillPerSec)
	case "EventSockThrottledPerSec":
		s = cfg.Render(c.EventSockThrottledPerSec)
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
	case "MemoryOOMGroup":
		s = cfg.Render(c.MemoryOOMGroup)
	case "MemorySwapCurrent":
		s = cfg.Render(c.MemorySwapCurrent)
	case "MemorySwapMax":
		s = cfg.Render(c.MemorySwapMax)
	case "MemoryZSwapCurrent":
		s = cfg.Render(c.MemoryZSwapCurrent)
	case "MemoryZSwapMax":
		s = cfg.Render(c.MemoryZSwapMax)
	case "CpuWeight":
		s = cfg.Render(c.CpuWeight)
	case "CpuMax":
		s = cfg.Render(c.CpuMax)
	case "CpuSetCpus":
		s = cfg.Render(c.CpuSetCpus)
	case "CpuSetCpusEffective":
		s = cfg.Render(c.CpuSetCpusEffective)
	case "CpuSetCpusExclusive":
		s = cfg.Render(c.CpuSetCpusExclusive)
	case "CpuSetCpusExclusiveEffective":
		s = cfg.Render(c.CpuSetCpusExclusiveEffective)
	case "TidsCurrent":
		s = cfg.Render(c.TidsCurrent)
	case "TidsMax":
		s = cfg.Render(c.TidsMax)
	case "RxPacketPerSec":
		s = cfg.Render(c.RxPacketPerSec)
	case "RxBytePerSec":
		s = cfg.Render(c.RxBytePerSec)
	case "TxPacketPerSec":
		s = cfg.Render(c.TxPacketPerSec)
	case "TxBytePerSec":
		s = cfg.Render(c.TxBytePerSec)
	default:
		s = "no " + field + " for cgroup stat"
	}
	return s
}

func (c *Cgroup) Collect(prev, curr *store.CgroupSample, interval int64) {
	if curr == nil {
		*c = Cgroup{
			Child: make(map[string]*Cgroup),
		}
		return
	}

	if prev == nil || curr.Name != prev.Name || curr.Inode != prev.Inode {
		prev = &store.CgroupSample{}
	}

	*c = Cgroup{
		FullPath:                     curr.FullPath,
		Name:                         curr.Name,
		Level:                        curr.Level,
		Inode:                        curr.Inode,
		IsExpand:                     true,
		Child:                        make(map[string]*Cgroup),
		Controllers:                  curr.Controllers,
		CgoupStat:                    curr.CgoupStat,
		UsagePercent:                 SubWithInterval(curr.UsageUsec, prev.UsageUsec, interval*1000000) * 100,
		UserPercent:                  SubWithInterval(curr.UserUsec, prev.UserUsec, interval*1000000) * 100,
		SystemPercent:                SubWithInterval(curr.SystemUsec, prev.SystemUsec, interval*1000000) * 100,
		NrPeriodsPerSec:              SubWithInterval(curr.NrPeriods, prev.NrPeriods, interval),
		NrThrottledPerSec:            SubWithInterval(curr.NrThrottled, prev.NrThrottled, interval),
		ThrottledPercent:             SubWithInterval(curr.ThrottledUsec, prev.ThrottledUsec, interval*1000000) * 100,
		NrBurstsPerSec:               SubWithInterval(curr.NrBursts, prev.NrBursts, interval),
		BurstPercent:                 SubWithInterval(curr.BurstUsec, prev.BurstUsec, interval*1000000) * 100,
		Anon:                         curr.Anon,
		File:                         curr.File,
		Kernel:                       curr.Kernel,
		KernelStack:                  curr.KernelStack,
		PageTables:                   curr.PageTables,
		SecPageTables:                curr.SecPageTables,
		PerCPU:                       curr.PerCPU,
		Sock:                         curr.Sock,
		Vmalloc:                      curr.Vmalloc,
		Shmem:                        curr.Shmem,
		Zswap:                        curr.Zswap,
		Zswapped:                     curr.Zswapped,
		FileMapped:                   curr.FileMapped,
		FileDirty:                    curr.FileDirty,
		FileWriteback:                curr.FileWriteback,
		SwapCached:                   curr.SwapCached,
		AnonThp:                      curr.AnonThp,
		FileThp:                      curr.FileThp,
		ShmemThp:                     curr.ShmemThp,
		InactiveAnon:                 curr.InactiveAnon,
		ActiveAnon:                   curr.ActiveAnon,
		InactiveFile:                 curr.InactiveFile,
		ActiveFile:                   curr.ActiveFile,
		Unevictable:                  curr.Unevictable,
		SlabReclaimable:              curr.SlabReclaimable,
		SlabUnreclaimable:            curr.SlabUnreclaimable,
		Slab:                         curr.Slab,
		WorkingsetRefaultPerSec:      SubWithInterval(curr.WorkingsetRefaultAnon, prev.WorkingsetRefaultAnon, interval),
		WorkingsetActivatePerSec:     SubWithInterval(curr.WorkingsetActivateAnon, prev.WorkingsetActivateAnon, interval),
		WorkingsetRestorePerSec:      SubWithInterval(curr.WorkingsetRestoreAnon, prev.WorkingsetRestoreAnon, interval),
		WorkingsetNodereclaimPerSec:  SubWithInterval(curr.WorkingsetNodereclaim, prev.WorkingsetNodereclaim, interval),
		PgscanPerSec:                 SubWithInterval(curr.Pgscan, prev.Pgscan, interval),
		PgstealPerSec:                SubWithInterval(curr.Pgsteal, prev.Pgsteal, interval),
		PgscanKswapdPerSec:           SubWithInterval(curr.PgscanKswapd, prev.PgscanKswapd, interval),
		PgscanDirectPerSec:           SubWithInterval(curr.PgscanDirect, prev.PgscanDirect, interval),
		PgscanKhugepagedPerSec:       SubWithInterval(curr.PgscanKhugepaged, prev.PgscanKhugepaged, interval),
		PgstealKswapdPerSec:          SubWithInterval(curr.PgstealKswapd, prev.PgstealKswapd, interval),
		PgstealDirectPerSec:          SubWithInterval(curr.PgstealDirect, prev.PgstealDirect, interval),
		PgstealKhugepagedPerSec:      SubWithInterval(curr.PgstealKhugepaged, prev.PgstealKhugepaged, interval),
		PgfaultPerSec:                SubWithInterval(curr.Pgfault, prev.Pgfault, interval),
		PgmajfaultPerSec:             SubWithInterval(curr.Pgmajfault, prev.Pgmajfault, interval),
		PgrefillPerSec:               SubWithInterval(curr.Pgrefill, prev.Pgrefill, interval),
		PgactivatePerSec:             SubWithInterval(curr.Pgactivate, prev.Pgactivate, interval),
		PgdeactivatePerSec:           SubWithInterval(curr.Pgdeactivate, prev.Pgdeactivate, interval),
		PglazyfreePerSec:             SubWithInterval(curr.Pglazyfree, prev.Pglazyfree, interval),
		PglazyfreedPerSec:            SubWithInterval(curr.Pglazyfreed, prev.Pglazyfreed, interval),
		ZswpInPerSec:                 SubWithInterval(curr.ZswpIn, prev.ZswpIn, interval),
		ZswpWbPerSec:                 SubWithInterval(curr.ZswpOut, prev.ZswpOut, interval),
		ZswpOutPerSec:                SubWithInterval(curr.ZswpWb, prev.ZswpWb, interval),
		ThpFaultAllocPerSec:          SubWithInterval(curr.ThpFaultAlloc, prev.ThpFaultAlloc, interval),
		ThpCollapseAllocPerSec:       SubWithInterval(curr.ThpCollapseAlloc, prev.ThpCollapseAlloc, interval),
		EventLowPerSec:               SubWithInterval(curr.MemoryEvents.Low, prev.MemoryEvents.Low, interval),
		EventHighPerSec:              SubWithInterval(curr.MemoryEvents.High, prev.MemoryEvents.High, interval),
		EventMaxPerSec:               SubWithInterval(curr.MemoryEvents.Max, prev.MemoryEvents.Max, interval),
		EventOomPerSec:               SubWithInterval(curr.MemoryEvents.Oom, prev.MemoryEvents.Oom, interval),
		EventOomKillPerSec:           SubWithInterval(curr.MemoryEvents.OomKill, prev.MemoryEvents.OomKill, interval),
		EventOomGroupKillPerSec:      SubWithInterval(curr.MemoryEvents.OomKill, prev.MemoryEvents.OomKill, interval),
		EventSockThrottledPerSec:     SubWithInterval(curr.MemoryEvents.OomKill, prev.MemoryEvents.OomKill, interval),
		RbytePerSec:                  math.NaN(),
		WbytePerSec:                  math.NaN(),
		RioPerSec:                    math.NaN(),
		WioPerSec:                    math.NaN(),
		DbytePerSec:                  math.NaN(),
		DioPerSec:                    math.NaN(),
		CPUSomePressure:              curr.CpuPressure.Some.Avg60,
		CPUFullPressure:              curr.CpuPressure.Full.Avg60,
		MemorySomePressure:           curr.MemoryPressure.Some.Avg60,
		MemoryFullPressure:           curr.MemoryPressure.Full.Avg60,
		IOSomePressure:               curr.IOPressure.Some.Avg60,
		IOFullPressure:               curr.IOPressure.Full.Avg60,
		MemoryCurrent:                curr.MemoryCurrent,
		MemoryLow:                    curr.MemoryLow,
		MemoryHigh:                   curr.MemoryHigh,
		MemoryMin:                    curr.MemoryMin,
		MemoryMax:                    curr.MemoryMax,
		MemoryOOMGroup:               curr.MemoryOOMGroup,
		MemorySwapCurrent:            curr.MemorySwapCurrent,
		MemorySwapMax:                curr.MemorySwapMax,
		MemoryZSwapCurrent:           curr.MemoryZSwapCurrent,
		MemoryZSwapMax:               curr.MemoryZSwapMax,
		CpuWeight:                    curr.CpuWeight,
		CpuMax:                       curr.CpuMax,
		CpuSetCpus:                   curr.CpuSetCpus,
		CpuSetCpusEffective:          curr.CpuSetCpusEffective,
		CpuSetCpusExclusive:          curr.CpuSetCpusExclusive,
		CpuSetCpusExclusiveEffective: curr.CpuSetCpusExclusiveEffective,
		TidsCurrent:                  curr.TidsCurrent,
		TidsMax:                      curr.TidsMax,
		RxPacketPerSec:               math.NaN(),
		RxBytePerSec:                 math.NaN(),
		TxPacketPerSec:               math.NaN(),
		TxBytePerSec:                 math.NaN(),
	}

	if len(curr.IOStats) > 0 {
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
	}

	if curr.RxByte != math.MaxUint64 {
		c.RxPacketPerSec = float64(curr.RxPacket-prev.RxPacket) / float64(interval)
		c.RxBytePerSec = float64(curr.RxByte-prev.RxByte) / float64(interval)
		c.TxPacketPerSec = float64(curr.TxPacket-prev.TxPacket) / float64(interval)
		c.TxBytePerSec = float64(curr.TxByte-prev.TxByte) / float64(interval)
	}

	for _, currChild := range curr.Child {
		prevChild := prev.Child[currChild.Name]

		child := &Cgroup{
			Child: make(map[string]*Cgroup),
		}
		child.Collect(&prevChild, &currChild, interval)
		c.Child[child.Name] = child
	}
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
		if !output.(bool) {
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

	if !isMatch && len(childReSult) == 0 {
		return []*Cgroup{}
	}
	cgs := []*Cgroup{c}
	if !c.IsExpand {
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
		case "Kernel":
			return childs[i].Kernel > childs[j].Kernel
		case "KernelStack":
			return childs[i].KernelStack > childs[j].KernelStack
		case "PageTables":
			return childs[i].PageTables > childs[j].PageTables
		case "SecPageTables":
			return childs[i].SecPageTables > childs[j].SecPageTables
		case "PerCPU":
			return childs[i].PerCPU > childs[j].PerCPU
		case "Sock":
			return childs[i].Sock > childs[j].Sock
		case "Vmalloc":
			return childs[i].Vmalloc > childs[j].Vmalloc
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
		case "SwapCached":
			return childs[i].SwapCached > childs[j].SwapCached
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
		case "Slab":
			return childs[i].Slab > childs[j].Slab
		case "WorkingsetRefaultPerSec":
			return childs[i].WorkingsetRefaultPerSec > childs[j].WorkingsetRefaultPerSec
		case "WorkingsetActivatePerSec":
			return childs[i].WorkingsetActivatePerSec > childs[j].WorkingsetActivatePerSec
		case "WorkingsetRestorePerSec":
			return childs[i].WorkingsetRestorePerSec > childs[j].WorkingsetRestorePerSec
		case "WorkingsetNodereclaimPerSec":
			return childs[i].WorkingsetNodereclaimPerSec > childs[j].WorkingsetNodereclaimPerSec
		case "PgscanPerSec":
			return childs[i].PgscanPerSec > childs[j].PgscanPerSec
		case "PgstealPerSec":
			return childs[i].PgstealPerSec > childs[j].PgstealPerSec
		case "PgscanKswapdPerSec":
			return childs[i].PgscanKswapdPerSec > childs[j].PgscanKswapdPerSec
		case "PgscanDirectPerSec":
			return childs[i].PgscanDirectPerSec > childs[j].PgscanDirectPerSec
		case "PgscanKhugepagedPerSec":
			return childs[i].PgscanKhugepagedPerSec > childs[j].PgscanKhugepagedPerSec
		case "PgstealKswapdPerSec":
			return childs[i].PgstealKswapdPerSec > childs[j].PgstealKswapdPerSec
		case "PgstealDirectPerSec":
			return childs[i].PgstealDirectPerSec > childs[j].PgstealDirectPerSec
		case "PgstealKhugepagedPerSec":
			return childs[i].PgstealKhugepagedPerSec > childs[j].PgstealKhugepagedPerSec
		case "PgfaultPerSec":
			return childs[i].PgfaultPerSec > childs[j].PgfaultPerSec
		case "PgmajfaultPerSec":
			return childs[i].PgmajfaultPerSec > childs[j].PgmajfaultPerSec
		case "PgrefillPerSec":
			return childs[i].PgrefillPerSec > childs[j].PgrefillPerSec
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
		case "ZswpWbPerSec":
			return childs[i].ZswpWbPerSec > childs[j].ZswpWbPerSec
		case "ThpFaultAllocPerSec":
			return childs[i].ThpFaultAllocPerSec > childs[j].ThpFaultAllocPerSec
		case "ThpCollapseAllocPerSec":
			return childs[i].ThpCollapseAllocPerSec > childs[j].ThpCollapseAllocPerSec
		case "EventLowPerSec":
			return childs[i].EventLowPerSec > childs[j].EventLowPerSec
		case "EventHighPerSec":
			return childs[i].EventHighPerSec > childs[j].EventHighPerSec
		case "EventMaxPerSec":
			return childs[i].EventMaxPerSec > childs[j].EventMaxPerSec
		case "EventOomPerSec":
			return childs[i].EventOomPerSec > childs[j].EventOomPerSec
		case "EventOomKillPerSec":
			return childs[i].EventOomKillPerSec > childs[j].EventOomKillPerSec
		case "EventOomGroupKillPerSec":
			return childs[i].EventOomGroupKillPerSec > childs[j].EventOomGroupKillPerSec
		case "EventSockThrottledPerSec":
			return childs[i].EventSockThrottledPerSec > childs[j].EventSockThrottledPerSec
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
		case "MemoryOOMGroup":
			return childs[i].MemoryOOMGroup > childs[j].MemoryOOMGroup
		case "MemorySwapCurrent":
			return childs[i].MemorySwapCurrent > childs[j].MemorySwapCurrent
		case "MemorySwapMax":
			return childs[i].MemorySwapMax > childs[j].MemorySwapMax
		case "MemoryZSwapCurrent":
			return childs[i].MemoryZSwapCurrent > childs[j].MemoryZSwapCurrent
		case "MemoryZSwapMax":
			return childs[i].MemoryZSwapMax > childs[j].MemoryZSwapMax
		case "CpuWeight":
			return childs[i].CpuWeight > childs[j].CpuWeight
		case "CpuMax":
			return childs[i].CpuMax > childs[j].CpuMax
		case "CpuSetCpus":
			return childs[i].CpuSetCpus > childs[j].CpuSetCpus
		case "CpuSetCpusEffective":
			return childs[i].CpuSetCpusEffective > childs[j].CpuSetCpusEffective
		case "CpuSetCpusExclusive":
			return childs[i].CpuSetCpusExclusive > childs[j].CpuSetCpusExclusive
		case "CpuSetCpusExclusiveEffective":
			return childs[i].CpuSetCpusExclusiveEffective > childs[j].CpuSetCpusExclusiveEffective
		case "TidsCurrent":
			return childs[i].TidsCurrent > childs[j].TidsCurrent
		case "TidsMax":
			return childs[i].TidsMax > childs[j].TidsMax
		case "RxPacketPerSec":
			return childs[i].RxPacketPerSec > childs[j].RxPacketPerSec
		case "RxBytePerSec":
			return childs[i].RxBytePerSec > childs[j].RxBytePerSec
		case "TxPacketPerSec":
			return childs[i].TxPacketPerSec > childs[j].TxPacketPerSec
		case "TxBytePerSec":
			return childs[i].TxBytePerSec > childs[j].TxBytePerSec
		}
		return false
	})

	if !descOrder {
		for i := 0; i < len(childs)/2; i++ {
			childs[i], childs[len(childs)-1-i] = childs[len(childs)-1-i], childs[i]
		}
	}
	return childs
}
