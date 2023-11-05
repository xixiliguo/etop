package model

import (
	"sort"

	"github.com/xixiliguo/etop/store"
)

var DefaultDiskFields = []string{
	"Disk", "Util",
	"ReadPerSec", "ReadBytePerSec", "WritePerSec", "WriteBytePerSec",
	"AvgIOSize", "AvgQueueLen", "InFlight", "AvgIOWait", "AvgIOTime",
}

type Disk struct {
	DeviceName             string
	ReadIOs                uint64
	ReadMerges             uint64
	ReadSectors            uint64
	ReadTicks              uint64
	WriteIOs               uint64
	WriteMerges            uint64
	WriteSectors           uint64
	WriteTicks             uint64
	IOsInProgress          uint64
	IOsTotalTicks          uint64
	WeightedIOTicks        uint64
	DiscardIOs             uint64
	DiscardMerges          uint64
	DiscardSectors         uint64
	DiscardTicks           uint64
	FlushRequestsCompleted uint64
	TimeSpentFlushing      uint64
	ReadPerSec             float64
	WritePerSec            float64
	DiscardPerSec          float64
	ReadBytePerSec         float64
	WriteBytePerSec        float64
	DiscardBytePerSec      float64
	ReadAvgIOSize          float64
	WriteAvgIOSize         float64
	DiscardAvgIOSize       float64
	AvgIOSize              float64
	ReadAvgWait            float64
	WriteAvgWait           float64
	DiscardAvgWait         float64
	AvgIOWait              float64
	AvgQueueLength         float64
	AvgIOTime              float64
	Util                   float64
}

type DiskMap map[string]Disk

func (d *Disk) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Disk":
		cfg = Field{"Disk", Raw, 0, "", 10, false}
	case "Util":
		cfg = Field{"Util", Raw, 1, "%", 10, false}
	case "Read":
		cfg = Field{"Read", Raw, 0, "", 10, false}
	case "ReadPerSec":
		cfg = Field{"Read/s", Raw, 0, "/s", 10, false}
	case "ReadBytePerSec":
		cfg = Field{"ReadByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "Write":
		cfg = Field{"Write", Raw, 0, "", 10, false}
	case "WritePerSec":
		cfg = Field{"Write/s", Raw, 0, "/s", 10, false}
	case "WriteBytePerSec":
		cfg = Field{"WriteByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "AvgIOSize":
		cfg = Field{"AvgIOSize", HumanReadableSize, 1, "", 10, false}
	case "AvgQueueLen":
		cfg = Field{"AvgQueueLen", Raw, 1, "", 10, false}
	case "InFlight":
		cfg = Field{"InFlight", Raw, 1, "", 10, false}
	case "AvgIOWait":
		cfg = Field{"AvgIOWait", Raw, 1, " ms", 10, false}
	case "AvgIOTime":
		cfg = Field{"AvgIOTime", Raw, 1, " ms", 10, false}
	}
	return cfg
}

func (d *Disk) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "Disk":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"Disk"}}
	case "Util":
		cfg = OpenMetricField{"Util", Gauge, "", "", []string{"Disk"}}
	case "Read":
		cfg = OpenMetricField{"Read", Gauge, "", "", []string{"Disk"}}
	case "ReadPerSec":
		cfg = OpenMetricField{"ReadPerSec", Gauge, "", "", []string{"Disk"}}
	case "ReadBytePerSec":
		cfg = OpenMetricField{"ReadBytePerSec", Gauge, "", "", []string{"Disk"}}
	case "Write":
		cfg = OpenMetricField{"Write", Gauge, "", "", []string{"Disk"}}
	case "WritePerSec":
		cfg = OpenMetricField{"WritePerSec", Gauge, "", "", []string{"Disk"}}
	case "WriteBytePerSec":
		cfg = OpenMetricField{"WriteBytePerSec", Gauge, "", "", []string{"Disk"}}
	case "AvgIOSize":
		cfg = OpenMetricField{"AvgIOSize", Gauge, "", "", []string{"Disk"}}
	case "AvgQueueLen":
		cfg = OpenMetricField{"AvgQueueLen", Gauge, "", "", []string{"Disk"}}
	case "InFlight":
		cfg = OpenMetricField{"InFlight", Gauge, "", "", []string{"Disk"}}
	case "AvgIOWait":
		cfg = OpenMetricField{"AvgIOWait", Gauge, "", "", []string{"Disk"}}
	case "AvgIOTime":
		cfg = OpenMetricField{"AvgIOTime", Gauge, "", "", []string{"Disk"}}
	}
	return cfg
}

