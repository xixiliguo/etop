package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/xixiliguo/etop/store"
)

var DefaultDiskFields = []string{"Disk", "Util",
	"Read/s", "ReadByte/s", "Write/s", "WriteByte/s",
	"AvgIOSize", "AvgQueueLen", "InFlight", "AvgIOWait", "AvgIOTime"}

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

func (d *Disk) GetRenderValue(config RenderConfig, field string) string {

	s := fmt.Sprintf("no %s for disk stat", field)
	switch field {
	case "Disk":
		s = config[field].Render(d.DeviceName)
	case "Util":
		s = config[field].Render(d.Util)
	case "Read":
		s = config[field].Render(d.ReadIOs)
	case "Read/s":
		s = config[field].Render(d.ReadPerSec)
	case "ReadByte/s":
		s = config[field].Render(d.ReadBytePerSec)
	case "Write":
		s = config[field].Render(d.WriteIOs)
	case "Write/s":
		s = config[field].Render(d.WritePerSec)
	case "WriteByte/s":
		s = config[field].Render(d.WriteBytePerSec)
	case "AvgIOSize":
		s = config[field].Render(d.AvgIOSize)
	case "AvgQueueLen":
		s = config[field].Render(d.AvgQueueLength)
	case "InFlight":
		s = config[field].Render(d.IOsInProgress)
	case "AvgIOWait":
		s = config[field].Render(d.AvgIOWait)
	case "AvgIOTime":
		s = config[field].Render(d.AvgIOTime)
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

func (diskMap DiskMap) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
	looptext:
		for _, d := range diskMap {
			row := strings.Builder{}
			row.WriteString(dateTime)
			for _, f := range opt.Fields {
				renderValue := d.GetRenderValue(config, f)
				if f == opt.SelectField && opt.Filter != nil {
					if opt.Filter.MatchString(renderValue) == false {
						continue looptext
					}
				}
				row.WriteString(" ")
				row.WriteString(renderValue)
			}
			row.WriteString("\n")

			opt.Output.WriteString(row.String())
		}
	case "json":
		t := []any{}
	loopjson:
		for _, d := range diskMap {
			row := make(map[string]string)
			row["Timestamp"] = dateTime
			for _, f := range opt.Fields {
				renderValue := d.GetRenderValue(config, f)
				if f == opt.SelectField && opt.Filter != nil {
					if opt.Filter.MatchString(renderValue) == false {
						continue loopjson
					}
				}
				row[config[f].Name] = renderValue
			}
			t = append(t, row)
		}
		b, _ := json.Marshal(t)
		opt.Output.Write(b)
	}
}
