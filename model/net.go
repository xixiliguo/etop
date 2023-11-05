package model

import (
	"sort"
	"strings"

	"github.com/xixiliguo/etop/store"
)

var DefaultNetDevFields = []string{
	"Name",
	"RxPacketPerSec", "TxPacketPerSec",
	"RxBytePerSec", "TxBytePerSec",
	"RxErrors", "RxDropped", "RxFIFO", "RxFrame",
	"TxErrors", "TxDropped", "TxFIFO", "TxCollisions",
}

type NetDev struct {
	Name           string
	RxBytes        uint64
	RxPackets      uint64
	RxErrors       uint64
	RxDropped      uint64
	RxFIFO         uint64
	RxFrame        uint64
	RxCompressed   uint64
	RxMulticast    uint64
	TxBytes        uint64
	TxPackets      uint64
	TxErrors       uint64
	TxDropped      uint64
	TxFIFO         uint64
	TxCollisions   uint64
	TxCarrier      uint64
	TxCompressed   uint64
	RxBytePerSec   float64
	RxPacketPerSec float64
	TxBytePerSec   float64
	TxPacketPerSec float64
}

type NetDevMap map[string]NetDev

func (n *NetDev) DefaultConfig(field string) Field {

	cfg := Field{}
	switch field {
	case "Name":
		cfg = Field{"Name", Raw, 0, "", 10, false}
	case "RxBytes":
		cfg = Field{"RxBytes", Raw, 0, "", 10, false}
	case "RxPackets":
		cfg = Field{"RxPackets", Raw, 0, "", 10, false}
	case "RxErrors":
		cfg = Field{"RxErrors", Raw, 0, "", 10, false}
	case "RxDropped":
		cfg = Field{"RxDropped", Raw, 0, "", 10, false}
	case "RxFIFO":
		cfg = Field{"RxFIFO", Raw, 0, "", 10, false}
	case "RxFrame":
		cfg = Field{"RxFrame", Raw, 0, "", 10, false}
	case "RxCompressed":
		cfg = Field{"RxCompressed", Raw, 0, "", 10, false}
	case "RxMulticast":
		cfg = Field{"RxMulticast", Raw, 0, "", 10, false}
	case "TxBytes":
		cfg = Field{"TxBytes", Raw, 0, "", 10, false}
	case "TxPackets":
		cfg = Field{"TxPackets", Raw, 0, "", 10, false}
	case "TxErrors":
		cfg = Field{"TxErrors", Raw, 0, "", 10, false}
	case "TxDropped":
		cfg = Field{"TxDropped", Raw, 0, "", 10, false}
	case "TxFIFO":
		cfg = Field{"TxFIFO", Raw, 0, "", 10, false}
	case "TxCollisions":
		cfg = Field{"TxCollisions", Raw, 0, "", 10, false}
	case "TxCarrier":
		cfg = Field{"TxCarrier", Raw, 0, "", 10, false}
	case "TxCompressed":
		cfg = Field{"TxCompressed", Raw, 0, "", 10, false}
	case "RxBytePerSec":
		cfg = Field{"RxByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "RxPacketPerSec":
		cfg = Field{"RxPacket/s", Raw, 1, "/s", 10, false}
	case "TxBytePerSec":
		cfg = Field{"TxByte/s", HumanReadableSize, 1, "/s", 10, false}
	case "TxPacketPerSec":
		cfg = Field{"TxPacket/s", Raw, 1, "/s", 10, false}
	}
	return cfg
}

func (n *NetDev) DefaultOMConfig(field string) OpenMetricField {

	cfg := OpenMetricField{}
	switch field {
	case "Name":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"Name"}}
	case "RxBytes":
		cfg = OpenMetricField{"RxBytes", Gauge, "", "", []string{"Name"}}
	case "RxPackets":
		cfg = OpenMetricField{"RxPackets", Gauge, "", "", []string{"Name"}}
	case "RxErrors":
		cfg = OpenMetricField{"RxErrors", Gauge, "", "", []string{"Name"}}
	case "RxDropped":
		cfg = OpenMetricField{"RxDropped", Gauge, "", "", []string{"Name"}}
	case "RxFIFO":
		cfg = OpenMetricField{"RxFIFO", Gauge, "", "", []string{"Name"}}
	case "RxFrame":
		cfg = OpenMetricField{"RxFrame", Gauge, "", "", []string{"Name"}}
	case "RxCompressed":
		cfg = OpenMetricField{"RxCompressed", Gauge, "", "", []string{"Name"}}
	case "RxMulticast":
		cfg = OpenMetricField{"RxMulticast", Gauge, "", "", []string{"Name"}}
	case "TxBytes":
		cfg = OpenMetricField{"TxBytes", Gauge, "", "", []string{"Name"}}
	case "TxPackets":
		cfg = OpenMetricField{"TxPackets", Gauge, "", "", []string{"Name"}}
	case "TxErrors":
		cfg = OpenMetricField{"TxErrors", Gauge, "", "", []string{"Name"}}
	case "TxDropped":
		cfg = OpenMetricField{"TxDropped", Gauge, "", "", []string{"Name"}}
	case "TxFIFO":
		cfg = OpenMetricField{"TxFIFO", Gauge, "", "", []string{"Name"}}
	case "TxCollisions":
		cfg = OpenMetricField{"TxCollisions", Gauge, "", "", []string{"Name"}}
	case "TxCarrier":
		cfg = OpenMetricField{"TxCarrier", Gauge, "", "", []string{"Name"}}
	case "TxCompressed":
		cfg = OpenMetricField{"TxCompressed", Gauge, "", "", []string{"Name"}}
	case "RxBytePerSec":
		cfg = OpenMetricField{"RxBytePerSec", Gauge, "", "", []string{"Name"}}
	case "RxPacketPerSec":
		cfg = OpenMetricField{"RxPacketPerSec", Gauge, "", "", []string{"Name"}}
	case "TxBytePerSec":
		cfg = OpenMetricField{"TxBytePerSec", Gauge, "", "", []string{"Name"}}
	case "TxPacketPerSec":
		cfg = OpenMetricField{"TxPacketPerSec", Gauge, "", "", []string{"Name"}}
	}
	return cfg
}

