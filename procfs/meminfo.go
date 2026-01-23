package procfs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

// Meminfo represents memory statistics.
type Meminfo struct {
	// Total usable ram (i.e. physical ram minus a few reserved
	// bits and the kernel binary code)
	MemTotal uint64
	// The sum of LowFree+HighFree
	MemFree uint64
	// An estimate of how much memory is available for starting
	// new applications, without swapping. Calculated from
	// MemFree, SReclaimable, the size of the file LRU lists, and
	// the low watermarks in each zone.  The estimate takes into
	// account that the system needs some page cache to function
	// well, and that not all reclaimable slab will be
	// reclaimable, due to items being in use. The impact of those
	// factors will vary from system to system.
	MemAvailable uint64
	// Relatively temporary storage for raw disk blocks shouldn't
	// get tremendously large (20MB or so)
	Buffers uint64
	Cached  uint64
	// Memory that once was swapped out, is swapped back in but
	// still also is in the swapfile (if memory is needed it
	// doesn't need to be swapped out AGAIN because it is already
	// in the swapfile. This saves I/O)
	SwapCached uint64
	// Memory that has been used more recently and usually not
	// reclaimed unless absolutely necessary.
	Active uint64
	// Memory which has been less recently used.  It is more
	// eligible to be reclaimed for other purposes
	Inactive     uint64
	ActiveAnon   uint64
	InactiveAnon uint64
	ActiveFile   uint64
	InactiveFile uint64
	Unevictable  uint64
	Mlocked      uint64
	// total amount of swap space available
	SwapTotal uint64
	// Memory which has been evicted from RAM, and is temporarily
	// on the disk
	SwapFree uint64
	// Memory which is waiting to get written back to the disk
	Dirty uint64
	// Memory which is actively being written back to the disk
	Writeback uint64
	// Non-file backed pages mapped into userspace page tables
	AnonPages uint64
	// files which have been mapped, such as libraries
	Mapped uint64
	Shmem  uint64
	// in-kernel data structures cache
	Slab uint64
	// Part of Slab, that might be reclaimed, such as caches
	SReclaimable uint64
	// Part of Slab, that cannot be reclaimed on memory pressure
	SUnreclaim  uint64
	KernelStack uint64
	// amount of memory dedicated to the lowest level of page
	// tables.
	PageTables uint64
	// NFS pages sent to the server, but not yet committed to
	// stable storage
	NFSUnstable uint64
	// Memory used for block device "bounce buffers"
	Bounce uint64
	// Memory used by FUSE for temporary writeback buffers
	WritebackTmp uint64
	// Based on the overcommit ratio ('vm.overcommit_ratio'),
	// this is the total amount of  memory currently available to
	// be allocated on the system. This limit is only adhered to
	// if strict overcommit accounting is enabled (mode 2 in
	// 'vm.overcommit_memory').
	// The CommitLimit is calculated with the following formula:
	// CommitLimit = ([total RAM pages] - [total huge TLB pages]) *
	//                overcommit_ratio / 100 + [total swap pages]
	// For example, on a system with 1G of physical RAM and 7G
	// of swap with a `vm.overcommit_ratio` of 30 it would
	// yield a CommitLimit of 7.3G.
	// For more details, see the memory overcommit documentation
	// in vm/overcommit-accounting.
	CommitLimit uint64
	// The amount of memory presently allocated on the system.
	// The committed memory is a sum of all of the memory which
	// has been allocated by processes, even if it has not been
	// "used" by them as of yet. A process which malloc()'s 1G
	// of memory, but only touches 300M of it will show up as
	// using 1G. This 1G is memory which has been "committed" to
	// by the VM and can be used at any time by the allocating
	// application. With strict overcommit enabled on the system
	// (mode 2 in 'vm.overcommit_memory'),allocations which would
	// exceed the CommitLimit (detailed above) will not be permitted.
	// This is useful if one needs to guarantee that processes will
	// not fail due to lack of memory once that memory has been
	// successfully allocated.
	CommittedAS uint64
	// total size of vmalloc memory area
	VmallocTotal uint64
	// amount of vmalloc area which is used
	VmallocUsed uint64
	// largest contiguous block of vmalloc area which is free
	VmallocChunk      uint64
	Percpu            uint64
	HardwareCorrupted uint64
	AnonHugePages     uint64
	ShmemHugePages    uint64
	ShmemPmdMapped    uint64
	CmaTotal          uint64
	CmaFree           uint64
	HugePagesTotal    uint64
	HugePagesFree     uint64
	HugePagesRsvd     uint64
	HugePagesSurp     uint64
	Hugepagesize      uint64
	DirectMap4k       uint64
	DirectMap2M       uint64
	DirectMap1G       uint64
}

