package store

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

var prefix = "└─ "

func draw(c CgroupSample) {
	indents := strings.Repeat("  ", c.Level)
	text := indents + prefix + c.Name
	fmt.Println(text)
	for _, child := range c.Child {
		draw(*child)
	}
}

func TestWalkCgroup(t *testing.T) {
	isCgroup2()
	t.Logf("%+v", isCgroup2())

	root := CgroupSample{
		Path:       "/",
		Name:       "/",
		Child:      make(map[string]*CgroupSample),
		IsNotExist: make(map[CgroupFile]struct{}),
	}
	root.collectStat()
	t.Log(root)
	err := walkCgroupNode(&root, slog.Default())
	t.Log(root, err)
	draw(root)
}
