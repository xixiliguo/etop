package model

import (
	"fmt"

	"github.com/xixiliguo/etop/store"
)

var DefaultNetProtocolFields = []string{"Name", "Sockets", "Memory"}

type NetProtocol struct {
	Name    string
	Sockets int64
	Memory  int64
}

type NetProtocolMap map[string]NetProtocol

func (n *NetProtocol) GetRenderValue(config RenderConfig, field string) string {

	s := fmt.Sprintf("no %s for netprotocol stat", field)
	switch field {
	case "Name":
		s = config[field].Render(n.Name)
	case "Sockets":
		s = config[field].Render(n.Sockets)
	case "Memory":
		s = config[field].Render(n.Memory)
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
