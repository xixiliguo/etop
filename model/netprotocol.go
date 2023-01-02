package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func (netProtocolMap NetProtocolMap) Dump(timeStamp int64, config RenderConfig, opt DumpOption) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	switch opt.Format {
	case "text":
		config.SetFixWidth(true)
	looptext:
		for _, n := range netProtocolMap {
			row := strings.Builder{}
			row.WriteString(dateTime)
			for _, f := range opt.Fields {
				renderValue := n.GetRenderValue(config, f)
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
		for _, n := range netProtocolMap {
			row := make(map[string]string)
			row["Timestamp"] = dateTime
			for _, f := range opt.Fields {
				renderValue := n.GetRenderValue(config, f)
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
