package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"

	"github.com/xixiliguo/etop/cgroupfs"
)

type Format int

const (
	Raw Format = iota
	HumanReadableSize
)

type Field struct {
	Name      string
	Format    Format
	Precision int
	Suffix    string
	Width     int
	FixWidth  bool
}

func appendReadableSize(dst []byte, fsize float64) []byte {
	unitMap := []string{" B", " KB", " MB", " GB", " TB", " PB"}
	i := 0
	unitsLimit := len(unitMap) - 1
	for fsize >= 1024 && i < unitsLimit {
		fsize = fsize / 1024
		i++
	}
	dst = strconv.AppendFloat(dst, fsize, 'f', 1, 64)
	dst = append(dst, unitMap[i]...)
	return dst
}

func (f Field) Render(value any) string {
	buf := make([]byte, 0, 16)

	addSuffix := true
	switch v := value.(type) {
	case uint64:
		if v == math.MaxUint64 {
			buf = append(buf, '-')
			addSuffix = false
			break
		}
		if v == cgroupfs.MaxCgroupPropertyUintValue {
			buf = append(buf, "max"...)
			addSuffix = false
			break
		}
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendUint(buf, uint64(v), 10)
		}
	case int:
		if f.Format == HumanReadableSize {
			f := float64(v)
			buf = appendReadableSize(buf, f)
		} else {
			buf = strconv.AppendInt(buf, int64(v), 10)
		}
	case float64:
		if math.IsNaN(v) {
			buf = append(buf, '-')
			addSuffix = false
			break
		}
		if f.Format == HumanReadableSize {
			buf = appendReadableSize(buf, v)
		} else {
			buf = strconv.AppendFloat(buf, v, 'f', f.Precision, 64)
		}
	case string:
		if v == cgroupfs.NoExistCgroupPropertyStrValue {
			v = "-"
		}
		if f.Suffix == "" && !f.FixWidth {
			return v
		}
		buf = append(buf, v...)
	case []string:
		if f.Suffix == "" && !f.FixWidth {
			return strings.Join(v, " ")
		}
		buf = append(buf, strings.Join(v, " ")...)
	default:
		return fmt.Sprintf("%T is unknown type", v)
	}

	if addSuffix {
		buf = append(buf, f.Suffix...)
	}

	if f.FixWidth {
		width := f.Width
		if len(f.Name) > width {
			width = len(f.Name)
		}

		if padding := width - len(buf); padding > 0 {

			cache := "                "
			if padding <= 16 {
				buf = append(buf,
					unsafe.Slice(unsafe.StringData(cache), len(cache))[:padding]...)
			} else {
				for padding != 0 {
					buf = append(buf, ' ')
					padding--
				}
			}
		}
	}
	return string(buf)
}

type FieldOpt struct {
	FixWidth bool
	Raw      bool
}

func (f Field) SetRawData() Field {
	f.Format = Raw
	f.Suffix = ""
	return f
}

func (f Field) SetFixWidth() Field {
	f.FixWidth = true
	return f
}

func (f *Field) ApplyOpt(opt FieldOpt) {
	if opt.FixWidth {
		f.FixWidth = true
	}
	if opt.Raw {
		f.Format = Raw
		f.Suffix = ""
	}
}
