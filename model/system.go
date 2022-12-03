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
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if err := p.CollectSampleByTime(t.Unix()); err != nil {
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
