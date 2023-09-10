package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xixiliguo/etop/store"
)

var DefaultNetDevFields = []string{
	"Name",
	"RxPacket/s", "TxPacket/s",
	"RxByte/s", "TxByte/s",
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

func (n *NetDev) GetRenderValue(config RenderConfig, field string) string {

	s := fmt.Sprintf("no %s for netdev stat", field)
	switch field {
	case "Name":
		s = config[field].Render(n.Name)
	case "RxBytes":
		s = config[field].Render(n.RxBytes)
	case "RxPackets":
		s = config[field].Render(n.RxPackets)
	case "RxErrors":
		s = config[field].Render(n.RxErrors)
	case "RxDropped":
		s = config[field].Render(n.RxDropped)
	case "RxFIFO":
		s = config[field].Render(n.RxFIFO)
	case "RxFrame":
		s = config[field].Render(n.RxFrame)
	case "RxCompressed":
		s = config[field].Render(n.RxCompressed)
	case "RxMulticast":
		s = config[field].Render(n.RxMulticast)
	case "TxBytes":
		s = config[field].Render(n.TxBytes)
	case "TxPackets":
		s = config[field].Render(n.TxPackets)
	case "TxErrors":
		s = config[field].Render(n.TxErrors)
	case "TxDropped":
		s = config[field].Render(n.TxDropped)
	case "TxFIFO":
		s = config[field].Render(n.TxFIFO)
	case "TxCollisions":
		s = config[field].Render(n.TxCollisions)
	case "TxCarrier":
		s = config[field].Render(n.TxCarrier)
	case "TxCompressed":
		s = config[field].Render(n.TxCompressed)
	case "RxByte/s":
		s = config[field].Render(n.RxBytePerSec)
	case "RxPacket/s":
		s = config[field].Render(n.RxPacketPerSec)
	case "TxByte/s":
		s = config[field].Render(n.TxBytePerSec)
	case "TxPacket/s":
		s = config[field].Render(n.TxPacketPerSec)
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
