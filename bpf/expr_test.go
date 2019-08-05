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

package bpf_test

import (
	"testing"

	"github.com/booster-proj/lsaddr/bpf"
)

type bpfmock string

func (f bpfmock) BPF() string {
	return string(f)
}

func TestJoin(t *testing.T) {
	tt := []struct {
		prev string
		in   string
		op   bpf.Operator
		out  string
	}{
		{
			prev: "",
			in:   "",
			op:   bpf.AND,
			out:  "",
		},
		{
			prev: "",
			in:   "(a)",
			op:   bpf.AND,
			out:  "(a)",
		},
		{
			prev: "foo",
			in:   "bar",
			op:   bpf.AND,
			out:  "foo and bar",
		},
		{
			prev: "(foo and bar)",
			in:   "baz",
			op:   bpf.OR,
			out:  "(foo and bar) or baz",
		},
	}

	for i, v := range tt {
		expr := bpf.Expr(v.prev).Join(v.op, bpfmock(v.in))
		if string(expr) != v.out {
			t.Fatalf("%d: unexpected expression: expected %v, found %v", i, v.out, expr)
		}
	}
}
