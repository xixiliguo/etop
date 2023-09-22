package model

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

var typToString = [2]string{"gauge", "counter"}

type OpenMetricField struct {
	Name   string
	Typ    MetricType
	Unit   string
	Help   string
	Labels []string
}

type OpenMetricRenderConfig map[string]OpenMetricField

var DefaultOMRenderConfig = make(map[string]OpenMetricRenderConfig)

var sysDefaultOMRenderConfig = make(OpenMetricRenderConfig)

var cpuDefaultOMRenderConfig = make(OpenMetricRenderConfig)

var memDefaultOMRenderConfig = make(OpenMetricRenderConfig)

var vmDefaultOMRenderConfig = make(OpenMetricRenderConfig)

var diskDefaultOMRenderConfig = make(OpenMetricRenderConfig)

func genSysDefaultOMConfig() {

	sysDefaultOMRenderConfig["Load1"] = OpenMetricField{
		Name:   "Load1",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Load5"] = OpenMetricField{
		Name:   "Load5",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Load15"] = OpenMetricField{
		Name:   "Load15",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Processes"] = OpenMetricField{
		Name:   "Processes",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Threads"] = OpenMetricField{
		Name:   "Threads",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ProcessesRunning"] = OpenMetricField{
		Name:   "Running",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ProcessesBlocked"] = OpenMetricField{
		Name:   "Blocked",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	sysDefaultOMRenderConfig["ClonePerSec"] = OpenMetricField{
		Name:   "ClonePerSec",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ContextSwitchPerSec"] = OpenMetricField{
		Name:   "ContextSwitchPerSec",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
}

func genCPUDefaultOMConfig() {
	cpuDefaultOMRenderConfig["User"] = OpenMetricField{
		Name:   "User",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["Nice"] = OpenMetricField{
		Name:   "Nice",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["System"] = OpenMetricField{
		Name:   "System",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["Idle"] = OpenMetricField{
		Name:   "Idle",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["Iowait"] = OpenMetricField{
		Name:   "Iowait",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["IRQ"] = OpenMetricField{
		Name:   "IRQ",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["SoftIRQ"] = OpenMetricField{
		Name:   "SoftIRQ",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["IRQ"] = OpenMetricField{
		Name:   "IRQ",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["Steal"] = OpenMetricField{
		Name:   "Steal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["Guest"] = OpenMetricField{
		Name:   "Guest",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}
	cpuDefaultOMRenderConfig["GuestNice"] = OpenMetricField{
		Name:   "GuestNice",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Index"},
	}

}

func genMEMDefaultOMConfig() {

	memDefaultOMRenderConfig["Total"] = OpenMetricField{
		Name:   "Total",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Free"] = OpenMetricField{
		Name:   "Free",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Avail"] = OpenMetricField{
		Name:   "Avail",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HSlab"] = OpenMetricField{
		Name:   "HSlab",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Buffer"] = OpenMetricField{
		Name:   "Buffer",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Cache"] = OpenMetricField{
		Name:   "Cache",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["MemTotal"] = OpenMetricField{
		Name:   "MemTotal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["MemFree"] = OpenMetricField{
		Name:   "MemFree",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["MemAvailable"] = OpenMetricField{
		Name:   "MemAvailable",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Buffers"] = OpenMetricField{
		Name:   "Buffers",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Cached"] = OpenMetricField{
		Name:   "Cached",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["SwapCached"] = OpenMetricField{
		Name:   "SwapCached",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Active"] = OpenMetricField{
		Name:   "Active",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Inactive"] = OpenMetricField{
		Name:   "Inactive",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["ActiveAnon"] = OpenMetricField{
		Name:   "ActiveAnon",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["InactiveAnon"] = OpenMetricField{
		Name:   "InactiveAnon",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["ActiveFile"] = OpenMetricField{
		Name:   "ActiveFile",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["InactiveFile"] = OpenMetricField{
		Name:   "InactiveFile",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Unevictable"] = OpenMetricField{
		Name:   "Unevictable",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Mlocked"] = OpenMetricField{
		Name:   "Mlocked",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["SwapTotal"] = OpenMetricField{
		Name:   "SwapTotal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["SwapFree"] = OpenMetricField{
		Name:   "SwapFree",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Dirty"] = OpenMetricField{
		Name:   "Dirty",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["Writeback"] = OpenMetricField{
		Name:   "Writeback",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["AnonPages"] = OpenMetricField{
		Name:   "AnonPages",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Mapped"] = OpenMetricField{
		Name:   "Mapped",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["Shmem"] = OpenMetricField{
		Name:   "Shmem",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Slab"] = OpenMetricField{
		Name:   "Slab",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["SReclaimable"] = OpenMetricField{
		Name:   "SReclaimable",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["SUnreclaim"] = OpenMetricField{
		Name:   "SUnreclaim",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["KernelStack"] = OpenMetricField{
		Name:   "KernelStack",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["PageTables"] = OpenMetricField{
		Name:   "PageTables",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["NFSUnstable"] = OpenMetricField{
		Name:   "NFSUnstable",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	memDefaultOMRenderConfig["Bounce"] = OpenMetricField{
		Name:   "Bounce",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["WritebackTmp"] = OpenMetricField{
		Name:   "WritebackTmp",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["CommitLimit"] = OpenMetricField{
		Name:   "CommitLimit",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["CommittedAS"] = OpenMetricField{
		Name:   "CommittedAS",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["VmallocTotal"] = OpenMetricField{
		Name:   "VmallocTotal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["VmallocUsed"] = OpenMetricField{
		Name:   "VmallocUsed",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["VmallocChunk"] = OpenMetricField{
		Name:   "VmallocChunk",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HardwareCorrupted"] = OpenMetricField{
		Name:   "HardwareCorrupted",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["AnonHugePages"] = OpenMetricField{
		Name:   "AnonHugePages",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["ShmemHugePages"] = OpenMetricField{
		Name:   "ShmemHugePages",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["ShmemPmdMapped"] = OpenMetricField{
		Name:   "ShmemPmdMapped",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["CmaTotal"] = OpenMetricField{
		Name:   "CmaTotal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["CmaFree"] = OpenMetricField{
		Name:   "CmaFree",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HugePagesTotal"] = OpenMetricField{
		Name:   "HugePagesTotal",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HugePagesFree"] = OpenMetricField{
		Name:   "HugePagesFree",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HugePagesRsvd"] = OpenMetricField{
		Name:   "HugePagesRsvd",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["HugePagesSurp"] = OpenMetricField{
		Name:   "HugePagesSurp",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["Hugepagesize"] = OpenMetricField{
		Name:   "Hugepagesize",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["DirectMap4k"] = OpenMetricField{
		Name:   "DirectMap4k",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["DirectMap2M"] = OpenMetricField{
		Name:   "DirectMap2M",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	memDefaultOMRenderConfig["DirectMap1G"] = OpenMetricField{
		Name:   "DirectMap1G",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
}

func genVmDefaultOMConfig() {

	vmDefaultOMRenderConfig["PageIn"] = OpenMetricField{
		Name:   "PageIn",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	vmDefaultOMRenderConfig["PageOut"] = OpenMetricField{
		Name:   "PageOut",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["SwapIn"] = OpenMetricField{
		Name:   "SwapIn",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["SwapOut"] = OpenMetricField{
		Name:   "SwapOut",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["PageScanKswapd"] = OpenMetricField{
		Name:   "PageScanKswapd",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["PageScanDirect"] = OpenMetricField{
		Name:   "PageScanDirect",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	vmDefaultOMRenderConfig["PageStealKswapd"] = OpenMetricField{
		Name:   "PageStealKswapd",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["PageStealDirect"] = OpenMetricField{
		Name:   "PageStealDirect",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	vmDefaultOMRenderConfig["OOMKill"] = OpenMetricField{
		Name:   "OOMKill",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
}

func genDiskDefaultOMConfig() {

	diskDefaultOMRenderConfig["Util"] = OpenMetricField{
		Name:   "Util",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}

	diskDefaultOMRenderConfig["Read"] = OpenMetricField{
		Name:   "Read",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["Read/s"] = OpenMetricField{
		Name:   "Read/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["ReadByte/s"] = OpenMetricField{
		Name:   "ReadByte/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["Write"] = OpenMetricField{
		Name:   "Write",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["Write/s"] = OpenMetricField{
		Name:   "Write/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["WriteByte/s"] = OpenMetricField{
		Name:   "WriteByte/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["Discard"] = OpenMetricField{
		Name:   "Discard",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["Discard/s"] = OpenMetricField{
		Name:   "Discard/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["DiscardByte/s"] = OpenMetricField{
		Name:   "DiscardByte/s",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["AvgIOSize"] = OpenMetricField{
		Name:   "AvgIOSize",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["AvgQueueLen"] = OpenMetricField{
		Name:   "AvgQueueLen",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["InFlight"] = OpenMetricField{
		Name:   "InFlight",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["AvgIOWait"] = OpenMetricField{
		Name:   "AvgIOWait",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}
	diskDefaultOMRenderConfig["AvgIOTime"] = OpenMetricField{
		Name:   "AvgIOTime",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{"Disk"},
	}

}

func init() {
	genSysDefaultOMConfig()
	genCPUDefaultOMConfig()
	genMEMDefaultOMConfig()
	genVmDefaultOMConfig()
	genDiskDefaultOMConfig()
	DefaultOMRenderConfig["system"] = sysDefaultOMRenderConfig
	DefaultOMRenderConfig["cpu"] = cpuDefaultOMRenderConfig
	DefaultOMRenderConfig["memory"] = memDefaultOMRenderConfig
	DefaultOMRenderConfig["vm"] = vmDefaultOMRenderConfig
	DefaultOMRenderConfig["disk"] = diskDefaultOMRenderConfig
}
