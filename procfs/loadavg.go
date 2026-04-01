// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package procfs

import (
	"fmt"
	"strconv"

	"github.com/xixiliguo/etop/internal/stringutil"
)

// LoadAvg represents an entry in /proc/loadavg.
type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func (fs FS) Load() (LoadAvg, error) {

	load := LoadAvg{}

	path := fs.path("loadavg")

	err := fs.processFile(path, func(i int, line string) error {
		var fields [5]string

		nFields := stringutil.FieldsN(line, fields[:])

		if nFields < 3 {
			return fmt.Errorf("unexpected line in loadavg: '%s'", line)
		}

		load.Load1, _ = strconv.ParseFloat(fields[0], 64)
		load.Load5, _ = strconv.ParseFloat(fields[1], 64)
		load.Load15, _ = strconv.ParseFloat(fields[2], 64)
		return nil
	})

	return load, err
}
