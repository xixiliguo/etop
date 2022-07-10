//go:build ingore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

var RenderAndField = []struct {
	Render string
	Field  string
}{
	{"PID", "PID"},
	{"COMM", "Comm"},
	{"STATE", "State"},
	{"PPID", "PPID"},
	{"THR", "NumThreads"},
	{"STARTTIME", "Starttime"},
	{"USERCPU", "UTime"},
	{"SYSCPU", "STime"},
	{"PRI", "Priority"},
	{"NICE", "Nice"},
	{"CPU", "CPUUsage"},
	{"MINFLT", "MinFlt"},
	{"MAJFLT", "MajFlt"},
	{"VSIZE", "VSize"},
	{"RSS", "RSS"},
	{"MEM", "MEMUsage"},
	{"RCHAR", "RChar"},
	{"WCHAR", "WChar"},
	{"SYSCR", "SyscR"},
	{"SYSCW", "SyscW"},
	{"READ", "ReadBytes"},
	{"WRITE", "WriteBytes"},
	{"WCANCEL", "CancelledWriteBytes"},
	{"DISK", "DiskUage"},
}

func main() {
	outputFile, _ := os.OpenFile("sortfunc.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer outputFile.Close()
	buff := bytes.NewBuffer(nil)
	buff.WriteString("package model\n")
	buff.WriteString("var SortMap = map[string]sortFunc{\n")
	for _, info := range RenderAndField {
		c := fmt.Sprintf("	%q:  SortBy%s,\n", info.Render, info.Field)
		buff.WriteString(c)
	}
	buff.WriteString("}\n")
	for _, info := range RenderAndField {
		c := fmt.Sprintf(`
			func SortBy%s(i, j Process) bool {
				return i.%s > j.%s
			}
			`, info.Field, info.Field, info.Field)
		buff.WriteString(c)
	}
	buff.Bytes()
	formarted, _ := format.Source(buff.Bytes())
	outputFile.Write(formarted)
}
