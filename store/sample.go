package store

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/xixiliguo/etop/cgroupfs"
	"github.com/xixiliguo/etop/procfs"
	"golang.org/x/sys/unix"
)

var bootTimeTick uint64

// Sample represent all system info and process info.
type Sample struct {
	TimeStamp    int64  // unix time when sample was generated
	SystemSample        // system information
	ProcSamples  PidMap // process information
	CgroupSample CgroupSample
}

type SystemSample struct {
	HostName      string
	KernelVersion string
	PageSize      int
	BootTimeTick  uint64
	procfs.LoadAvg
	procfs.Stat
	procfs.Meminfo
	procfs.VmStat
	NetDevStats procfs.NetDev
	DiskStats   procfs.DiskStat
	procfs.NetProtocolStats
	SoftNetStats []procfs.SoftnetStat
}

type PidMap map[int]ProcSample

func (p PidMap) mergeWithExitProcess(e *ExitProcess) {
	e.Lock()
	for pid, s := range e.Samples {
		if _, ok := p[pid]; ok {
		} else {
			p[pid] = s
		}
		delete(e.Samples, pid)
	}
	e.Unlock()
}

type ProcSample struct {
	procfs.ProcStat
	procfs.ProcIO
	procfs.ProcSchedstat
	CmdLine  string
	Cgroup   string
	EndTime  uint64
	ExitCode uint64
}

func NewSample() Sample {
	s := Sample{
		TimeStamp: 0,
		SystemSample: SystemSample{
			NetDevStats:      make(procfs.NetDev),
			DiskStats:        make(procfs.DiskStat),
			NetProtocolStats: make(procfs.NetProtocolStats),
		},
		ProcSamples: make(PidMap),
	}
	return s
}

func (s *Sample) Reset() {

	clear(s.NetDevStats)
	clear(s.DiskStats)
	clear(s.NetProtocolStats)
	clear(s.ProcSamples)
	return
}

func (s *Sample) Marshal() ([]byte, error) {
	return cbor.Marshal(s)
}

func (s *Sample) Unmarshal(b []byte) error {
	return cbor.Unmarshal(b, s)
}

func CollectSampleFromSys(s *Sample, exit *ExitProcess, log *slog.Logger) error {

	//collect one sample
	var (
		err error
	)
	s.TimeStamp = time.Now().Unix()
	u := unix.Utsname{}
	unix.Uname(&u)

	newFS := procfs.NewFS("")

	s.HostName = unix.ByteSliceToString(u.Nodename[:])
	s.KernelVersion = unix.ByteSliceToString(u.Release[:])
	s.PageSize = os.Getpagesize()
	s.BootTimeTick = bootTimeTick

	if s.LoadAvg, err = newFS.Load(); err != nil {
		return err
	}

	if s.Stat, err = newFS.Stat(); err != nil {
		return err
	}

	if s.Meminfo, err = newFS.Meminfo(); err != nil {
		return err
	}

	if s.VmStat, err = newFS.VmStat(); err != nil {
		return err
	}

	if s.NetDevStats, err = newFS.NetDev(); err != nil {
		return err
	}

	if s.NetProtocolStats, err = newFS.NetProtocols(); err != nil {
		return err
	}

	if s.SoftNetStats, err = newFS.NetSoftnetStat(); err != nil {
		return err
	}

	if s.DiskStats, err = newFS.DiskStat(); err != nil {
		return err
	}

	err = newFS.EachProc(func(proc procfs.Proc) error {
		p := ProcSample{}
		var err error
		var pathErr *os.PathError
		if p.ProcStat, err = proc.Stat(); err != nil && !errors.As(err, &pathErr) {

			return err
		}
		if p.ProcIO, err = proc.IO(); err != nil && !errors.As(err, &pathErr) {

			return err
		}
		if p.ProcSchedstat, err = proc.Schedstat(); err != nil && !errors.As(err, &pathErr) {

			return err
		}
		if p.CmdLine, err = proc.CmdLine(); err != nil && !errors.As(err, &pathErr) {

			return err
		}
		if isCgroup2() {
			if p.Cgroup, err = proc.Cgroup(); err != nil && !errors.As(err, &pathErr) {
				return err
			}
		}
		s.ProcSamples[p.PID] = p
		return nil
	})

	if err != nil {
		return err
	}

	if exit != nil {
		s.ProcSamples.mergeWithExitProcess(exit)
	}

	// collect cgroupv2 if enabled
	if isCgroup2() {

		cgRoot := cgroupfs.NewCgroup("/", "/")
		s.CgroupSample, err = walkCgroupNode(0, cgRoot)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	ts := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_REALTIME, &ts)
	ts1 := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_BOOTTIME, &ts1)
	bootTimeTick = uint64(ts.Nano()-ts1.Nano()) / 10000000
}
