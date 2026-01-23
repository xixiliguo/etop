package procfs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

// NetProtocolStats stores the contents from /proc/net/protocols.
type NetProtocolStats map[string]NetProtocolStatLine

// NetProtocolStatLine contains a single line parsed from /proc/net/protocols. We
// only care about the first six columns as the rest are not likely to change
// and only serve to provide a set of capabilities for each protocol.
type NetProtocolStatLine struct {
	Name       string // 0 The name of the protocol
	Size       uint64 // 1 The size, in bytes, of a given protocol structure. e.g. sizeof(struct tcp_sock) or sizeof(struct unix_sock)
	Sockets    uint64 // 2 Number of sockets in use by this protocol
	Memory     int    // 3 Number of 4KB pages allocated by all sockets of this protocol
	Pressure   string // 4 This is either yes, no, or NI (not implemented). For the sake of simplicity we treat NI as not experiencing memory pressure.
	MaxHeader  uint64 // 5 Protocol specific max header size
	Slab       string // 6 Indicates whether or not memory is allocated from the SLAB
	ModuleName string // 7 The name of the module that implemented this protocol or "kernel" if not from a module
}

// NetProtocols reads stats from /proc/net/protocols and returns a map of
// PortocolStatLine entries. As of this writing no official Linux Documentation
// exists, however the source is fairly self-explanatory and the format seems
// stable since its introduction in 2.6.12-rc2
// Linux 2.6.12-rc2 - https://elixir.bootlin.com/linux/v2.6.12-rc2/source/net/core/sock.c#L1452
// Linux 5.10 - https://elixir.bootlin.com/linux/v5.10.4/source/net/core/sock.c#L3586
func (fs FS) NetProtocols() (NetProtocolStats, error) {
	netPros := NetProtocolStats{}

	path := fs.path("net/protocols")
	f, err := os.Open(path)
	if err != nil {
		return netPros, err
	}
	defer f.Close()

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		netPro := NetProtocolStatLine{}
		var fields [27]string
		nFields := stringutil.FieldsN(line, fields[:])
		if nFields < 27 {
			return fmt.Errorf("unexpected line in /proc/net/protocols: '%s'", line)
		}
		netPro.Name = string(fields[0])
		netPro.Size, _ = strconv.ParseUint(fields[1], 10, 64)
		netPro.Sockets, _ = strconv.ParseUint(fields[2], 10, 64)
		netPro.Memory, _ = strconv.Atoi(fields[3])
		netPro.Pressure = string(fields[4])
		netPro.MaxHeader, _ = strconv.ParseUint(fields[5], 10, 64)
		netPro.Slab = string(fields[6])
		netPro.ModuleName = string(fields[7])
		netPros[netPro.Name] = netPro
		return nil
	})
	return netPros, err
}
