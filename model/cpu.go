package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var DefaultCPUFields = []string{
	"Index", "User", "Nice",
	"System", "Idle", "Iowait", "IRQ",
	"SoftIRQ", "Steal", "Guest", "GuestNice",
}

type CPU struct {
	Index     string
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

func (c *CPU) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "Index":
		cfg = Field{"Index", Raw, 0, "", 10, false}
	case "User":
		cfg = Field{"User", Raw, 1, "%", 10, false}
	case "Nice":
		cfg = Field{"Nice", Raw, 1, "%", 10, false}
	case "System":
		cfg = Field{"System", Raw, 1, "%", 10, false}
	case "Idle":
		cfg = Field{"Idle", Raw, 1, "%", 10, false}
	case "Iowait":
		cfg = Field{"Iowait", Raw, 1, "%", 10, false}
	case "IRQ":
		cfg = Field{"IRQ", Raw, 1, "%", 10, false}
	case "SoftIRQ":
		cfg = Field{"SoftIRQ", Raw, 1, "%", 10, false}
	case "Steal":
		cfg = Field{"Steal", Raw, 1, "%", 10, false}
	case "Guest":
		cfg = Field{"Guest", Raw, 1, "%", 10, false}
	case "GuestNice":
		cfg = Field{"GuestNice", Raw, 1, "%", 10, false}
	}
	return cfg
}

func (c *CPU) DefaultOMConfig(field string) OpenMetricField {
	cfg := OpenMetricField{}
	switch field {
	case "Index":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"Index"}}
	case "User":
		cfg = OpenMetricField{"User", Gauge, "", "", []string{"Index"}}
	case "Nice":
		cfg = OpenMetricField{"Nice", Gauge, "", "", []string{"Index"}}
	case "System":
		cfg = OpenMetricField{"System", Gauge, "", "", []string{"Index"}}
	case "Idle":
		cfg = OpenMetricField{"Idle", Gauge, "", "", []string{"Index"}}
	case "Iowait":
		cfg = OpenMetricField{"Iowait", Gauge, "", "", []string{"Index"}}
	case "IRQ":
		cfg = OpenMetricField{"IRQ", Gauge, "", "", []string{"Index"}}
	case "SoftIRQ":
		cfg = OpenMetricField{"SoftIRQ", Gauge, "", "", []string{"Index"}}
	case "Steal":
		cfg = OpenMetricField{"Steal", Gauge, "", "", []string{"Index"}}
	case "Guest":
		cfg = OpenMetricField{"Guest", Gauge, "", "", []string{"Index"}}
	case "GuestNice":
		cfg = OpenMetricField{"GuestNice", Gauge, "", "", []string{"Index"}}
	}
	return cfg
}

func (c *CPU) GetRenderValue(field string, opt FieldOpt) string {

	cfg := c.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Index":
		s = cfg.Render(c.Index)
	case "User":
		s = cfg.Render(c.User)
	case "Nice":
		s = cfg.Render(c.Nice)
	case "System":
		s = cfg.Render(c.System)
	case "Idle":
		s = cfg.Render(c.Idle)
	case "Iowait":
		s = cfg.Render(c.Iowait)
	case "IRQ":
		s = cfg.Render(c.IRQ)
	case "SoftIRQ":
		s = cfg.Render(c.SoftIRQ)
	case "Steal":
		s = cfg.Render(c.Steal)
	case "Guest":
		s = cfg.Render(c.Guest)
	case "GuestNice":
		s = cfg.Render(c.GuestNice)
	default:
		s = "no " + field + " for cpu stat"
	}
	return s
}

type CPUSlice []CPU

func (cpus *CPUSlice) Collect(prev, curr *store.Sample) {

	*cpus = (*cpus)[:0]

	c := calcCpuUsage(prev.CPUTotal, curr.CPUTotal)
	c.Index = "total"

	*cpus = append(*cpus, c)

	indexs := []int64{}
	for i := range curr.CPU {
		indexs = append(indexs, i)
	}
	sort.Slice(indexs, func(i, j int) bool {
		return indexs[i] < indexs[j]
	})

	for _, i := range indexs {
		c := calcCpuUsage(prev.CPU[i], curr.CPU[i])
		c.Index = fmt.Sprintf("%d", i)
		*cpus = append(*cpus, c)
	}
}

func (cpus *CPUSlice) GetOtelMetric(timeStamp int64, sm *metricdata.ScopeMetrics) {

	sm.Scope = instrumentation.Scope{Name: "cpu", Version: "0.0.1"}
	m := metricdata.Metrics{
		Name: "cpu.usage",
	}
	data := metricdata.Gauge[float64]{}

	for _, c := range *cpus {
		idx := attribute.String("cpu", c.Index)
		data.DataPoints = append(data.DataPoints, []metricdata.DataPoint[float64]{
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "user")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.User,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "nice")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.Nice,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "system")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.System,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "idle")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.Idle,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "iowait")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.Iowait,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "irq")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.IRQ,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "soft_irq")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.SoftIRQ,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "steal")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.Steal,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "guest")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.Guest,
			},
			{
				Attributes: attribute.NewSet(idx, attribute.String("mode", "guest_nice")),
				Time:       time.Unix(timeStamp, 0),
				Value:      c.GuestNice,
			},
		}...)
	}
	m.Data = data
	sm.Metrics = append(sm.Metrics, m)
}

func calcCpuUsage(prev, curr procfs.CPUStat) CPU {

	c := CPU{}

	user := curr.User - prev.User
	nice := curr.Nice - prev.Nice
	system := curr.System - prev.System
	idle := curr.Idle - prev.Idle
	iowait := curr.Iowait - prev.Iowait
	irq := curr.IRQ - prev.IRQ
	softIRQ := curr.SoftIRQ - prev.SoftIRQ
	steal := curr.Steal - prev.Steal
	guest := curr.Guest - prev.Guest
	guestNice := curr.GuestNice - prev.GuestNice

	total := user + nice + system + idle + iowait + irq + softIRQ + steal + guest + guestNice
	c.User = user * 100 / total
	c.Nice = nice * 100 / total
	c.System = system * 100 / total
	c.Iowait = iowait * 100 / total
	c.IRQ = irq * 100 / total
	c.SoftIRQ = softIRQ * 100 / total
	c.Steal = steal * 100 / total
	c.Guest = guest * 100 / total
	c.GuestNice = guestNice * 100 / total
	c.Idle = idle * 100 / total

	return c
}
