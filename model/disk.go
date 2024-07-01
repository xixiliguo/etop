package model

import (
	"sort"
	"time"

	"github.com/xixiliguo/etop/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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
		if d.ReadIOs != 0 {
			d.ReadAvgIOSize = float64(d.ReadSectors*512) / float64((d.ReadIOs))
			d.ReadAvgWait = float64(d.ReadTicks) / float64((d.ReadIOs))
		}
		if d.WriteIOs != 0 {
			d.WriteAvgIOSize = float64(d.WriteSectors*512) / float64((d.WriteIOs))
			d.WriteAvgWait = float64(d.WriteTicks) / float64((d.WriteIOs))
		}
		if d.DiscardIOs != 0 {
			d.DiscardAvgIOSize = float64(d.DiscardSectors*512) / float64((d.DiscardIOs))
			d.DiscardAvgWait = float64(d.DiscardTicks) / float64((d.DiscardIOs))
		}
		if d.ReadIOs+d.WriteIOs+d.DiscardIOs != 0 {
			d.AvgIOSize = float64(d.ReadSectors+d.WriteSectors+d.DiscardSectors) * 512 / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))
			d.AvgIOWait = float64(d.ReadTicks+d.WriteTicks+d.DiscardTicks) / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))
			d.AvgIOTime = float64(d.IOsTotalTicks) / float64((d.ReadIOs + d.WriteIOs + d.DiscardIOs))
		}

		d.AvgQueueLength = float64(d.WeightedIOTicks) / 1000 / float64(interval)

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

func (diskMap DiskMap) GetOtelMetric(timeStamp int64, sm *metricdata.ScopeMetrics) {

	sm.Scope = instrumentation.Scope{Name: "disk", Version: "0.0.1"}
	diskIO := metricdata.Metrics{
		Name: "disk.io",
		Data: metricdata.Gauge[float64]{},
	}
	diskIOData := metricdata.Gauge[float64]{}
	diskByte := metricdata.Metrics{
		Name: "disk.byte",
	}
	diskByteData := metricdata.Gauge[float64]{}
	diskUsage := metricdata.Metrics{
		Name: "disk.usage",
	}
	diskUsageData := metricdata.Gauge[float64]{}
	diskIOSize := metricdata.Metrics{
		Name: "disk.io.size",
	}
	diskIOSizeData := metricdata.Gauge[float64]{}
	diskQueueLen := metricdata.Metrics{
		Name: "disk.io.queue.len",
	}
	diskQueueLenData := metricdata.Gauge[float64]{}
	diskFlight := metricdata.Metrics{
		Name: "disk.flight",
	}
	diskFlightData := metricdata.Gauge[float64]{}
	diskIOWait := metricdata.Metrics{
		Name: "disk.iowait",
	}
	diskIOWaitData := metricdata.Gauge[float64]{}
	diskIOTime := metricdata.Metrics{
		Name: "disk.iotime",
	}
	diskIOTimeData := metricdata.Gauge[float64]{}
	for _, n := range diskMap.GetKeys() {
		disk := diskMap[n]
		name := attribute.String("disk", n)

		diskIOData.DataPoints = append(diskIOData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "read")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.ReadPerSec,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "write")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.WritePerSec,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "discard")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.DiscardPerSec,
			},
		}...)

		diskByteData.DataPoints = append(diskByteData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "read")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.ReadBytePerSec,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "write")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.WriteBytePerSec,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "discard")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.DiscardBytePerSec,
			},
		}...)

		diskUsageData.DataPoints = append(diskUsageData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.Util,
			},
		}...)

		diskIOSizeData.DataPoints = append(diskIOSizeData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "read")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.ReadAvgIOSize,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "write")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.WriteAvgIOSize,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "discard")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.DiscardAvgIOSize,
			},
		}...)

		diskQueueLenData.DataPoints = append(diskQueueLenData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.AvgQueueLength,
			},
		}...)

		diskFlightData.DataPoints = append(diskFlightData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name),
				Time:       time.Unix(timeStamp, 0),
				Value:      float64(disk.IOsInProgress),
			},
		}...)

		diskIOWaitData.DataPoints = append(diskIOWaitData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "read")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.ReadAvgWait,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "write")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.WriteAvgWait,
			},
			{
				Attributes: attribute.NewSet(name, attribute.String("action", "discard")),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.DiscardAvgWait,
			},
		}...)

		diskIOTimeData.DataPoints = append(diskIOTimeData.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(name),
				Time:       time.Unix(timeStamp, 0),
				Value:      disk.AvgIOTime,
			},
		}...)

	}
	diskIO.Data = diskIOData
	diskByte.Data = diskByteData
	diskUsage.Data = diskUsageData
	diskIOSize.Data = diskIOSizeData
	diskQueueLen.Data = diskQueueLenData
	diskFlight.Data = diskFlightData
	diskIOWait.Data = diskIOWaitData
	diskIOTime.Data = diskIOTimeData
	sm.Metrics = append(sm.Metrics, diskIO, diskByte, diskUsage, diskIOSize, diskQueueLen, diskFlight, diskIOWait, diskIOTime)
}
