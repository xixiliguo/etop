package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs/blockdevice"
	"github.com/xixiliguo/etop/store"
)

func TestDiskGetRenderValue(t *testing.T) {

	n := Disk{
		DeviceName:    "vda",
		ReadIOs:       1,
		WriteIOs:      2,
		IOsTotalTicks: 100,
	}

	tests := []struct {
		field string
		want  string
	}{
		{
			field: "Disk",
			want:  "vda",
		},
		{
			field: "Read",
			want:  "1",
		},
		{
			field: "Write",
			want:  "2",
		},
	}
	for _, tt := range tests {

		if got := n.GetRenderValue(tt.field, FieldOpt{}); got != tt.want {
			t.Errorf("Disk.GetRenderValue() = %v, want %v", got, tt.want)
		}

	}
}

func TestDiskCollect(t *testing.T) {

	prev := &store.Sample{
		TimeStamp: 0,
		SystemSample: store.SystemSample{
			DiskStats: map[string]blockdevice.Diskstats{
				"vda": {
					Info: blockdevice.Info{
						MajorNumber: 0,
						MinorNumber: 0,
						DeviceName:  "vda",
					},
					IOStats: blockdevice.IOStats{
						ReadIOs:  1,
						WriteIOs: 1,
					},
					IoStatsCount: 0,
				},
			},
		},
	}

	curr := &store.Sample{
		TimeStamp: 2,
		SystemSample: store.SystemSample{
			DiskStats: map[string]blockdevice.Diskstats{
				"vda": {
					Info: blockdevice.Info{
						MajorNumber: 0,
						MinorNumber: 0,
						DeviceName:  "vda",
					},
					IOStats: blockdevice.IOStats{
						ReadIOs:    10,
						WriteIOs:   10,
						DiscardIOs: 1,
					},
					IoStatsCount: 0,
				},
			},
		},
	}

	testCase := DiskMap{
		"vda": {
			DeviceName:    "vda",
			ReadIOs:       9,
			WriteIOs:      9,
			DiscardIOs:    1,
			ReadPerSec:    4.5,
			WritePerSec:   4.5,
			DiscardPerSec: 0.5,
		},
	}

	re := DiskMap{}

	re.Collect(prev, curr)

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(),
	}

	if cmp.Equal(testCase, re, opts...) == false {
		t.Errorf("%s", cmp.Diff(testCase, re, opts...))
	}
}

func TestDiskMapGetKeys(t *testing.T) {
	tests := []struct {
		diskMap DiskMap
		want    []string
	}{
		{
			diskMap: DiskMap{
				"vda": {},
				"vdb": {},
			},
			want: []string{"vda", "vdb"},
		},
		{
			diskMap: DiskMap{
				"sdb": {},
				"sda": {},
			},
			want: []string{"sda", "sdb"},
		},
		{
			diskMap: DiskMap{
				"vda": {},
				"sda": {},
				"sdb": {},
			},
			want: []string{"sda", "sdb", "vda"},
		},
	}
	for _, tt := range tests {
		got := tt.diskMap.GetKeys()
		if cmp.Equal(tt.want, got) == false {
			t.Errorf("%s", cmp.Diff(tt.want, got))
		}
	}
}
