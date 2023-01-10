package store

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"golang.org/x/sys/unix"
)

// Sample represent all system info and process info.
type Sample struct {
	TimeStamp    int64              // unix time when sample was generated
	SystemSample                    // system information
	ProcSamples  map[int]ProcSample // process information
}

type SystemSample struct {
	HostName      string
	KernelVersion string
	PageSize      int
	procfs.LoadAvg
	procfs.Stat
	procfs.Meminfo
	VmStat
	NetDevStats map[string]procfs.NetDevLine
	DiskStats   map[string]blockdevice.Diskstats
	NetStat
	procfs.NetProtocolStats
	SoftNetStats []procfs.SoftnetStat
}

type NetStat struct {
	Ip       procfs.Ip
	Icmp     procfs.Icmp
	IcmpMsg  procfs.IcmpMsg
	Tcp      procfs.Tcp
	Udp      procfs.Udp
	UdpLite  procfs.UdpLite
	Ip6      procfs.Ip6
	Icmp6    procfs.Icmp6
	Udp6     procfs.Udp6
	UdpLite6 procfs.UdpLite6
	TcpExt   procfs.TcpExt
	IpExt    procfs.IpExt
}

type VmStat struct {
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

func newVmstat(file string) (VmStat, error) {
	f, err := os.Open(file)
	if err != nil {
		return VmStat{}, err
	}
	defer f.Close()

	vmStat := VmStat{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 {
			return VmStat{}, fmt.Errorf("malformed vmstat line: %q", s.Text())
		}

		v, err := strconv.ParseUint(fields[1], 0, 64)
		if err != nil {
			return VmStat{}, err
		}

		if fields[0] == "pgpgin" {
			vmStat.PageIn = v
			continue
		}
		if fields[0] == "pgpgout" {
			vmStat.PageOut = v
			continue
		}
		if fields[0] == "pswpin" {
			vmStat.SwapIn = v
			continue
		}
		if fields[0] == "pswpout" {
			vmStat.SwapOut = v
			continue
		}
		if fields[0] == "pswpout" {
			vmStat.SwapOut = v
			continue
		}
		if fields[0] == "pgscan_direct_throttle" {
			continue
		}
		if strings.HasPrefix(fields[0], "pgscan_kswapd") {
			vmStat.PageScanKswapd += v
			continue
		}
		if strings.HasPrefix(fields[0], "pgscan_direct") {
			vmStat.PageScanDirect += v
			continue
		}
		if strings.HasPrefix(fields[0], "pgsteal_kswapd") {
			vmStat.PageStealKswapd += v
			continue
		}
		if strings.HasPrefix(fields[0], "pgsteal_direct") {
			vmStat.PageStealDirect += v
			continue
		}
		if fields[0] == "oom_kill" {
			vmStat.OOMKill = v
			continue
		}
	}
	return vmStat, s.Err()
}

type ProcSample struct {
	procfs.ProcStat
	procfs.ProcIO
}

func NewSample() Sample {
	s := Sample{
		TimeStamp: 0,
		SystemSample: SystemSample{
			NetDevStats:      make(map[string]procfs.NetDevLine),
			DiskStats:        make(map[string]blockdevice.Diskstats),
			NetProtocolStats: make(map[string]procfs.NetProtocolStatLine),
		},
		ProcSamples: map[int]ProcSample{},
	}
	return s
}

func (s *Sample) Marshal() ([]byte, error) {

	b, err := cbor.Marshal(s)
	if err != nil {
		return nil, err
	}
	enc, _ := zstd.NewWriter(nil)
	return enc.EncodeAll(b, make([]byte, 0, len(b))), nil
}

func (s *Sample) Unmarshal(b []byte) error {

	dec, _ := zstd.NewReader(nil)
	uncompressed, err := dec.DecodeAll(b, make([]byte, 0, len(b)))
	if err != nil {
		return err
	}

	if err = cbor.Unmarshal(uncompressed, s); err != nil {
		return nil
	}
	return nil
}

func CollectSampleFromSys(s *Sample) error {
	//collect one sample
	var (
		fs     procfs.FS
		diskFS blockdevice.FS
		err    error
	)
	s.TimeStamp = time.Now().Unix()
	u := unix.Utsname{}
	unix.Uname(&u)
	s.HostName = string(u.Nodename[:])
	s.KernelVersion = string(u.Release[:])
	s.PageSize = os.Getpagesize()
	if fs, err = procfs.NewFS("/proc"); err != nil {
		return err
	}
	if diskFS, err = blockdevice.NewFS("/proc", "/sys"); err != nil {
		return err
	}

	if avg, err := fs.LoadAvg(); err != nil {
		return err
	} else {
		s.LoadAvg = *avg
	}

	if s.Stat, err = fs.Stat(); err != nil {
		return err
	}

	if s.Meminfo, err = fs.Meminfo(); err != nil {
		return err
	}

	if s.VmStat, err = newVmstat("/proc/vmstat"); err != nil {
		return err
	}

	if s.NetDevStats, err = fs.NetDev(); err != nil {
		return err
	}

	if s.NetProtocolStats, err = fs.NetProtocols(); err != nil {
		return err
	}

	p, _ := fs.NewProc(1)
	if snmp, err := p.Snmp(); err != nil {
		return err
	} else {
		s.Ip = snmp.Ip
		s.Icmp = snmp.Icmp
		s.IcmpMsg = snmp.IcmpMsg
		s.Tcp = snmp.Tcp
		s.Udp = snmp.Udp
		s.UdpLite = snmp.UdpLite

	}
	if snmp6, err := p.Snmp6(); err != nil {
		return err
	} else {
		s.Ip6 = snmp6.Ip6
		s.Icmp6 = snmp6.Icmp6
		s.Udp6 = snmp6.Udp6
		s.UdpLite6 = snmp6.UdpLite6
	}
	if netStat, err := p.Netstat(); err != nil {
		return err
	} else {
		s.TcpExt = netStat.TcpExt
		s.IpExt = netStat.IpExt
	}

	if s.SoftNetStats, err = fs.NetSoftnetStat(); err != nil {
		return err
	}

	if diskStats, err := diskFS.ProcDiskstats(); err != nil {
		return err
	} else {
		deviceNames := make(map[string]bool)
		if bds, err := diskFS.SysBlockDevices(); err != nil {
			return err
		} else {
			for _, db := range bds {
				deviceNames[db] = true
			}
		}
		for _, diskStat := range diskStats {
			if deviceNames[diskStat.DeviceName] {
				s.DiskStats[diskStat.DeviceName] = diskStat
			}
		}
	}
	procs := make(procfs.Procs, 0, 1024)
	if procs, err = fs.AllProcs(); err != nil {
		return err
	}
	for _, proc := range procs {
		p := ProcSample{}
		if p.ProcStat, err = proc.Stat(); err != nil {
			continue
		}
		if p.ProcIO, err = proc.IO(); err != nil {
			continue
		}
		s.ProcSamples[p.PID] = p
	}
	return nil
}
