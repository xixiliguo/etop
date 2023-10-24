package store

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"github.com/xixiliguo/etop/util"
)

func TestLocalStoreopenFile(t *testing.T) {
	dir := t.TempDir()
	local, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}
	//1697760000 mean 2013/10/20
	err = local.openFile(1697760000, false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "index_01697760000"))

	err = local.openFile(1697760000, false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "data_01697760000"))
	err = local.openFile(1697760000, false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile(1697760000, false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile(1697760000, true)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}
	if filepath.Base(local.Index.Name()) != "index_01697760000" {
		t.Fatalf("index file got: %s, but want %s\n", filepath.Base(local.Index.Name()), "index_20221127")
	}

	if filepath.Base(local.Data.Name()) != "data_01697760000" {
		t.Fatalf("data file got: %s, but want %s\n", filepath.Base(local.Data.Name()), "data_20221127")
	}
}

func TestNextSample(t *testing.T) {
	dir := t.TempDir()
	writeStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 8),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
	)
	if err != nil {
		t.Fatalf("now readStore: %s\n", err)
	}
	defer readStore.Close()

	s := NewSample()

	d := NewSample()
	noExist := NewSample()

	if err := readStore.NextSample(1, &d); err != ErrOutOfRange {
		t.Fatalf("read sample: %s\n", err)
	}

	if err := writeStore.CollectSample(&s); err != nil {
		t.Fatalf("collect sample: %s\n", err)
	}
	if _, err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}

	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.NextSample(1, &d); err != nil {
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
	if _, err := writeStore.WriteSample(&s); err != nil {
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
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 8),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
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
	if _, err := writeStore.WriteSample(&s); err != nil {
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
	if _, err := writeStore.WriteSample(&s); err != nil {
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

func TestGetSampleFromFileWithMultipleDataFormat(t *testing.T) {
	dir := t.TempDir()

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
				Stat: procfs.Stat{
					IRQ: []uint64{1, 2, 3},
				},
				Meminfo:     procfs.Meminfo{},
				NetDevStats: make(map[string]procfs.NetDevLine),
				DiskStats:   make(map[string]blockdevice.Diskstats),
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
			TimeStamp: 2,
			SystemSample: SystemSample{
				HostName:      "abc",
				KernelVersion: "xyz",
				PageSize:      4096,
				LoadAvg:       procfs.LoadAvg{Load1: 1, Load5: 5, Load15: 10},
				Stat:          procfs.Stat{},
				Meminfo:       procfs.Meminfo{},
				NetDevStats:   make(map[string]procfs.NetDevLine),
				DiskStats:     make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: map[int]ProcSample{
				0: {ProcStat: procfs.ProcStat{
					PID:   0,
					Comm:  "systemd",
					State: "",
					PPID:  0,
				},
					ProcIO: procfs.ProcIO{},
				},
				1: {ProcStat: procfs.ProcStat{
					PID:   1,
					Comm:  "test",
					State: "Sleeping",
					PPID:  9999,
				},
					ProcIO: procfs.ProcIO{},
				},
			},
		},
		{
			TimeStamp: 3,
			SystemSample: SystemSample{
				HostName:      "abc3",
				KernelVersion: "xyz3",
				PageSize:      4096,
				LoadAvg:       procfs.LoadAvg{Load1: 9, Load5: 99, Load15: 999},
				Stat:          procfs.Stat{BootTime: 99},
				Meminfo:       procfs.Meminfo{},
				NetDevStats:   make(map[string]procfs.NetDevLine),
				DiskStats:     make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: map[int]ProcSample{
				0: {ProcStat: procfs.ProcStat{
					PID:   0,
					Comm:  "etop",
					State: "X",
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
		{
			TimeStamp: 9997,
			SystemSample: SystemSample{
				HostName:      "",
				KernelVersion: "5.10",
				PageSize:      1024,
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
					State: "R",
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
					ProcIO: procfs.ProcIO{RChar: 1},
				},
			},
		},
		{
			TimeStamp: 9998,
			SystemSample: SystemSample{
				HostName:      "xyz",
				KernelVersion: "4.19",
				PageSize:      1024,
				LoadAvg:       procfs.LoadAvg{Load1: 1, Load5: 5, Load15: 10},
				Stat:          procfs.Stat{BootTime: 1},
				Meminfo:       procfs.Meminfo{},
				NetDevStats:   make(map[string]procfs.NetDevLine),
				DiskStats:     make(map[string]blockdevice.Diskstats),
			},
			ProcSamples: map[int]ProcSample{
				0: {ProcStat: procfs.ProcStat{
					PID:   0,
					Comm:  "systemd",
					State: "S",
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
				999: {ProcStat: procfs.ProcStat{
					PID:   999,
					Comm:  "test",
					State: "Sleeping",
					PPID:  1,
				},
					ProcIO: procfs.ProcIO{},
				},
			},
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

	if writeStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(NoCompress, 0),
	); err == nil {
		writeStore.WriteSample(&testCases[0])
		writeStore.Close()
	} else {
		t.Fatalf("new writeStore: %s\n", err)
	}

	if writeStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 2),
	); err == nil {
		writeStore.WriteSample(&testCases[1])
		writeStore.WriteSample(&testCases[2])
		writeStore.WriteSample(&testCases[3])
		writeStore.Close()
	} else {
		t.Fatalf("new writeStore: %s\n", err)
	}

	if writeStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompress, 0),
	); err == nil {
		writeStore.WriteSample(&testCases[4])
		writeStore.Close()
	} else {
		t.Fatalf("new writeStore: %s\n", err)
	}
	if writeStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 8),
	); err == nil {
		writeStore.WriteSample(&testCases[5])
		writeStore.WriteSample(&testCases[6])
		writeStore.Close()
	} else {
		t.Fatalf("new writeStore: %s\n", err)
	}

	res := []Sample{}
	if readStore, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
	); err == nil {
		for i := 0; i < 7; i++ {
			s := NewSample()
			t.Logf("%+v, %+v", readStore.idxs, readStore.curIdx)
			if err := readStore.NextSample(1, &s); err != nil {
				t.Fatalf("readStore: %s\n", err)
			}

			res = append(res, s)
		}

		readStore.Close()
	} else {
		t.Fatalf("now readStore: %s\n", err)
	}

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
	}
	if cmp.Equal(testCases, res, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(testCases, res, opts...))
	}

}

