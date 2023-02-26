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
	{"Pid", "Pid"},
	{"Comm", "Comm"},
	{"State", "State"},
	{"Ppid", "Ppid"},
	{"Thr", "NumThreads"},
	{"StartTime", "StartTime"},
	{"UserCPU", "UTime"},
	{"SysCPU", "STime"},
	{"Pri", "Priority"},
	{"Nice", "Nice"},
	{"CPU", "CPUUsage"},
	{"Minflt", "MinFlt"},
	{"Majflt", "MajFlt"},
	{"Vsize", "VSize"},
	{"RSS", "RSS"},
	{"Mem", "MemUsage"},
	{"Rchar", "RChar"},
	{"Wchar", "WChar"},
	{"Rchar/s", "RCharPerSec"},
	{"Wchar/s", "WCharPerSec"},
	{"Syscr", "SyscR"},
	{"Syscw", "SyscW"},
	{"Syscr/s", "SyscRPerSec"},
	{"Syscw/s", "SyscWPerSec"},
	{"Read", "ReadBytes"},
	{"Write", "WriteBytes"},
	{"Wcancel", "CancelledWriteBytes"},
	{"R/s", "ReadBytesPerSec"},
	{"W/s", "WriteBytesPerSec"},
	{"CW/s", "CancelledWriteBytesPerSec"},
	{"Disk", "DiskUage"},
}

func main() {
	outputFile, _ := os.OpenFile("sortfunc.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer outputFile.Close()
	buff := bytes.NewBuffer(nil)
	buff.WriteString(`// Code generated by model/sort.go; DO NOT EDIT.

	package model
	
	`)
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
