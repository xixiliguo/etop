package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/xixiliguo/etop/store"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Model struct {
	Mode  string
	Store store.Store
	log   *slog.Logger
	Prev  store.Sample
	Curr  store.Sample
	Sys   System
	CPUs  CPUSlice
	MEM
	Vm
	Disks DiskMap
	Nets  NetDevMap
	NetStat
	NetProtocols NetProtocolMap
	Softnets     SoftnetSlice
	Processes    ProcessMap
	Cgroup
}

func NewSysModel(s *store.LocalStore, log *slog.Logger) (*Model, error) {
	p := &Model{
		Mode:         "report",
		Store:        s,
		log:          log,
		Prev:         store.NewSample(),
		Curr:         store.NewSample(),
		Sys:          System{},
		CPUs:         []CPU{},
		MEM:          MEM{},
		Vm:           Vm{},
		Disks:        make(DiskMap),
		Nets:         make(NetDevMap),
		NetProtocols: make(NetProtocolMap),
		Softnets:     []Softnet{},
		Processes:    make(ProcessMap),
		Cgroup:       Cgroup{},
	}
	return p, nil
}

func (s *Model) CollectLiveSample(exit *store.ExitProcess) error {

	s.Prev = s.Curr
	s.Curr = store.NewSample()
	if err := store.CollectSampleFromSys(&s.Curr, exit, s.log); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectNext() error {

	n := store.NewSample()
	if err := s.Store.NextSample(1, &n); err != nil {
		return err
	}
	s.Prev = s.Curr
	s.Curr = n
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Info("skip one sample since system reboot")
		return s.CollectNext()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectPrev() error {
	n := store.NewSample()
	if err := s.Store.NextSample(-2, &n); err != nil {
		return err
	}
	s.Prev = n
	s.Curr = store.NewSample()
	if err := s.Store.NextSample(1, &s.Curr); err != nil {
		return err
	}
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Info("skip one sample since system reboot")
		return s.CollectPrev()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectSampleByTime(timeStamp int64) error {
	s.Curr = store.NewSample()
	if err := s.Store.JumpSampleByTimeStamp(timeStamp, &s.Curr); err != nil {
		return err
	}

	s.Prev = store.NewSample()
	if err := s.Store.NextSample(-1, &s.Prev); err != nil {
		return err
	}
	s.Curr = store.NewSample()
	if err := s.Store.NextSample(1, &s.Curr); err != nil {
		return err
	}
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Info("skip one sample since system reboot")
		return s.CollectNext()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectField() {

	s.Sys.Collect(&s.Prev, &s.Curr)
	s.CPUs.Collect(&s.Prev, &s.Curr)
	s.MEM.Collect(&s.Prev, &s.Curr)
	s.Vm.Collect(&s.Prev, &s.Curr)
	s.Disks.Collect(&s.Prev, &s.Curr)
	s.Nets.Collect(&s.Prev, &s.Curr)
	s.NetStat.Collect(&s.Prev, &s.Curr)
	s.NetProtocols.Collect(&s.Prev, &s.Curr)
	s.Softnets.Collect(&s.Prev, &s.Curr)
	s.Sys.Processes, s.Sys.Threads = s.Processes.Collect(&s.Prev, &s.Curr)
	s.Cgroup.Collect(s.Prev.CgroupSample, s.Curr.CgroupSample, s.Curr.TimeStamp-s.Prev.TimeStamp)
}

type DumpOption struct {
	Begin           int64
	End             int64
	Module          string
	Output          *os.File
	Format          string
	Fields          []string
	FilterText      string
	FilterProgram   *vm.Program
	SortField       string
	DescendingOrder bool
	Top             int
	DisableTitle    bool
	RepeatTitle     int
	RawData         bool
}

func (s *Model) Dump(opt DumpOption) error {

	for _, c := range opt.Fields {
		name, _ := getNameAndWidthOfField(opt.Module, c)
		if name == "" {
			return fmt.Errorf("%s is not available field for module %s", c, opt.Module)
		}
	}
	if err := verifyFilterText(&opt); err != nil {
		return err
	}

	switch opt.Format {
	case "text":
		return s.dumpText(opt)
	case "json":
		return s.dumpJson(opt)
	default:
		return fmt.Errorf("no support output format: %s", opt.Format)
	}
}

type DumpOtelOption struct {
	Begin  int64
	End    int64
	Output string
}

func (s *Model) DumpToOtel(opt DumpOtelOption) (err error) {
	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}
	ctx := context.Background()
	var export metric.Exporter
	progress := false
	if opt.Output == "http" {
		export, err = otlpmetrichttp.New(ctx)
		progress = true
	} else {
		export, err = stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
		)
	}

	// export, err := stdoutmetric.New(
	// 	stdoutmetric.WithPrettyPrint(),
	// )
	if err != nil {
		return err
	}

	cnt := 0
	rm := metricdata.ResourceMetrics{
		Resource: resource.NewSchemaless(
			semconv.ServiceName("etop"),
			semconv.ServiceInstanceID(s.Curr.HostName),
		),
	}

	for opt.End >= s.Curr.TimeStamp {
		var sm metricdata.ScopeMetrics
		s.CPUs.GetOtelMetric(s.Curr.TimeStamp, &sm)
		s.MEM.GetOtelMetric(s.Curr.TimeStamp, &sm)
		s.Disks.GetOtelMetric(s.Curr.TimeStamp, &sm)
		s.Nets.GetOtelMetric(s.Curr.TimeStamp, &sm)
		s.Processes.GetOtelMetric(s.Curr.TimeStamp, &sm)
		rm.ScopeMetrics = append(rm.ScopeMetrics, sm)
		if err := s.CollectNext(); err != nil {
			if err == store.ErrOutOfRange {
				break
			}
			return err
		}

		if err := export.Export(ctx, &rm); err != nil {
			return err
		}
		cnt++
		if progress {
			fmt.Printf("Dump Otel data from %s to %s: %d samples\033[0K\r",
				time.Unix(opt.Begin, 0).Format(time.RFC3339),
				time.Unix(opt.End, 0).Format(time.RFC3339), cnt)
		}
		rm.ScopeMetrics = rm.ScopeMetrics[:0]
	}
	if progress {
		fmt.Printf("\n")
	}
	return export.Shutdown(ctx)
}

func verifyFilterText(opt *DumpOption) (err error) {
	if opt.FilterText == "" {
		return nil
	}
	var s any
	switch opt.Module {
	case "system":
		s = &System{}
	case "cpu":
		s = &CPU{}
	case "memory":
		s = &MEM{}
	case "vm":
		s = &Vm{}
	case "disk":
		s = &Disk{}
	case "netdev":
		s = &NetDev{}
	case "network":
		s = &NetStat{}
	case "networkprotocol":
		s = &NetProtocol{}
	case "softnet":
		s = &Softnet{}
	case "process":
		s = &Process{}
	case "cgroup":
		s = &Cgroup{}
	}
	opt.FilterProgram, err = expr.Compile(opt.FilterText, expr.Env(s), expr.AsBool())
	return err
}

func isFilter(opt DumpOption, m Render) bool {
	if opt.FilterProgram != nil {
		output, _ := expr.Run(opt.FilterProgram, m)
		return output.(bool)
	}
	return true
}

func getNameAndWidthOfField(module string, f string) (string, int) {
	var s Render
	switch module {
	case "system":
		s = &System{}
	case "cpu":
		s = &CPU{}
	case "memory":
		s = &MEM{}
	case "vm":
		s = &Vm{}
	case "disk":
		s = &Disk{}
	case "netdev":
		s = &NetDev{}
	case "network":
		s = &NetStat{}
	case "networkprotocol":
		s = &NetProtocol{}
	case "softnet":
		s = &Softnet{}
	case "process":
		s = &Process{}
	case "cgroup":
		s = &Cgroup{}
	}
	cfg := s.DefaultConfig(f)
	return cfg.Name, cfg.Width
}

func (s *Model) dumpText(opt DumpOption) error {

	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}

	title := fmt.Sprintf("%-25s", "TimeStamp")
	for _, c := range opt.Fields {
		name, width := getNameAndWidthOfField(opt.Module, c)
		if len(name) > width {
			width = len(name)
		}
		title += fmt.Sprintf(" %-*s", width, name)
	}
	title += "\n"
	if !opt.DisableTitle {
		opt.Output.WriteString(title)
	}
	cnt := 0
	for opt.End >= s.Curr.TimeStamp {
		if !opt.DisableTitle && opt.RepeatTitle != 0 && cnt%opt.RepeatTitle == 0 {
			opt.Output.WriteString(title)
		}
		switch opt.Module {
		case "system":
			dumpText(s.Curr.TimeStamp, opt, &s.Sys)
		case "cpu":
			for _, c := range s.CPUs {
				dumpText(s.Curr.TimeStamp, opt, &c)
			}
		case "memory":
			dumpText(s.Curr.TimeStamp, opt, &s.MEM)
		case "vm":
			dumpText(s.Curr.TimeStamp, opt, &s.Vm)
		case "disk":
			for _, disk := range s.Disks.GetKeys() {
				d := s.Disks[disk]
				dumpText(s.Curr.TimeStamp, opt, &d)
			}
		case "netdev":
			for _, dev := range s.Nets.GetKeys() {
				n := s.Nets[dev]
				dumpText(s.Curr.TimeStamp, opt, &n)
			}
		case "network":
			dumpText(s.Curr.TimeStamp, opt, &s.NetStat)
		case "networkprotocol":
			for _, n := range s.NetProtocols {
				dumpText(s.Curr.TimeStamp, opt, &n)
			}
		case "softnet":
			for _, soft := range s.Softnets {
				dumpText(s.Curr.TimeStamp, opt, &soft)
			}
		case "process":
			processList := s.Processes.Iterate(nil, opt.SortField, opt.DescendingOrder)
			cnt := 0
			for _, p := range processList {
				dumpText(s.Curr.TimeStamp, opt, &p)
				cnt++
				if opt.Top > 0 && opt.Top == cnt {
					break
				}
			}
		case "cgroup":
			dumpTextForCgroup(s.Curr.TimeStamp, opt, s.Cgroup)
		}
		if err := s.CollectNext(); err != nil {
			if err == store.ErrOutOfRange {
				return nil
			}
			return err
		}
		cnt++
	}
	return nil

}

func (s *Model) dumpJson(opt DumpOption) error {

	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}

	opt.Output.WriteString("[\n")
	first := true
	for opt.End >= s.Curr.TimeStamp {
		if first {
			first = false
		} else {
			opt.Output.WriteString(",\n\n")
		}
		switch opt.Module {
		case "system":
			if isFilter(opt, &s.Sys) {
				dumpJson(s.Curr.TimeStamp, opt, &s.Sys)
			}
		case "cpu":
			opt.Output.WriteString("[")
			first := true
			for _, c := range s.CPUs {
				if isFilter(opt, &c) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &c)
				}
			}
			opt.Output.WriteString("]")
		case "memory":
			if isFilter(opt, &s.MEM) {
				dumpJson(s.Curr.TimeStamp, opt, &s.MEM)
			}
		case "vm":
			if isFilter(opt, &s.Vm) {
				dumpJson(s.Curr.TimeStamp, opt, &s.Vm)
			}
		case "disk":
			opt.Output.WriteString("[")
			first := true
			for _, disk := range s.Disks.GetKeys() {
				d := s.Disks[disk]
				if isFilter(opt, &d) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &d)
				}
			}
			opt.Output.WriteString("]")
		case "netdev":
			opt.Output.WriteString("[")
			first := true
			for _, dev := range s.Nets.GetKeys() {
				n := s.Nets[dev]
				if isFilter(opt, &n) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &n)
				}
			}
			opt.Output.WriteString("]")
		case "network":
			if isFilter(opt, &s.NetStat) {
				dumpJson(s.Curr.TimeStamp, opt, &s.NetStat)
			}
		case "networkprotocol":
			opt.Output.WriteString("[")
			first := true
			for _, n := range s.NetProtocols {
				if isFilter(opt, &n) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &n)
				}
			}
			opt.Output.WriteString("]")
		case "softnet":
			opt.Output.WriteString("[")
			first := true
			for _, soft := range s.Softnets {
				if isFilter(opt, &soft) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &soft)
				}
			}
			opt.Output.WriteString("]")
		case "process":
			processList := s.Processes.Iterate(nil, opt.SortField, opt.DescendingOrder)
			cnt := 0
			opt.Output.WriteString("[")
			first := true
			for _, p := range processList {
				if isFilter(opt, &p) {
					if first {
						first = false
					} else {
						opt.Output.WriteString(",\n")
					}
					dumpJson(s.Curr.TimeStamp, opt, &p)
					cnt++
					if opt.Top > 0 && opt.Top == cnt {
						break
					}
				}
			}
			opt.Output.WriteString("]")
		case "cgroup":
			re := dumpJsonForCgroup(s.Curr.TimeStamp, opt, s.Cgroup)
			b, _ := json.Marshal(re)
			opt.Output.Write(b)
		}
		if err := s.CollectNext(); err != nil {
			if err == store.ErrOutOfRange {
				opt.Output.WriteString("]\n")
				return nil
			}
			return err
		}
	}
	opt.Output.WriteString("\n]\n")
	return nil

}
