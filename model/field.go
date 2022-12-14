package model

import (
	"fmt"

	"github.com/xixiliguo/etop/util"
)

type Format int

const (
	Raw Format = iota
	HumanReadableSize
)

type Field struct {
	Name      string
	Format    Format
	Precision int
	Suffix    string
	Width     int64
	FixWidth  bool
}

func (f Field) Render(value any) string {
	s := ""
	switch v := value.(type) {
	case uint64:
		if f.Format == HumanReadableSize {
			s = util.GetHumanSize(v)
		} else {
			s = fmt.Sprintf("%d", v)
		}
	case int:
		if f.Format == HumanReadableSize {
			s = util.GetHumanSize(v)
		} else {
			s = fmt.Sprintf("%d", v)
		}
	case float64:
		s = fmt.Sprintf("%.*f", f.Precision, v)
	case string:
		s = v
	default:
		return fmt.Sprintf("%T is unknown type", v)
	}

	s += f.Suffix

	if f.FixWidth == true {
		s = fmt.Sprintf("%*s", f.Width, s)
	}
	return s
}

type RenderConfig map[string]Field

func (renderConfig RenderConfig) Update(s string, f Field) {
	renderConfig[s] = f
}

func (renderConfig RenderConfig) SetFixWidth(fixWidth bool) {
	for k, v := range renderConfig {
		v.FixWidth = fixWidth
		renderConfig[k] = v
	}
}

func (renderConfig RenderConfig) SetRawData() {
	for k, v := range renderConfig {
		v.Format = Raw
		v.Suffix = ""
		renderConfig[k] = v
	}
}

var DefaultRenderConfig = make(map[string]RenderConfig)

var cpuDefaultRenderConfig = make(RenderConfig)

var memDefaultRenderConfig = make(RenderConfig)

func genCPUDefaultConfig() {
	cpuDefaultRenderConfig["Index"] = Field{
		Name:      "Index",
		Format:    Raw,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}

	cpuDefaultRenderConfig["User"] = Field{
		Name:      "User",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}

	cpuDefaultRenderConfig["Nice"] = Field{
		Name:      "Nice",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["System"] = Field{
		Name:      "System",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["Idle"] = Field{
		Name:      "Idle",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["Iowait"] = Field{
		Name:      "Iowait",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["IRQ"] = Field{
		Name:      "IRQ",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["SoftIRQ"] = Field{
		Name:      "SoftIRQ",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["Steal"] = Field{
		Name:      "Steal",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}
	cpuDefaultRenderConfig["Guest"] = Field{
		Name:      "Guest",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}

	cpuDefaultRenderConfig["GuestNice"] = Field{
		Name:      "GuestNice",
		Format:    Raw,
		Precision: 1,
		Suffix:    "%",
		Width:     10,
		FixWidth:  false,
	}

}

func genMEMDefaultConfig() {
	memDefaultRenderConfig["Total"] = Field{
		Name:      "Total",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["Free"] = Field{
		Name:      "Free",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Avail"] = Field{
		Name:      "Avail",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["HSlab"] = Field{
		Name:      "Slab",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Buffer"] = Field{
		Name:      "Buffer",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Cache"] = Field{
		Name:      "Cache",
		Format:    HumanReadableSize,
		Precision: 0,
		Suffix:    "",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["MemTotal"] = Field{
		Name:      "MemTotal",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["MemFree"] = Field{
		Name:      "MemFree",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["MemAvailable"] = Field{
		Name:      "MemAvailable",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Buffers"] = Field{
		Name:      "Buffers",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Cached"] = Field{
		Name:      "Cached",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["SwapCached"] = Field{
		Name:      "SwapCached",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Active"] = Field{
		Name:      "Active",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Inactive"] = Field{
		Name:      "Inactive",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["ActiveAnon"] = Field{
		Name:      "ActiveAnon",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["InactiveAnon"] = Field{
		Name:      "InactiveAnon",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["ActiveFile"] = Field{
		Name:      "ActiveFile",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["InactiveFile"] = Field{
		Name:      "InactiveFile",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Unevictable"] = Field{
		Name:      "Unevictable",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Mlocked"] = Field{
		Name:      "Mlocked",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["SwapTotal"] = Field{
		Name:      "SwapTotal",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["SwapFree"] = Field{
		Name:      "SwapFree",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Dirty"] = Field{
		Name:      "Dirty",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Writeback"] = Field{
		Name:      "Writeback",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["AnonPages"] = Field{
		Name:      "AnonPages",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Mapped"] = Field{
		Name:      "Mapped",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["Shmem"] = Field{
		Name:      "Shmem",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Slab"] = Field{
		Name:      "Slab",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["SReclaimable"] = Field{
		Name:      "SReclaimable",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["SUnreclaim"] = Field{
		Name:      "SUnreclaim",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["KernelStack"] = Field{
		Name:      "KernelStack",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["PageTables"] = Field{
		Name:      "PageTables",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["NFSUnstable"] = Field{
		Name:      "NFSUnstable",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["Bounce"] = Field{
		Name:      "Bounce",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["WritebackTmp"] = Field{
		Name:      "WritebackTmp",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["CommitLimit"] = Field{
		Name:      "CommitLimit",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["CommittedAS"] = Field{
		Name:      "CommittedAS",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["VmallocTotal"] = Field{
		Name:      "VmallocTotal",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["VmallocUsed"] = Field{
		Name:      "VmallocUsed",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["VmallocChunk"] = Field{
		Name:      "VmallocChunk",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["HardwareCorrupted"] = Field{
		Name:      "HardwareCorrupted",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["AnonHugePages"] = Field{
		Name:      "AnonHugePages",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["ShmemHugePages"] = Field{
		Name:      "ShmemHugePages",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["ShmemPmdMapped"] = Field{
		Name:      "ShmemPmdMapped",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["CmaTotal"] = Field{
		Name:      "CmaTotal",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["CmaFree"] = Field{
		Name:      "CmaFree",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["HugePagesTotal"] = Field{
		Name:      "HugePagesTotal",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["HugePagesFree"] = Field{
		Name:      "HugePagesFree",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["HugePagesRsvd"] = Field{
		Name:      "HugePagesRsvd",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["HugePagesSurp"] = Field{
		Name:      "HugePagesSurp",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["Hugepagesize"] = Field{
		Name:      "Hugepagesize",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}

	memDefaultRenderConfig["DirectMap4k"] = Field{
		Name:      "DirectMap4k",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["DirectMap2M"] = Field{
		Name:      "DirectMap2M",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
	memDefaultRenderConfig["DirectMap1G"] = Field{
		Name:      "DirectMap1G",
		Format:    Raw,
		Precision: 0,
		Suffix:    " KB",
		Width:     10,
		FixWidth:  false,
	}
}
func init() {
	genCPUDefaultConfig()
	genMEMDefaultConfig()
	DefaultRenderConfig["cpu"] = cpuDefaultRenderConfig
	DefaultRenderConfig["memory"] = memDefaultRenderConfig
}
