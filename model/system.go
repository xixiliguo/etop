package model

import (
	"log"
	"time"

	"github.com/xixiliguo/etop/store"
)

type System struct {
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
		Store:       s,
		log:         log,
		CPUs:        []CPU{},
		MEM:         MEM{},
		ProcessList: []Process{},
	}
	if err := p.CollectNext(); err != nil {
		return nil, err
	}
	return p, nil
}

func NewSysModelWithLive(log *log.Logger) (*System, error) {
	p := &System{
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
	s.Curr.CurrTime = time.Now().Unix()
	s.CollectField()
	return nil
}

func (s *System) CollectNext() error {

	if err := s.CollectSampleByStep(0, &s.Prev); err != nil {
		return err
	}
	if err := s.CollectSampleByStep(1, &s.Curr); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *System) CollectPrev() error {

	if err := s.CollectSampleByStep(-2, &s.Prev); err != nil {
		return err
	}

	if err := s.CollectSampleByStep(1, &s.Curr); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *System) CollectSampleByTime(value string) error {

	if err := s.Store.ChangeIndex(value); err != nil {
		return err
	}
	if err := s.CollectSampleByStep(-1, &s.Prev); err != nil {
		return err
	}

	if err := s.CollectSampleByStep(1, &s.Curr); err != nil {
		return err
	}
	s.CollectField()
	return nil
}

func (s *System) CollectSampleByStep(step int, sample *store.Sample) error {

	if err := s.Store.AdjustIndex(step); err != nil {
		return err
	}

	if err := s.Store.ReadSample(sample); err != nil {
		return err
	}
	return nil
}

func (s *System) CollectField() {

	s.CPUs.Collect(&s.Prev, &s.Curr)
	s.MEM.Collect(&s.Prev, &s.Curr)
	s.Disks.Collect(&s.Prev, &s.Curr)
	s.Nets.Collect(&s.Prev, &s.Curr)
	s.Prcesses, s.Threads = s.ProcessList.Collect(&s.Prev, &s.Curr)
	s.Clones = (s.Curr.ProcessCreated - s.Prev.ProcessCreated) / uint64(s.Curr.CurrTime-s.Prev.CurrTime)
	s.ContextSwitch = (s.Curr.ContextSwitches - s.Prev.ContextSwitches) / uint64(s.Curr.CurrTime-s.Prev.CurrTime)

}
