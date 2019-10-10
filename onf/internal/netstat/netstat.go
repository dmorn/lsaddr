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

package netstat

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/jecoz/lsaddr/onf/internal"
)

type ActiveConnection struct {
	Proto   string
	SrcAddr net.Addr
	DstAddr net.Addr
	State   string
	Pid     int
}

// ParseOutput expects "r" to contain the output of
// a ``netstat -nao'' call. The output is splitted into lines, and
// each line that ``ParseActiveConnection'' is able to Unmarshal is
// appended to the final output.
// Returns an error only if reading from "r" produces an error
// different from ``io.EOF''.
func ParseOutput(r io.Reader) ([]ActiveConnection, error) {
	set := []ActiveConnection{}
	err := internal.ScanLines(r, func(line string) error {
		af, err := ParseActiveConnection(line)
		if err != nil {
			log.Printf("skipping netstat active connection \"%s\": %v", line, err)
			return nil
		}
		set = append(set, *af)
		return nil
	})
	return set, err
}

// ParseActiveConnection expectes "line" to be a single line output from
// ``netstat -nao'' call. The line is unmarshaled into an ``ActiveConnection''
// only if is splittable by " " into a slice of at least 4 items. "line" should
// not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
//
// "line" examples:
// "  TCP    0.0.0.0:5357           0.0.0.0:0              LISTENING       4"
// "  UDP    [::1]:62261            *:*                                    1036"
func ParseActiveConnection(line string) (*ActiveConnection, error) {
	chunks, err := internal.ChunkLine(line, " ", 4)
	if err != nil {
		return nil, err
	}

	proto := chunks[0]

	var src, dst net.Addr
	src, err = internal.ParseNetAddr(proto, chunks[1])
	dst, err = internal.ParseNetAddr(proto, chunks[2])
	if err != nil && src == dst {
		// We where not able to parse an address. We consider this
		// connection not usable.
		return nil, fmt.Errorf("unable to parse addresses: %w", err)
	}

	ac := &ActiveConnection{
		Proto:   proto,
		SrcAddr: src,
		DstAddr: dst,
	}
	hasState := len(chunks) > 4
	pidIndex := 3
	if hasState {
		pidIndex = 4
		ac.State = chunks[3]
	}
	pid, err := strconv.Atoi(chunks[pidIndex])
	if err != nil {
		return nil, fmt.Errorf("error parsing pid: %w", err)
	}
	ac.Pid = pid

	return ac, nil
}
