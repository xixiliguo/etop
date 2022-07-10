package model

var SortMap = map[string]sortFunc{
	"PID":       SortByPID,
	"COMM":      SortByComm,
	"STATE":     SortByState,
	"PPID":      SortByPPID,
	"THR":       SortByNumThreads,
	"STARTTIME": SortByStarttime,
	"USERCPU":   SortByUTime,
	"SYSCPU":    SortBySTime,
	"PRI":       SortByPriority,
	"NICE":      SortByNice,
	"CPU":       SortByCPUUsage,
	"MINFLT":    SortByMinFlt,
	"MAJFLT":    SortByMajFlt,
	"VSIZE":     SortByVSize,
	"RSS":       SortByRSS,
	"MEM":       SortByMEMUsage,
	"RCHAR":     SortByRChar,
	"WCHAR":     SortByWChar,
	"SYSCR":     SortBySyscR,
	"SYSCW":     SortBySyscW,
	"READ":      SortByReadBytes,
	"WRITE":     SortByWriteBytes,
	"WCANCEL":   SortByCancelledWriteBytes,
	"DISK":      SortByDiskUage,
}

func SortByPID(i, j Process) bool {
	return i.PID > j.PID
}

func SortByComm(i, j Process) bool {
	return i.Comm > j.Comm
}

func SortByState(i, j Process) bool {
	return i.State > j.State
}

func SortByPPID(i, j Process) bool {
	return i.PPID > j.PPID
}

func SortByNumThreads(i, j Process) bool {
	return i.NumThreads > j.NumThreads
}

func SortByStarttime(i, j Process) bool {
	return i.Starttime > j.Starttime
}

func SortByUTime(i, j Process) bool {
	return i.UTime > j.UTime
}

func SortBySTime(i, j Process) bool {
	return i.STime > j.STime
}

func SortByPriority(i, j Process) bool {
	return i.Priority > j.Priority
}

func SortByNice(i, j Process) bool {
	return i.Nice > j.Nice
}

func SortByCPUUsage(i, j Process) bool {
	return i.CPUUsage > j.CPUUsage
}

func SortByMinFlt(i, j Process) bool {
	return i.MinFlt > j.MinFlt
}

func SortByMajFlt(i, j Process) bool {
	return i.MajFlt > j.MajFlt
}

func SortByVSize(i, j Process) bool {
	return i.VSize > j.VSize
}

func SortByRSS(i, j Process) bool {
	return i.RSS > j.RSS
}

func SortByMEMUsage(i, j Process) bool {
	return i.MEMUsage > j.MEMUsage
}

func SortByRChar(i, j Process) bool {
	return i.RChar > j.RChar
}

func SortByWChar(i, j Process) bool {
	return i.WChar > j.WChar
}

func SortBySyscR(i, j Process) bool {
	return i.SyscR > j.SyscR
}

func SortBySyscW(i, j Process) bool {
	return i.SyscW > j.SyscW
}

func SortByReadBytes(i, j Process) bool {
	return i.ReadBytes > j.ReadBytes
}

func SortByWriteBytes(i, j Process) bool {
	return i.WriteBytes > j.WriteBytes
}

func SortByCancelledWriteBytes(i, j Process) bool {
	return i.CancelledWriteBytes > j.CancelledWriteBytes
}

func SortByDiskUage(i, j Process) bool {
	return i.DiskUage > j.DiskUage
}
