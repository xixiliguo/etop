package store

import (
	"bytes"
	"encoding/binary"
	"os"
	"sync"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/prometheus/procfs"
	"golang.org/x/exp/slog"
	"golang.org/x/sys/unix"
)

//go:generate bpf2go -cc clang -cflags $BPF_CFLAGS -type event bpf exit.bpf.c -- -I./include

type ExitProcess struct {
	sync.Mutex
	Samples map[int]ProcSample
	log     *slog.Logger
}

func NewExitProcess(log *slog.Logger) *ExitProcess {
	return &ExitProcess{
		Samples: make(PidMap),
		log:     log,
	}
}

func (e *ExitProcess) Collect() {

	pageSize := uint(os.Getpagesize())

	if err := rlimit.RemoveMemlock(); err != nil {
		e.log.Error("remove Memlock: ", err)
		return
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		e.log.Error("loading objects: ", err)
		return
	}
	defer objs.Close()

	kp, err := link.Tracepoint("sched", "sched_process_exit", objs.HandleExit, nil)
	if err != nil {
		e.log.Error("opening tracepoint: ", err)
		return
	}
	defer kp.Close()

	rd, err := perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		e.log.Error("creating event reader: ", err)
		return
	}
	defer rd.Close()

	var event bpfEvent
	var record perf.Record
	for {
		err := rd.ReadInto(&record)
		if err != nil {
			e.log.Error("reading from reader: ", err)
			return
		}

		if record.LostSamples != 0 {
			e.log.Info("perf event lost %d samples", record.LostSamples)
			continue
		}

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			e.log.Error("parsing ringbuf event: ", err)
			continue
		}

		e.Lock()
		e.Samples[int(event.Pid)] = ProcSample{
			ProcStat: procfs.ProcStat{
				PID:                 int(event.Pid),
				Comm:                unix.ByteSliceToString(event.Comm[:]),
				State:               "X",
				PPID:                int(event.Ppid),
				PGRP:                0,
				Session:             0,
				TTY:                 0,
				TPGID:               0,
				Flags:               0,
				MinFlt:              uint(event.MinFlt),
				CMinFlt:             0,
				MajFlt:              uint(event.MajFlt),
				CMajFlt:             0,
				UTime:               uint(event.Utime),
				STime:               uint(event.Stime),
				CUTime:              0,
				CSTime:              0,
				Priority:            int(event.Priority),
				Nice:                int(event.Nice),
				NumThreads:          1,
				Starttime:           event.StartTime,
				VSize:               uint(event.VssPages) * pageSize,
				RSS:                 int(event.RssPages),
				RSSLimit:            0,
				Processor:           uint(event.OnCpu),
				RTPriority:          0,
				Policy:              0,
				DelayAcctBlkIOTicks: event.DelayacctBlkioTicks,
			},
			ProcIO: procfs.ProcIO{
				RChar:               event.Rchar,
				WChar:               event.Wchar,
				SyscR:               event.Syscr,
				SyscW:               event.Syscw,
				ReadBytes:           event.IoReadBytes,
				WriteBytes:          event.IoWriteBytes,
				CancelledWriteBytes: int64(event.CancelledWriteBytes),
			},
		}
		e.Unlock()
	}
}
