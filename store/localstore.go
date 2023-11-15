package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/klauspost/compress/zstd"
	"github.com/xixiliguo/etop/util"
)

var (
	ErrOutOfRange                   = errors.New("data is out of range")
	ErrIndexCorrupt                 = errors.New("corrupt index")
	ErrDataCorrupt                  = errors.New("corrupt data")
	MinimumFreeSpaceForStore uint64 = 500 * (1 << 20) // 500MB
	ShardTime                       = int64(24 * 60 * 60)
)

const (
	NoCompress = uint32(1 << iota)
	ZstdCompress
	ZstdCompressWithDict
)
const (
	CompressModeShift = 0
	CompressModeLen   = 8
	DictOffsetShift   = 8
	DictOffsetLen     = 24
	MaxDictOffset     = 1<<DictOffsetLen - 1
)

// Index represent short info for one sample, so that speed up to find
// specifed sample.
type Index struct {
	TimeStamp int64  // unix time when sample was generated
	Offset    int64  // offset of file where sample is exist
	Len       int64  // length of one sample
	Flag      uint32 // compress mode and dict offset
	CRC       uint32 // crc for whole index, except CRC self
}

var sizeIndex = binary.Size(Index{})

func mask(len uint32) uint32 {
	return (1 << len) - 1
}

func readBits(value, len, shift uint32) uint32 {
	return (value >> shift) & mask(len)
}

func writeBits(value, len, shift, new uint32) uint32 {
	value &^= mask(len) << shift
	value |= (new & mask(len)) << shift
	return value
}

func (idx *Index) SetCompressMode(m uint32, offset uint32) {
	idx.Flag = writeBits(idx.Flag, CompressModeLen, CompressModeShift, m)
	idx.Flag = writeBits(idx.Flag, DictOffsetLen, DictOffsetShift, offset)
}

func (idx *Index) CompressMode() (uint32, uint32) {
	mode := readBits(idx.Flag, CompressModeLen, CompressModeShift)
	offset := readBits(idx.Flag, DictOffsetLen, DictOffsetShift)
	return mode, offset
}

func (idx *Index) Marshal() []byte {
	b := make([]byte, sizeIndex)
	binary.LittleEndian.PutUint64(b[0:], uint64(idx.TimeStamp))
	binary.LittleEndian.PutUint64(b[8:], uint64(idx.Offset))
	binary.LittleEndian.PutUint64(b[16:], uint64(idx.Len))
	binary.LittleEndian.PutUint32(b[24:], idx.Flag)
	binary.LittleEndian.PutUint32(b[28:], idx.CRC)
	return b
}

