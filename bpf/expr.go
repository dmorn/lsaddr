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

import "strings"

// Expr represents a BPF expression. It's zero value is ready to use.
type Expr string

// BPFer is a wrapper around the BPF function.
type BPFer interface {
	// Implementers should return a valid Berkeley Packet Filter
	// representation of themselves.
	BPF() string
}

// Operator restrics the
type Operator string

const (
	AND Operator = "and"
	OR           = "or"
)

// Join returns a copy of `e`, with all `ff` filters appended
// to the original expression, using `op` as separator.
// Callers have to ensure that operator precedence is preserved.
func (e Expr) Join(op Operator, f BPFer) Expr {
	raw := join(op, string(e), f.BPF())
	return Expr(raw)
}

// NewReader returns an io.Reader implementation, which will read
// the BPF expression from `e`. Later modifications of `e` will not
// affect the content of the reader.
func (e Expr) NewReader() *strings.Reader {
	return strings.NewReader(string(e) + "\n")
}

func join(op Operator, a, b string) string {
	if len(a) > 0 {
		return strings.Join([]string{a, string(op), b}, " ")
	}
	return b
}
