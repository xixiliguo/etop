package fileutil

import (
	"bytes"
	"io"
	"sync"

	"github.com/xixiliguo/etop/internal/stringutil"
)

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

var bufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1024)
		return &buf
	},
}

var bufLinePool = sync.Pool{
	New: func() any {
		buf := make([]byte, 3)
		return &buf
	},
}

func ProcessFile(r io.Reader, fn func(data string) error) error {

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)

	*bufPtr = (*bufPtr)[:0]

	for {
		n, err := r.Read((*bufPtr)[len(*bufPtr):cap(*bufPtr)])
		*bufPtr = (*bufPtr)[:len(*bufPtr)+n]
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if len(*bufPtr) == cap(*bufPtr) {
			// Add more capacity (let append pick how much).
			*bufPtr = append(*bufPtr, 0)[:len(*bufPtr)]
		}
	}
	stop := len(*bufPtr)
	for ; stop > 0; stop-- {
		if asciiSpace[(*bufPtr)[stop-1]] != 1 {
			break
		}
	}

	return fn(stringutil.ToString((*bufPtr)[:stop]))
}

type lineFunc func(i int, line string) error

func ProcessFileLine(r io.Reader, fn lineFunc) error {

	bufPtr := bufLinePool.Get().(*[]byte)
	defer bufLinePool.Put(bufPtr)

	*bufPtr = (*bufPtr)[:0]

	var (
		token   []byte
		lineIdx = 0
	)

	for {
		n, err := r.Read((*bufPtr)[len(*bufPtr):cap(*bufPtr)])
		token = (*bufPtr)[:len(*bufPtr)+n]
		if n > 0 {
			for i := bytes.IndexByte(token, '\n'); i >= 0; {
				if err := fn(lineIdx, stringutil.ToString(dropCR(token[:i]))); err != nil {
					return err
				}
				lineIdx++
				token = token[i+1:]
				i = bytes.IndexByte(token, '\n')
			}
			*bufPtr = append((*bufPtr)[:0], token...)
			if len((*bufPtr)) >= cap((*bufPtr)) {
				*bufPtr = append(*bufPtr, 0)[:len(*bufPtr)]
			}
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	if len(token) != 0 {
		if err := fn(lineIdx, stringutil.ToString(dropCR(token))); err != nil {
			return err
		}
	}
	return nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
