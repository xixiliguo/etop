package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs"
	"github.com/xixiliguo/etop/store"
)

func TestNetDevGetRenderValue(t *testing.T) {

	n := NetDev{
		Name:      "eth0",
		RxBytes:   1,
		RxPackets: 2,
		RxErrors:  3,
		RxDropped: 4,
		RxFIFO:    5,
	}

	tests := []struct {
		config RenderConfig
		field  string
		want   string
	}{
		{
			config: netDevDefaultRenderConfig,
			field:  "Name",
			want:   "eth0",
		},
		{
			config: netDevDefaultRenderConfig,
			field:  "RxBytes",
			want:   "1",
		},
		{
			config: netDevDefaultRenderConfig,
			field:  "abc",
			want:   "no abc for netdev stat",
		},
	}
	for _, tt := range tests {

		if got := n.GetRenderValue(tt.config, tt.field); got != tt.want {
			t.Errorf("NetDev.GetRenderValue() = %v, want %v", got, tt.want)
		}

	}
}

func TestNetDevCollect(t *testing.T) {

	prev := &store.Sample{
		TimeStamp: 0,
		SystemSample: store.SystemSample{
			NetDevStats: map[string]procfs.NetDevLine{
				"eth0": {
					Name:    "eth0",
					RxBytes: 1,
					RxFIFO:  0,
				},
				"ethx": {
					Name:    "ethx",
					RxBytes: 1,
					RxFIFO:  0,
				},
			},
		},
	}

	curr := &store.Sample{
		TimeStamp: 2,
		SystemSample: store.SystemSample{

			NetDevStats: map[string]procfs.NetDevLine{
				"eth0": {
					Name:    "eth0",
					RxBytes: 10,
					RxFIFO:  0,
				},
				"eth1": {
					Name:    "eth1",
					RxBytes: 20,
					RxFIFO:  1,
				},
			},
		},
	}

	testCase := NetDevMap{
		"eth0": {
			Name:         "eth0",
			RxBytes:      9,
			RxBytePerSec: 4.5,
			RxFIFO:       0,
		},
		"eth1": {
			Name:         "eth1",
			RxBytes:      20,
			RxBytePerSec: 10,
			RxFIFO:       1,
		},
	}

	re := NetDevMap{
		"xxxx": {
			Name:         "xxxx",
			RxBytes:      111,
			RxBytePerSec: 11,
			RxFIFO:       1,
		},
	}

	re.Collect(prev, curr)

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(),
	}

	if cmp.Equal(testCase, re, opts...) == false {
		t.Errorf("%s", cmp.Diff(testCase, re, opts...))
	}
}

func TestNetDevMapGetKeys(t *testing.T) {
	tests := []struct {
		netMap NetDevMap
		want   []string
	}{
		{
			netMap: NetDevMap{
				"lo":   {},
				"eth0": {},
			},
			want: []string{"eth0"},
		},
		{
			netMap: NetDevMap{
				"lo":    {},
				"eth0":  {},
				"eth1":  {},
				"eth11": {},
			},
			want: []string{"eth0", "eth1", "eth11"},
		},
		{
			netMap: NetDevMap{
				"lo":      {},
				"eth0":    {},
				"bond0":   {},
				"docker0": {},
			},
			want: []string{"bond0", "docker0", "eth0"},
		},
	}
	for _, tt := range tests {
		got := tt.netMap.GetKeys()
		if cmp.Equal(tt.want, got) == false {
			t.Errorf("%s", cmp.Diff(tt.want, got))
		}
	}
}
