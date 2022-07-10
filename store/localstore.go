package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"golang.org/x/sys/unix"
)

var (
	DefaultPath = "/var/log/etop"
	EOUTOFRANGE = errors.New("already out of data range")
)

type SystemSample struct {
	HostName      string
	KernelVersion string
	procfs.LoadAvg
	procfs.Stat
	procfs.Meminfo
	NetInfo   []procfs.NetDevLine
	DiskStats []blockdevice.Diskstats
}

type ProcSample struct {
	procfs.ProcStat
	procfs.ProcIO
}

type Sample struct {
	CurrTime int64
	SystemSample
	ProcSamples []ProcSample
}

type Store interface {
	AdjustIndex(step int) error
	ReadSample(s *Sample) error
	ChangeIndex(value string) error
}

type LocalStore struct {
	Path            string
	DataName        string
	Data            *os.File
	Log             *log.Logger
	NextIndexOffset int
	DataOffset      int
	readOnly        bool
	idxs            []Index
	curIdx          int
	dec             *zstd.Decoder
	enc             *zstd.Encoder
}

func NewLocalStore(path string, log *log.Logger) (*LocalStore, error) {
	if path == "" {
		path = DefaultPath
	}
	if _, err := os.Stat(path); err != nil {
		err := os.Mkdir(path, 0750)
		if err != nil {
			return nil, err
		}
	}

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	dataName := fmt.Sprintf("etop_%s", time.Now().Format("20060102"))

	data, err := os.OpenFile(filepath.Join(path, dataName), flags, 0644)
	if err != nil {
		return nil, err
	}
	initOffset := 0
	if info, err := data.Stat(); err == nil {
		initOffset = int(info.Size())
	} else {
		return nil, err
	}

	dec, _ := zstd.NewReader(nil)
	enc, _ := zstd.NewWriter(nil)
	localStore := &LocalStore{
		Path:            path,
		DataName:        dataName,
		Data:            data,
		Log:             log,
		NextIndexOffset: initOffset,
		readOnly:        false,
		dec:             dec,
		enc:             enc,
	}
	return localStore, nil
}

func NewLocalStoreWithReadOnly(fileName string, log *log.Logger) (*LocalStore, error) {

	path := DefaultPath
	dataName := fileName
	if filepath.IsAbs(fileName) {
		path, dataName = filepath.Split(fileName)
	}
	if dataName == "" {
		dataName = fmt.Sprintf("etop_%s", time.Now().Format("20060102"))
	}

	f, err := os.OpenFile(filepath.Join(path, dataName), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	dec, _ := zstd.NewReader(nil)
	enc, _ := zstd.NewWriter(nil)
	localStore := &LocalStore{
		Path:     path,
		Data:     f,
		DataName: dataName,
		Log:      log,
		readOnly: true,
		idxs:     []Index{},
		curIdx:   0,
		dec:      dec,
		enc:      enc,
	}

	return localStore, nil
}

type Index struct {
	Time   int64
	Offset int
	Len    int
}

func (idx *Index) Marshal() []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint64(b[0:], uint64(idx.Time))
	binary.LittleEndian.PutUint64(b[8:], uint64(idx.Offset))
	binary.LittleEndian.PutUint64(b[16:], uint64(idx.Len))
	return b[:]
}

func (idx *Index) Unmarshal(b []byte) {
	idx.Time = int64(binary.LittleEndian.Uint64(b[0:]))
	idx.Offset = int(binary.LittleEndian.Uint64(b[8:]))
	idx.Len = int(binary.LittleEndian.Uint64(b[16:]))
}

func (local *LocalStore) Close() error {
	return local.Data.Close()
}

func (local *LocalStore) ChangeIndex(value string) error {

	f := strings.TrimPrefix(local.DataName, "etop_") + " " + value

	t, err := time.ParseInLocation("20060102 15:04", f, time.Local)
	if err != nil {
		return err
	}

	secs := t.Unix()
	if secs < local.idxs[1].Time {
		return EOUTOFRANGE
	}

	if secs > local.idxs[len(local.idxs)-1].Time {
		for {
			end, err := local.tryExpandIndex()
			if err != nil {
				return err
			}
			if secs < local.idxs[len(local.idxs)-1].Time {
				break
			}
			if end == true {
				return EOUTOFRANGE
			}
		}
	}

	n := sort.Search(len(local.idxs), func(i int) bool {
		return local.idxs[i].Time >= secs
	})

	if n == len(local.idxs) {
		return fmt.Errorf("maybe bug")
	}

	local.curIdx = n

	return nil
}

func (local *LocalStore) AdjustIndex(step int) error {
	// local.Log.Printf("adjust %d step, %v", step, local.idxs)
	newIdx := local.curIdx + step
	if newIdx < 0 {
		return EOUTOFRANGE
	}

	if newIdx >= len(local.idxs) {
		for {
			end, err := local.tryExpandIndex()
			if err != nil {
				return err
			}
			if newIdx < len(local.idxs) {
				break
			}
			if end == true {
				return EOUTOFRANGE
			}
		}
	}
	local.curIdx = newIdx
	return nil
}

