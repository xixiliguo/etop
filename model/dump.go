package model

import (
	"bytes"
	"encoding/json"
	"strconv"
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

func dumpOpenMetric(timeStamp int64, OMConfig OpenMetricRenderConfig, config RenderConfig, opt DumpOption, m Render) {

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	for _, f := range opt.Fields {
		if _, ok := OMConfig[f]; !ok {
			continue
		}
		renderValue := m.GetRenderValue(config, f)
		if f == opt.SelectField && opt.Filter != nil {
			if opt.Filter.MatchString(renderValue) == false {
				return
			}
		}
		cfg := OMConfig[f]
		buf.WriteString("# TYPE" + " " + cfg.Name + " " + typToString[cfg.Typ] + "\n")
		if cfg.Unit != "" {
			buf.WriteString("# UNIT" + " " + cfg.Name + " " + cfg.Unit + "\n")
		}
		if cfg.Help != "" {
			buf.WriteString("# HELP" + " " + cfg.Name + " " + cfg.Help + "\n")
		}
		buf.WriteString(cfg.Name)
		if len(cfg.Labels) != 0 {
			buf.WriteString("{")
			first := true
			for _, l := range cfg.Labels {
				if first == true {
					first = false
				} else {
					buf.WriteString(",")
				}
				buf.WriteString(l + "=" + `"` + m.GetRenderValue(config, l) + `"`)
			}
			buf.WriteString("}")
		}

		buf.WriteString(" " + renderValue)
		buf.WriteString(" " + strconv.FormatInt(timeStamp, 10))
		buf.WriteString("\n")
	}
	buf.WriteTo(opt.Output)
}
