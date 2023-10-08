package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/prometheus/procfs"
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
		msg := fmt.Sprintf("remove Memlock: %s", err)
		e.log.Error(msg)
		return
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		msg := fmt.Sprintf("loading objects: %s", err)
		e.log.Error(msg)
		return
	}
	defer objs.Close()
	btf.FlushKernelSpec()

	kp, err := link.Tracepoint("sched", "sched_process_exit", objs.HandleExit, nil)
	if err != nil {
		msg := fmt.Sprintf("opening tracepoint: %s", err)
		e.log.Error(msg)
		return
	}
	defer kp.Close()

	rd, err := perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		msg := fmt.Sprintf("creating event reader: %s", err)
		e.log.Error(msg)
		return
	}
	defer rd.Close()

	var event bpfEvent
	var record perf.Record
	for {
		err := rd.ReadInto(&record)
		if err != nil {
			msg := fmt.Sprintf("reading from reader: %s", err)
			e.log.Error(msg)
			return
		}

		if record.LostSamples != 0 {
			msg := fmt.Sprintf("perf event lost %d samples", record.LostSamples)
			e.log.Info(msg)
			continue
		}

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			msg := fmt.Sprintf("parsing ringbuf event: %s", err)
			e.log.Error(msg)
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
				NumThreads:          int(event.NumThreads),
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
			EndTime:  event.EndTime,
			ExitCode: event.ExitCode,
		}
		e.Unlock()
	}
}
