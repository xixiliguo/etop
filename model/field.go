package model

import (
	"fmt"
	"strconv"
	"unsafe"
)

type Format int

const (
	Raw Format = iota
	HumanReadableSize
)

type Field struct {
	Name      string
	Format    Format
	Precision int
	Suffix    string
	Width     int
	FixWidth  bool
}

func appendReadableSize(dst []byte, fsize float64) []byte {
	unitMap := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := 0
	unitsLimit := len(unitMap) - 1
	for fsize >= 1024 && i < unitsLimit {
		fsize = fsize / 1024
		i++
	}
	dst = strconv.AppendFloat(dst, fsize, 'f', 1, 64)
	dst = append(dst, unitMap[i]...)
	return dst
}

func (f Field) Render(value any) string {
	buf := make([]byte, 0, 16)
	switch v := value.(type) {
	case uint64:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendUint(buf, uint64(v), 10)
		}
	case uint:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendUint(buf, uint64(v), 10)
		}
	case uint32:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendUint(buf, uint64(v), 10)
		}
	case int:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendInt(buf, int64(v), 10)
		}
	case int64:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendInt(buf, int64(v), 10)
		}
	case float64:
		if f.Format == HumanReadableSize {
			buf = appendReadableSize(buf, v)
		} else {
			buf = strconv.AppendFloat(buf, v, 'f', f.Precision, 64)
		}
	case string:
		if f.Suffix == "" && f.FixWidth == false {
			return v
		}
		buf = append(buf, v...)
	default:
		return fmt.Sprintf("%T is unknown type", v)
	}

	buf = append(buf, f.Suffix...)

	if f.FixWidth == true {
		width := f.Width
		if len(f.Name) > width {
			width = len(f.Name)
		}

		if padding := width - len(buf); padding > 0 {

			cache := "                "
			if padding <= 16 {
				buf = append(buf,
					unsafe.Slice(unsafe.StringData(cache), len(cache))[:padding]...)
			} else {
				for padding != 0 {
					buf = append(buf, ' ')
					padding--
				}
			}
		}
	}
	return string(buf)
}

type RenderConfig map[string]Field

func (renderConfig RenderConfig) Update(s string, f Field) {
	renderConfig[s] = f
}

func (renderConfig RenderConfig) SetFixWidth(fixWidth bool) {
	for k, v := range renderConfig {
		v.FixWidth = fixWidth
		renderConfig[k] = v
	}
}

func (renderConfig RenderConfig) SetRawData() {
	for k, v := range renderConfig {
		v.Format = Raw
		v.Suffix = ""
		renderConfig[k] = v
	}
}

var DefaultRenderConfig = make(map[string]RenderConfig)

var sysDefaultRenderConfig = make(RenderConfig)

var cpuDefaultRenderConfig = make(RenderConfig)

var memDefaultRenderConfig = make(RenderConfig)

var vmDefaultRenderConfig = make(RenderConfig)

var diskDefaultRenderConfig = make(RenderConfig)

var netDevDefaultRenderConfig = make(RenderConfig)

var netStatDefaultRenderConfig = make(RenderConfig)

var netProtocolDefaultRenderConfig = make(RenderConfig)

var softnetDefaultRenderConfig = make(RenderConfig)

var processDefaultRenderConfig = make(RenderConfig)

