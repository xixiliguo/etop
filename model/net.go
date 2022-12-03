package model

import (
	"fmt"

	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/util"
)

type Net struct {
	Name            string
	RxBytes         uint64
	RxPackets       uint64
	RxErrors        uint64
	RxDropped       uint64
	RxFIFO          uint64
	RxFrame         uint64
	RxCompressed    uint64
	RxMulticast     uint64
	TxBytes         uint64
	TxPackets       uint64
	TxErrors        uint64
	TxDropped       uint64
	TxFIFO          uint64
	TxCollisions    uint64
	TxCarrier       uint64
	TxCompressed    uint64
	RxBytesPerSec   uint64
	RxPacketsPerSec uint64
	TxBytesPerSec   uint64
	TxPacketsPerSec uint64
}

func (n *Net) GetRenderValue(field string) string {
	switch field {
	case "Name":
		return fmt.Sprintf("%s", n.Name)
	case "RxBytes":
		return fmt.Sprintf("%d", n.RxBytes)
	case "RxPackets":
		return fmt.Sprintf("%d", n.RxPackets)
	case "RxErrors":
		return fmt.Sprintf("%d", n.RxErrors)
	case "RxDropped":
		return fmt.Sprintf("%d", n.RxDropped)
	case "RxFIFO":
		return fmt.Sprintf("%d", n.RxFIFO)
	case "RxFrame":
		return fmt.Sprintf("%d", n.RxFrame)
	case "RxCompressed":
		return fmt.Sprintf("%d", n.RxCompressed)
	case "RxMulticast":
		return fmt.Sprintf("%d", n.RxMulticast)
	case "TxBytes":
		return fmt.Sprintf("%d", n.TxBytes)
	case "TxPackets":
		return fmt.Sprintf("%d", n.TxPackets)
	case "TxErrors":
		return fmt.Sprintf("%d", n.TxErrors)
	case "TxDropped":
		return fmt.Sprintf("%d", n.TxDropped)
	case "TxFIFO":
		return fmt.Sprintf("%d", n.TxFIFO)
	case "TxCollisions":
		return fmt.Sprintf("%d", n.TxCollisions)
	case "TxCarrier":
		return fmt.Sprintf("%d", n.TxCarrier)
	case "TxCompressed":
		return fmt.Sprintf("%d", n.TxCompressed)
	case "R/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(n.RxBytesPerSec))
	case "Rp/s":
		return fmt.Sprintf("%d/s", n.RxPackets)
	case "T/s":
		return fmt.Sprintf("%s/s", util.GetHumanSize(n.TxBytesPerSec))
	case "Tp/s":
		return fmt.Sprintf("%d/s", n.TxPacketsPerSec)
	}
	return ""
}

type NetSlice []Net

func (nets *NetSlice) Collect(prev, curr *store.Sample) {
	*nets = (*nets)[:0]
	prevMap := make(map[string]procfs.NetDevLine)
	for i := 0; i < len(prev.NetStats); i++ {
		prevMap[prev.NetStats[i].Name] = prev.NetStats[i]
	}
	for i := 0; i < len(curr.NetStats); i++ {

		new := curr.NetStats[i]
		old := prevMap[new.Name]
		interval := uint64(curr.TimeStamp) - uint64(prev.TimeStamp)

		n := Net{
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
		n.RxBytesPerSec = n.RxBytes / interval
		n.RxPacketsPerSec = n.RxPackets / interval
		n.TxBytesPerSec = n.TxBytes / interval
		n.TxPacketsPerSec = n.TxPackets / interval

		*nets = append(*nets, n)
	}

}
