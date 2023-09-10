package model

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"
)

type Render interface {
	GetRenderValue(config RenderConfig, field string) string
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

func dumpText(timeStamp int64, config RenderConfig, opt DumpOption, m Render) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	buf.WriteString(dateTime)
	for _, f := range opt.Fields {
		renderValue := m.GetRenderValue(config, f)
		if f == opt.SelectField && opt.Filter != nil {
			if opt.Filter.MatchString(renderValue) == false {
				return
			}
		}
		buf.WriteString(" ")
		buf.WriteString(renderValue)
	}
	buf.WriteString("\n")
	buf.WriteTo(opt.Output)
}

func dumpJson(timeStamp int64, config RenderConfig, opt DumpOption, m Render) {

	dateTime := time.Unix(timeStamp, 0).Format(time.RFC3339)
	buf := bufferMapPool.Get().(map[string]string)
	clear(buf)
	defer bufferMapPool.Put(buf)

	buf["Timestamp"] = dateTime
	for _, f := range opt.Fields {
		renderValue := m.GetRenderValue(config, f)
		if f == opt.SelectField && opt.Filter != nil {
			if opt.Filter.MatchString(renderValue) == false {
				return
			}
		}
		buf[config[f].Name] = renderValue
	}
	b, _ := json.Marshal(buf)
	opt.Output.Write(b)
}
