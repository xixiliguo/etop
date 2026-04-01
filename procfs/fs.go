package procfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"syscall"

	"github.com/xixiliguo/etop/internal/fileutil"
	"github.com/xixiliguo/etop/internal/stringutil"
	"golang.org/x/sys/unix"
)

const (
	// DefaultProcMountPoint is the common mount point of the proc filesystem.
	DefaultProcMountPoint = "/proc"

	// DefaultSysMountPoint is the common mount point of the sys filesystem.
	DefaultSysMountPoint = "/sys"

	// DefaultConfigfsMountPoint is the common mount point of the configfs.
	DefaultConfigfsMountPoint = "/sys/kernel/config"
)

const userHZ = 100

type FS struct {
	mountPoint string
	bufName    []byte
	bufData    []byte
}

func NewFS(mount string) *FS {
	fs := &FS{
		mountPoint: DefaultProcMountPoint,
		bufName:    make([]byte, 0, 64),
		bufData:    make([]byte, 0, 1024),
	}
	if mount != "" {
		fs.mountPoint = mount
	}

	return fs
}

func (fs *FS) path(file string) string {
	fs.bufName = fs.bufName[:0]
	fs.bufName = append(fs.bufName, fs.mountPoint...)

	fs.bufName = append(fs.bufName, "/"...)
	fs.bufName = append(fs.bufName, file...)

	return stringutil.ToString(fs.bufName)
}

func (fs *FS) Proc(pid int) Proc {
	return Proc{PID: pid, fs: fs}
}

func (fs *FS) EachProc(fn func(proc Proc) error) error {

	return fileutil.SubDirWalk(fs.mountPoint, func(name string) error {

		pid, err := strconv.ParseInt(name, 10, 64)
		if err != nil {
			return nil
		}
		err = fn(Proc{PID: int(pid), fs: fs})
		if err != nil {
			if _, ok := errors.AsType[syscall.Errno](err); ok {
				return nil
			}
			return err
		}
		return nil
	})
}

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR {
			return err
		}
	}
}

func (fs *FS) processFile(name string, fn func(i int, line string) error) error {

	var (
		fd  int
		err error
	)

	ignoringEINTR(func() error {
		fd, err = unix.Open(name, unix.O_RDONLY|unix.O_CLOEXEC, 0)
		return err
	})

	if err != nil {
		return fmt.Errorf("fs open %s: %+w", name, err)
	}
	defer unix.Close(fd)

	fs.bufData = fs.bufData[:0]

	for {
		var (
			n int
			e error
		)

		ignoringEINTR(func() error {
			n, e = unix.Read(fd, fs.bufData[len(fs.bufData):cap(fs.bufData)])
			return e
		})

		if e != nil {
			if e == io.EOF {
				break
			}
			return fmt.Errorf("fs read %s: %+w", name, e)
		}
		if n == 0 {
			break
		}
		fs.bufData = fs.bufData[:len(fs.bufData)+n]
		if len(fs.bufData) == cap(fs.bufData) {
			fs.bufData = append(fs.bufData, 0)[:len(fs.bufData)]
		}
	}

	i := 0
	for sub := range bytes.FieldsFuncSeq(fs.bufData, func(r rune) bool { return r == '\n' || r == '\r' }) {
		line := stringutil.ToString(sub)
		if err := fn(i, line); err != nil {
			return err
		}
		i++
	}
	return nil
}
