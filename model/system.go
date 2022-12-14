package model

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/xixiliguo/etop/store"
)

type System struct {
	Config        map[string]RenderConfig
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
	Disks       DiskSlice
	Nets        NetSlice
	ProcessList ProcessSlice
}

func NewSysModel(s *store.LocalStore, log *log.Logger) (*System, error) {
	p := &System{
		Config:      DefaultRenderConfig,
		Store:       s,
		log:         log,
		CPUs:        []CPU{},
		MEM:         MEM{},
		ProcessList: []Process{},
	}
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if err := p.CollectSampleByTime(t.Unix()); err != nil {
		return nil, err
	}
	return p, nil
}

func NewSysModelWithLive(log *log.Logger) (*System, error) {
	p := &System{
		Config:      DefaultRenderConfig,
		log:         log,
		CPUs:        []CPU{},
		MEM:         MEM{},
		ProcessList: []Process{},
	}
	return p, nil
}

func (s *System) CollectLiveSample() error {

	s.Prev = s.Curr
	s.Curr = store.Sample{}
	if err := store.CollectSampleFromSys(&s.Curr); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *System) CollectNext() error {
	next := store.Sample{}
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

func (s *System) CollectPrev() error {

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

func (s *System) CollectSampleByTime(timeStamp int64) error {

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

func (s *System) CollectField() {

	s.CPUs.Collect(&s.Prev, &s.Curr)
	s.MEM.Collect(&s.Prev, &s.Curr)
	s.Disks.Collect(&s.Prev, &s.Curr)
	s.Nets.Collect(&s.Prev, &s.Curr)
	s.Prcesses, s.Threads = s.ProcessList.Collect(&s.Prev, &s.Curr)
	s.Clones = (s.Curr.ProcessCreated - s.Prev.ProcessCreated) / uint64(s.Curr.TimeStamp-s.Prev.TimeStamp)
	s.ContextSwitch = (s.Curr.ContextSwitches - s.Prev.ContextSwitches) / uint64(s.Curr.TimeStamp-s.Prev.TimeStamp)

}

type DumpOption struct {
	Begin        int64
	End          int64
	Module       string
	Output       *os.File
	Format       string
	Fields       []string
	SelectField  string
	Filter       *regexp.Regexp
	DisableTitle bool
	RepeatTitle  int
	RawData      bool
}

func (s *System) Dump(opt DumpOption) error {

	switch opt.Format {
	case "text":
		return s.dumpText(s.Config[opt.Module], opt)
	case "json":
		return s.dumpJson(s.Config[opt.Module], opt)
	default:
		return fmt.Errorf("no support output format: %s", opt.Format)
	}
}

func (s *System) dumpText(config RenderConfig, opt DumpOption) error {

	if err := s.CollectSampleByTime(opt.Begin); err != nil {
		return err
	}

	title := fmt.Sprintf("%25s", "TimeStamp")
	for _, c := range opt.Fields {
		title += fmt.Sprintf("%*s", config[c].Width, config[c].Name)
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

func (s *System) dumpJson(config RenderConfig, opt DumpOption) error {

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