func (local *LocalStore) tryExpandIndex() (end bool, err error) {
	info, err := local.Data.Stat()
	if err != nil {
		return false, err
	}
	nextIdx := 0
	if len(local.idxs) != 0 {
		lastIdx := local.idxs[len(local.idxs)-1]
		nextIdx = lastIdx.Offset + lastIdx.Len
	}

	i := 0
	for ; int64(nextIdx) < info.Size() && i < 3600; i++ {
		local.Data.Seek(int64(nextIdx), os.SEEK_SET)
		var meta [24]byte
		var idx Index
		n, err := local.Data.Read(meta[:])
		if err != nil {
			return false, err
		}
		if n != 24 {
			return false, fmt.Errorf("expect 24, but get %d", n)
		}
		idx.Unmarshal(meta[:])
		local.idxs = append(local.idxs, idx)
		nextIdx = idx.Offset + idx.Len
	}
	if int64(nextIdx) < info.Size() {
		return false, nil
	}
	return true, nil
}
func (local *LocalStore) CollectSampleFromSys(s *Sample) error {
	return CollectSampleFromSys(s)
}

func CollectSampleFromSys(s *Sample) error {
	//collect one sample
	var (
		fs     procfs.FS
		diskFS blockdevice.FS
		err    error
	)
	u := unix.Utsname{}
	unix.Uname(&u)
	s.HostName = string(u.Nodename[:])
	s.KernelVersion = string(u.Release[:])
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

	if s.SystemSample.Stat, err = fs.Stat(); err != nil {
		return err
	}

	if s.SystemSample.Meminfo, err = fs.Meminfo(); err != nil {
		return err
	}

	if netDev, err := fs.NetDev(); err == nil {
		for _, v := range netDev {
			s.NetInfo = append(s.NetInfo, v)
		}
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
				s.SystemSample.DiskStats = append(s.SystemSample.DiskStats, diskStat)
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
		s.ProcSamples = append(s.ProcSamples, p)
	}
	return nil
}

func (local *LocalStore) SetDestFile(fileName string) error {
	data, err := os.OpenFile(filepath.Join(local.Path, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	initOffset := 0
	if info, err := data.Stat(); err == nil {
		initOffset = int(info.Size())
	} else {
		return err
	}
	if err := local.Data.Close(); err != nil {
		return err
	}
	local.DataName = fileName
	local.Data = data
	local.NextIndexOffset = initOffset
	return nil
}

func (local *LocalStore) WriteLoop(interval time.Duration) error {
	local.Log.Printf("start to write sample every %s to %s", interval.String(), local.Data.Name())
	for {
		start := time.Now()
		if dstFile := "etop_" + start.Format("20060102"); dstFile > local.DataName {
			local.Log.Printf("change sample file from %s to %s", local.DataName, dstFile)
			if err := local.SetDestFile(dstFile); err != nil {
				return err
			}
		}
		s := &Sample{}
		if err := local.CollectSampleFromSys(s); err != nil {
			return err
		}
		s.CurrTime = time.Now().Unix()
		if err := local.WriteSample(s); err != nil {
			return err
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

func (local *LocalStore) WriteSample(s *Sample) error {

	b, err := cbor.Marshal(s)
	if err != nil {
		return err
	}
	compressed := local.enc.EncodeAll(b, make([]byte, 0, len(b)))

	idx := Index{
		Time:   s.CurrTime,
		Offset: local.NextIndexOffset + 24,
		Len:    len(compressed),
	}

	iovs := make([][]byte, 2)
	iovs[0] = idx.Marshal()
	iovs[1] = compressed

	n, err := unix.Writev(int(local.Data.Fd()), iovs)

	if err != nil {
		return err
	}
	if n != 24+len(compressed) {
		return fmt.Errorf("want to write %d bytes, but get %d", 24+len(compressed), n)
	}

	local.NextIndexOffset += n
	return nil
}

func (local *LocalStore) ReadSample(s *Sample) error {

	idx := local.idxs[local.curIdx]

	data := make([]byte, idx.Len)
	_, err := local.Data.Seek(int64(idx.Offset), os.SEEK_SET)
	if err != nil {
		return err
	}
	n, err := local.Data.Read(data)
	if err != nil {
		return err
	}
	if n != idx.Len {
		return fmt.Errorf("read sample, want %d bytes, but get %d", idx.Len, n)
	}

	uncompressed, err := local.dec.DecodeAll(data, make([]byte, 0, len(data)))
	if err != nil {
		return err
	}

	err = cbor.Unmarshal(uncompressed, s)
	if err != nil {
		return err
	}
	return nil
}
