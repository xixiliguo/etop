package model

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

var omTypToString = [2]string{"gauge", "counter"}

type OpenMetricField struct {
	Name   string
	Typ    MetricType
	Unit   string
	Help   string
	Labels []string
}
