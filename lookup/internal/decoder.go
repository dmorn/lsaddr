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

// UnmarshalName unmarshals `lsof`'s name field, which by default is in the form:
// [46][protocol][@hostname|hostaddr][:service|port]
// but we're disabling hostname conversion with the ``-n'' option
// and  port conversion with the ``-P'' option, so the output
// in printed in the more decodable format: ``addr:port->addr:port''.
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
	ll := []*OpenFile{}
	err := scanLines(r, func(line string) {
		f, err := UnmarshalLsofLine(line)
		if err != nil {
			// Skip this line
			return
		}
		ll = append(ll, f)
	})
	return ll, err
}

// UnmarshalLsofLine expectes "line" to be a single line output from
// ``lsof -i -n -P'' call. The line is unmarshaled into an ``OpenFile''
// only if is splittable by " " into a slice of at least 9 items. "line" should
// not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
//
// "line" examples:
// "postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051"
// "Dropbox     614 danielmorandini  247u  IPv4 0x25c5bf09a393d583      0t0  TCP 192.168.0.61:58282->162.125.18.133:https (ESTABLISHED)"
func UnmarshalLsofLine(line string) (*OpenFile, error) {
	chunks, err := chunkLine(line, " ", 9)
	if err != nil {
		return nil, err
	}

	f := &OpenFile{
		Command: chunks[0],
		Pid:     chunks[1],
		User:    chunks[2],
		Fd:      chunks[3],
		Type:    chunks[4],
		Device:  chunks[5],
		Node:    chunks[7],
		Name:    chunks[8],
	}
	if len(chunks) >= 10 {
		f.State = chunks[9]
	}
	return f, nil
}

// Netstat

// DecodeNetstatOutput expects "r" to contain the output of
// a ``netstat -ano'' call. The output is splitted into lines, and
// each line that ``UnmarshalNetstatLine'' is able to Unmarshal is
// appended to the final output.
// As of ``DecodeLsofOutput'', this function returns an error only
// if reading from "r" produces an error different from ``io.EOF''.
func DecodeNetstatOutput(r io.Reader) ([]*OpenFile, error) {
	ll := []*OpenFile{}
	err := scanLines(r, func(line string) {
		f, err := UnmarshalNetstatLine(line)
		if err != nil {
			// Skip this line
			return
		}
		ll = append(ll, f)
	})
	return ll, err
}

// UnmarshalNetstatLine expectes "line" to be a single line output from
// ``netstat -ano'' call. The line is unmarshaled into an ``OpenFile''
// only if is splittable by " " into a slice of at least 4 items. "line" should
// not end with a "\n" delimitator, otherwise it will end up in the last
// unmarshaled item.
//
// "line" examples:
// "  TCP    0.0.0.0:5357           0.0.0.0:0              LISTENING       4"
// "  UDP    [::1]:62261            *:*                                    1036"
func UnmarshalNetstatLine(line string) (*OpenFile, error) {
	chunks, err := chunkLine(line, " ", 4)
	if err != nil {
		return nil, err
	}

	from := chunks[1]
	to := chunks[2]
	if !isValidAddress(from) {
		return nil, fmt.Errorf("unrecognised source ip address: %s", from)
	}
	if !isValidAddress(to) {
		return nil, fmt.Errorf("unrecognised destination ip address: %s", to)
	}

	f := &OpenFile{
		Node: chunks[0],
		Name: from + "->" + to,
	}
	if len(chunks) == 4 {
		// It means that the state field is missing
		f.Pid = chunks[3]
	} else {
		f.State = chunks[3]
		f.Pid = chunks[4]
	}

	return f, nil
}

// Tasklist

type Task struct {
	Pid string
	Image string
}

func (t *Task) String() string {
	return fmt.Sprintf("{Image: %s, Pid: %v}", t.Image, t.Pid)
}

func DecodeTasklistOutput(r io.Reader) ([]*Task, error) {
	ll := []*Task{}
	delim := "="
	headerTrimmed := false
	err := scanLines(r, func(line string) {
		if !headerTrimmed {
			if strings.HasPrefix(line, delim) && strings.HasSuffix(line, delim) {
				// This is the header delimiter!
				headerTrimmed = true
				return
			}
			// Still in the header
			return
		}

		t, err := UnmarshalTasklistLine(line)
		if err != nil {
			// Skip this line
			return
		}
		ll = append(ll, t)
	})
	return ll, err
}

func UnmarshalTasklistLine(line string) (*Task, error) {
	chunks, err := chunkLine(line, " ", 5)
	if err != nil {
		return nil, err
	}
	return &Task{
		Image: chunks[0],
		Pid: chunks[1],
	}, nil
}

// Plist

// ExtractAppName is used to find the value of the "CFBundleExecutable" key.
// "r" is expected to be an ".plist" encoded file.
func ExtractAppName(r io.Reader) (string, error) {
	rs, ok := r.(io.ReadSeeker)
	if !ok {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return "", err
		}
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

// Private helpers

func chunkLine(line string, sep string, min int) ([]string, error) {
	items := strings.Split(line, sep)
	chunks := make([]string, 0, len(items))
	for _, v := range items {
		if v == "" {
			continue
		}
		chunks = append(chunks, v)
	}
	n := len(chunks)
	if n < min {
		return chunks, fmt.Errorf("unable to chunk line: expected at least %d items, found %d: line \"%s\"", min, n, chunks)
	}

	return chunks, nil
}

func scanLines(r io.Reader, f func(string)) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		f(line)
	}
	return scanner.Err()
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