func (idx *Index) Unmarshal(b []byte) {
	idx.TimeStamp = int64(binary.LittleEndian.Uint64(b[0:]))
	idx.Offset = int64(binary.LittleEndian.Uint64(b[8:]))
	idx.Len = int64(binary.LittleEndian.Uint64(b[16:]))
	idx.Flag = binary.LittleEndian.Uint32(b[24:])
	idx.CRC = binary.LittleEndian.Uint32(b[28:])
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

func WithPathAndLogger(path string, log *slog.Logger) Option {
	return func(local *LocalStore) error {
		local.Path = path
		local.Log = log
		return nil
	}
}

func WithWriteOnly(mode, chunk uint32) Option {
	return func(local *LocalStore) error {
		local.writeOnly = true
		local.mode = mode
		local.chunk = chunk
		local.handelSignal()
		return nil
	}
}

func WithExitProcess(log *slog.Logger) Option {
	return func(local *LocalStore) error {
		local.exit = NewExitProcess(log)
		go local.exit.Collect()
		return nil
	}
}

// LocalStore represent local store, which consist of index and data files.
// All files was stored into Path (default: /var/log/etop).
// file format: index_{shard}, data_{shard}
// shard is unix time of 00:00 utc every day
// It should work either readonly mode or writeonly mode.
type LocalStore struct {
	Path       string   // path which index and data file is stored
	Index      *os.File // current active index file
	Data       *os.File // current active data file
	Log        *slog.Logger
	DataOffset int64 // file offset which next sample was written to
	writeOnly  bool
	mode       uint32 // compress mode
	// increment after writing one sample
	// reset to 0 when opening new file or initializing instance of LocalStore
	next    uint32
	curDict int64         // timestamp of sample which was used as dict
	chunk   uint32        // a number of adjacent samples as one group
	enc     *zstd.Encoder // zstd encoder
	dec     *zstd.Decoder // zstd decoder
	encDict *zstd.Encoder // zstd encoder with dict
	decDict *zstd.Decoder // zstd decoder with dict
	sync.Mutex
	closed   bool
	closeSig chan os.Signal
	idxs     []Index // all indexs from Path
	shard    int64
	curIdx   int
	exit     *ExitProcess
}

func NewLocalStore(opts ...Option) (*LocalStore, error) {
	local := &LocalStore{}
	for _, opt := range opts {
		if err := opt(local); err != nil {
			return nil, err
		}
	}

	// 0 is valid, so cannot set to 0 as initial value
	local.shard = -1
	// so that nextSample(1, &s) after NewLocalStore can get first sample
	local.curIdx = -1

	if local.mode != NoCompress {
		local.enc, _ = zstd.NewWriter(
			nil,
			zstd.WithLowerEncoderMem(true),
			zstd.WithEncoderConcurrency(1),
		)
		local.dec, _ = zstd.NewReader(
			nil,
			zstd.WithDecoderLowmem(true),
			zstd.WithDecoderConcurrency(1),
		)
	}

	if local.writeOnly == true {
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

func calcshard(sec int64) int64 {
	return sec - sec%ShardTime
}

// getIndexFrames walk path and get data from all index file
// if no any data, return ErrNoIndexFile
func getIndexFrames(path string) ([]Index, error) {

	idxFiles := []string{}
	if ds, err := os.ReadDir(path); err == nil {
		for _, d := range ds {
			if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "index_") {
				idxFiles = append(idxFiles, d.Name())
			}
		}
	} else {
		return nil, err
	}

	idxs := []Index{}
	for _, file := range idxFiles {

		if f, err := os.Open(filepath.Join(path, file)); err == nil {
			buf := make([]byte, sizeIndex)
			for {
				n, err := f.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				if n != sizeIndex {
					return nil, fmt.Errorf("request %d byte, but got %d from %s", sizeIndex, n, f.Name())
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

func (local *LocalStore) handelSignal() {
	local.closeSig = make(chan os.Signal, 1)
	signal.Notify(local.closeSig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		s := <-local.closeSig
		local.Lock()
		defer local.Unlock()
		local.closed = true
		msg := fmt.Sprintf("recevice %s signal, exiting", s.String())
		local.Log.Info(msg)
	}()
	return
}

func (local *LocalStore) FileStatInfo() (result string, err error) {
	suffixs, indexSize, dataSize, err := getIndexAndDataInfo(local.Path)
	if err != nil {
		return
	}
	if len(suffixs) == 0 {
		return fmt.Sprintf("no data at %s", local.Path), nil
	}

	result += fmt.Sprintf("%-5s: %s\n", "Path", local.Path)
	result += fmt.Sprintf("%-5s: %d files %s\n", "Index", len(suffixs), util.GetHumanSize(indexSize))
	result += fmt.Sprintf("%-5s: %d files %s\n", "Data", len(suffixs), util.GetHumanSize(dataSize))

	start := time.Unix(local.idxs[0].TimeStamp, 0).Format(time.RFC3339)
	end := time.Unix(local.idxs[len(local.idxs)-1].TimeStamp, 0).Format(time.RFC3339)
	result += fmt.Sprintf("%d samples from %s to %s", len(local.idxs), start, end)
	return
}

// getDataInfo walk path and get info of index file and data file
// return suffixs array and total size
func getIndexAndDataInfo(path string) (shards []int64, indexSize int64, dataSize int64, err error) {
	indexNum := 0
	dataNum := 0

	if ds, err := os.ReadDir(path); err == nil {
		for _, d := range ds {
			if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "index_") {
				indexNum++
				if info, err := d.Info(); err == nil {
					indexSize += info.Size()
				} else {
					return shards, indexSize, dataSize, err
				}
			}
			if d.Type().IsRegular() && strings.HasPrefix(d.Name(), "data_") {
				dataNum++
				var shard int64
				if _, err := fmt.Sscanf(d.Name(), "data_0%d", &shard); err != nil {
					return shards, indexSize, dataSize, err
				}
				shards = append(shards, shard)
				if info, err := d.Info(); err == nil {
					dataSize += info.Size()
				} else {
					return shards, indexSize, dataSize, err
				}
			}
		}
	} else {
		return shards, indexSize, dataSize, err
	}

	if indexNum != dataNum {
		return shards, indexSize, dataSize,
			fmt.Errorf("%d index files is not equal to %d data files", indexNum, dataNum)
	}
	sort.Slice(shards, func(i, j int) bool {
		return shards[i] < shards[j]
	})
	return
}

func (local *LocalStore) openFile(shard int64, writeonly bool) (err error) {

	flags := os.O_RDONLY
	if writeonly == true {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	lock := syscall.LOCK_EX | syscall.LOCK_NB

	idxPath := filepath.Join(local.Path, fmt.Sprintf("index_%011d", shard))
	if local.Index, err = os.OpenFile(idxPath, flags, 0644); err != nil {
		return err
	}
	if writeonly == true && syscall.Flock(int(local.Index.Fd()), lock) != nil {
		return fmt.Errorf("can not acquire lock for file %s", idxPath)
	}

	dataPath := filepath.Join(local.Path, fmt.Sprintf("data_%011d", shard))
	if local.Data, err = os.OpenFile(dataPath, flags, 0644); err != nil {
		local.Index.Close()
		return err
	}
	if writeonly == true && syscall.Flock(int(local.Data.Fd()), lock) != nil {
		return fmt.Errorf("can not acquire lock for file %s", dataPath)
	}

	// should be new group for zstd dict compress when opening new file
	// reset next to 0
	local.next = 0
	// data file may already exist and take file size as initail DataOffset
	if info, err := local.Data.Stat(); err == nil {
		local.DataOffset = info.Size()
	} else {
		return err
	}

	local.shard = shard
	return nil
}

func (local *LocalStore) changeFile(shard int64, writeOnly bool) error {

	local.Index.Close()
	local.Data.Close()
	return local.openFile(shard, writeOnly)
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

	shard := calcshard(local.idxs[target].TimeStamp)

	if shard != local.shard {
		if err := local.changeFile(shard, false); err != nil {
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
		local.Log.Warn("None or only one sample", slog.String("path", local.Path))
		return ErrOutOfRange
	}
	if target >= len(local.idxs) {
		target = len(local.idxs) - 1
	}
	if target == 0 {
		target = 1
	}

	shard := calcshard(local.idxs[target].TimeStamp)

	if shard != local.shard {
		if err := local.changeFile(shard, false); err != nil {
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
	var err error
	if err = local.getDataBytes(idx, &buff); err != nil {
		return err
	}

	mode, offset := idx.CompressMode()
	if mode == ZstdCompress || (mode == ZstdCompressWithDict && offset == 0) {
		if buff, err = local.dec.DecodeAll(buff, make([]byte, 0, len(buff))); err != nil {
			return err
		}
	} else if mode == ZstdCompressWithDict {
		dictIdxPos := target - int(offset)
		dictIdx := local.idxs[dictIdxPos]
		dictMode, dictOffSet := dictIdx.CompressMode()
		if dictMode != ZstdCompressWithDict || dictOffSet != 0 {
			return fmt.Errorf("dict of %+v is not compress mode: %+v", idx, dictIdx)
		}
		if dictIdx.TimeStamp != local.curDict {
			dictBuff := make([]byte, dictIdx.Len)
			if err = local.getDataBytes(dictIdx, &dictBuff); err != nil {
				return err
			}
			if dictBuff, err = local.dec.DecodeAll(dictBuff, make([]byte, 0, len(dictBuff))); err != nil {
				return err
			}
			opts := []zstd.DOption{
				zstd.WithDecoderDictRaw(0, dictBuff),
				zstd.WithDecoderLowmem(true),
				zstd.WithDecoderConcurrency(1),
			}
			if local.decDict, err = zstd.NewReader(nil, opts...); err != nil {
				return err
			}
			local.curDict = dictIdx.TimeStamp
		}
		if buff, err = local.decDict.DecodeAll(buff, make([]byte, 0, len(buff))); err != nil {
			return err
		}
	}

	if err := sample.Unmarshal(buff); err != nil {
		return err
	}

	return nil
}

func (local *LocalStore) getDataBytes(idx Index, buff *[]byte) error {
	if idx.CRC != crc32.ChecksumIEEE((*[32]byte)(unsafe.Pointer(&idx))[:28]) {
		return fmt.Errorf("%s timestamp %d: %w", local.Index.Name(), idx.TimeStamp, ErrIndexCorrupt)
	}

	n, err := local.Data.ReadAt(*buff, idx.Offset)
	if err != nil {
		return fmt.Errorf("offset %d read %d bytes at %s: %w", idx.Offset, idx.Len, local.Data.Name(), err)
	}

	if n != int(idx.Len) {
		return fmt.Errorf("got %d bytes, but want %d from %s\n", n, idx.Len, local.Data.Name())
	}
	return nil
}

func (local *LocalStore) CollectSample(s *Sample) error {
	return CollectSampleFromSys(s, local.exit, local.Log)
}

func (local *LocalStore) WriteSample(s *Sample) (bool, error) {

	newSuffix := false
	shard := calcshard(s.TimeStamp)
	if shard != local.shard {
		if err := local.changeFile(shard, true); err != nil {
			return newSuffix, err
		}
		msg := fmt.Sprintf("switch and write data into %s", local.Data.Name())
		local.Log.Info(msg)
		newSuffix = true
	}

	var err error

	dataBytes := []byte{}
	offset := uint32(0)
	if dataBytes, err = s.Marshal(); err != nil {
		return newSuffix, err
	}

	if local.mode == ZstdCompress {
		dataBytes = local.enc.EncodeAll(dataBytes, make([]byte, 0, len(dataBytes)))
	} else if local.mode == ZstdCompressWithDict {
		offset = local.next % local.chunk
		if offset == 0 {
			opts := []zstd.EOption{
				zstd.WithEncoderDictRaw(0, dataBytes),
				zstd.WithLowerEncoderMem(true),
				zstd.WithEncoderConcurrency(1),
			}
			if local.encDict, err = zstd.NewWriter(nil, opts...); err != nil {
				return newSuffix, err
			}
			dataBytes = local.enc.EncodeAll(dataBytes, make([]byte, 0, len(dataBytes)))
		} else {
			dataBytes = local.encDict.EncodeAll(dataBytes, make([]byte, 0, len(dataBytes)))
		}
	}

	if info, err := local.Data.Stat(); err != nil {
		return newSuffix, err
	} else {
		if s := info.Size(); s != local.DataOffset {
			msg := fmt.Sprintf("Data file length mismatch, expect %d, but got %d", local.DataOffset, s)
			local.Log.Error(msg)
			local.DataOffset = s
		}
	}

	idx := Index{
		TimeStamp: s.TimeStamp,
		Offset:    local.DataOffset,
		Len:       int64(len(dataBytes)),
	}
	idx.SetCompressMode(local.mode, offset)

	if _, err = local.Data.Write(dataBytes); err != nil {
		return newSuffix, err
	}

	idx.CRC = crc32.ChecksumIEEE((*[32]byte)(unsafe.Pointer(&idx))[:28])

	if _, err = local.Index.Write(idx.Marshal()); err != nil {
		return newSuffix, err
	}

	local.next++
	local.DataOffset += int64(len(dataBytes))
	return newSuffix, nil
}

type WriteOption struct {
	Interval   time.Duration
	RetainDay  int
	RetainSize int64
}

func (local *LocalStore) WriteLoop(opt WriteOption) error {

	interval := opt.Interval
	msg := fmt.Sprintf("start to collect sample every %s",
		interval.String())
	local.Log.Info(msg)
	isSkip := 0
	for {
		var shouldClose bool
		local.Lock()
		shouldClose = local.closed
		local.Unlock()
		if shouldClose == true {
			local.Close()
			return nil
		}
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
				msg := fmt.Sprintf("resume to write sample (%d skipped)", isSkip)
				local.Log.Info(msg)
				isSkip = 0
			}

			newSuffix, err := local.WriteSample(&s)
			if err != nil {
				return err
			}
			if newSuffix == true {
				// it is time to check if clean old data or not.
				// oldest first.
				local.CleanOldFiles(opt)

			}
		} else {
			if isSkip == 0 {
				msg := fmt.Sprintf("filesystem free space %s below %s, write sample skipped",
					util.GetHumanSize(statInfo.Bavail*uint64(statInfo.Bsize)),
					util.GetHumanSize(MinimumFreeSpaceForStore))
				local.Log.Warn(msg)
			}
			isSkip++
		}

		collectDuration := time.Now().Sub(start)
		if collectDuration > 500*time.Millisecond {
			msg := fmt.Sprintf("write sample take %s (larger than 500 ms)", collectDuration.String())
			local.Log.Warn(msg)
		}
		sleepDuration := time.Duration(1 * time.Second)
		if interval-collectDuration > 1*time.Second {
			sleepDuration = interval - collectDuration
		}
		time.Sleep(sleepDuration)
	}
}

func (local *LocalStore) CleanOldFiles(opt WriteOption) {
	shards, _, _, err := getIndexAndDataInfo(local.Path)
	if err != nil {
		msg := fmt.Sprintf("get index and data files: %s", err)
		local.Log.Warn(msg)
		return
	}

	unixTime := time.Now().AddDate(0, 0, -opt.RetainDay).Unix()
	oldestKeepShard := calcshard(unixTime)
	for _, shard := range shards {
		if shard < oldestKeepShard {
			local.DeleteSingleFile(shard)
		}
	}

	for {
		shards, idxSize, dataSize, err := getIndexAndDataInfo(local.Path)
		if err != nil {
			msg := fmt.Sprintf("get index and data files: %s", err)
			local.Log.Warn(msg)
			return
		}
		if idxSize+dataSize > opt.RetainSize {
			currShard := calcshard(time.Now().Unix())
			if len(shards) != 0 && shards[0] < currShard {
				s := shards[0]
				shards = shards[1:]
				local.DeleteSingleFile(s)
			} else {
				local.Log.Info("file of today was not permitted to delete")
				break
			}
		} else {
			break
		}
	}
}

func (local *LocalStore) DeleteSingleFile(shard int64) {
	for _, prefix := range []string{"index", "data"} {
		file := fmt.Sprintf("%s_%011d", prefix, shard)
		absFile := filepath.Join(local.Path, file)
		if err := os.Remove(absFile); err != nil {
			local.Log.Warn(fmt.Sprintf("%s", err))
		} else {
			msg := fmt.Sprintf("delete oldest %s", absFile)
			local.Log.Info(msg)
		}
	}
}
