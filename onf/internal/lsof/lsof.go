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

package lsof

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jecoz/lsaddr/onf/internal"
	"gopkg.in/pipe.v2"
)

type OpenFile struct {
	Command string
	Pid     int
	User    string
	Fd      string
	Type    string
	Device  string
	State   string   // (ENSTABLISHED), (LISTEN), ...
	SrcAddr net.Addr // Source address
	DstAddr net.Addr // Destination address
}

func Run() ([]OpenFile, error) {
	acc := []OpenFile{}
	p := pipe.Exec("lsof", "-i", "-n", "-P")
	out, err := pipe.OutputTimeout(p, time.Millisecond*100)
	if err != nil {
		return acc, fmt.Errorf("unable to run lsof: %w", err)
	}
	buf := bytes.NewBuffer(out)
	return ParseOutput(buf)
}

// ParseOutput expects "r" to contain the output of
// an ``lsof -i -n -P'' call. The output is splitted into each new line,
// and each line that ``ParseOpenFile'' is able to parse
// is appended to the final output.
// Returns an error only if reading from "r" produces an error
// different from ``io.EOF''.
func ParseOutput(r io.Reader) ([]OpenFile, error) {
	set := []OpenFile{}
	err := internal.ScanLines(r, func(line string) error {
		onf, err := ParseOpenFile(line)
		if err != nil {
			log.Printf("skipping onf \"%s\": %v", line, err)
			return nil
		}
		set = append(set, *onf)
		return nil
	})
	return set, err
}

// ParseOpenFile expectes "line" to be a single line output from
// ``lsof -i -n -P'' call. The line is unmarshaled into an ``OpenFile''
// only if is splittable by " " into a slice of at least 9 items. "line" should
// not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
//
// "line" examples:
// "postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051"
// "Dropbox     614 danielmorandini  247u  IPv4 0x25c5bf09a393d583      0t0  TCP 192.168.0.61:58282->162.125.18.133:https (ESTABLISHED)"
func ParseOpenFile(line string) (*OpenFile, error) {
	chunks, err := internal.ChunkLine(line, " ", 9)
	if err != nil {
		return nil, err
	}
	pid, err := strconv.Atoi(chunks[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing pid: %w", err)
	}

	onf := &OpenFile{
		Command: chunks[0],
		Pid:     pid,
		User:    chunks[2],
		Fd:      chunks[3],
		Type:    chunks[4],
		Device:  chunks[5],
	}
	src, dst, err := ParseName(chunks[7], chunks[8])
	if err != nil {
		return nil, fmt.Errorf("error parsing name: %w", err)
	}
	onf.SrcAddr = src
	onf.DstAddr = dst
	if len(chunks) >= 10 {
		onf.State = chunks[9]
	}

	return onf, nil
}

// ParseName parses `lsof`'s name field, which by default is in the form:
// [46][protocol][@hostname|hostaddr][:service|port]
// but we're disabling hostname conversion with the ``-n'' option
// and port conversion with the ``-P'' option, so the output
// in printed in the more decodable format: ``addr:port->addr:port''.
func ParseName(node, name string) (net.Addr, net.Addr, error) {
	chunks := strings.Split(name, "->")
	if len(chunks) == 0 {
		return nil, nil, fmt.Errorf("unable to split name by ->")
	}
	src, err := internal.ParseNetAddr(node, chunks[0])
	if err != nil {
		return nil, nil, err
	}
	if len(chunks) == 1 {
		return src, addr{}, nil
	}
	dst, err := internal.ParseNetAddr(node, chunks[1])
	if err != nil {
		return nil, nil, err
	}
	return src, dst, nil
}

// addr is a net.Addr implementation.
type addr struct {
	addr string
	net  string
}

func (a addr) String() string {
	return a.addr
}

func (a addr) Network() string {
	return a.net
}
