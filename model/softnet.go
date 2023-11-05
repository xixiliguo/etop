package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultSoftnetFields = []string{"CPU", "Processed", "Dropped",
	"TimeSqueezed", "CPUCollision", "ReceivedRps", "FlowLimitCount", "SoftnetBacklogLen"}

type Softnet struct {
	CPU               uint32
	Processed         uint32
	Dropped           uint32
	TimeSqueezed      uint32
	CPUCollision      uint32
	ReceivedRps       uint32
	FlowLimitCount    uint32
	SoftnetBacklogLen uint32
}

func (softnet *Softnet) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "CPU":
		cfg = Field{"CPU", Raw, 0, "", 10, false}
	case "Processed":
		cfg = Field{"Processed", Raw, 0, "", 10, false}
	case "Dropped":
		cfg = Field{"Dropped", Raw, 0, "", 10, false}
	case "TimeSqueezed":
		cfg = Field{"TimeSqueezed", Raw, 0, "", 10, false}
	case "CPUCollision":
		cfg = Field{"CPUCollision", Raw, 0, "", 10, false}
	case "ReceivedRps":
		cfg = Field{"ReceivedRps", Raw, 0, "", 10, false}
	case "FlowLimitCount":
		cfg = Field{"FlowLimitCount", Raw, 0, "", 10, false}
	case "SoftnetBacklogLen":
		cfg = Field{"SoftnetBacklogLen", Raw, 0, "", 10, false}
	}
	return cfg
}

func (softnet *Softnet) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "CPU":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"CPU"}}
	case "Processed":
		cfg = OpenMetricField{"Processed", Gauge, "", "", []string{"CPU"}}
	case "Dropped":
		cfg = OpenMetricField{"Dropped", Gauge, "", "", []string{"CPU"}}
	case "TimeSqueezed":
		cfg = OpenMetricField{"TimeSqueezed", Gauge, "", "", []string{"CPU"}}
	case "CPUCollision":
		cfg = OpenMetricField{"CPUCollision", Gauge, "", "", []string{"CPU"}}
	case "ReceivedRps":
		cfg = OpenMetricField{"ReceivedRps", Gauge, "", "", []string{"CPU"}}
	case "FlowLimitCount":
		cfg = OpenMetricField{"FlowLimitCount", Gauge, "", "", []string{"CPU"}}
	case "SoftnetBacklogLen":
		cfg = OpenMetricField{"SoftnetBacklogLen", Gauge, "", "", []string{"CPU"}}
	}
	return cfg
}

func (softnet *Softnet) GetRenderValue(field string, opt FieldOpt) string {
	cfg := softnet.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "CPU":
		s = cfg.Render(softnet.CPU)
	case "Processed":
		s = cfg.Render(softnet.Processed)
	case "Dropped":
		s = cfg.Render(softnet.Dropped)
	case "TimeSqueezed":
		s = cfg.Render(softnet.TimeSqueezed)
	case "CPUCollision":
		s = cfg.Render(softnet.CPUCollision)
	case "ReceivedRps":
		s = cfg.Render(softnet.ReceivedRps)
	case "FlowLimitCount":
		s = cfg.Render(softnet.FlowLimitCount)
	case "SoftnetBacklogLen":
		s = cfg.Render(softnet.SoftnetBacklogLen)
	default:
		s = "no " + field + " for softnet stat"
	}
	return s
}

type SoftnetSlice []Softnet

func (softnets *SoftnetSlice) Collect(prev, curr *store.Sample) {

	*softnets = (*softnets)[:0]

	for i, new := range curr.SoftNetStats {
		s := Softnet{
			CPU:               uint32(i),
			Processed:         new.Processed,
			Dropped:           new.Dropped,
			TimeSqueezed:      new.TimeSqueezed,
			CPUCollision:      new.CPUCollision,
			ReceivedRps:       new.ReceivedRps,
			FlowLimitCount:    new.FlowLimitCount,
			SoftnetBacklogLen: new.SoftnetBacklogLen,
		}
		if i < len(prev.SoftNetStats) {
			old := prev.SoftNetStats[i]
			s = Softnet{
				CPU:               uint32(i),
				Processed:         new.Processed - old.Processed,
				Dropped:           new.Dropped - old.Dropped,
				TimeSqueezed:      new.TimeSqueezed - old.TimeSqueezed,
				CPUCollision:      new.CPUCollision - old.CPUCollision,
				ReceivedRps:       new.ReceivedRps - old.ReceivedRps,
				FlowLimitCount:    new.FlowLimitCount - old.FlowLimitCount,
				SoftnetBacklogLen: new.SoftnetBacklogLen - old.SoftnetBacklogLen,
			}
		}

		*softnets = append(*softnets, s)
	}
}
