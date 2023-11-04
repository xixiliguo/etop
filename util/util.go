package util

import (
	"archive/tar"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var unitMap = []string{" B", " KB", " MB", " GB", " TB", " PB"}

func GetHumanSize[T int | int64 | uint64 | uint | uint32 | float64](size T) string {
	if size < 0 {
		return "-1 B"
	}
	fsize := float64(size)
	i := 0
	unitsLimit := len(unitMap) - 1
	for fsize >= 1024 && i < unitsLimit {
		fsize = fsize / 1024
		i++
	}
	return strconv.FormatFloat(fsize, 'f', 1, 64) + unitMap[i]
}

func ConvertToUnixTime(s string) (timeStamp int64, err error) {

	if d, err := time.ParseDuration(strings.TrimSpace(s)); err == nil {
		return time.Now().Add(-d).Unix(), nil
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		return t.Unix(), nil
	}

	if t, err := time.ParseInLocation("01-02 15:04", s, time.Local); err == nil {
		y := time.Now().Year()
		return time.Date(y, t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local).Unix(), nil
	}

	if t, err := time.ParseInLocation("15:04", s, time.Local); err == nil {
		y, m, d := time.Now().Date()
		return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, time.Local).Unix(), nil
	}
	return timeStamp, fmt.Errorf("cannot parse %s: not support format", s)
}

func ExtractFileFromTar(tarFileName string) (string, error) {
	f, err := os.Open(tarFileName)
	if err != nil {
		return "", err
	}
	tempPath, err := os.MkdirTemp("", "etop-snapshot")
	if err != nil {
		return "", err
	}
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", err
		}

		subFile, err := os.OpenFile(filepath.Join(tempPath, hdr.Name),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(subFile, tr); err != nil {
			return "", err
		}
		subFile.Close()
	}
	return tempPath, nil
}

func ArchiveToTarFile(path string, tarFileName string) error {
	tarFile, err := os.OpenFile(tarFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(tarFile)
	if entrys, err := os.ReadDir(path); err != nil {
		return err
	} else {
		for _, e := range entrys {
			fullPath := filepath.Join(path, e.Name())

			size := int64(0)
			if info, err := os.Stat(fullPath); err != nil {
				return err
			} else {
				size = info.Size()
			}
			src, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer src.Close()
			now := time.Now()
			hdr := &tar.Header{
				Name:       e.Name(),
				Mode:       0644,
				Size:       size,
				ModTime:    now,
				AccessTime: now,
				ChangeTime: now,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			n, err := io.CopyN(tw, src, size)
			if err != nil {
				return err
			}
			if n != size {
				return fmt.Errorf("%s expect written %d bytes, but %d bytes", e.Name(), size, n)
			}
			src.Close()
			os.Remove(src.Name())
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}

func CreateLogger(w io.Writer, onlyMsg bool) *slog.Logger {
	hOpt := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}
	if onlyMsg {
		hOpt = &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}
	}
	th := slog.NewTextHandler(w, hOpt)
	return slog.New(th)
}
