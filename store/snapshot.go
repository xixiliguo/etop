package store

import (
	"fmt"
	"os"
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
		WithSetExitProcess(local.Log),
	)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	sample := NewSample()
	if err := local.JumpSampleByTimeStamp(begin, &sample); err != nil {
		return "", err
	}

	for sample.TimeStamp <= end {
		if _, err := dest.WriteSample(&sample); err != nil {
			return "", err
		}
		sample = NewSample()
		if err := local.NextSample(1, &sample); err == ErrOutOfRange {
			break
		} else if err != nil {
			return "", err
		}
	}

	local.Close()
	dest.Close()

	tarFileName := fmt.Sprintf("snapshot_%s_%s",
		time.Unix(begin, 0).Format("200601021504"),
		time.Unix(end, 0).Format("200601021504"))

	if err := util.ArchiveToTarFile(dest.Path, tarFileName); err != nil {
		return "", err
	}
	return tarFileName, nil
}
