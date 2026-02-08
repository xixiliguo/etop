package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/xixiliguo/etop/procfs"
	"golang.org/x/sys/unix"
)

//go:generate go tool bpf2go -cc clang -cflags $BPF_CFLAGS -type event process exit.bpf.c -- -I./include

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

	objs := processObjects{}
	if err := loadProcessObjects(&objs, nil); err != nil {
		msg := fmt.Sprintf("loading objects: %s", err)
		e.log.Error(msg)
		return
	}
	defer objs.Close()
	btf.FlushKernelSpec()
	debug.FreeOSMemory()

	kp, err := link.Kprobe("acct_process", objs.HandleExit, nil)
	if err != nil {
		msg := fmt.Sprintf("opening kprobe: %s", err)
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

	var event processEvent
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
				State:               procfs.Dead,
				PPID:                int(event.Ppid),
				PGRP:                0,
				Session:             0,
				TTY:                 0,
				TPGID:               0,
				Flags:               0,
				MinFlt:              event.MinFlt,
				CMinFlt:             0,
				MajFlt:              event.MajFlt,
				CMajFlt:             0,
				UTime:               event.Utime,
				STime:               event.Stime,
				CUTime:              0,
				CSTime:              0,
				Priority:            int(event.Priority),
				Nice:                int(event.Nice),
				NumThreads:          int(event.NumThreads),
				Starttime:           event.StartTime,
				VSize:               event.VssPages * uint64(pageSize),
				RSS:                 event.RssPages,
				RSSLimit:            0,
				Processor:           int(event.OnCpu),
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
				CancelledWriteBytes: event.CancelledWriteBytes,
			},
			EndTime:  event.EndTime,
			ExitCode: uint64(event.ExitCode),
		}
		e.Unlock()
	}
}