func (d *Disk) GetRenderValue(field string, opt FieldOpt) string {

	cfg := d.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Disk":
		s = cfg.Render(d.DeviceName)
	case "Util":
		s = cfg.Render(d.Util)
	case "Read":
		s = cfg.Render(d.ReadIOs)
	case "ReadPerSec":
		s = cfg.Render(d.ReadPerSec)
	case "ReadBytePerSec":
		s = cfg.Render(d.ReadBytePerSec)
	case "Write":
		s = cfg.Render(d.WriteIOs)
	case "WritePerSec":
		s = cfg.Render(d.WritePerSec)
	case "WriteBytePerSec":
		s = cfg.Render(d.WriteBytePerSec)
	case "AvgIOSize":
		s = cfg.Render(d.AvgIOSize)
	case "AvgQueueLen":
		s = cfg.Render(d.AvgQueueLength)
	case "InFlight":
		s = cfg.Render(d.IOsInProgress)
	case "AvgIOWait":
		s = cfg.Render(d.AvgIOWait)
	case "AvgIOTime":
		s = cfg.Render(d.AvgIOTime)
	default:
		s = "no " + field + " for disk stat"
	}
	return s
}

func (diskMap DiskMap) Collect(prev, curr *store.Sample) {

	for k := range diskMap {
		delete(diskMap, k)
	}
	for name := range curr.DiskStats {
		new := curr.DiskStats[name]
		old := prev.DiskStats[name]
		interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)
		d := Disk{
			DeviceName:             new.DeviceName,
			ReadIOs:                new.ReadIOs - old.ReadIOs,
			ReadMerges:             new.ReadMerges - old.ReadMerges,
			ReadSectors:            new.ReadSectors - old.ReadSectors,
			ReadTicks:              new.ReadTicks - old.ReadTicks,
			WriteIOs:               new.WriteIOs - old.WriteIOs,
			WriteMerges:            new.WriteMerges - old.WriteMerges,
			WriteSectors:           new.WriteSectors - old.WriteSectors,
			WriteTicks:             new.WriteTicks - old.WriteTicks,
			IOsInProgress:          new.IOsInProgress,
			IOsTotalTicks:          new.IOsTotalTicks - old.IOsTotalTicks,
			WeightedIOTicks:        new.WeightedIOTicks - old.WeightedIOTicks,
			DiscardIOs:             new.DiscardIOs - old.DiscardIOs,
			DiscardMerges:          new.DiscardMerges - old.DiscardMerges,
			DiscardSectors:         new.DiscardSectors - old.DiscardSectors,
			DiscardTicks:           new.DiscardTicks - old.DiscardTicks,
			FlushRequestsCompleted: new.FlushRequestsCompleted - old.FlushRequestsCompleted,
			TimeSpentFlushing:      new.TimeSpentFlushing - old.TimeSpentFlushing,
		}

		d.ReadPerSec = float64(d.ReadIOs) / float64(interval)
		d.WritePerSec = float64(d.WriteIOs) / float64(interval)
		d.DiscardPerSec = float64(d.DiscardIOs) / float64(interval)

		d.ReadBytePerSec = float64(d.ReadSectors*512) / float64(interval)
		d.WriteBytePerSec = float64(d.WriteSectors*512) / float64(interval)
		d.DiscardBytePerSec = float64(d.DiscardBytePerSec*512) / float64(interval)

		d.ReadAvgIOSize = float64(d.ReadSectors*512) / float64((d.ReadIOs))
		d.WriteAvgIOSize = float64(d.WriteSectors*512) / float64((d.WriteIOs))
		d.DiscardAvgIOSize = float64(d.DiscardSectors*512) / float64((d.DiscardIOs))
		d.AvgIOSize = float64(d.ReadSectors+d.WriteSectors+d.DiscardSectors) * 512 / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))

		d.ReadAvgWait = float64(d.ReadTicks) / float64((d.ReadIOs))
		d.WriteAvgWait = float64(d.WriteTicks) / float64((d.WriteIOs))
		d.DiscardAvgWait = float64(d.DiscardTicks) / float64((d.DiscardIOs))
		d.AvgIOWait = float64(d.ReadTicks+d.WriteTicks+d.DiscardTicks) / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))

		d.AvgQueueLength = float64(d.WeightedIOTicks) / 1000 / float64(interval)

		d.AvgIOTime = float64(d.IOsTotalTicks) / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))
		d.Util = float64(d.IOsTotalTicks) / 1000 * 100 / float64(interval)
		diskMap[name] = d
	}

}

func (diskMap DiskMap) GetKeys() []string {

	keys := []string{}
	for k := range diskMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
