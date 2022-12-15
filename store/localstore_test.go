package store

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
)

func TestIndexMarshalAndUnmarshal(t *testing.T) {
	testCases := []Index{
		{
			TimeStamp: 0,
			Offset:    0,
			Len:       0,
		},
		{
			TimeStamp: 1,
			Offset:    2,
			Len:       3,
		},
		{
			TimeStamp: 99999,
			Offset:    999,
			Len:       0,
		},
	}
	for _, testCase := range testCases {
		re := Index{}
		re.Unmarshal(testCase.Marshal())
		if cmp.Equal(testCase, re) == false {
			t.Errorf("got %+v\nwant %+v\n", re, testCase)
		}
	}
}

func TestSampleMarshalAndUnmarshal(t *testing.T) {
	m1 := uint64(9)
	m2 := uint64(8)
	m3 := uint64(1)
	testCases := []Sample{
		{
			TimeStamp: 0,
			SystemSample: SystemSample{
				HostName:      "",
				KernelVersion: "",
				PageSize:      0,
				LoadAvg:       procfs.LoadAvg{},
				Stat:          procfs.Stat{},
				Meminfo:       procfs.Meminfo{},
				NetDevStats:   make(map[string]procfs.NetDevLine),
				DiskStats:     make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: make(map[int]ProcSample),
		},
		{
			TimeStamp: 1,
			SystemSample: SystemSample{
				HostName:      "abc",
				KernelVersion: "5.14",
				PageSize:      0,
				LoadAvg:       procfs.LoadAvg{},
				Stat:          procfs.Stat{},
				Meminfo: procfs.Meminfo{
					MemTotal:     &m1,
					MemFree:      &m2,
					MemAvailable: &m3,
				},
				NetDevStats: nil,
				DiskStats:   make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: nil,
		},
		{
			TimeStamp: 9999,
			SystemSample: SystemSample{
				HostName:      "",
				KernelVersion: "",
				PageSize:      4096,
				LoadAvg:       procfs.LoadAvg{},
				Stat:          procfs.Stat{},
				Meminfo:       procfs.Meminfo{},
				NetDevStats:   make(map[string]procfs.NetDevLine),
				DiskStats:     make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: map[int]ProcSample{
				0: {ProcStat: procfs.ProcStat{
					PID:   0,
					Comm:  "",
					State: "",
					PPID:  0,
				},
					ProcIO: procfs.ProcIO{},
				},
				1: {ProcStat: procfs.ProcStat{
					PID:   1,
					Comm:  "test",
					State: "Sleeping",
					PPID:  0,
				},
					ProcIO: procfs.ProcIO{},
				},
			},
		},
	}
	for _, testCase := range testCases {
		var b []byte
		var err error
		if b, err = testCase.Marshal(); err != nil {
			t.Errorf("%+v Marshal: %s", testCase, err)
		}
		re := Sample{}
		if err = re.Unmarshal(b); err != nil {
			t.Errorf("%+v Unmarshal: %s", b, err)
		}
		opts := []cmp.Option{
			cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
		}
		if cmp.Equal(testCase, re, opts...) == false {
			t.Errorf("got %+v\nwant %+v\n", re, testCase)
		}
	}
}

func BenchmarkSampleMarshal(b *testing.B) {
	testCase := NewSample()
	if err := CollectSampleFromSys(&testCase); err != nil {
		b.Fatalf("collect sample: %s", err)
	}
	re, _ := testCase.Marshal()
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		testCase.Marshal()
	}
	b.Logf("Marshal to %d bytes", len(re))
}

func TestLocalStoreopenFile(t *testing.T) {
	dir := t.TempDir()
	local, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	err = local.openFile("20221127", false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "index_20221127"))

	err = local.openFile("20221127", false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "data_20221127"))
	err = local.openFile("20221127", false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile("20221127", false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile("20221127", true)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}
	if filepath.Base(local.Index.Name()) != "index_20221127" {
		t.Fatalf("index file got: %s, but want %s\n", filepath.Base(local.Index.Name()), "index_20221127")
	}

	if filepath.Base(local.Data.Name()) != "data_20221127" {
		t.Fatalf("data file got: %s, but want %s\n", filepath.Base(local.Data.Name()), "data_20221127")
	}
}

func TestNextSample(t *testing.T) {
	dir := t.TempDir()
	writeStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
	)
	if err != nil {
		t.Fatalf("now readStore: %s\n", err)
	}
	defer readStore.Close()

	s := NewSample()

	d := NewSample()
	noExist := NewSample()

	if err := readStore.NextSample(0, &d); err != ErrOutOfRange {
		t.Fatalf("read sample: %s\n", err)
	}

	if err := writeStore.CollectSample(&s); err != nil {
		t.Fatalf("collect sample: %s\n", err)
	}
	if err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}

	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.NextSample(0, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.NextSample(1, &noExist); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got no error: %s\n", err)
	}
	if err := readStore.NextSample(-1, &noExist); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got no error: %s\n", err)
	}

	if err := readStore.NextSample(0, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	s.TimeStamp += 1
	if err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write new sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.NextSample(1, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}
}

func TestJumpSampleByTimeStamp(t *testing.T) {
	dir := t.TempDir()
	writeStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
	)
	if err != nil {
		t.Fatalf("now readStore: %s\n", err)
	}
	defer readStore.Close()

	d := NewSample()
	c := NewSample()

	if err := readStore.JumpSampleByTimeStamp(123, &c); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	// write 1st sample
	s := NewSample()
	if err := writeStore.CollectSample(&s); err != nil {
		t.Fatalf("collect sample: %s\n", err)
	}
	if err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	// shoud ignore 1st sample, because no sample can compare with it
	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp-1, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp+1, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	// write 2nd sample
	s.TimeStamp += 1
	if err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp-1, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp+1, &c); err != nil {
		t.Logf("%d %+v\n", s.TimeStamp, readStore.idxs)
		t.Fatalf("read sample: %s\n", err)
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, c, opts...))
	}
}
