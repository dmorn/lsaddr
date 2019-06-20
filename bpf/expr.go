// Copyright © 2019 booster authors
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

package bpf

import (
	"net"
	"strings"
)

type Expr struct {
	strings.Builder
}

func NewExpr() *Expr {
	return &Expr{}
}

func (e *Expr) Host(addrs []net.Addr) *Expr {
	seen := make(map[string]bool)
	acc := make([]string, 0, len(addrs))
	for _, v := range addrs {
		host, _, err := net.SplitHostPort(v.String())
		if err != nil {
			continue
		}

		if _, ok := seen[host]; !ok {
			acc = append(acc, host)
			seen[host] = true
		}
	}

	e.WriteString("host ")
	e.WriteString(strings.Join(acc, " or "))
	return e
}

func (e *Expr) NewReader() *strings.Reader {
	return strings.NewReader(e.String() + "\n")
}

func (e *Expr) WriteString(s string) (int, error) {
	if len(e.String()) > 0 {
		s = " " + s
	}
	return e.Builder.WriteString(s)
}
