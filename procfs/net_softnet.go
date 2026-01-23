package procfs

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

// SoftnetStat contains a single row of data from /proc/net/softnet_stat.
type SoftnetStat struct {
	// Number of processed packets.
	Processed uint64
	// Number of dropped packets.
	Dropped uint64
	// Number of times processing packets ran out of quota.
	TimeSqueezed uint64
	// Number of collision occur while obtaining device lock while transmitting.
	CPUCollision uint64
	// Number of times cpu woken up received_rps.
	ReceivedRps uint64
	// number of times flow limit has been reached.
	FlowLimitCount uint64
	// Softnet backlog status.
	SoftnetBacklogLen uint64
	// CPU id owning this softnet_data.
	Index uint64
}

// NetSoftnetStat reads data from /proc/net/softnet_stat.
func (fs FS) NetSoftnetStat() ([]SoftnetStat, error) {

	softNetStats := []SoftnetStat{}

	path := fs.path("net/softnet_stat")
	f, err := os.Open(path)
	if err != nil {
		return softNetStats, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {

		softnetStat := SoftnetStat{
			Processed:         math.MaxUint64,
			Dropped:           math.MaxUint64,
			TimeSqueezed:      math.MaxUint64,
			CPUCollision:      math.MaxUint64,
			ReceivedRps:       math.MaxUint64,
			FlowLimitCount:    math.MaxUint64,
			SoftnetBacklogLen: math.MaxUint64,
			Index:             math.MaxUint64,
		}
		var fields [13]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 9 {
			return fmt.Errorf("unexpected line in /proc/net/softnet_stat: '%s'", line)
		}

		if nFields >= 9 {
			softnetStat.Processed, _ = strconv.ParseUint(fields[0], 16, 64)
			softnetStat.Dropped, _ = strconv.ParseUint(fields[1], 16, 64)
			softnetStat.TimeSqueezed, _ = strconv.ParseUint(fields[2], 16, 64)
			softnetStat.CPUCollision, _ = strconv.ParseUint(fields[8], 16, 64)
		}
		if nFields >= 10 {
			softnetStat.ReceivedRps, _ = strconv.ParseUint(fields[9], 16, 64)
		}

		if nFields >= 11 {
			softnetStat.FlowLimitCount, _ = strconv.ParseUint(fields[10], 16, 64)
		}

		if nFields >= 13 {
			softnetStat.SoftnetBacklogLen, _ = strconv.ParseUint(fields[11], 16, 64)
			softnetStat.Index, _ = strconv.ParseUint(fields[12], 16, 64)
		} else {
			softnetStat.Index = uint64(i)
		}
		softNetStats = append(softNetStats, softnetStat)
		return nil
	})

	return softNetStats, nil
}
