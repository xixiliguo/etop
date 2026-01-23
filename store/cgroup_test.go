package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xixiliguo/etop/cgroupfs"
)

var prefix = "└─ "

func draw(c CgroupSample) {
	indents := strings.Repeat("  ", c.Level)
	text := indents + prefix + c.Name
	fmt.Println(text)
	for _, child := range c.Child {
		draw(child)
	}
}

func TestWalkCgroup(t *testing.T) {
	isCgroup2()
	t.Logf("%+v", isCgroup2())
	cgRoot := cgroupfs.NewCgroup("/", "/")

	sample, err := walkCgroupNode(0, cgRoot)
	t.Log(sample, err)
	draw(sample)
}
