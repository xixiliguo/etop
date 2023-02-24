package store

import (
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
	realData := NewSample()
	CollectSampleFromSys(&realData, nil)
	testCases = append(testCases, realData)
	for i, testCase := range testCases {
		var b []byte
		var err error
		if b, err = testCase.Marshal(); err != nil {
			t.Errorf("%+v Marshal: %s", testCase, err)
		}
		re := NewSample()
		if err = re.Unmarshal(b); err != nil {
			t.Fatalf("%+v Unmarshal: %s", b, err)
		}
		opts := []cmp.Option{
			cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
		}
		if cmp.Equal(testCase, re, opts...) == false {
			t.Errorf("case%d: %s", i, cmp.Diff(testCase, re, opts...))
		}
	}
}

func BenchmarkSampleMarshal(b *testing.B) {

	testCase := NewSample()
	CollectSampleFromSys(&testCase, nil)

	b.ReportAllocs()
	b.ResetTimer()
	size := 0
	for n := 0; n < b.N; n++ {
		re, _ := testCase.Marshal()
		size += len(re)
	}
	b.ReportMetric(float64(size)/float64(b.N), "size/op")
}
