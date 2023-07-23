package util

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var unitMap = []string{"B", "KB", "MB", "GB", "TB", "PB"}

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
	return fmt.Sprintf("%.1f %s", fsize, unitMap[i])
}

func ConvertToTime(s string) (timeStamp int64, err error) {
	if strings.HasSuffix(s, "ago") {
		d, err := time.ParseDuration(strings.TrimSpace(strings.TrimSuffix(s, "ago")))
		if err != nil {
			return timeStamp, err
		}
		return time.Now().Add(-d).Unix(), nil
	}

	t, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local)
	if err == nil {
		return t.Unix(), nil
	}

	t, err = time.ParseInLocation("01-02 15:04", s, time.Local)
	if err == nil {
		return t.AddDate(time.Now().Year(), 0, 0).Unix(), nil
	}
	t, err = time.ParseInLocation("15:04", s, time.Local)
	if err != nil {
		return timeStamp, fmt.Errorf("cannot parse %s: not support format", s)
	}
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, time.Local).Unix(), nil
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

func ArchiveToTarFile(subFiles []string, tarFileName string) error {
	tarFile, err := os.OpenFile(tarFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(tarFile)

	for _, file := range subFiles {
		base := filepath.Base(file)
		size := int64(0)
		if info, err := os.Stat(file); err != nil {
			return err
		} else {
			size = info.Size()
		}
		src, err := os.Open(file)
		if err != nil {
			return err
		}
		defer src.Close()

		hdr := &tar.Header{
			Name: base,
			Mode: 0644,
			Size: size,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		n, err := io.CopyN(tw, src, size)
		if err != nil {
			return err
		}
		if n != size {
			return fmt.Errorf("%s expect written %d bytes, but %d bytes", file, size, n)
		}
		src.Close()
		os.Remove(src.Name())
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}
