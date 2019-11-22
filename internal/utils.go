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

package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func ChunkLine(line string, sep string, min int) ([]string, error) {
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

func ScanLines(r io.Reader, f func(string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		if err := f(line); err != nil {
			return err
		}
	}
	return scanner.Err()
}

type uncheckedAddr struct {
	net, addr string
}

func (a uncheckedAddr) Network() string { return a.net }
func (a uncheckedAddr) String() string  { return a.addr }

func ParseNetAddr(network, addr string) (net.Addr, error) {
	network = strings.ToLower(network)
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(host); ip == nil {
		return nil, fmt.Errorf("%v is not a valid ip address", host)
	}
	return uncheckedAddr{
		net:  network,
		addr: addr,
	}, nil
}
