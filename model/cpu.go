package model

import (
	"fmt"

	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
)

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

func (c *CPU) GetRenderValue(field string) string {
	switch field {
	case "Index":
		return fmt.Sprintf("%s", c.Index)
	case "User":
		return fmt.Sprintf("%.1f%%", c.User)
	case "Nice":
		return fmt.Sprintf("%.1f%%", c.Nice)
	case "System":
		return fmt.Sprintf("%.1f%%", c.System)
	case "Idle":
		return fmt.Sprintf("%.1f%%", c.Idle)
	case "Iowait":
		return fmt.Sprintf("%.1f%%", c.Iowait)
	case "IRQ":
		return fmt.Sprintf("%.1f%%", c.IRQ)
	case "SoftIRQ":
		return fmt.Sprintf("%.1f%%", c.SoftIRQ)
	case "Steal":
		return fmt.Sprintf("%.1f%%", c.Steal)
	case "Guest":
		return fmt.Sprintf("%.1f%%", c.Guest)
	case "GuestNice":
		return fmt.Sprintf("%.1f%%", c.GuestNice)
	}
	return ""
}

type CPUSlice []CPU

func (cpus *CPUSlice) Collect(prev, curr *store.Sample) {

	*cpus = (*cpus)[:0]

	c := calcCpuUsage(prev.CPUTotal, curr.CPUTotal)
	c.Index = "total"

	*cpus = append(*cpus, c)

	prevMap := make(map[int]procfs.CPUStat)
	for i, v := range prev.CPU {
		prevMap[i] = v
	}

	for i := 0; i < len(curr.CPU); i++ {
		c := calcCpuUsage(prevMap[i], curr.CPU[i])
		c.Index = fmt.Sprintf("%d", i)

		*cpus = append(*cpus, c)
	}

}

func calcCpuUsage(prev, curr procfs.CPUStat) CPU {

	c := CPU{}

	empty := procfs.CPUStat{}
	if curr == empty {
		return c
	}

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
	c.Idle = 100 - c.User - c.Nice - c.System - c.Iowait - c.IRQ - c.SoftIRQ - c.Steal - c.Guest - c.GuestNice

	return c
}
