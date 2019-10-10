// Copyright © 2019 Jecoz
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package tasklist

import (
	"bytes"
	"testing"
	"reflect"
)

const tasklistExample = `
Image Name                     PID Session Name        Session#    Mem Usage
========================= ======== ================ =========== ============
System Idle Process              0 Services                   0          4 K
System                           4 Services                   0     15,376 K
smss.exe                       296 Services                   0      1,008 K
csrss.exe                      380 Services                   0      4,124 K
wininit.exe                    452 Services                   0      4,828 K
services.exe                   588 Services                   0      6,284 K
lsass.exe                      596 Services                   0     12,600 K
svchost.exe                    688 Services                   0     17,788 K
svchost.exe                    748 Services                   0      8,980 K
svchost.exe                    888 Services                   0     21,052 K
svchost.exe                    904 Services                   0     21,200 K
svchost.exe                    940 Services                   0     52,336 K
WUDFHost.exe                   464 Services                   0      6,128 K
svchost.exe                   1036 Services                   0     14,524 K
svchost.exe                   1044 Services                   0     27,488 K
svchost.exe                   1104 Services                   0     28,428 K
WUDFHost.exe                  1240 Services                   0      6,888 K
`

func TestParseOutput(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBufferString(tasklistExample)
	ll, err := ParseOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ll) != 17 {
		t.Fatalf("Unexpected ll length: wanted 17, found: %d: %v", len(ll), ll)
	}
}

func TestParseTask(t *testing.T) {
	t.Parallel()
	tt := []struct {
		line  string
		segs  []int
		image string
		pid   int
	}{
		{
			line:  "svchost.exe                    940 Services                   0     52,336 K",
			segs:  []int{25, 8, 16, 11, 12},
			image: "svchost.exe",
			pid:   940,
		},
		{
			line:  "System Idle Process              0 Services                   0          4 K",
			segs:  []int{25, 8, 16, 11, 12},
			image: "System Idle Process",
			pid:   0,
		},
	}

	for i, v := range tt {
		task, err := ParseTask(v.line, v.segs)
		if err != nil {
			t.Fatalf("%d: unexpected error: %v", i, err)
		}
		assert(t, v.image, task.Image)
		assert(t, v.pid, task.Pid)
	}
}

func assert(t *testing.T, exp, x interface{}) {
	if !reflect.DeepEqual(exp, x) {
		t.Fatalf("Assert failed: expected %v, found %v", exp, x)
	}
}