func genSysDefaultConfig() {

	sysDefaultRenderConfig["Load1"] = Field{"Load1", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["Load5"] = Field{"Load5", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["Load15"] = Field{"Load15", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["Processes"] = Field{"Process", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["Threads"] = Field{"Thread", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["ProcessesRunning"] = Field{"Running", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["ProcessesBlocked"] = Field{"Blocked", Raw, 0, "", 10, false}
	sysDefaultRenderConfig["ClonePerSec"] = Field{"Clone", Raw, 1, " /s", 10, false}
	sysDefaultRenderConfig["ContextSwitchPerSec"] = Field{"CtxSw", Raw, 1, " /s", 10, false}

}

func genCPUDefaultConfig() {
	cpuDefaultRenderConfig["Index"] = Field{"Index", Raw, 0, "", 10, false}
	cpuDefaultRenderConfig["User"] = Field{"User", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["Nice"] = Field{"Nice", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["System"] = Field{"System", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["Idle"] = Field{"Idle", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["Iowait"] = Field{"Iowait", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["IRQ"] = Field{"IRQ", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["SoftIRQ"] = Field{"SoftIRQ", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["Steal"] = Field{"Steal", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["Guest"] = Field{"Guest", Raw, 1, "%", 10, false}
	cpuDefaultRenderConfig["GuestNice"] = Field{"GuestNice", Raw, 1, "%", 10, false}

}

func genMEMDefaultConfig() {
	memDefaultRenderConfig["Total"] = Field{"Total", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["Free"] = Field{"Free", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["Avail"] = Field{"Avail", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["HSlab"] = Field{"Slab", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["Buffer"] = Field{"Buffer", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["Cache"] = Field{"Cache", HumanReadableSize, 0, "", 10, false}
	memDefaultRenderConfig["MemTotal"] = Field{"MemTotal", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["MemFree"] = Field{"MemFree", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["MemAvailable"] = Field{"MemAvailable", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Buffers"] = Field{"Buffers", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Cached"] = Field{"Cached", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["SwapCached"] = Field{"SwapCached", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Active"] = Field{"Active", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Inactive"] = Field{"Inactive", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["ActiveAnon"] = Field{"ActiveAnon", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["InactiveAnon"] = Field{"InactiveAnon", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["ActiveFile"] = Field{"ActiveFile", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["InactiveFile"] = Field{"InactiveFile", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Unevictable"] = Field{"Unevictable", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Mlocked"] = Field{"Mlocked", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["SwapTotal"] = Field{"SwapTotal", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["SwapFree"] = Field{"SwapFree", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Dirty"] = Field{"Dirty", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Writeback"] = Field{"Writeback", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["AnonPages"] = Field{"AnonPages", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Mapped"] = Field{"Mapped", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Shmem"] = Field{"Shmem", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Slab"] = Field{"Slab", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["SReclaimable"] = Field{"SReclaimable", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["SUnreclaim"] = Field{"SUnreclaim", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["KernelStack"] = Field{"KernelStack", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["PageTables"] = Field{"PageTables", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["NFSUnstable"] = Field{"NFSUnstable", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Bounce"] = Field{"Bounce", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["WritebackTmp"] = Field{"WritebackTmp", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["CommitLimit"] = Field{"CommitLimit", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["CommittedAS"] = Field{"CommittedAS", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["VmallocTotal"] = Field{"VmallocTotal", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["VmallocUsed"] = Field{"VmallocUsed", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["VmallocChunk"] = Field{"VmallocChunk", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["HardwareCorrupted"] = Field{"HardwareCorrupted", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["AnonHugePages"] = Field{"AnonHugePages", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["ShmemHugePages"] = Field{"ShmemHugePages", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["ShmemPmdMapped"] = Field{"ShmemPmdMapped", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["CmaTotal"] = Field{"CmaTotal", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["CmaFree"] = Field{"CmaFree", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["HugePagesTotal"] = Field{"HugePagesTotal", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["HugePagesFree"] = Field{"HugePagesFree", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["HugePagesRsvd"] = Field{"HugePagesRsvd", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["HugePagesSurp"] = Field{"HugePagesSurp", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["Hugepagesize"] = Field{"Hugepagesize", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["DirectMap4k"] = Field{"DirectMap4k", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["DirectMap2M"] = Field{"DirectMap2M", Raw, 0, " KB", 10, false}
	memDefaultRenderConfig["DirectMap1G"] = Field{"DirectMap1G", Raw, 0, " KB", 10, false}
}

func genVmDefaultConfig() {

	vmDefaultRenderConfig["PageIn"] = Field{"PageIn", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["PageOut"] = Field{"PageOut", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["SwapIn"] = Field{"SwapIn", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["SwapOut"] = Field{"SwapOut", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["PageScanKswapd"] = Field{"PageScanKswapd", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["PageScanDirect"] = Field{"PageScanDirect", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["PageStealKswapd"] = Field{"PageStealKswapd", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["PageStealDirect"] = Field{"PageStealDirect", Raw, 0, "", 10, false}
	vmDefaultRenderConfig["OOMKill"] = Field{"OOMKill", Raw, 0, "", 10, false}
}

func genDiskDefaultConfig() {

	diskDefaultRenderConfig["Disk"] = Field{"Disk", Raw, 0, "", 10, false}
	diskDefaultRenderConfig["Util"] = Field{"Util", Raw, 1, "%", 10, false}

	diskDefaultRenderConfig["Read"] = Field{"Read", Raw, 0, "", 10, false}
	diskDefaultRenderConfig["Read/s"] = Field{"Read/s", Raw, 0, "/s", 10, false}
	diskDefaultRenderConfig["ReadByte/s"] = Field{"ReadByte/s", HumanReadableSize, 1, "/s", 10, false}

	diskDefaultRenderConfig["Write"] = Field{"Write", Raw, 0, "", 10, false}
	diskDefaultRenderConfig["Write/s"] = Field{"Write/s", Raw, 0, "/s", 10, false}
	diskDefaultRenderConfig["WriteByte/s"] = Field{"WriteByte/s", HumanReadableSize, 1, "/s", 10, false}

	diskDefaultRenderConfig["Discard"] = Field{"Discard", Raw, 0, "", 10, false}
	diskDefaultRenderConfig["Discard/s"] = Field{"Discard/s", Raw, 0, "/s", 10, false}
	diskDefaultRenderConfig["DiscardByte/s"] = Field{"DiscardByte/s", HumanReadableSize, 1, "/s", 10, false}

	diskDefaultRenderConfig["AvgIOSize"] = Field{"AvgIOSize", HumanReadableSize, 1, "", 10, false}
	diskDefaultRenderConfig["AvgQueueLen"] = Field{"AvgQueueLen", Raw, 1, "", 10, false}
	diskDefaultRenderConfig["InFlight"] = Field{"InFlight", Raw, 1, "", 10, false}
	diskDefaultRenderConfig["AvgIOWait"] = Field{"AvgIOWait", Raw, 1, " ms", 10, false}
	diskDefaultRenderConfig["AvgIOTime"] = Field{"AvgIOTime", Raw, 1, " ms", 10, false}

}

func genNetDevDefaultConfig() {
	netDevDefaultRenderConfig["Name"] = Field{"Name", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxBytes"] = Field{"RxBytes", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxPackets"] = Field{"RxPackets", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxErrors"] = Field{"RxErrors", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxDropped"] = Field{"RxDropped", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxFIFO"] = Field{"RxFIFO", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxFrame"] = Field{"RxFrame", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxCompressed"] = Field{"RxCompressed", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxMulticast"] = Field{"RxMulticast", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxBytes"] = Field{"TxBytes", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxPackets"] = Field{"TxPackets", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxErrors"] = Field{"TxErrors", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxDropped"] = Field{"TxDropped", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxFIFO"] = Field{"TxFIFO", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxCollisions"] = Field{"TxCollisions", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxCarrier"] = Field{"TxCarrier", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["TxCompressed"] = Field{"TxCompressed", Raw, 0, "", 10, false}
	netDevDefaultRenderConfig["RxByte/s"] = Field{"RxByte/s", HumanReadableSize, 1, "/s", 10, false}
	netDevDefaultRenderConfig["RxPacket/s"] = Field{"Rp/s", Raw, 1, "/s", 10, false}
	netDevDefaultRenderConfig["TxByte/s"] = Field{"TxByte/s", HumanReadableSize, 1, "/s", 10, false}
	netDevDefaultRenderConfig["TxPacket/s"] = Field{"Tp/s", Raw, 1, "/s", 10, false}
}

func genNetStatDefaultConfig() {
	netStatDefaultRenderConfig["IpInReceives"] = Field{"IpInReceives", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpInHdrErrors"] = Field{"IpInHdrErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpInAddrErrors"] = Field{"IpInAddrErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpForwDatagrams"] = Field{"IpForwDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpInUnknownProtos"] = Field{"IpInUnknownProtos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpInDiscards"] = Field{"IpInDiscards", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpInDelivers"] = Field{"IpInDelivers", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpOutRequests"] = Field{"IpOutRequests", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpOutDiscards"] = Field{"IpOutDiscards", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpOutNoRoutes"] = Field{"IpOutNoRoutes", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpReasmTimeout"] = Field{"IpReasmTimeout", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpReasmReqds"] = Field{"IpReasmReqds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpReasmOKs"] = Field{"IpReasmOKs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpReasmFails"] = Field{"IpReasmFails", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpFragOKs"] = Field{"IpFragOKs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpFragFails"] = Field{"IpFragFails", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpFragCreates"] = Field{"IpFragCreates", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInMsgs"] = Field{"IcmpInMsgs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInErrors"] = Field{"IcmpInErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInCsumErrors"] = Field{"IcmpInCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInDestUnreachs"] = Field{"IcmpInDestUnreachs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInTimeExcds"] = Field{"IcmpInTimeExcds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInParmProbs"] = Field{"IcmpInParmProbs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInSrcQuenchs"] = Field{"IcmpInSrcQuenchs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInRedirects"] = Field{"IcmpInRedirects", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInEchos"] = Field{"IcmpInEchos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInEchoReps"] = Field{"IcmpInEchoReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInTimestamps"] = Field{"IcmpInTimestamps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInTimestampReps"] = Field{"IcmpInTimestampReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInAddrMasks"] = Field{"IcmpInAddrMasks", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInAddrMaskReps"] = Field{"IcmpInAddrMaskReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutMsgs"] = Field{"IcmpOutMsgs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutErrors"] = Field{"IcmpOutErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutDestUnreachs"] = Field{"IcmpOutDestUnreachs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutTimeExcds"] = Field{"IcmpOutTimeExcds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutParmProbs"] = Field{"IcmpOutParmProbs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutSrcQuenchs"] = Field{"IcmpOutSrcQuenchs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutRedirects"] = Field{"IcmpOutRedirects", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutEchos"] = Field{"IcmpOutEchos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutEchoReps"] = Field{"IcmpOutEchoReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutTimestamps"] = Field{"IcmpOutTimestamps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutTimestampReps"] = Field{"IcmpOutTimestampReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutAddrMasks"] = Field{"IcmpOutAddrMasks", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutAddrMaskReps"] = Field{"IcmpOutAddrMaskReps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpInType3"] = Field{"IcmpInType3", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IcmpOutType3"] = Field{"IcmpOutType3", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpActiveOpens"] = Field{"TcpActiveOpens", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpPassiveOpens"] = Field{"TcpPassiveOpens", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpAttemptFails"] = Field{"TcpAttemptFails", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpEstabResets"] = Field{"TcpEstabResets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpCurrEstab"] = Field{"TcpCurrEstab", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpInSegs"] = Field{"TcpInSegs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpOutSegs"] = Field{"TcpOutSegs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpRetransSegs"] = Field{"TcpRetransSegs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpInErrs"] = Field{"TcpInErrs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpOutRsts"] = Field{"TcpOutRsts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpInCsumErrors"] = Field{"TcpInCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpInDatagrams"] = Field{"UdpInDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpNoPorts"] = Field{"UdpNoPorts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpInErrors"] = Field{"UdpInErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpOutDatagrams"] = Field{"UdpOutDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpRcvbufErrors"] = Field{"UdpRcvbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpSndbufErrors"] = Field{"UdpSndbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpInCsumErrors"] = Field{"UdpInCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpIgnoredMulti"] = Field{"UdpIgnoredMulti", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteInDatagrams"] = Field{"UdpLiteInDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteNoPorts"] = Field{"UdpLiteNoPorts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteInErrors"] = Field{"UdpLiteInErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteOutDatagrams"] = Field{"UdpLiteOutDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteRcvbufErrors"] = Field{"UdpLiteRcvbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteSndbufErrors"] = Field{"UdpLiteSndbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteInCsumErrors"] = Field{"UdpLiteInCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLiteIgnoredMulti"] = Field{"UdpLiteIgnoredMulti", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InReceives"] = Field{"Ip6InReceives", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InHdrErrors"] = Field{"Ip6InHdrErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InTooBigErrors"] = Field{"Ip6InTooBigErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InNoRoutes"] = Field{"Ip6InNoRoutes", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InAddrErrors"] = Field{"Ip6InAddrErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InUnknownProtos"] = Field{"Ip6InUnknownProtos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InTruncatedPkts"] = Field{"Ip6InTruncatedPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InDiscards"] = Field{"Ip6InDiscards", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InDelivers"] = Field{"Ip6InDelivers", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutForwDatagrams"] = Field{"Ip6OutForwDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutRequests"] = Field{"Ip6OutRequests", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutDiscards"] = Field{"Ip6OutDiscards", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutNoRoutes"] = Field{"Ip6OutNoRoutes", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6ReasmTimeout"] = Field{"Ip6ReasmTimeout", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6ReasmReqds"] = Field{"Ip6ReasmReqds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6ReasmOKs"] = Field{"Ip6ReasmOKs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6ReasmFails"] = Field{"Ip6ReasmFails", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6FragOKs"] = Field{"Ip6FragOKs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6FragFails"] = Field{"Ip6FragFails", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6FragCreates"] = Field{"Ip6FragCreates", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InMcastPkts"] = Field{"Ip6InMcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutMcastPkts"] = Field{"Ip6OutMcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InOctets"] = Field{"Ip6InOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutOctets"] = Field{"Ip6OutOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InMcastOctets"] = Field{"Ip6InMcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutMcastOctets"] = Field{"Ip6OutMcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InBcastOctets"] = Field{"Ip6InBcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6OutBcastOctets"] = Field{"Ip6OutBcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InNoECTPkts"] = Field{"Ip6InNoECTPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InECT1Pkts"] = Field{"Ip6InECT1Pkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InECT0Pkts"] = Field{"Ip6InECT0Pkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Ip6InCEPkts"] = Field{"Ip6InCEPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InMsgs"] = Field{"Icmp6InMsgs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InErrors"] = Field{"Icmp6InErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutMsgs"] = Field{"Icmp6OutMsgs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutErrors"] = Field{"Icmp6OutErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InCsumErrors"] = Field{"Icmp6InCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InDestUnreachs"] = Field{"Icmp6InDestUnreachs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InPktTooBigs"] = Field{"Icmp6InPktTooBigs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InTimeExcds"] = Field{"Icmp6InTimeExcds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InParmProblems"] = Field{"Icmp6InParmProblems", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InEchos"] = Field{"Icmp6InEchos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InEchoReplies"] = Field{"Icmp6InEchoReplies", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InGroupMembQueries"] = Field{"Icmp6InGroupMembQueries", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InGroupMembResponses"] = Field{"Icmp6InGroupMembResponses", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InGroupMembReductions"] = Field{"Icmp6InGroupMembReductions", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InRouterSolicits"] = Field{"Icmp6InRouterSolicits", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InRouterAdvertisements"] = Field{"Icmp6InRouterAdvertisements", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InNeighborSolicits"] = Field{"Icmp6InNeighborSolicits", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InNeighborAdvertisements"] = Field{"Icmp6InNeighborAdvertisements", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InRedirects"] = Field{"Icmp6InRedirects", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InMLDv2Reports"] = Field{"Icmp6InMLDv2Reports", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutDestUnreachs"] = Field{"Icmp6OutDestUnreachs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutPktTooBigs"] = Field{"Icmp6OutPktTooBigs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutTimeExcds"] = Field{"Icmp6OutTimeExcds", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutParmProblems"] = Field{"Icmp6OutParmProblems", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutEchos"] = Field{"Icmp6OutEchos", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutEchoReplies"] = Field{"Icmp6OutEchoReplies", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutGroupMembQueries"] = Field{"Icmp6OutGroupMembQueries", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutGroupMembResponses"] = Field{"Icmp6OutGroupMembResponses", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutGroupMembReductions"] = Field{"Icmp6OutGroupMembReductions", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutRouterSolicits"] = Field{"Icmp6OutRouterSolicits", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutRouterAdvertisements"] = Field{"Icmp6OutRouterAdvertisements", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutNeighborSolicits"] = Field{"Icmp6OutNeighborSolicits", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutNeighborAdvertisements"] = Field{"Icmp6OutNeighborAdvertisements", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutRedirects"] = Field{"Icmp6OutRedirects", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutMLDv2Reports"] = Field{"Icmp6OutMLDv2Reports", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InType1"] = Field{"Icmp6InType1", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InType134"] = Field{"Icmp6InType134", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InType135"] = Field{"Icmp6InType135", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InType136"] = Field{"Icmp6InType136", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6InType143"] = Field{"Icmp6InType143", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutType133"] = Field{"Icmp6OutType133", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutType135"] = Field{"Icmp6OutType135", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutType136"] = Field{"Icmp6OutType136", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Icmp6OutType143"] = Field{"Icmp6OutType143", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6InDatagrams"] = Field{"Udp6InDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6NoPorts"] = Field{"Udp6NoPorts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6InErrors"] = Field{"Udp6InErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6OutDatagrams"] = Field{"Udp6OutDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6RcvbufErrors"] = Field{"Udp6RcvbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6SndbufErrors"] = Field{"Udp6SndbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6InCsumErrors"] = Field{"Udp6InCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["Udp6IgnoredMulti"] = Field{"Udp6IgnoredMulti", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6InDatagrams"] = Field{"UdpLite6InDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6NoPorts"] = Field{"UdpLite6NoPorts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6InErrors"] = Field{"UdpLite6InErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6OutDatagrams"] = Field{"UdpLite6OutDatagrams", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6RcvbufErrors"] = Field{"UdpLite6RcvbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6SndbufErrors"] = Field{"UdpLite6SndbufErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["UdpLite6InCsumErrors"] = Field{"UdpLite6InCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtSyncookiesSent"] = Field{"TcpExtSyncookiesSent", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtSyncookiesRecv"] = Field{"TcpExtSyncookiesRecv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtSyncookiesFailed"] = Field{"TcpExtSyncookiesFailed", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtEmbryonicRsts"] = Field{"TcpExtEmbryonicRsts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtPruneCalled"] = Field{"TcpExtPruneCalled", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtRcvPruned"] = Field{"TcpExtRcvPruned", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtOfoPruned"] = Field{"TcpExtOfoPruned", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtOutOfWindowIcmps"] = Field{"TcpExtOutOfWindowIcmps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtLockDroppedIcmps"] = Field{"TcpExtLockDroppedIcmps", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtArpFilter"] = Field{"TcpExtArpFilter", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTW"] = Field{"TcpExtTW", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTWRecycled"] = Field{"TcpExtTWRecycled", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTWKilled"] = Field{"TcpExtTWKilled", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtPAWSActive"] = Field{"TcpExtPAWSActive", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtPAWSEstab"] = Field{"TcpExtPAWSEstab", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtDelayedACKs"] = Field{"TcpExtDelayedACKs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtDelayedACKLocked"] = Field{"TcpExtDelayedACKLocked", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtDelayedACKLost"] = Field{"TcpExtDelayedACKLost", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtListenOverflows"] = Field{"TcpExtListenOverflows", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtListenDrops"] = Field{"TcpExtListenDrops", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHPHits"] = Field{"TcpExtTCPHPHits", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPPureAcks"] = Field{"TcpExtTCPPureAcks", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHPAcks"] = Field{"TcpExtTCPHPAcks", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRenoRecovery"] = Field{"TcpExtTCPRenoRecovery", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackRecovery"] = Field{"TcpExtTCPSackRecovery", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSACKReneging"] = Field{"TcpExtTCPSACKReneging", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSACKReorder"] = Field{"TcpExtTCPSACKReorder", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRenoReorder"] = Field{"TcpExtTCPRenoReorder", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPTSReorder"] = Field{"TcpExtTCPTSReorder", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFullUndo"] = Field{"TcpExtTCPFullUndo", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPPartialUndo"] = Field{"TcpExtTCPPartialUndo", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKUndo"] = Field{"TcpExtTCPDSACKUndo", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPLossUndo"] = Field{"TcpExtTCPLossUndo", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPLostRetransmit"] = Field{"TcpExtTCPLostRetransmit", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRenoFailures"] = Field{"TcpExtTCPRenoFailures", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackFailures"] = Field{"TcpExtTCPSackFailures", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPLossFailures"] = Field{"TcpExtTCPLossFailures", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastRetrans"] = Field{"TcpExtTCPFastRetrans", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSlowStartRetrans"] = Field{"TcpExtTCPSlowStartRetrans", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPTimeouts"] = Field{"TcpExtTCPTimeouts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPLossProbes"] = Field{"TcpExtTCPLossProbes", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPLossProbeRecovery"] = Field{"TcpExtTCPLossProbeRecovery", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRenoRecoveryFail"] = Field{"TcpExtTCPRenoRecoveryFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackRecoveryFail"] = Field{"TcpExtTCPSackRecoveryFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRcvCollapsed"] = Field{"TcpExtTCPRcvCollapsed", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKOldSent"] = Field{"TcpExtTCPDSACKOldSent", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKOfoSent"] = Field{"TcpExtTCPDSACKOfoSent", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKRecv"] = Field{"TcpExtTCPDSACKRecv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKOfoRecv"] = Field{"TcpExtTCPDSACKOfoRecv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortOnData"] = Field{"TcpExtTCPAbortOnData", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortOnClose"] = Field{"TcpExtTCPAbortOnClose", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortOnMemory"] = Field{"TcpExtTCPAbortOnMemory", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortOnTimeout"] = Field{"TcpExtTCPAbortOnTimeout", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortOnLinger"] = Field{"TcpExtTCPAbortOnLinger", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAbortFailed"] = Field{"TcpExtTCPAbortFailed", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMemoryPressures"] = Field{"TcpExtTCPMemoryPressures", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMemoryPressuresChrono"] = Field{"TcpExtTCPMemoryPressuresChrono", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSACKDiscard"] = Field{"TcpExtTCPSACKDiscard", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKIgnoredOld"] = Field{"TcpExtTCPDSACKIgnoredOld", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDSACKIgnoredNoUndo"] = Field{"TcpExtTCPDSACKIgnoredNoUndo", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSpuriousRTOs"] = Field{"TcpExtTCPSpuriousRTOs", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMD5NotFound"] = Field{"TcpExtTCPMD5NotFound", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMD5Unexpected"] = Field{"TcpExtTCPMD5Unexpected", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMD5Failure"] = Field{"TcpExtTCPMD5Failure", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackShifted"] = Field{"TcpExtTCPSackShifted", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackMerged"] = Field{"TcpExtTCPSackMerged", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSackShiftFallback"] = Field{"TcpExtTCPSackShiftFallback", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPBacklogDrop"] = Field{"TcpExtTCPBacklogDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtPFMemallocDrop"] = Field{"TcpExtPFMemallocDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMinTTLDrop"] = Field{"TcpExtTCPMinTTLDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPDeferAcceptDrop"] = Field{"TcpExtTCPDeferAcceptDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtIPReversePathFilter"] = Field{"TcpExtIPReversePathFilter", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPTimeWaitOverflow"] = Field{"TcpExtTCPTimeWaitOverflow", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPReqQFullDoCookies"] = Field{"TcpExtTCPReqQFullDoCookies", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPReqQFullDrop"] = Field{"TcpExtTCPReqQFullDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRetransFail"] = Field{"TcpExtTCPRetransFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRcvCoalesce"] = Field{"TcpExtTCPRcvCoalesce", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPRcvQDrop"] = Field{"TcpExtTCPRcvQDrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPOFOQueue"] = Field{"TcpExtTCPOFOQueue", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPOFODrop"] = Field{"TcpExtTCPOFODrop", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPOFOMerge"] = Field{"TcpExtTCPOFOMerge", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPChallengeACK"] = Field{"TcpExtTCPChallengeACK", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSYNChallenge"] = Field{"TcpExtTCPSYNChallenge", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenActive"] = Field{"TcpExtTCPFastOpenActive", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenActiveFail"] = Field{"TcpExtTCPFastOpenActiveFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenPassive"] = Field{"TcpExtTCPFastOpenPassive", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenPassiveFail"] = Field{"TcpExtTCPFastOpenPassiveFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenListenOverflow"] = Field{"TcpExtTCPFastOpenListenOverflow", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenCookieReqd"] = Field{"TcpExtTCPFastOpenCookieReqd", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFastOpenBlackhole"] = Field{"TcpExtTCPFastOpenBlackhole", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSpuriousRtxHostQueues"] = Field{"TcpExtTCPSpuriousRtxHostQueues", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtBusyPollRxPackets"] = Field{"TcpExtBusyPollRxPackets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPAutoCorking"] = Field{"TcpExtTCPAutoCorking", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPFromZeroWindowAdv"] = Field{"TcpExtTCPFromZeroWindowAdv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPToZeroWindowAdv"] = Field{"TcpExtTCPToZeroWindowAdv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPWantZeroWindowAdv"] = Field{"TcpExtTCPWantZeroWindowAdv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPSynRetrans"] = Field{"TcpExtTCPSynRetrans", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPOrigDataSent"] = Field{"TcpExtTCPOrigDataSent", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHystartTrainDetect"] = Field{"TcpExtTCPHystartTrainDetect", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHystartTrainCwnd"] = Field{"TcpExtTCPHystartTrainCwnd", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHystartDelayDetect"] = Field{"TcpExtTCPHystartDelayDetect", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPHystartDelayCwnd"] = Field{"TcpExtTCPHystartDelayCwnd", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedSynRecv"] = Field{"TcpExtTCPACKSkippedSynRecv", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedPAWS"] = Field{"TcpExtTCPACKSkippedPAWS", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedSeq"] = Field{"TcpExtTCPACKSkippedSeq", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedFinWait2"] = Field{"TcpExtTCPACKSkippedFinWait2", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedTimeWait"] = Field{"TcpExtTCPACKSkippedTimeWait", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPACKSkippedChallenge"] = Field{"TcpExtTCPACKSkippedChallenge", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPWinProbe"] = Field{"TcpExtTCPWinProbe", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPKeepAlive"] = Field{"TcpExtTCPKeepAlive", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMTUPFail"] = Field{"TcpExtTCPMTUPFail", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPMTUPSuccess"] = Field{"TcpExtTCPMTUPSuccess", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["TcpExtTCPWqueueTooBig"] = Field{"TcpExtTCPWqueueTooBig", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInNoRoutes"] = Field{"IpExtInNoRoutes", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInTruncatedPkts"] = Field{"IpExtInTruncatedPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInMcastPkts"] = Field{"IpExtInMcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtOutMcastPkts"] = Field{"IpExtOutMcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInBcastPkts"] = Field{"IpExtInBcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtOutBcastPkts"] = Field{"IpExtOutBcastPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInOctets"] = Field{"IpExtInOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtOutOctets"] = Field{"IpExtOutOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInMcastOctets"] = Field{"IpExtInMcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtOutMcastOctets"] = Field{"IpExtOutMcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInBcastOctets"] = Field{"IpExtInBcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtOutBcastOctets"] = Field{"IpExtOutBcastOctets", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInCsumErrors"] = Field{"IpExtInCsumErrors", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInNoECTPkts"] = Field{"IpExtInNoECTPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInECT1Pkts"] = Field{"IpExtInECT1Pkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInECT0Pkts"] = Field{"IpExtInECT0Pkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtInCEPkts"] = Field{"IpExtInCEPkts", Raw, 0, "", 10, false}
	netStatDefaultRenderConfig["IpExtReasmOverlaps"] = Field{"IpExtReasmOverlaps", Raw, 0, "", 10, false}
}

func genNetProtocolDefaultConfig() {

	netProtocolDefaultRenderConfig["Name"] = Field{"Name", Raw, 0, "", 10, false}
	netProtocolDefaultRenderConfig["Sockets"] = Field{"Sockets", Raw, 0, "", 10, false}
	netProtocolDefaultRenderConfig["Memory"] = Field{"Memory", HumanReadableSize, 0, "", 10, false}

}

func genSoftnetDefaultConfig() {
	softnetDefaultRenderConfig["CPU"] = Field{"CPU", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["Processed"] = Field{"Processed", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["Dropped"] = Field{"Dropped", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["TimeSqueezed"] = Field{"TimeSqueezed", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["CPUCollision"] = Field{"CPUCollision", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["ReceivedRps"] = Field{"ReceivedRps", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["FlowLimitCount"] = Field{"FlowLimitCount", Raw, 0, "", 10, false}
	softnetDefaultRenderConfig["SoftnetBacklogLen"] = Field{"SoftnetBacklogLen", Raw, 0, "", 10, false}
}

func genProcessDefaultConfig() {

	processDefaultRenderConfig["Pid"] = Field{"Pid", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Comm"] = Field{"Comm", Raw, 0, "", 16, false}
	processDefaultRenderConfig["State"] = Field{"State", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Ppid"] = Field{"Ppid", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Thr"] = Field{"Thr", Raw, 0, "", 10, false}
	processDefaultRenderConfig["StartTime"] = Field{"StartTime", Raw, 0, "", 10, false}
	processDefaultRenderConfig["OnCPU"] = Field{"OnCPU", Raw, 0, "", 10, false}
	processDefaultRenderConfig["CmdLine"] = Field{"CmdLine", Raw, 0, "", 10, false}

	processDefaultRenderConfig["UserCPU"] = Field{"UserCPU", Raw, 1, "%", 10, false}
	processDefaultRenderConfig["SysCPU"] = Field{"SysCPU", Raw, 1, "%", 10, false}
	processDefaultRenderConfig["Pri"] = Field{"Pri", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Nice"] = Field{"Nice", Raw, 0, "", 10, false}
	processDefaultRenderConfig["CPU"] = Field{"CPU", Raw, 1, "%", 10, false}

	processDefaultRenderConfig["Minflt"] = Field{"Minflt", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Majflt"] = Field{"Majflt", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Vsize"] = Field{"Vsize", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["RSS"] = Field{"RSS", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["Mem"] = Field{"Mem", Raw, 1, "%", 10, false}

	processDefaultRenderConfig["Rchar"] = Field{"Rchar", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["Wchar"] = Field{"Wchar", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["Rchar/s"] = Field{"Rchar/s", HumanReadableSize, 1, "/s", 10, false}
	processDefaultRenderConfig["Wchar/s"] = Field{"Wchar/s", HumanReadableSize, 1, "/s", 10, false}
	processDefaultRenderConfig["Syscr"] = Field{"Syscr", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Syscw"] = Field{"Syscw", Raw, 0, "", 10, false}
	processDefaultRenderConfig["Syscr/s"] = Field{"Syscr/s", Raw, 1, "/s", 10, false}
	processDefaultRenderConfig["Syscw/s"] = Field{"Syscw/s", Raw, 1, "/s", 10, false}
	processDefaultRenderConfig["Read"] = Field{"Read", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["Write"] = Field{"Write", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["Wcancel"] = Field{"Wcancel", HumanReadableSize, 0, "", 10, false}
	processDefaultRenderConfig["R/s"] = Field{"R/s", HumanReadableSize, 1, "/s", 10, false}
	processDefaultRenderConfig["W/s"] = Field{"W/s", HumanReadableSize, 1, "/s", 10, false}
	processDefaultRenderConfig["CW/s"] = Field{"CW/s", HumanReadableSize, 1, "/s", 10, false}
	processDefaultRenderConfig["Disk"] = Field{"Disk", Raw, 1, "%", 10, false}

}

func init() {
	genSysDefaultConfig()
	genCPUDefaultConfig()
	genMEMDefaultConfig()
	genVmDefaultConfig()
	genDiskDefaultConfig()
	genNetDevDefaultConfig()
	genNetStatDefaultConfig()
	genNetProtocolDefaultConfig()
	genProcessDefaultConfig()
	genSoftnetDefaultConfig()
	DefaultRenderConfig["system"] = sysDefaultRenderConfig
	DefaultRenderConfig["cpu"] = cpuDefaultRenderConfig
	DefaultRenderConfig["memory"] = memDefaultRenderConfig
	DefaultRenderConfig["vm"] = vmDefaultRenderConfig
	DefaultRenderConfig["disk"] = diskDefaultRenderConfig
	DefaultRenderConfig["netdev"] = netDevDefaultRenderConfig
	DefaultRenderConfig["network"] = netStatDefaultRenderConfig
	DefaultRenderConfig["networkprotocol"] = netProtocolDefaultRenderConfig
	DefaultRenderConfig["softnet"] = softnetDefaultRenderConfig
	DefaultRenderConfig["process"] = processDefaultRenderConfig
}
