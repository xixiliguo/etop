package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultNetProtocolFields = []string{"Name", "Sockets", "Memory", "Pressure"}

type NetProtocol struct {
	Name     string
	Sockets  int64
	Memory   int64
	Pressure int
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
	case "Pressure":
		cfg = Field{"Pressure", Raw, 0, "", 10, false}
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
	case "Pressure":
		s = cfg.Render(n.Pressure)
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
			Name:     v.Name,
			Sockets:  v.Sockets,
			Memory:   memory,
			Pressure: v.Pressure,
		}
	}
}
