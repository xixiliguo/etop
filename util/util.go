package util

import (
	"fmt"
	"strings"
	"time"
)

var unitMap = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func GetHumanSize[T int | uint64](size T) string {
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
		return time.Now().Truncate(d).Unix(), nil
	}

	t, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local)
	if err == nil {
		return t.Unix(), nil
	}
	t, err = time.ParseInLocation("15:04", s, time.Local)
	if err != nil {
		return timeStamp, err
	}
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, time.Local).Unix(), nil
}
