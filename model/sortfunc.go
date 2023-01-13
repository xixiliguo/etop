package model

var SortMap = map[string]sortFunc{
	"Pid":       SortByPid,
	"Comm":      SortByComm,
	"State":     SortByState,
	"Ppid":      SortByPpid,
	"Thr":       SortByNumThreads,
	"StartTime": SortByStartTime,
	"UserCPU":   SortByUTime,
	"SysCPU":    SortBySTime,
	"Pri":       SortByPriority,
	"Nice":      SortByNice,
	"CPU":       SortByCPUUsage,
	"Minflt":    SortByMinFlt,
	"Majflt":    SortByMajFlt,
	"Vsize":     SortByVSize,
	"RSS":       SortByRSS,
	"Mem":       SortByMemUsage,
	"Rchar":     SortByRChar,
	"Wchar":     SortByWChar,
	"Syscr":     SortBySyscR,
	"Syscw":     SortBySyscW,
	"Read":      SortByReadBytes,
	"Write":     SortByWriteBytes,
	"R/s":       SortByReadBytesPerSec,
	"W/s":       SortByWriteBytesPerSec,
	"Wcancel":   SortByCancelledWriteBytes,
	"Disk":      SortByDiskUage,
}

func SortByPid(i, j Process) bool {
	return i.Pid > j.Pid
}

func SortByComm(i, j Process) bool {
	return i.Comm > j.Comm
}

func SortByState(i, j Process) bool {
	return i.State > j.State
}

func SortByPpid(i, j Process) bool {
	return i.Ppid > j.Ppid
}

func SortByNumThreads(i, j Process) bool {
	return i.NumThreads > j.NumThreads
}

func SortByStartTime(i, j Process) bool {
	return i.StartTime > j.StartTime
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

func SortByMemUsage(i, j Process) bool {
	return i.MemUsage > j.MemUsage
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

func SortByReadBytesPerSec(i, j Process) bool {
	return i.ReadBytesPerSec > j.ReadBytesPerSec
}

func SortByWriteBytesPerSec(i, j Process) bool {
	return i.WriteBytesPerSec > j.WriteBytesPerSec
}

func SortByCancelledWriteBytes(i, j Process) bool {
	return i.CancelledWriteBytes > j.CancelledWriteBytes
}

func SortByDiskUage(i, j Process) bool {
	return i.DiskUage > j.DiskUage
}
