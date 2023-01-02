package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"github.com/xixiliguo/etop/util"
	"golang.org/x/sys/unix"
)

var (
	DefaultPath                     = "/var/log/etop"
	ErrOutOfRange                   = errors.New("data is out of range")
	MinimumFreeSpaceForStore uint64 = 500 * (1 << 20) // 500MB
)

type SystemSample struct {
	HostName      string
	KernelVersion string
	PageSize      int
	procfs.LoadAvg
	procfs.Stat
	procfs.Meminfo
	NetDevStats map[string]procfs.NetDevLine
	DiskStats   map[string]blockdevice.Diskstats
	procfs.ProcSnmp
	procfs.ProcSnmp6
	procfs.ProcNetstat
	procfs.NetProtocolStats
	SoftNetStats []procfs.SoftnetStat
}

type ProcSample struct {
	procfs.ProcStat
	procfs.ProcIO
}

// Sample represent all system info and process info.
type Sample struct {
	TimeStamp    int64              // unix time when sample was generated
	SystemSample                    // system information
	ProcSamples  map[int]ProcSample // process information
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

// Index represent short info for one sample, so that speed up to find
// specifed sample.
type Index struct {
	TimeStamp int64 // unix time when sample was generated
	Offset    int64 // offset of file where sample is exist
	Len       int64 // length of one sample
}

func (idx *Index) Marshal() []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:], uint64(idx.TimeStamp))
	binary.LittleEndian.PutUint64(b[8:], uint64(idx.Offset))
	binary.LittleEndian.PutUint64(b[16:], uint64(idx.Len))
	return b[:]
}

func (idx *Index) Unmarshal(b []byte) {
	idx.TimeStamp = int64(binary.LittleEndian.Uint64(b[0:]))
	idx.Offset = int64(binary.LittleEndian.Uint64(b[8:]))
	idx.Len = int64(binary.LittleEndian.Uint64(b[16:]))
}

// Store is interface which operate file/network which include sample data.
type Store interface {
	// if step is positive , advance the corresponding number of steps
	// otherwise back, then get sample
	NextSample(step int, sample *Sample) error
	// Get sample which collect at timestamp
	JumpSampleByTimeStamp(timestamp int64, sample *Sample) error
}

type Option func(*LocalStore) error

func WithSetDefault(path string, log *log.Logger) Option {
	return func(local *LocalStore) error {
		if path == "" {
			path = DefaultPath
		}
		local.Path = path
		local.Log = log
		return nil
	}
}

func WithWriteOnly() Option {
	return func(local *LocalStore) error {
		local.writeOnly = true
		return nil
	}
}

func WithSetPath(path string) Option {
	return func(local *LocalStore) error {
		local.Path = path
		return nil
	}
}

// LocalStore represent local store, which consist of index and data files.
// All files was stored into Path (default: /var/log/etop).
// file format: index_{suffix}, data_{suffix}  suffix: yyyymmdd
// It should work either readonly mode or writeonly mode.
type LocalStore struct {
	Path       string   // path which index and data file is stored
	Index      *os.File // current active index file
	Data       *os.File // current active data file
	Log        *log.Logger
	DataOffset int64 // file offset which next sample was written to
	writeOnly  bool
	idxs       []Index // all indexs from Path
	suffix     string
	curIdx     int
}

func NewLocalStore(opts ...Option) (*LocalStore, error) {
	local := &LocalStore{}
	for _, opt := range opts {
		if err := opt(local); err != nil {
			return nil, err
		}
	}

	if local.writeOnly == true {
		suffix := time.Now().Format("20060102")
		if err := local.openFile(suffix, true); err != nil {
			return nil, err
		}
		return local, nil
	}

	// readonly mode
	if idxs, err := getIndexFrames(local.Path); err != nil {
		return local, err
	} else {
		local.idxs = idxs
	}

	return local, nil
}

