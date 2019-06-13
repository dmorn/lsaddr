// +build darwin

// Copyright © 2019 booster authors
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

package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"howett.net/plist"
)

// ExtractAppName is used to find the value of the "CFBundleExecutable" key.
// "r" is expected to be an ".plist" encoded file.
func ExtractAppName(r io.Reader) (string, error) {
	rs, ok := r.(io.ReadSeeker)
	if !ok {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return "", err
		}
		fmt.Printf("Buffer length: %d\n", buf.Len())
		rs = bytes.NewReader(buf.Bytes())
	}

	var data struct {
		Name string `plist:"CFBundleExecutable"`
	}
	if err := plist.NewDecoder(rs).Decode(&data); err != nil {
		return "", err
	}

	return data.Name, nil
}

type OpenFile struct {
	Command string
	Pid     string
	User    string
	Fd      string
	Type    string
	Device  string
	Node    string // contains L4 proto
	Name    string // contains src->dst addresses
	State   string // (ENSTABLISHED), (LISTEN), ...
}

// By default, ``lsof'' "name" output is of the form:
// [46][protocol][@hostname|hostaddr][:service|port]
// but we're disabling hostname conversion with the ``-n'' option
// and  port conversion with the ``-P'' option, so the output
// looks like what you expect: ``addr:port->addr:port''.
func (f *OpenFile) UnmarshalName() (net.Addr, net.Addr) {
	chunks := strings.Split(f.Name, "->")
	if len(chunks) == 0 {
		return addr{}, addr{}
	}
	src := addr{net: strings.ToLower(f.Node), addr: chunks[0]}
	if len(chunks) == 1 {
		return src, addr{}
	}

	return src, addr{net: strings.ToLower(f.Node), addr: chunks[1]}
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

func DecodeLsofOutput(r io.Reader) ([]*OpenFile, error) {
	ll := []*OpenFile{}
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return ll, nil
			}
			return ll, err
		}
		line = strings.Trim(line, "\n")
		f, err := UnmarshalLsofLine(line)
		if err != nil {
			// Skip this line
			continue
		}
		ll = append(ll, f)
	}
}

func UnmarshalLsofLine(line string) (*OpenFile, error) {
	chunks := strings.Split(line, " ")
	l := make([]string, 0, len(chunks))
	for _, v := range chunks {
		if v != "" {
			l = append(l, v)
		}
	}
	if len(l) < 9 {
		return nil, fmt.Errorf("unrecognised open file line: expected 10 items, found %d: line \"%s\"", len(l), l)
	}

	f := &OpenFile{
		Command: l[0],
		Pid:     l[1],
		User:    l[2],
		Fd:      l[3],
		Type:    l[4],
		Device:  l[5],
		Node:    l[7],
		Name:    l[8],
	}
	if len(l) >= 10 {
		f.State = l[9]
	}
	return f, nil
}
