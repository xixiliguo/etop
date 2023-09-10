package model

import (
	"fmt"

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

func (v *Vm) GetRenderValue(config RenderConfig, field string) string {
	s := fmt.Sprintf("no %s for vm stat", field)
	switch field {
	case "PageIn":
		s = config[field].Render(v.PageIn)
	case "PageOut":
		s = config[field].Render(v.PageOut)
	case "SwapIn":
		s = config[field].Render(v.SwapIn)
	case "SwapOut":
		s = config[field].Render(v.SwapOut)
	case "PageScanKswapd":
		s = config[field].Render(v.PageScanKswapd)
	case "PageScanDirect":
		s = config[field].Render(v.PageScanDirect)
	case "PageStealKswapd":
		s = config[field].Render(v.PageStealKswapd)
	case "PageStealDirect":
		s = config[field].Render(v.PageStealDirect)
	case "OOMKill":
		s = config[field].Render(v.OOMKill)
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
	return
}
