// Copyright © 2019 Jecoz
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
	"log"
	"fmt"
	"strings"
	"io"
	"bytes"
	"strconv"

	"github.com/jecoz/lsaddr/onf/internal"
)

type Task struct {
	Pid   int
	Image string
}

// ParseOutput expects "r" to contain the output of
// a ``tasklist'' call. The output is splitted into lines, and
// each line that ``UnmarshakTasklistLine'' is able to Unmarshal is
// appended to the final output, with the expections of the first lines
// that come before the separator line composed by only "=". Those lines
// are considered part of the "header".
//
// As of ``ParseTask'', this function returns an error only
// if reading from "r" produces an error different from ``io.EOF''.
func ParseOutput(r io.Reader) ([]Task, error) {
	ll := []Task{}
	delim := "="
	headerTrimmed := false
	segLengths := []int{}
	err := internal.ScanLines(r, func(line string) error {
		if !headerTrimmed {
			if strings.HasPrefix(line, delim) && strings.HasSuffix(line, delim) {
				headerTrimmed = true
				// This is the header delimiter!
				chunks, err := internal.ChunkLine(line, " ", 5)
				if err != nil {
					return fmt.Errorf("unexpected header format: %w", err)
				}
				for _, v := range chunks {
					segLengths = append(segLengths, len(v))
				}
			}
			// Still in the header
			return nil
		}

		t, err := ParseTask(line, segLengths)
		if err != nil {
			log.Printf("skipping tasklist line \"%s\": %v", line, err)
			return nil
		}
		ll = append(ll, *t)
		return nil
	})
	return ll, err
}

// ParseTask expectes "line" to be a single line output from
// ``tasklist'' call. The line is unmarshaled into a ``Task'' and the operation
// is performed by readying bytes equal to "segLengths"[i], in order. "segLengths"
// should be computed using the header delimitator and counting the number of
// "=" in each segment of the header (split it by " ")
//
// "line" should not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
// The "header" lines (see below) should not be passed to this function.
//
// Example header:
// Image Name                     PID Session Name        Session#    Mem Usage
// ========================= ======== ================ =========== ============
//
// Example line:
// svchost.exe                    940 Services                   0     52,336 K
func ParseTask(line string, segLengths []int) (*Task, error) {
	buf := bytes.NewBufferString(line)
	p := make([]byte, 32)

	var image, pidRaw string
	for i, v := range segLengths[:2] {
		n, err := buf.Read(p[:v+1])
		if err != nil {
			return nil, fmt.Errorf("unable to read tasklist chunk: %w", err)
		}
		s := strings.Trim(string(p[:n]), " ")
		switch i {
		case 0:
			image = s
		case 1:
			pidRaw = s
		default:
		}
	}
	if image == "" {
		return nil, fmt.Errorf("couldn't decode image from line")
	}
	if pidRaw == "" {
		return nil, fmt.Errorf("couldn't decode pid from line")
	}
	pid, err := strconv.Atoi(pidRaw)
	if err != nil {
		return nil, fmt.Errorf("error parsing pid: %w", err)
	}

	return &Task{
		Image: image,
		Pid:   pid,
	}, nil
}