func (fs FS) Meminfo() (Meminfo, error) {
	m := Meminfo{}

	path := fs.path("meminfo")
	f, err := os.Open(path)
	if err != nil {
		return m, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		var val uint64
		var err error

		var fields [3]string
		nFields := stringutil.FieldsN(line, fields[:])

		val, err = strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return err
		}

		switch nFields {
		case 2:
		case 3:
			// Unit present in optional 3rd field, convert it to
			// bytes. The only unit supported within the Linux
			// kernel is `kB`.
			if fields[2] != "kB" {
				return fmt.Errorf("Unsupported unit in optional 3rd field in meminfo: %q", line)
			}
		default:
			return fmt.Errorf("unexpected line in stat: '%s'", line)
		}

		switch fields[0] {
		case "MemTotal:":
			m.MemTotal = val
		case "MemFree:":
			m.MemFree = val
		case "MemAvailable:":
			m.MemAvailable = val
		case "Buffers:":
			m.Buffers = val
		case "Cached:":
			m.Cached = val
		case "SwapCached:":
			m.SwapCached = val
		case "Active:":
			m.Active = val
		case "Inactive:":
			m.Inactive = val
		case "Active(anon):":
			m.ActiveAnon = val
		case "Inactive(anon):":
			m.InactiveAnon = val
		case "Active(file):":
			m.ActiveFile = val
		case "Inactive(file):":
			m.InactiveFile = val
		case "Unevictable:":
			m.Unevictable = val
		case "Mlocked:":
			m.Mlocked = val
		case "SwapTotal:":
			m.SwapTotal = val
		case "SwapFree:":
			m.SwapFree = val
		case "Dirty:":
			m.Dirty = val
		case "Writeback:":
			m.Writeback = val
		case "AnonPages:":
			m.AnonPages = val
		case "Mapped:":
			m.Mapped = val
		case "Shmem:":
			m.Shmem = val
		case "Slab:":
			m.Slab = val
		case "SReclaimable:":
			m.SReclaimable = val
		case "SUnreclaim:":
			m.SUnreclaim = val
		case "KernelStack:":
			m.KernelStack = val
		case "PageTables:":
			m.PageTables = val
		case "NFS_Unstable:":
			m.NFSUnstable = val
		case "Bounce:":
			m.Bounce = val
		case "WritebackTmp:":
			m.WritebackTmp = val
		case "CommitLimit:":
			m.CommitLimit = val
		case "Committed_AS:":
			m.CommittedAS = val
		case "VmallocTotal:":
			m.VmallocTotal = val
		case "VmallocUsed:":
			m.VmallocUsed = val
		case "VmallocChunk:":
			m.VmallocChunk = val
		case "Percpu:":
			m.Percpu = val
		case "HardwareCorrupted:":
			m.HardwareCorrupted = val
		case "AnonHugePages:":
			m.AnonHugePages = val
		case "ShmemHugePages:":
			m.ShmemHugePages = val
		case "ShmemPmdMapped:":
			m.ShmemPmdMapped = val
		case "CmaTotal:":
			m.CmaTotal = val
		case "CmaFree:":
			m.CmaFree = val
		case "HugePages_Total:":
			m.HugePagesTotal = val
		case "HugePages_Free:":
			m.HugePagesFree = val
		case "HugePages_Rsvd:":
			m.HugePagesRsvd = val
		case "HugePages_Surp:":
			m.HugePagesSurp = val
		case "Hugepagesize:":
			m.Hugepagesize = val
		case "DirectMap4k:":
			m.DirectMap4k = val
		case "DirectMap2M:":
			m.DirectMap2M = val
		case "DirectMap1G:":
			m.DirectMap1G = val
		}
		return nil
	})

	return m, err
}
