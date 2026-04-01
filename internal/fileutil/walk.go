package fileutil

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/xixiliguo/etop/internal/stringutil"
	"golang.org/x/sys/unix"
)

func readInt(b []byte, off, size uintptr) (u uint64, ok bool) {
	if len(b) < int(off+size) {
		return 0, false
	}

	switch size {
	case 1:
		return uint64(b[off]), true
	case 2:
		return uint64(binary.NativeEndian.Uint16(b[off:])), true
	case 4:
		return uint64(binary.NativeEndian.Uint32(b[off:])), true
	case 8:
		return uint64(binary.NativeEndian.Uint64(b[off:])), true
	default:
		panic("syscall: readInt with unsupported size")
	}
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(syscall.Dirent{}.Ino), unsafe.Sizeof(syscall.Dirent{}.Ino))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(syscall.Dirent{}.Reclen), unsafe.Sizeof(syscall.Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	reclen, ok := direntReclen(buf)
	if !ok {
		return 0, false
	}
	return reclen - uint64(unsafe.Offsetof(syscall.Dirent{}.Name)), true
}

func direntType(buf []byte) os.FileMode {
	off := unsafe.Offsetof(syscall.Dirent{}.Type)
	if off >= uintptr(len(buf)) {
		return ^os.FileMode(0) // unknown
	}
	typ := buf[off]
	switch typ {
	case syscall.DT_BLK:
		return os.ModeDevice
	case syscall.DT_CHR:
		return os.ModeDevice | os.ModeCharDevice
	case syscall.DT_DIR:
		return os.ModeDir
	case syscall.DT_FIFO:
		return os.ModeNamedPipe
	case syscall.DT_LNK:
		return os.ModeSymlink
	case syscall.DT_REG:
		return 0
	case syscall.DT_SOCK:
		return os.ModeSocket
	}
	return ^os.FileMode(0) // unknown
}

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR {
			return err
		}
	}
}

var walkBufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 8192)
		return &buf
	},
}

func SubDirWalk(path string, walkFunc func(subDir string) error) error {

	var (
		fd  int
		err error
	)

	ignoringEINTR(func() error {
		fd, err = unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
		return err
	})

	if err != nil {
		return fmt.Errorf("open %s: %+w", path, err)
	}
	defer unix.Close(fd)

	bufPtr := walkBufPool.Get().(*[]byte)
	defer walkBufPool.Put(bufPtr)

	bufp := 0
	nbuf := 0
	for {
		if bufp >= nbuf {
			bufp = 0
			var errno error

			ignoringEINTR(func() error {
				nbuf, errno = unix.ReadDirent(fd, *bufPtr)
				return errno
			})

			if errno != nil {
				return fmt.Errorf("readDirent %s: %+w", path, err)
			}
			if nbuf <= 0 {
				return nil // EOF
			}
		}

		buf := (*bufPtr)[bufp:nbuf]
		reclen, ok := direntReclen(buf)
		if !ok || reclen > uint64(len(buf)) {
			break
		}
		rec := buf[:reclen]
		bufp += int(reclen)

		const namoff = uint64(unsafe.Offsetof(syscall.Dirent{}.Name))
		namlen, ok := direntNamlen(rec)
		if !ok || namoff+namlen > uint64(len(rec)) {
			break
		}
		name := rec[namoff : namoff+namlen]
		for i, c := range name {
			if c == 0 {
				name = name[:i]
				break
			}
		}

		recName := stringutil.ToString(name)
		// Check for useless names
		if recName == "." || recName == ".." {
			continue
		}

		recType := direntType(rec)
		if recType.IsDir() {
			if err := walkFunc(recName); err != nil {
				return err
			}
		}
	}
	return nil
}
