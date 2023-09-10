package model

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"

	"github.com/xixiliguo/etop/store"
)

type Model struct {
	Config map[string]RenderConfig
	Mode   string
	Store  store.Store
	log    *slog.Logger
	Prev   store.Sample
	Curr   store.Sample
	Sys    System
	CPUs   CPUSlice
	MEM
	Vm
	Disks DiskMap
	Nets  NetDevMap
	NetStat
	NetProtocols NetProtocolMap
	Softnets     SoftnetSlice
	Processes    ProcessMap
}

func NewSysModel(s *store.LocalStore, log *slog.Logger) (*Model, error) {
	p := &Model{
		Config:       DefaultRenderConfig,
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
	}
	return p, nil
}

func (s *Model) CollectLiveSample(exit *store.ExitProcess) error {

	s.Prev = s.Curr
	s.Curr = store.NewSample()
	if err := store.CollectSampleFromSys(&s.Curr, exit); err != nil {
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

}

type DumpOption struct {
	Begin          int64
	End            int64
	Module         string
	Output         *os.File
	Format         string
	Fields         []string
	SelectField    string
	Filter         *regexp.Regexp
	SortField      string
	AscendingOrder bool
	Top            int
	DisableTitle   bool
	RepeatTitle    int
	RawData        bool
}

func (s *Model) Dump(opt DumpOption) error {

	switch opt.Format {
	case "text":
		return s.dumpText(s.Config[opt.Module], opt)
	case "json":
		return s.dumpJson(s.Config[opt.Module], opt)
	default:
		return fmt.Errorf("no support output format: %s", opt.Format)
	}
}

func (s *Model) dumpText(config RenderConfig, opt DumpOption) error {

	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}

	title := fmt.Sprintf("%-25s", "TimeStamp")
	for _, c := range opt.Fields {
		width := config[c].Width
		if len(config[c].Name) > width {
			width = len(config[c].Name)
		}
		title += fmt.Sprintf(" %-*s", width, config[c].Name)
	}
	title += "\n"
	if opt.DisableTitle == false {
		opt.Output.WriteString(title)
	}
	config.SetFixWidth(true)
	cnt := 0
	for opt.End >= s.Curr.TimeStamp {
		if opt.DisableTitle == false && opt.RepeatTitle != 0 && cnt%opt.RepeatTitle == 0 {
			opt.Output.WriteString(title)
		}
		switch opt.Module {
		case "system":
			dumpText(s.Curr.TimeStamp, config, opt, &s.Sys)
		case "cpu":
			for _, c := range s.CPUs {
				dumpText(s.Curr.TimeStamp, config, opt, &c)
			}
		case "memory":
			dumpText(s.Curr.TimeStamp, config, opt, &s.MEM)
		case "vm":
			dumpText(s.Curr.TimeStamp, config, opt, &s.Vm)
		case "disk":
			for _, disk := range s.Disks.GetKeys() {
				d := s.Disks[disk]
				dumpText(s.Curr.TimeStamp, config, opt, &d)
			}
		case "netdev":
			for _, dev := range s.Nets.GetKeys() {
				n := s.Nets[dev]
				dumpText(s.Curr.TimeStamp, config, opt, &n)
			}
		case "network":
			dumpText(s.Curr.TimeStamp, config, opt, &s.NetStat)
		case "networkprotocol":
			for _, n := range s.NetProtocols {
				dumpText(s.Curr.TimeStamp, config, opt, &n)
			}
		case "softnet":
			for _, soft := range s.Softnets {
				dumpText(s.Curr.TimeStamp, config, opt, &soft)
			}
		case "process":
			processList := []Process{}
			for _, p := range s.Processes {
				processList = append(processList, p)
			}

			sort.SliceStable(processList, func(i, j int) bool {
				return SortMap[opt.SortField](processList[i], processList[j])
			})
			if opt.AscendingOrder == true {
				for i := 0; i < len(processList)/2; i++ {
					processList[i], processList[len(processList)-1-i] = processList[len(processList)-1-i], processList[i]
				}
			}
			cnt := 0
			for _, p := range processList {
				dumpText(s.Curr.TimeStamp, config, opt, &p)
				cnt++
				if opt.Top > 0 && opt.Top == cnt {
					break
				}
			}
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

func (s *Model) dumpJson(config RenderConfig, opt DumpOption) error {

	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}

	opt.Output.WriteString("[\n")
	first := true
	for opt.End >= s.Curr.TimeStamp {
		if first == true {
			first = false
		} else {
			opt.Output.WriteString(",\n\n")
		}
		switch opt.Module {
		case "system":
			dumpJson(s.Curr.TimeStamp, config, opt, &s.Sys)
		case "cpu":
			opt.Output.WriteString("[")
			first := true
			for _, c := range s.CPUs {
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &c)
			}
			opt.Output.WriteString("]")
		case "memory":
			dumpJson(s.Curr.TimeStamp, config, opt, &s.MEM)
		case "vm":
			dumpJson(s.Curr.TimeStamp, config, opt, &s.Vm)
		case "disk":
			opt.Output.WriteString("[")
			first := true
			for _, disk := range s.Disks.GetKeys() {
				d := s.Disks[disk]
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &d)
			}
			opt.Output.WriteString("]")
		case "netdev":
			opt.Output.WriteString("[")
			first := true
			for _, dev := range s.Nets.GetKeys() {
				n := s.Nets[dev]
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &n)
			}
			opt.Output.WriteString("]")
		case "network":
			dumpJson(s.Curr.TimeStamp, config, opt, &s.NetStat)
		case "networkprotocol":
			opt.Output.WriteString("[")
			first := true
			for _, n := range s.NetProtocols {
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &n)
			}
			opt.Output.WriteString("]")
		case "softnet":
			opt.Output.WriteString("[")
			first := true
			for _, soft := range s.Softnets {
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &soft)
			}
			opt.Output.WriteString("]")
		case "process":
			processList := []Process{}
			for _, p := range s.Processes {
				processList = append(processList, p)
			}

			sort.SliceStable(processList, func(i, j int) bool {
				return SortMap[opt.SortField](processList[i], processList[j])
			})
			if opt.AscendingOrder == true {
				for i := 0; i < len(processList)/2; i++ {
					processList[i], processList[len(processList)-1-i] = processList[len(processList)-1-i], processList[i]
				}
			}
			cnt := 0
			opt.Output.WriteString("[")
			first := true
			for _, p := range processList {
				if first == true {
					first = false
				} else {
					opt.Output.WriteString(",\n")
				}
				dumpJson(s.Curr.TimeStamp, config, opt, &p)
				cnt++
				if opt.Top > 0 && opt.Top == cnt {
					break
				}
			}
			opt.Output.WriteString("]")
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
