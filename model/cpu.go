package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
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

func (c *CPU) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for cpu stat", field)
	switch field {
	case "Index":
		s = config[field].Render(c.Index)
	case "User":
		s = config[field].Render(c.User)
	case "Nice":
		s = config[field].Render(c.Nice)
	case "System":
		s = config[field].Render(c.System)
	case "Idle":
		s = config[field].Render(c.Idle)
	case "Iowait":
		s = config[field].Render(c.Iowait)
	case "IRQ":
		s = config[field].Render(c.IRQ)
	case "SoftIRQ":
		s = config[field].Render(c.SoftIRQ)
	case "Steal":
		s = config[field].Render(c.Steal)
	case "Guest":
		s = config[field].Render(c.Guest)
	case "GuestNice":
		s = config[field].Render(c.GuestNice)
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

func (cpus *CPUSlice) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
	looptext:
		for _, c := range *cpus {
			row := strings.Builder{}
			row.WriteString(dateTime)
			for _, f := range opt.Fields {
				renderValue := c.GetRenderValue(config, f)
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
		for _, c := range *cpus {
			row := make(map[string]string)
			row["Timestamp"] = dateTime
			for _, f := range opt.Fields {
				renderValue := c.GetRenderValue(config, f)
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