// getIndexFrames walk path and get data from all index file
// if no any data, return ErrNoIndexFile
func getIndexFrames(path string) ([]Index, error) {

	idxFiles := []string{}
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "index_") {
			idxFiles = append(idxFiles, d.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	idxs := []Index{}
	for _, file := range idxFiles {

		if f, err := os.Open(filepath.Join(path, file)); err == nil {
			buf := make([]byte, 24)
			for {
				n, err := f.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil || n != 24 {
					return nil, err
				}
				index := Index{}
				index.Unmarshal(buf)
				idxs = append(idxs, index)
			}

		} else {
			return nil, err
		}
	}
	sort.Slice(idxs, func(i, j int) bool {
		return idxs[i].TimeStamp < idxs[j].TimeStamp
	})
	return idxs, nil
}

func (local *LocalStore) FileStatInfo() (result string, err error) {
	suffixs, indexSize, dataSize, err := getIndexAndDataInfo(local.Path)
	if err != nil {
		return
	}
	result += fmt.Sprintf("%-5s: %d files %s\n", "Index", len(suffixs), util.GetHumanSize(indexSize))
	result += fmt.Sprintf("%-5s: %d files %s\n", "Data", len(suffixs), util.GetHumanSize(dataSize))

	start := time.Unix(local.idxs[0].TimeStamp, 0).Format(time.RFC3339)
	end := time.Unix(local.idxs[len(local.idxs)-1].TimeStamp, 0).Format(time.RFC3339)
	result += fmt.Sprintf("%d samples from %s to %s", len(local.idxs), start, end)
	return
}

// getDataInfo walk path and get info of index file and data file
// return suffixs array and total size
func getIndexAndDataInfo(path string) (suffixs []string, indexSize int64, dataSize int64, err error) {
	indexNum := 0
	dataNum := 0
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "index_") {
			indexNum++
			if info, err := d.Info(); err == nil {
				indexSize += info.Size()
			} else {
				return err
			}
		}
		if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "data_") {
			dataNum++
			suffixs = append(suffixs, strings.TrimLeft(d.Name(), "data_"))
			if info, err := d.Info(); err == nil {
				dataSize += info.Size()
			} else {
				return err
			}
		}
		return nil
	})
	if indexNum != dataNum {
		return suffixs, indexSize, dataSize,
			fmt.Errorf("%d index files is not equal to %d data files", indexNum, dataNum)
	}
	sort.Slice(suffixs, func(i, j int) bool {
		return suffixs[i] < suffixs[j]
	})
	return
}

