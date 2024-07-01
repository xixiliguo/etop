package model

import (
	"github.com/xixiliguo/etop/store"
)

var DefaultVmFields = []string{"PageIn", "PageOut",
	"SwapIn", "SwapOut",
	"PageScanKswapd", "PageScanDirect",
	"PageStealKswapd", "PageStealDirect", "OOMKill"}

type Vm struct {
	PageIn          uint64
	PageOut         uint64
	SwapIn          uint64
	SwapOut         uint64
	PageScanKswapd  uint64
	PageScanDirect  uint64
	PageStealKswapd uint64
	PageStealDirect uint64
	OOMKill         uint64
}

func (v *Vm) DefaultConfig(field string) Field {
	cfg := Field{}
	switch field {
	case "PageIn":
		cfg = Field{"PageIn", Raw, 0, "", 10, false}
	case "PageOut":
		cfg = Field{"PageOut", Raw, 0, "", 10, false}
	case "SwapIn":
		cfg = Field{"SwapIn", Raw, 0, "", 10, false}
	case "SwapOut":
		cfg = Field{"SwapOut", Raw, 0, "", 10, false}
	case "PageScanKswapd":
		cfg = Field{"PageScanKswapd", Raw, 0, "", 10, false}
	case "PageScanDirect":
		cfg = Field{"PageScanDirect", Raw, 0, "", 10, false}
	case "PageStealKswapd":
		cfg = Field{"PageStealKswapd", Raw, 0, "", 10, false}
	case "PageStealDirect":
		cfg = Field{"PageStealDirect", Raw, 0, "", 10, false}
	case "OOMKill":
		cfg = Field{"OOMKill", Raw, 0, "", 10, false}
	}
	return cfg
}

func (v *Vm) GetRenderValue(field string, opt FieldOpt) string {
	cfg := v.DefaultConfig(field)
	cfg.ApplyOpt(opt)
	s := ""
	switch field {
	case "PageIn":
		s = cfg.Render(v.PageIn)
	case "PageOut":
		s = cfg.Render(v.PageOut)
	case "SwapIn":
		s = cfg.Render(v.SwapIn)
	case "SwapOut":
		s = cfg.Render(v.SwapOut)
	case "PageScanKswapd":
		s = cfg.Render(v.PageScanKswapd)
	case "PageScanDirect":
		s = cfg.Render(v.PageScanDirect)
	case "PageStealKswapd":
		s = cfg.Render(v.PageStealKswapd)
	case "PageStealDirect":
		s = cfg.Render(v.PageStealDirect)
	case "OOMKill":
		s = cfg.Render(v.OOMKill)
	default:
		s = "no " + field + " for vm stat"
	}
	return s
}

func (v *Vm) Collect(prev, curr *store.Sample) {

	v.PageIn = (curr.PageIn - prev.PageIn) * 1024 / uint64(curr.PageSize)
	v.PageOut = curr.PageOut - prev.PageOut
	v.SwapIn = curr.SwapIn - prev.SwapIn
	v.SwapOut = curr.SwapOut - prev.SwapOut
	v.PageScanKswapd = curr.PageScanKswapd - prev.PageScanKswapd
	v.PageScanDirect = curr.PageScanDirect - prev.PageScanDirect
	v.PageStealKswapd = curr.PageStealKswapd - prev.PageStealKswapd
	v.PageStealDirect = curr.PageStealDirect - prev.PageStealDirect
	v.OOMKill = curr.OOMKill - prev.OOMKill
}