func (n *NetDev) GetRenderValue(field string, opt FieldOpt) string {

	cfg := n.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Name":
		s = cfg.Render(n.Name)
	case "RxBytes":
		s = cfg.Render(n.RxBytes)
	case "RxPackets":
		s = cfg.Render(n.RxPackets)
	case "RxErrors":
		s = cfg.Render(n.RxErrors)
	case "RxDropped":
		s = cfg.Render(n.RxDropped)
	case "RxFIFO":
		s = cfg.Render(n.RxFIFO)
	case "RxFrame":
		s = cfg.Render(n.RxFrame)
	case "RxCompressed":
		s = cfg.Render(n.RxCompressed)
	case "RxMulticast":
		s = cfg.Render(n.RxMulticast)
	case "TxBytes":
		s = cfg.Render(n.TxBytes)
	case "TxPackets":
		s = cfg.Render(n.TxPackets)
	case "TxErrors":
		s = cfg.Render(n.TxErrors)
	case "TxDropped":
		s = cfg.Render(n.TxDropped)
	case "TxFIFO":
		s = cfg.Render(n.TxFIFO)
	case "TxCollisions":
		s = cfg.Render(n.TxCollisions)
	case "TxCarrier":
		s = cfg.Render(n.TxCarrier)
	case "TxCompressed":
		s = cfg.Render(n.TxCompressed)
	case "RxBytePerSec":
		s = cfg.Render(n.RxBytePerSec)
	case "RxPacketPerSec":
		s = cfg.Render(n.RxPacketPerSec)
	case "TxBytePerSec":
		s = cfg.Render(n.TxBytePerSec)
	case "TxPacketPerSec":
		s = cfg.Render(n.TxPacketPerSec)
	default:
		s = "no " + field + " for netdev stat"
	}
	return s
}

func (netMap NetDevMap) Collect(prev, curr *store.Sample) {
	for k := range netMap {
		delete(netMap, k)
	}
	for name := range curr.NetDevStats {
		new := curr.NetDevStats[name]
		old := prev.NetDevStats[name]
		interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)
		n := NetDev{
			Name:         new.Name,
			RxBytes:      new.RxBytes - old.RxBytes,
			RxPackets:    new.RxPackets - old.RxPackets,
			RxErrors:     new.RxErrors - old.RxErrors,
			RxDropped:    new.RxDropped - old.RxDropped,
			RxFIFO:       new.RxFIFO - old.RxFIFO,
			RxFrame:      new.RxFrame - old.RxFrame,
			RxCompressed: new.RxCompressed - old.RxCompressed,
			RxMulticast:  new.RxMulticast - old.RxMulticast,
			TxBytes:      new.TxBytes - old.TxBytes,
			TxPackets:    new.TxPackets - old.TxPackets,
			TxErrors:     new.TxErrors - old.TxErrors,
			TxDropped:    new.TxDropped - old.TxDropped,
			TxFIFO:       new.TxFIFO - old.TxFIFO,
			TxCollisions: new.TxCollisions - old.TxCollisions,
			TxCarrier:    new.TxCarrier - old.TxCarrier,
			TxCompressed: new.TxCompressed - old.TxCompressed,
		}
		n.RxBytePerSec = float64(n.RxBytes) / float64(interval)
		n.RxPacketPerSec = float64(n.RxPackets) / float64(interval)
		n.TxBytePerSec = float64(n.TxBytes) / float64(interval)
		n.TxPacketPerSec = float64(n.TxPackets) / float64(interval)
		netMap[name] = n
	}
}

func (netMap NetDevMap) GetKeys() []string {

	keys := []string{}
	for k := range netMap {
		kk := k
		if k == "lo" {
			continue
		}
		keys = append(keys, kk)
	}
	sort.Slice(keys, func(i, j int) bool {
		if strings.HasPrefix(keys[i], "eth") && strings.HasPrefix(keys[j], "eth") == false {
			return true
		}
		if strings.HasPrefix(keys[j], "eth") && strings.HasPrefix(keys[i], "eth") == false {
			return false
		}
		return keys[i] < keys[j]
	})
	return keys
}
