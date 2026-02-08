package store

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go tool bpf2go -cc clang -cflags $BPF_CFLAGS -type cgroup_net_stat cgroup cgroup.bpf.c -- -I./include

type CgroupNetStat struct {
	Stats *ebpf.Map
	links []link.Link
	log   *slog.Logger
}

func NewCgroupNetStat(log *slog.Logger) *CgroupNetStat {
	return &CgroupNetStat{
		log: log,
	}
}

func (c *CgroupNetStat) Collect() {

	if err := rlimit.RemoveMemlock(); err != nil {
		msg := fmt.Sprintf("remove Memlock: %s", err)
		c.log.Error(msg)
		return
	}

	objs := cgroupObjects{}
	if err := loadCgroupObjects(&objs, nil); err != nil {
		msg := fmt.Sprintf("loading objects: %s", err)
		c.log.Error(msg)
		return
	}

	c.Stats = objs.cgroupMaps.CgroupNetStats

	linkIngres, err := link.AttachCgroup(link.CgroupOptions{
		Path:    "/sys/fs/cgroup",
		Attach:  ebpf.AttachCGroupInetIngress,
		Program: objs.CountIngressPackets,
	})

	if err != nil {
		log.Fatal(err)
	}
	c.links = append(c.links, linkIngres)

	linkEgress, err := link.AttachCgroup(link.CgroupOptions{
		Path:    "/sys/fs/cgroup",
		Attach:  ebpf.AttachCGroupInetEgress,
		Program: objs.CountEgressPackets,
	})

	if err != nil {
		log.Fatal(err)
	}
	c.links = append(c.links, linkEgress)
}

func (c *CgroupNetStat) NetStat(id uint64) (cgroupCgroupNetStat, error) {

	var result cgroupCgroupNetStat
	err := c.Stats.Lookup(&id, &result)
	return result, err
}
