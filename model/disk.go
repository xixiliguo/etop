package model

import (
	"fmt"

	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/util"
)

type Disk struct {
	DeviceName       string
	ReadIOs          uint64
	ReadMerges       uint64
	ReadSectors      uint64
	ReadTicks        uint64
	WriteIOs         uint64
	WriteMerges      uint64
	WriteSectors     uint64
	WriteTicks       uint64
	IOsInProgress    uint64
	IOsTotalTicks    uint64
	WeightedIOTicks  uint64
	ReadBytesPerSec  uint64
	WriteBytesPerSec uint64
	Await            float64
	Avio             float64
	Busy             uint64
}

type DiskMap map[string]Disk

func (d *Disk) GetRenderValue(field string) string {
	switch field {
	case "Disk":
		return fmt.Sprintf("%s", d.DeviceName)
	case "Busy":
		return fmt.Sprintf("%d%%", d.Busy)
	case "Read":
		return fmt.Sprintf("%d", d.ReadIOs)
	case "R/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(d.ReadBytesPerSec))
	case "Write":
		return fmt.Sprintf("%d", d.WriteIOs)
	case "W/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(d.WriteBytesPerSec))
	case "Await":
		return fmt.Sprintf("%.2fms", d.Await)
	case "Avio":
		return fmt.Sprintf("%.2fms", d.Avio)
	}
	return ""
}

func (diskMap DiskMap) Collect(prev, curr *store.Sample) {

	for k := range diskMap {
		delete(diskMap, k)
	}
	for name, _ := range curr.DiskStats {
		new := curr.DiskStats[name]
		old := prev.DiskStats[name]
		interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)
		d := Disk{
			DeviceName:      new.DeviceName,
			ReadIOs:         new.ReadIOs - old.ReadIOs,
			ReadMerges:      new.ReadMerges - old.ReadMerges,
			ReadSectors:     new.ReadSectors - old.ReadSectors,
			ReadTicks:       new.ReadTicks - old.ReadTicks,
			WriteIOs:        new.WriteIOs - old.WriteIOs,
			WriteMerges:     new.WriteMerges - old.WriteMerges,
			WriteSectors:    new.WriteSectors - old.WriteSectors,
			WriteTicks:      new.WriteTicks - old.WriteTicks,
			IOsTotalTicks:   new.IOsTotalTicks - old.IOsTotalTicks,
			WeightedIOTicks: new.WeightedIOTicks - old.WeightedIOTicks,
		}
		d.ReadBytesPerSec = d.ReadSectors * 512 / interval
		d.WriteBytesPerSec = d.WriteSectors * 512 / interval
		d.Await = float64(d.WeightedIOTicks) / float64((d.ReadIOs + d.WriteIOs))
		d.Avio = float64(d.IOsTotalTicks) / float64((d.ReadIOs + d.WriteIOs))
		d.Busy = d.WeightedIOTicks / 1000 * 100 / interval
		diskMap[name] = d
	}

}
