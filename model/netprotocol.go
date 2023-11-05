package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultNetProtocolFields = []string{"Name", "Sockets", "Memory"}

type NetProtocol struct {
	Name    string
	Sockets int64
	Memory  int64
}

type NetProtocolMap map[string]NetProtocol

func (n *NetProtocol) DefaultConfig(field string) Field {

	cfg := Field{}
	switch field {
	case "Name":
		cfg = Field{"Name", Raw, 0, "", 10, false}
	case "Sockets":
		cfg = Field{"Sockets", Raw, 0, "", 10, false}
	case "Memory":
		cfg = Field{"Memory", HumanReadableSize, 0, "", 10, false}
	}
	return cfg
}

func (n *NetProtocol) DefaultOMConfig(field string) OpenMetricField {

	cfg := OpenMetricField{}
	switch field {
	case "Name":
		cfg = OpenMetricField{"", Gauge, "", "", []string{"Name"}}
	case "Sockets":
		cfg = OpenMetricField{"Sockets", Gauge, "", "", []string{"Name"}}
	case "Memory":
		cfg = OpenMetricField{"Memory", Gauge, "", "", []string{"Name"}}
	}
	return cfg
}

func (n *NetProtocol) GetRenderValue(field string, opt FieldOpt) string {

	cfg := n.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "Name":
		s = cfg.Render(n.Name)
	case "Sockets":
		s = cfg.Render(n.Sockets)
	case "Memory":
		s = cfg.Render(n.Memory)
	default:
		s = "no " + field + " for netprotocol stat"
	}
	return s
}

func (netProtocolMap NetProtocolMap) Collect(prev, curr *store.Sample) {
	for k := range netProtocolMap {
		delete(netProtocolMap, k)
	}
	for name, v := range curr.NetProtocolStats {
		memory := v.Memory
		if v.Memory >= 0 {
			memory = v.Memory * int64(curr.PageSize)
		}
		netProtocolMap[name] = NetProtocol{
			Name:    v.Name,
			Sockets: v.Sockets,
			Memory:  memory,
		}
	}
}