func BenchmarkWriteSample(b *testing.B) {

	testCase := NewSample()
	src, _ := os.ReadFile("testdata/sample.data")
	if err := testCase.Unmarshal(src); err != nil {
		b.Fatalf("testCase: %s", err)
	}
	devNull, _ := os.Open(os.DevNull)

	b.Run("nocompress", func(b *testing.B) {
		dir := b.TempDir()
		if writeStore, err := NewLocalStore(
			WithPathAndLogger(dir, util.CreateLogger(devNull, true)),
			WithWriteOnly(NoCompress, 0),
		); err == nil {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				writeStore.WriteSample(&testCase)
			}
			writeStore.Close()

			b.ReportMetric(float64(len(src)), "before/op")
			b.ReportMetric(float64(writeStore.DataOffset/int64(b.N)), "after/op")
		} else {
			b.Fatalf("new writeStore: %s\n", err)
		}

	})

	b.Run("compress", func(b *testing.B) {
		dir := b.TempDir()
		if writeStore, err := NewLocalStore(
			WithPathAndLogger(dir, util.CreateLogger(devNull, true)),
			WithWriteOnly(ZstdCompress, 0),
		); err == nil {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				writeStore.WriteSample(&testCase)
			}
			writeStore.Close()

			b.ReportMetric(float64(len(src)), "before/op")
			b.ReportMetric(float64(writeStore.DataOffset/int64(b.N)), "after/op")
		} else {
			b.Fatalf("new writeStore: %s\n", err)
		}

	})

	b.Run("compresswithdict", func(b *testing.B) {
		dir := b.TempDir()
		if writeStore, err := NewLocalStore(
			WithPathAndLogger(dir, util.CreateLogger(devNull, true)),
			WithWriteOnly(ZstdCompressWithDict, 1024),
		); err == nil {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				writeStore.WriteSample(&testCase)
			}
			writeStore.Close()

			b.ReportMetric(float64(len(src)), "before/op")
			b.ReportMetric(float64(writeStore.DataOffset/int64(b.N)), "after/op")
		} else {
			b.Fatalf("new writeStore: %s\n", err)
		}

	})

}

func getDirAndFilesName(path string) ([]string, []string) {
	entry, _ := os.ReadDir(path)
	dirs := []string{}
	files := []string{}
	for _, e := range entry {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		} else {
			files = append(files, e.Name())
		}
	}
	sort.Strings(dirs)
	sort.Strings(files)
	return dirs, files
}

func TestCleanOldFilesByDays(t *testing.T) {

	dir := t.TempDir()
	local, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 8),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	expect := []string{}

	now := time.Now()
	for i := 0; i < 5; i++ {
		shard := calcshard(now.AddDate(0, 0, -i).Unix())
		if i < 3 {
			expect = append(expect, fmt.Sprintf("data_%011d", shard))
			expect = append(expect, fmt.Sprintf("index_%011d", shard))
		}
		local.changeFile(shard, true)
	}
	sort.Strings(expect)

	local.CleanOldFiles(WriteOption{
		RetainDay:  2,
		RetainSize: 9999,
	})
	_, reFiles := getDirAndFilesName(local.Path)
	if cmp.Equal(expect, reFiles) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(expect, reFiles))
	}

}

func TestCleanOldFilesBySize(t *testing.T) {

	dir := t.TempDir()
	local, err := NewLocalStore(
		WithPathAndLogger(dir, slog.Default()),
		WithWriteOnly(ZstdCompressWithDict, 8),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	expect := []string{}

	now := time.Now()
	for i := 0; i < 5; i++ {
		shard := calcshard(now.AddDate(0, 0, -i).Unix())
		if i == 0 {
			expect = append(expect, fmt.Sprintf("data_%011d", shard))
			expect = append(expect, fmt.Sprintf("index_%011d", shard))
		}
		idx, err := os.Create(path.Join(local.Path, fmt.Sprintf("index_%011d", shard)))
		if err != nil {
			t.Fatal(err)
		}

		if err := idx.Truncate(1 << 20); err != nil {
			t.Fatal(err)
		}
		data, err := os.Create(path.Join(local.Path, fmt.Sprintf("data_%011d", shard)))
		if err != nil {
			t.Fatal(err)
		}

		if err := data.Truncate(5 << 20); err != nil {
			t.Fatal(err)
		}
	}
	sort.Strings(expect)

	local.CleanOldFiles(WriteOption{
		RetainDay:  9999,
		RetainSize: 6 << 20,
	})
	_, reFiles := getDirAndFilesName(local.Path)
	if cmp.Equal(expect, reFiles) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(expect, reFiles))
	}

}
