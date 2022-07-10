package model

import (
	"fmt"

	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
)

type CPU struct {
	Index     string
	User      int
	Nice      int
	System    int
	Idle      int
	Iowait    int
	IRQ       int
	SoftIRQ   int
	Steal     int
	Guest     int
	GuestNice int
}

func (c *CPU) GetRenderValue(field string) string {
	switch field {
	case "Index":
		return fmt.Sprintf("%s", c.Index)
	case "User":
		return fmt.Sprintf("%d%%", c.User)
	case "Nice":
		return fmt.Sprintf("%d%%", c.Nice)
	case "System":
		return fmt.Sprintf("%d%%", c.System)
	case "Idle":
		return fmt.Sprintf("%d%%", c.Idle)
	case "Iowait":
		return fmt.Sprintf("%d%%", c.Iowait)
	case "IRQ":
		return fmt.Sprintf("%d%%", c.IRQ)
	case "SoftIRQ":
		return fmt.Sprintf("%d%%", c.SoftIRQ)
	case "Steal":
		return fmt.Sprintf("%d%%", c.Steal)
	case "Guest":
		return fmt.Sprintf("%d%%", c.Guest)
	case "GuestNice":
		return fmt.Sprintf("%d%%", c.GuestNice)
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
	c.User = int(user * 100 / total)
	c.Nice = int(nice * 100 / total)
	c.System = int(system * 100 / total)
	c.Iowait = int(iowait * 100 / total)
	c.IRQ = int(irq * 100 / total)
	c.SoftIRQ = int(softIRQ * 100 / total)
	c.Steal = int(steal * 100 / total)
	c.Guest = int(guest * 100 / total)
	c.GuestNice = int(guestNice * 100 / total)
	c.Idle = 100 - c.User - c.Nice - c.System - c.Iowait - c.IRQ - c.SoftIRQ - c.Steal - c.Guest - c.GuestNice

	return c
}
