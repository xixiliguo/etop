package procfs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xixiliguo/etop/internal/stringutil"
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
}

func NewFS(mount string) *FS {
	fs := &FS{
		mountPoint: DefaultProcMountPoint,
		bufName:    make([]byte, 0, 1024),
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
		pid, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			continue
		}
		err = fn(Proc{PID: int(pid), fs: fs})
		if err != nil {
			return err
		}
	}
	return nil
}