func (local *LocalStore) openFile(suffix string, writeonly bool) (err error) {

	flags := os.O_RDONLY
	if writeonly == true {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	idxPath := filepath.Join(local.Path, fmt.Sprintf("index_%s", suffix))
	if local.Index, err = os.OpenFile(idxPath, flags, 0644); err != nil {
		return err
	}

	dataPath := filepath.Join(local.Path, fmt.Sprintf("data_%s", suffix))
	if local.Data, err = os.OpenFile(dataPath, flags, 0644); err != nil {
		local.Index.Close()
		return err
	}

	// data file may already exist and take file size as initail DataOffset
	if info, err := local.Data.Stat(); err == nil {
		local.DataOffset = info.Size()
	} else {
		return err
	}

	local.suffix = suffix
	return nil
}

func (local *LocalStore) changeFile(suffix string, writeOnly bool) error {

	local.Index.Close()
	local.Data.Close()
	return local.openFile(suffix, writeOnly)
}

func (local *LocalStore) Close() error {
	if err := local.Index.Close(); err != nil {
		return err
	}
	return local.Data.Close()
}

func (local *LocalStore) NextSample(step int, sample *Sample) error {
	target := local.curIdx + step
	if target < 0 {
		return ErrOutOfRange
	}
	if target >= len(local.idxs) {
		// try to read file again and check if have new data available
		if idxs, err := getIndexFrames(local.Path); err != nil {
			return err
		} else {
			local.idxs = idxs
			if target >= len(local.idxs) {
				return ErrOutOfRange
			}
		}
	}

	suffix := time.Unix(local.idxs[target].TimeStamp, 0).Format("20060102")

	if suffix != local.suffix {
		if err := local.changeFile(suffix, false); err != nil {
			return err
		}
	}

	if err := local.getSample(target, sample); err != nil {
		return err
	}

	local.curIdx = target
	return nil
}

// JumpSampleByTimeStamp get sample by specific timestamp (unix time)
// if no, search the nearest one.
// ignore 1st sample for better handle corner case.
func (local *LocalStore) JumpSampleByTimeStamp(timestamp int64, sample *Sample) error {

	target := sort.Search(len(local.idxs), func(i int) bool {
		return local.idxs[i].TimeStamp >= timestamp
	})

	if target >= len(local.idxs) {
		// try to read file again and check if have new data available
		if idxs, err := getIndexFrames(local.Path); err != nil {
			return err
		} else {
			local.idxs = idxs
			target = sort.Search(len(local.idxs), func(i int) bool {
				return local.idxs[i].TimeStamp >= timestamp
			})
		}
	}
	if len(local.idxs) == 0 || len(local.idxs) == 1 {
		return ErrOutOfRange
	}
	if target >= len(local.idxs) {
		target = len(local.idxs) - 1
	}
	if target == 0 {
		target = 1
	}

	suffix := time.Unix(local.idxs[target].TimeStamp, 0).Format("20060102")

	if suffix != local.suffix {
		if err := local.changeFile(suffix, false); err != nil {
			return err
		}
	}

	if err := local.getSample(target, sample); err != nil {
		return err
	}

	local.curIdx = target
	return nil
}

func (local *LocalStore) getSample(target int, sample *Sample) error {
	idx := local.idxs[target]
	buff := make([]byte, idx.Len)
	n, err := local.Data.ReadAt(buff, idx.Offset)
	if err != nil {
		return err
	}

	if n != int(idx.Len) {
		return fmt.Errorf("got %d bytes, but want %d\n", n, idx.Len)
	}

	if err := sample.Unmarshal(buff); err != nil {
		return err
	}
	return nil
}

func (local *LocalStore) CollectSample(s *Sample) error {
	return CollectSampleFromSys(s)
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
		s.ProcSnmp = snmp
	}
	if snmp6, err := p.Snmp6(); err != nil {
		return err
	} else {
		s.ProcSnmp6 = snmp6
	}
	if netStat, err := p.Netstat(); err != nil {
		return err
	} else {
		s.ProcNetstat = netStat
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

func (local *LocalStore) WriteSample(s *Sample) error {

	var err error

	compressed := []byte{}
	if compressed, err = s.Marshal(); err != nil {
		return err
	}

	idx := Index{
		TimeStamp: s.TimeStamp,
		Offset:    local.DataOffset,
		Len:       int64(len(compressed)),
	}
	_, err = local.Index.Write(idx.Marshal())
	if err != nil {
		return err
	}
	_, err = local.Data.Write(compressed)

	if err != nil {
		return err
	}

	local.DataOffset += int64(len(compressed))
	return nil
}

type WriteOption struct {
	Interval   time.Duration
	RetainDay  int
	RetainSize int64
}

func (local *LocalStore) WriteLoop(opt WriteOption) error {

	interval := opt.Interval

	local.Log.Printf("start to write sample every %s to %s", interval.String(), local.Data.Name())
	isSkip := 0
	for {
		start := time.Now()
		s := NewSample()
		if err := local.CollectSample(&s); err != nil {
			return err
		}

		statInfo := syscall.Statfs_t{}
		if err := syscall.Statfs(local.Path, &statInfo); err != nil {
			return err
		}
		if statInfo.Bavail*uint64(statInfo.Bsize) > MinimumFreeSpaceForStore {
			if isSkip != 0 {
				local.Log.Printf("resume to write sample (%d skipped)", isSkip)
				isSkip = 0
			}

			suffix := time.Unix(s.TimeStamp, 0).Format("20060102")
			if suffix != local.suffix {
				if err := local.changeFile(suffix, true); err != nil {
					return err
				}
				local.Log.Printf("switch and write data into data_%s", suffix)
				// it is time to check if clean old data or not.
				// oldest first.
				if suffixs, _, size, err := getIndexAndDataInfo(local.Path); err != nil {
					return err
				} else {
					if len(suffixs) > 0 {
						if len(suffixs)-1 > opt.RetainDay || size > opt.RetainSize {
							err := os.Remove("index_" + suffixs[0])
							if err != nil {
								return err
							}
							err = os.Remove("data_" + suffixs[0])
							if err != nil {
								return err
							}
							local.Log.Printf("total historical files: %d %s",
								len(suffixs)-1, util.GetHumanSize(size))
							local.Log.Printf("delete oldest data_%s", suffix)
						}
					}
				}
			}

			if err := local.WriteSample(&s); err != nil {
				return err
			}
		} else {
			if isSkip == 0 {
				local.Log.Printf("filesystem free space %s below %s, write sample skipped",
					util.GetHumanSize(statInfo.Bavail*uint64(statInfo.Bsize)),
					util.GetHumanSize(MinimumFreeSpaceForStore))
			}
			isSkip++
		}

		collectDuration := time.Now().Sub(start)
		if collectDuration > 500*time.Millisecond {
			local.Log.Printf("write sample take %s (larger than 500 ms)", collectDuration.String())
		}
		sleepDuration := time.Duration(1 * time.Second)
		if interval-collectDuration > 1*time.Second {
			sleepDuration = interval - collectDuration
		}
		time.Sleep(sleepDuration)
	}
}
