package store

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs"
)

func TestLocalStoreopenFile(t *testing.T) {
	dir := t.TempDir()
	local, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	err = local.openFile("20221127", false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "index_20221127"))

	err = local.openFile("20221127", false)
	if err == nil || errors.Is(err, os.ErrNotExist) == false {
		t.Fatalf("open no-exist file: %s\n", err)
	}
	_, _ = os.Create(filepath.Join(dir, "data_20221127"))
	err = local.openFile("20221127", false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile("20221127", false)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}

	err = local.changeFile("20221127", true)
	if err != nil {
		t.Fatalf("open exist file: %s\n", err)
	}
	if filepath.Base(local.Index.Name()) != "index_20221127" {
		t.Fatalf("index file got: %s, but want %s\n", filepath.Base(local.Index.Name()), "index_20221127")
	}

	if filepath.Base(local.Data.Name()) != "data_20221127" {
		t.Fatalf("data file got: %s, but want %s\n", filepath.Base(local.Data.Name()), "data_20221127")
	}
}

func TestNextSample(t *testing.T) {
	dir := t.TempDir()
	writeStore, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
	)
	if err != nil {
		t.Fatalf("now readStore: %s\n", err)
	}
	defer readStore.Close()

	s := NewSample()

	d := NewSample()
	noExist := NewSample()

	if err := readStore.NextSample(0, &d); err != ErrOutOfRange {
		t.Fatalf("read sample: %s\n", err)
	}

	if err := writeStore.CollectSample(&s); err != nil {
		t.Fatalf("collect sample: %s\n", err)
	}
	if _, err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}

	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.NextSample(0, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.NextSample(1, &noExist); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got no error: %s\n", err)
	}
	if err := readStore.NextSample(-1, &noExist); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got no error: %s\n", err)
	}

	if err := readStore.NextSample(0, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	s.TimeStamp += 1
	if _, err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write new sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.NextSample(1, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}
}

func TestJumpSampleByTimeStamp(t *testing.T) {
	dir := t.TempDir()
	writeStore, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
	)
	if err != nil {
		t.Fatalf("now readStore: %s\n", err)
	}
	defer readStore.Close()

	d := NewSample()
	c := NewSample()

	if err := readStore.JumpSampleByTimeStamp(123, &c); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	// write 1st sample
	s := NewSample()
	if err := writeStore.CollectSample(&s); err != nil {
		t.Fatalf("collect sample: %s\n", err)
	}
	if _, err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	// shoud ignore 1st sample, because no sample can compare with it
	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp-1, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp+1, &d); err != ErrOutOfRange {
		t.Fatalf("read sample should fail, but got: %s\n", err)
	}

	// write 2nd sample
	s.TimeStamp += 1
	if _, err := writeStore.WriteSample(&s); err != nil {
		t.Fatalf("write sample: %s\n", err)
	}
	writeStore.Index.Sync()
	writeStore.Data.Sync()

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(Sample{}, procfs.ProcStat{}),
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp-1, &d); err != nil {
		t.Fatalf("read sample: %s\n", err)
	}

	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, d, opts...))
	}

	if err := readStore.JumpSampleByTimeStamp(s.TimeStamp+1, &c); err != nil {
		t.Logf("%d %+v\n", s.TimeStamp, readStore.idxs)
		t.Fatalf("read sample: %s\n", err)
	}
	if cmp.Equal(s, d, opts...) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(s, c, opts...))
	}
}

func getDirAndFilesName(path string) ([]string, []string) {
	entry, _ := os.ReadDir(path)
	dirs := []string{}
	files := []string{}
	for _, e := range entry {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		} else {
			files = append(files, e.Name())
		}
	}
	sort.Strings(dirs)
	sort.Strings(files)
	return dirs, files
}

func TestCleanOldFilesByDays(t *testing.T) {

	dir := t.TempDir()
	local, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	expect := []string{}

	now := time.Now()
	for i := 0; i < 5; i++ {
		suffix := now.AddDate(0, 0, -i).Format("20060102")
		if i < 3 {
			expect = append(expect, "data_"+suffix)
			expect = append(expect, "index_"+suffix)
		}
		local.changeFile(suffix, true)
	}
	sort.Strings(expect)

	local.CleanOldFiles(WriteOption{
		RetainDay:  2,
		RetainSize: 9999,
	})
	_, reFiles := getDirAndFilesName(local.Path)
	if cmp.Equal(expect, reFiles) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(expect, reFiles))
	}

}

func TestCleanOldFilesBySize(t *testing.T) {

	dir := t.TempDir()
	local, err := NewLocalStore(
		WithSetDefault(dir, slog.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("NewLocalStore: %s\n", err)

	}

	expect := []string{}

	now := time.Now()
	for i := 0; i < 5; i++ {
		suffix := now.AddDate(0, 0, -i).Format("20060102")
		if i == 0 {
			expect = append(expect, "data_"+suffix)
			expect = append(expect, "index_"+suffix)
		}
		idx, err := os.Create(path.Join(local.Path, "index_"+suffix))
		if err != nil {
			t.Fatal(err)
		}

		if err := idx.Truncate(1 << 20); err != nil {
			t.Fatal(err)
		}
		data, err := os.Create(path.Join(local.Path, "data_"+suffix))
		if err != nil {
			t.Fatal(err)
		}

		if err := data.Truncate(5 << 20); err != nil {
			t.Fatal(err)
		}
	}
	sort.Strings(expect)

	local.CleanOldFiles(WriteOption{
		RetainDay:  9999,
		RetainSize: 6 << 20,
	})
	_, reFiles := getDirAndFilesName(local.Path)
	if cmp.Equal(expect, reFiles) == false {
		t.Fatalf("data should be the same\n%s\n", cmp.Diff(expect, reFiles))
	}

}
