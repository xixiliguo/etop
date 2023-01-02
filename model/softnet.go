package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func (softnet *Softnet) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for softnet stat", field)
	switch field {
	case "CPU":
		s = config[field].Render(softnet.CPU)
	case "Processed":
		s = config[field].Render(softnet.Processed)
	case "Dropped":
		s = config[field].Render(softnet.Dropped)
	case "TimeSqueezed":
		s = config[field].Render(softnet.TimeSqueezed)
	case "CPUCollision":
		s = config[field].Render(softnet.CPUCollision)
	case "ReceivedRps":
		s = config[field].Render(softnet.ReceivedRps)
	case "FlowLimitCount":
		s = config[field].Render(softnet.FlowLimitCount)
	case "SoftnetBacklogLen":
		s = config[field].Render(softnet.SoftnetBacklogLen)
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

func (softnets *SoftnetSlice) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
	looptext:
		for _, s := range *softnets {
			row := strings.Builder{}
			row.WriteString(dateTime)
			for _, f := range opt.Fields {
				renderValue := s.GetRenderValue(config, f)
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
		for _, s := range *softnets {
			row := make(map[string]string)
			row["Timestamp"] = dateTime
			for _, f := range opt.Fields {
				renderValue := s.GetRenderValue(config, f)
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
