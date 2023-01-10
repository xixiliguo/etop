package store

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/prometheus/procfs"
)

func TestLocalStoreopenFile(t *testing.T) {
	dir := t.TempDir()
	local, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
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
		WithSetDefault(dir, log.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
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
		WithSetDefault(dir, log.Default()),
		WithWriteOnly(),
	)
	if err != nil {
		t.Fatalf("new writeStore: %s\n", err)
	}
	defer writeStore.Close()

	readStore, err := NewLocalStore(
		WithSetDefault(dir, log.Default()),
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
