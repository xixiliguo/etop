package util

import "fmt"

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
