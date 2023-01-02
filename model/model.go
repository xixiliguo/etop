package model

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/xixiliguo/etop/store"
)

type Model struct {
	Config        map[string]RenderConfig
	Mode          string
	Store         store.Store
	log           *log.Logger
	Prev          store.Sample
	Curr          store.Sample
	Prcesses      int
	Threads       int
	Clones        uint64
	ContextSwitch uint64
	CPUs          CPUSlice
	MEM
	Disks DiskMap
	Nets  NetDevMap
	NetStat
	NetProtocols NetProtocolMap
	Softnets     SoftnetSlice
	Processes    ProcessMap
}

func NewSysModel(s *store.LocalStore, log *log.Logger) (*Model, error) {
	p := &Model{
		Config:       DefaultRenderConfig,
		Mode:         "report",
		Store:        s,
		log:          log,
		Prev:         store.NewSample(),
		Curr:         store.NewSample(),
		CPUs:         []CPU{},
		MEM:          MEM{},
		Disks:        make(DiskMap),
		Nets:         make(NetDevMap),
		NetProtocols: make(NetProtocolMap),
		Softnets:     []Softnet{},
		Processes:    make(ProcessMap),
	}
	return p, nil
}

func (s *Model) CollectLiveSample() error {

	s.Prev = s.Curr
	s.Curr = store.NewSample()
	if err := store.CollectSampleFromSys(&s.Curr); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectNext() error {
	next := store.NewSample()
	if err := s.Store.NextSample(1, &next); err != nil {
		return err
	}
	s.Prev = s.Curr
	s.Curr = next
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Printf("skip one sample since system reboot")
		return s.CollectNext()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectPrev() error {

	if err := s.Store.NextSample(-2, &s.Prev); err != nil {
		return err
	}

	if err := s.Store.NextSample(1, &s.Curr); err != nil {
		return err
	}
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Printf("skip one sample since system reboot")
		return s.CollectPrev()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectSampleByTime(timeStamp int64) error {

	if err := s.Store.JumpSampleByTimeStamp(timeStamp, &s.Curr); err != nil {
		return err
	}
	if err := s.Store.NextSample(-1, &s.Prev); err != nil {
		return err
	}

	if err := s.Store.NextSample(1, &s.Curr); err != nil {
		return err
	}
	if s.Curr.BootTime != s.Prev.BootTime {
		//system ever reboot, skip one sample
		s.log.Printf("skip one sample since system reboot")
		return s.CollectNext()
	}
	s.CollectField()
	return nil
}

func (s *Model) CollectField() {

	s.CPUs.Collect(&s.Prev, &s.Curr)
	s.MEM.Collect(&s.Prev, &s.Curr)
	s.Disks.Collect(&s.Prev, &s.Curr)
	s.Nets.Collect(&s.Prev, &s.Curr)
	s.NetStat.Collect(&s.Prev, &s.Curr)
	s.NetProtocols.Collect(&s.Prev, &s.Curr)
	s.Softnets.Collect(&s.Prev, &s.Curr)
	s.Prcesses, s.Threads = s.Processes.Collect(&s.Prev, &s.Curr)
	s.Clones = (s.Curr.ProcessCreated - s.Prev.ProcessCreated) / uint64(s.Curr.TimeStamp-s.Prev.TimeStamp)
	s.ContextSwitch = (s.Curr.ContextSwitches - s.Prev.ContextSwitches) / uint64(s.Curr.TimeStamp-s.Prev.TimeStamp)

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
	cnt := 0
	for opt.End >= s.Curr.TimeStamp {
		if opt.DisableTitle == false && opt.RepeatTitle != 0 && cnt%opt.RepeatTitle == 0 {
			opt.Output.WriteString(title)
		}
		switch opt.Module {
		case "cpu":
			s.CPUs.Dump(s.Curr.TimeStamp, config, opt)
		case "memory":
			s.MEM.Dump(s.Curr.TimeStamp, config, opt)
		case "disk":
			s.Disks.Dump(s.Curr.TimeStamp, config, opt)
		case "netdev":
			s.Nets.Dump(s.Curr.TimeStamp, config, opt)
		case "network":
			s.NetStat.Dump(s.Curr.TimeStamp, config, opt)
		case "networkprotocol":
			s.NetProtocols.Dump(s.Curr.TimeStamp, config, opt)
		case "softnet":
			s.Softnets.Dump(s.Curr.TimeStamp, config, opt)
		case "process":
			s.Processes.Dump(s.Curr.TimeStamp, config, opt)
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

	opt.Output.WriteString("[")
	first := true
	for opt.End >= s.Curr.TimeStamp {
		if first == true {
			first = false
		} else {
			opt.Output.WriteString("\n,")
		}
		switch opt.Module {
		case "cpu":
			s.CPUs.Dump(s.Curr.TimeStamp, config, opt)
		case "memory":
			s.MEM.Dump(s.Curr.TimeStamp, config, opt)
		case "disk":
			s.Disks.Dump(s.Curr.TimeStamp, config, opt)
		case "netdev":
			s.Nets.Dump(s.Curr.TimeStamp, config, opt)
		case "network":
			s.NetStat.Dump(s.Curr.TimeStamp, config, opt)
		case "networkprotocol":
			s.NetProtocols.Dump(s.Curr.TimeStamp, config, opt)
		case "softnet":
			s.Softnets.Dump(s.Curr.TimeStamp, config, opt)
		case "process":
			s.Processes.Dump(s.Curr.TimeStamp, config, opt)
		}
		if err := s.CollectNext(); err != nil {
			if err == store.ErrOutOfRange {
				opt.Output.WriteString("]\n")
				return nil
			}
			return err
		}
	}
	opt.Output.WriteString("]\n")
	return nil

}
