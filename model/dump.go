package model

import (
	"bytes"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

type Render interface {
	DefaultConfig(field string) Field
	GetRenderValue(field string, opt FieldOpt) string
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var bufferMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

func dumpText(timeStamp int64, opt DumpOption, m Render) {

	if !isFilter(opt, m) {
		return
	}

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	buf.WriteString(dateTime)
	for _, f := range opt.Fields {
		renderValue := m.GetRenderValue(f, FieldOpt{
			FixWidth: true,
			Raw:      opt.RawData,
		})
		buf.WriteString(" ")
		buf.WriteString(renderValue)
	}
	buf.WriteString("\n")
	buf.WriteTo(opt.Output)
}

func dumpTextForCgroup(timeStamp int64, opt DumpOption, c Cgroup) {

	dumpText(timeStamp, opt, &c)

	names := []string{}
	for _, child := range c.Child {
		names = append(names, child.Name)
	}
	sort.Strings(names)
	for _, name := range names {
		dumpTextForCgroup(timeStamp, opt, *c.Child[name])
	}
}

func dumpJson(timeStamp int64, opt DumpOption, m Render) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	buf := bufferMapPool.Get().(map[string]string)
	clear(buf)
	defer bufferMapPool.Put(buf)

	buf["Timestamp"] = dateTime
	for _, f := range opt.Fields {
		renderValue := m.GetRenderValue(f, FieldOpt{
			Raw: opt.RawData,
		})
		buf[m.DefaultConfig(f).Name] = renderValue
	}
	b, _ := json.Marshal(buf)
	opt.Output.Write(b)
}

func dumpJsonForCgroup(timeStamp int64, opt DumpOption, c Cgroup) map[string]any {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	bufMap := map[string]any{}

	bufMap["Timestamp"] = dateTime
	for _, f := range opt.Fields {
		renderValue := c.GetRenderValue(f, FieldOpt{
			Raw: opt.RawData,
		})
		bufMap[c.DefaultConfig(f).Name] = renderValue
	}

	childs := []any{}

	names := []string{}
	for _, child := range c.Child {
		names = append(names, child.Name)
	}
	sort.Strings(names)

	for _, name := range names {
		child := c.Child[name]
		childs = append(childs, dumpJsonForCgroup(timeStamp, opt, *child))
	}

	bufMap["Child"] = childs
	return bufMap
}
