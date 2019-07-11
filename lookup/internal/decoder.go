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
	"fmt"
	"io"
	"net"
	"strings"
)

// Lsof section

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

func (f *OpenFile) String() string {
	return fmt.Sprintf("{Pid: %s, Proto: %s, Conn: %s}", f.Pid, f.Node, f.Name)
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

// DecodeLsofOutput expects "r" to contain the output of
// an ``lsof -i -n -P'' call. The output is splitted into each new line,
// and each line that ``UnmarshalLsofLine'' is able to Unmarshal
// is appended to the final output.
// Returns an error only if reading from "r" produces an error
// different from ``io.EOF''.
func DecodeLsofOutput(r io.Reader) ([]*OpenFile, error) {
	return scanLines(r, UnmarshalLsofLine)
}

// UnmarshalLsofLine expectes "line" to be a single line output from
// ``lsof -i -n -P'' call. The line is unmarshaled into an ``OpenFile''
// only if is splittable by " " into a slice of 9 items. "line" should
// not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
//
// "line" example:
// "postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051"
func UnmarshalLsofLine(line string) (*OpenFile, error) {
	chunks := strings.Split(line, " ")
	l := make([]string, 0, len(chunks))
	for _, v := range chunks {
		if v == "" {
			continue
		}
		l = append(l, v)
	}
	if len(l) < 9 {
		return nil, fmt.Errorf("unrecognised open file line: expected at least 9 items, found %d: line \"%s\"", len(l), l)
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

// Netstat

func DecodeNetstatOutput(r io.Reader) ([]*OpenFile, error) {
	return scanLines(r, UnmarshalNetstatLine)
}

func UnmarshalNetstatLine(line string) (*OpenFile, error) {
	chunks := strings.Split(line, " ")
	l := make([]string, 0, len(chunks))
	for _, v := range chunks {
		if v == "" {
			continue
		}
		l = append(l, v)
	}
	if len(l) < 4 {
		return nil, fmt.Errorf("unerecognised open file line: expected 4 items, found %d: line \"%s\"", len(l), l)
	}

	from := l[1]
	to := l[2]
	if !isValidAddress(from) {
		return nil, fmt.Errorf("unrecognised source ip address: %s", from)
	}
	if !isValidAddress(to) {
		return nil, fmt.Errorf("unrecognised destination ip address: %s", to)
	}

	f := &OpenFile{
		Node: l[0],
		Name: from+"->"+to,
	}
	if len(l) == 4 {
		// It means that the state field is missing
		f.Pid = l[3]
	} else {
		f.State = l[3]
		f.Pid = l[4]
	}

	return f, nil
}

// Private helpers

type lineUnmarshalerFunc func(string) (*OpenFile, error)

func scanLines(r io.Reader, f lineUnmarshalerFunc) ([]*OpenFile, error) {
	ll := []*OpenFile{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		f, err := f(line)
		if err != nil {
			// Skip this line
			continue
		}
		ll = append(ll, f)
	}
	return ll, scanner.Err()
}

func isValidAddress(s string) bool {
	_, _, err := net.SplitHostPort(s)
	return err == nil
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

