package procfs

import (
	"fmt"
	"strconv"

	"github.com/xixiliguo/etop/internal/stringutil"
)

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

func (fs FS) VmStat() (VmStat, error) {

	vmStat := VmStat{}

	path := fs.path("vmstat")

	err := fs.processFile(path, func(i int, line string) error {
		var fields [2]string
		nFields := stringutil.FieldsN(line, fields[:])

		if nFields < 2 {
			return fmt.Errorf("unexpected line in vmstat: '%s'", line)
		}

		v, err := strconv.ParseUint(fields[1], 0, 64)
		if err != nil {
			return err
		}

		switch fields[0] {
		case "pgpgin":
			vmStat.PageIn = v
		case "pgpgout":
			vmStat.PageOut = v
		case "pswpin":
			vmStat.SwapIn = v
		case "pswpout":
			vmStat.SwapOut = v
		// case "pgscan_direct_throttle":
		case "pgscan_kswapd":
			vmStat.PageScanKswapd = v
		case "pgscan_direct":
			vmStat.PageScanDirect = v
		case "pgsteal_kswapd":
			vmStat.PageStealKswapd = v
		case "pgsteal_direct":
			vmStat.PageStealDirect = v
		case "oom_kill":
			vmStat.OOMKill = v
		}
		return nil
	})

	return vmStat, err
}
