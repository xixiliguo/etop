package procfs

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

type DiskStat map[string]DiskStatLine

type DiskStatLine struct {
	MajorNumber uint64
	MinorNumber uint64
	DeviceName  string
	// ReadIOs is the number of reads completed successfully.
	ReadIOs uint64
	// ReadMerges is the number of reads merged.  Reads and writes
	// which are adjacent to each other may be merged for efficiency.
	ReadMerges uint64
	// ReadSectors is the total number of sectors read successfully.
	ReadSectors uint64
	// ReadTicks is the total number of milliseconds spent by all reads.
	ReadTicks uint64
	// WriteIOs is the total number of writes completed successfully.
	WriteIOs uint64
	// WriteMerges is the number of reads merged.
	WriteMerges uint64
	// WriteSectors is the total number of sectors written successfully.
	WriteSectors uint64
	// WriteTicks is the total number of milliseconds spent by all writes.
	WriteTicks uint64
	// IOsInProgress is number of I/Os currently in progress.
	IOsInProgress uint64
	// IOsTotalTicks is the number of milliseconds spent doing I/Os.
	// This field increases so long as IosInProgress is nonzero.
	IOsTotalTicks uint64
	// WeightedIOTicks is the weighted number of milliseconds spent doing I/Os.
	// This can also be used to estimate average queue wait time for requests.
	WeightedIOTicks uint64
	// DiscardIOs is the total number of discards completed successfully.
	DiscardIOs uint64
	// DiscardMerges is the number of discards merged.
	DiscardMerges uint64
	// DiscardSectors is the total number of sectors discarded successfully.
	DiscardSectors uint64
	// DiscardTicks is the total number of milliseconds spent by all discards.
	DiscardTicks uint64
	// FlushRequestsCompleted is the total number of flush request completed successfully.
	FlushRequestsCompleted uint64
	// TimeSpentFlushing is the total number of milliseconds spent flushing.
	TimeSpentFlushing uint64
	Scheduler         string
	NrRequests        uint64
	ReadAheadKb       uint64
	QueueNum          uint64
}

func (fs FS) DiskStat() (DiskStat, error) {
	diskStat := DiskStat{}

	path := fs.path("diskstats")
	f, err := os.Open(path)
	if err != nil {
		return diskStat, err
	}
	defer f.Close()

	sysBlock := NewSysBlocFS("")
	sysDisks := map[string]BlockDevStat{}

	err = sysBlock.EachBlockDev(func(b BlockDev) error {
		stat, err := b.BlockDevStat()
		if err != nil {
			return err
		}
		sysDisks[b.Name] = stat
		return nil
	})

	if err != nil {
		return diskStat, err
	}

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		var fields [20]string
		nFields := stringutil.FieldsN(line, fields[:])

		if nFields < 14 {
			return fmt.Errorf("unexpected line in diskstats: '%s'", line)
		}

		disk := DiskStatLine{
			DiscardIOs:             math.MaxUint64,
			DiscardMerges:          math.MaxUint64,
			DiscardSectors:         math.MaxUint64,
			DiscardTicks:           math.MaxUint64,
			FlushRequestsCompleted: math.MaxUint64,
			TimeSpentFlushing:      math.MaxUint64,
		}

		disk.MajorNumber, _ = strconv.ParseUint(fields[0], 10, 64)
		disk.MinorNumber, _ = strconv.ParseUint(fields[1], 10, 64)
		disk.DeviceName = strings.Clone(fields[2])

		if s, ok := sysDisks[disk.DeviceName]; !ok {
			return nil
		} else {
			disk.Scheduler = s.Scheduler
			disk.NrRequests = s.NrRequests
			disk.ReadAheadKb = s.ReadAheadKb
			disk.QueueNum = s.QueueNum
		}

		disk.ReadIOs, _ = strconv.ParseUint(fields[3], 10, 64)
		disk.ReadMerges, _ = strconv.ParseUint(fields[4], 10, 64)
		disk.ReadSectors, _ = strconv.ParseUint(fields[5], 10, 64)
		disk.ReadTicks, _ = strconv.ParseUint(fields[6], 10, 64)

		disk.WriteIOs, _ = strconv.ParseUint(fields[7], 10, 64)
		disk.WriteMerges, _ = strconv.ParseUint(fields[8], 10, 64)
		disk.WriteSectors, _ = strconv.ParseUint(fields[9], 10, 64)
		disk.WriteTicks, _ = strconv.ParseUint(fields[10], 10, 64)

		disk.IOsInProgress, _ = strconv.ParseUint(fields[11], 10, 64)
		disk.IOsTotalTicks, _ = strconv.ParseUint(fields[12], 10, 64)
		disk.WeightedIOTicks, _ = strconv.ParseUint(fields[13], 10, 64)

		if nFields < 18 {
			disk.DiscardIOs, _ = strconv.ParseUint(fields[14], 10, 64)
			disk.DiscardMerges, _ = strconv.ParseUint(fields[15], 10, 64)
			disk.DiscardSectors, _ = strconv.ParseUint(fields[16], 10, 64)
			disk.DiscardTicks, _ = strconv.ParseUint(fields[17], 10, 64)
		}

		if nFields < 20 {
			disk.FlushRequestsCompleted, _ = strconv.ParseUint(fields[18], 10, 64)
			disk.TimeSpentFlushing, _ = strconv.ParseUint(fields[19], 10, 64)
		}

		diskStat[disk.DeviceName] = disk

		return nil
	})

	return diskStat, err
}
