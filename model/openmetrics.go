package model

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

var typToString = [2]string{"gauge", "counter"}

type OpenMetricField struct {
	Name   string
	Typ    MetricType
	Unit   string
	Help   string
	Labels []string
}

type OpenMetricRenderConfig map[string]OpenMetricField

var DefaultOMRenderConfig = make(map[string]OpenMetricRenderConfig)

var sysDefaultOMRenderConfig = make(OpenMetricRenderConfig)

func genSysDefaultOMConfig() {

	sysDefaultOMRenderConfig["Load1"] = OpenMetricField{
		Name:   "Load1",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Load5"] = OpenMetricField{
		Name:   "Load5",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Load15"] = OpenMetricField{
		Name:   "Load15",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Processes"] = OpenMetricField{
		Name:   "Processes",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["Threads"] = OpenMetricField{
		Name:   "Threads",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ProcessesRunning"] = OpenMetricField{
		Name:   "Running",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ProcessesBlocked"] = OpenMetricField{
		Name:   "Blocked",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}

	sysDefaultOMRenderConfig["ClonePerSec"] = OpenMetricField{
		Name:   "ClonePerSec",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
	sysDefaultOMRenderConfig["ContextSwitchPerSec"] = OpenMetricField{
		Name:   "ContextSwitchPerSec",
		Typ:    Gauge,
		Unit:   "",
		Help:   "",
		Labels: []string{},
	}
}

func init() {
	genSysDefaultOMConfig()
	DefaultOMRenderConfig["system"] = sysDefaultOMRenderConfig
}
