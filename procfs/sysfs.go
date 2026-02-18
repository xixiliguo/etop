package procfs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
)

const (

	// DefaultSysMountPoint is the common mount point of the sys filesystem.
	DefaultSysBlockMountPoint = "/sys/block"
)

type SysBlockFS struct {
	mountPoint string
	bufName    []byte
}

func NewSysBlocFS(mount string) *SysBlockFS {
	fs := &SysBlockFS{
		mountPoint: DefaultSysBlockMountPoint,
		bufName:    make([]byte, 0, 1024),
	}
	if mount != "" {
		fs.mountPoint = mount
	}

	return fs
}

func (fs *SysBlockFS) path(file string) string {
	fs.bufName = fs.bufName[:0]
	fs.bufName = append(fs.bufName, fs.mountPoint...)

	fs.bufName = append(fs.bufName, "/"...)
	fs.bufName = append(fs.bufName, file...)

	return stringutil.ToString(fs.bufName)
}

func (fs *SysBlockFS) BlockDev(name string) BlockDev {
	return BlockDev{Name: name, fs: fs}
}

func (fs *SysBlockFS) EachBlockDev(fn func(b BlockDev) error) error {

	d, err := os.Open(fs.mountPoint)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("Cannot read file: %s: %w", fs.mountPoint, err)
	}

	for _, n := range names {
		err = fn(BlockDev{Name: n, fs: fs})
		if err != nil {
			return err
		}
	}
	return nil
}

type BlockDev struct {
	Name string
	fs   *SysBlockFS
}

type BlockDevStat struct {
	Scheduler   string
	NrRequests  uint64
	ReadAheadKb uint64
	QueueNum    uint64
}

func (b BlockDev) path(file string) string {
	b.fs.bufName = b.fs.bufName[:0]
	b.fs.bufName = append(b.fs.bufName, b.fs.mountPoint...)
	b.fs.bufName = append(b.fs.bufName, "/"...)
	b.fs.bufName = append(b.fs.bufName, b.Name...)

	b.fs.bufName = append(b.fs.bufName, "/"...)
	b.fs.bufName = append(b.fs.bufName, file...)
	return stringutil.ToString(b.fs.bufName)
}

func (b BlockDev) Scheduler() (string, error) {
	path := b.path("queue/scheduler")
	res := ""
	f, err := os.Open(path)
	if err != nil {
		return res, err
	}

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		if i == 0 {
			var (
				l = strings.Index(line, "[")
				r = strings.LastIndex(line, "]")
			)

			if l < 0 || r < 0 {
				return fmt.Errorf("unexpected format, couldn't extract scheduler %q", line)
			}
			res = strings.Clone(line[l+1 : r])
		}
		return nil
	})
	return res, err
}

func (b BlockDev) NrRequests() (uint64, error) {
	path := b.path("queue/nr_requests")

	res := uint64(0)

	f, err := os.Open(path)
	if err != nil {
		return res, err
	}

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		if i == 0 {
			res, err = strconv.ParseUint(line, 10, 64)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return res, err
}

func (b BlockDev) ReadAheadKb() (uint64, error) {
	path := b.path("queue/read_ahead_kb")

	res := uint64(0)

	f, err := os.Open(path)
	if err != nil {
		return res, err
	}

	err = fileutil.ProcessFileLine(f, func(i int, line string) error {
		if i == 0 {
			res, err = strconv.ParseUint(line, 10, 64)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return res, err
}

func (b BlockDev) QueueNum() (uint64, error) {
	path := b.path("mq")

	res := uint64(0)

	if dirs, err := os.ReadDir(path); err != nil {
		return res, err
	} else {
		res = uint64(len(dirs))
	}
	return res, nil
}

func (b BlockDev) BlockDevStat() (BlockDevStat, error) {

	stat := BlockDevStat{}
	var err error
	if stat.Scheduler, err = b.Scheduler(); err != nil {
		return stat, err
	}
	if stat.NrRequests, err = b.NrRequests(); err != nil {
		return stat, err
	}
	if stat.ReadAheadKb, err = b.ReadAheadKb(); err != nil {
		return stat, err
	}
	if stat.QueueNum, err = b.QueueNum(); err != nil {
		return stat, err
	}
	return stat, err
}
