package util

import (
	"fmt"
	"time"
)

var unitMap = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func GetHumanSize[T int | uint64](size T) string {
	i := 0
	unitsLimit := len(unitMap) - 1
	for size >= 1024 && i < unitsLimit {
		size = size / 1024
		i++
	}
	return fmt.Sprintf("%d%s", size, unitMap[i])
}

func ConvertToTime(s string) (t time.Time, err error) {
	t, err = time.ParseInLocation("2006-01-02 15:04", s, time.Local)
	if err == nil {
		return
	}
	t, err = time.ParseInLocation("15:04", s, time.Local)
	if err != nil {
		return
	}
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, time.Local), nil
}
