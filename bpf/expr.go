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

type Operator string

// Supported operators
const (
	AND  Operator = "and"
	OR            = "or"
	NOOP          = ""
)

type Dir string

const (
	SRC   Dir = "src"
	DST       = "dst"
	NODIR     = ""
)

// Expr represents a BPF expression. It's zero value is ready to use.
type Expr string

// Join returns a new expression, made of the conjunction of the
// caller with `r`, wrapped in an Expr.
// Callers have to ensure that operator precedence is preserved.
func (l Expr) Join(r string) Expr {
	raw := join(NOOP, string(l), r)
	return Expr(raw)
}

// And works as `Join`, but uses "and" to join the two expressions.
func (l Expr) And(r string) Expr {
	raw := join(AND, string(l), r)
	return Expr(raw)
}

// Or is the same as `And`, but with "or".
func (l Expr) Or(r string) Expr {
	raw := join(OR, string(l), r)
	return Expr(raw)
}

// Wrap surrounds `e` with ().
func (e Expr) Wrap() Expr {
	return Expr("(" + string(e) + ")")
}

// FromAddr returns a BPF from a network address, plus direction
// information. Use NODIR to make a filter that matches both src and
// dst packets.
func FromAddr(d Dir, addr net.Addr) Expr {
	if addr.String() == "" {
		return Expr("")
	}

	expr := Expr(addr.Network()) // <udp, tcp>
	addrExprRaw := string(fromAddr(addr))
	if len(addrExprRaw) == 0 {
		return expr
	}

	if d == NODIR {
		return expr.And(addrExprRaw)
	}
	return expr.And(string(d)).Join(addrExprRaw) // and <src, dst>
}

func fromAddr(addr net.Addr) Expr {
	var expr Expr
	host, port, err := net.SplitHostPort(addr.String())

	valid := func(s string) bool {
		return s != "" && s != "*"
	}

	switch {
	case err != nil:
		return expr.Join("host " + addr.String())
	case !valid(host) && !valid(port):
		return expr
	case !valid(host):
		return expr.Join("port " + port)
	case !valid(port):
		return expr.Join("host " + host)
	default:
		return expr.Join("host " + host).And("port " + port)
	}
}

// NewReader returns an io.Reader implementation, which will read
// the BPF expression from `e`. Later modifications of `e` will not
// affect the content of the reader.
func (e Expr) NewReader() *strings.Reader {
	return strings.NewReader(string(e) + "\n")
}

func join(op Operator, a, b string) string {
	// validate input
	if len(a) == 0 && b != "()" {
		return b
	}
	if len(b) == 0 || b == "()" {
		return a
	}

	switch op {
	case NOOP:
		return strings.Join([]string{a, b}, " ")
	default:
		return strings.Join([]string{a, string(op), b}, " ")
	}
}
