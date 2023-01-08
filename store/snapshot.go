package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xixiliguo/etop/util"
)

func (local *LocalStore) Snapshot(begin int64, end int64) (string, error) {

	tempPath, err := os.MkdirTemp("", "etop-snapshot")
	if err != nil {
		return "", err
	}

	dest, err := NewLocalStore(
		WithSetDefault(tempPath, local.Log),
		WithWriteOnly(),
	)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	sample := NewSample()
	if err := local.JumpSampleByTimeStamp(begin, &sample); err != nil {
		return "", err
	}
	files := []string{}

	for sample.TimeStamp <= end {
		if newSuffix, err := dest.WriteSample(&sample); err != nil {
			return "", err
		} else if newSuffix == true {
			files = append(files,
				filepath.Join(dest.Path, "index_"+dest.suffix),
				filepath.Join(dest.Path, "data_"+dest.suffix),
			)
		}
		sample = NewSample()
		if err := local.NextSample(1, &sample); err == ErrOutOfRange {
			break
		} else if err != nil {
			return "", err
		}
	}

	tarFileName := fmt.Sprintf("snapshot_%s_%s",
		time.Unix(begin, 0).Format("200601021504"),
		time.Unix(end, 0).Format("200601021504"))

	if err := util.ArchiveToTarFile(files, tarFileName); err != nil {
		return "", err
	}
	return tarFileName, nil
}
