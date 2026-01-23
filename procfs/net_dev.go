package procfs

import (
	"os"
	"strconv"
	"strings"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

// NetDevLine is single line parsed from /proc/net/dev or /proc/[pid]/net/dev.
type NetDevLine struct {
	Name         string
	RxBytes      uint64
	RxPackets    uint64
	RxErrors     uint64
	RxDropped    uint64
	RxFIFO       uint64
	RxFrame      uint64
	RxCompressed uint64
	RxMulticast  uint64
	TxBytes      uint64
	TxPackets    uint64
	TxErrors     uint64
	TxDropped    uint64
	TxFIFO       uint64
	TxCollisions uint64
	TxCarrier    uint64
	TxCompressed uint64
}

type NetDev map[string]NetDevLine

func (fs FS) NetDev() (NetDev, error) {
	netDev := NetDev{}

	path := fs.path("net/dev")
	f, err := os.Open(path)
	if err != nil {
		return netDev, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		var err error

		var fields [17]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 17 {
			return nil
		}
		idx := strings.Index(fields[0], ":")
		if idx == -1 {
			return nil
		}

		name := strings.Clone(fields[0][:idx])
		devLine := NetDevLine{}

		devLine.Name = name
		// RX
		devLine.RxBytes, err = strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxPackets, err = strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxErrors, err = strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxDropped, err = strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxFIFO, err = strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxFrame, err = strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxCompressed, err = strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return err
		}
		devLine.RxMulticast, err = strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return err
		}

		// TX
		devLine.TxBytes, err = strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxPackets, err = strconv.ParseUint(fields[10], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxErrors, err = strconv.ParseUint(fields[11], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxDropped, err = strconv.ParseUint(fields[12], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxFIFO, err = strconv.ParseUint(fields[13], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxCollisions, err = strconv.ParseUint(fields[14], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxCarrier, err = strconv.ParseUint(fields[15], 10, 64)
		if err != nil {
			return err
		}
		devLine.TxCompressed, err = strconv.ParseUint(fields[16], 10, 64)
		if err != nil {
			return err
		}
		netDev[name] = devLine
		return nil
	})

	return netDev, err
}
