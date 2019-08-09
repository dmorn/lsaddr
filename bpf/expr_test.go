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
	"net/url"
	"testing"

	"github.com/booster-proj/lsaddr/bpf"
)

func TestJoin(t *testing.T) {
	tt := []struct {
		prev string
		in   string
		out  string
	}{
		{prev: "", in: "", out: ""},
		{prev: "", in: "(a)", out: "(a)"},
		{prev: "", in: "()", out: ""},
		{prev: "foo", in: "bar", out: "foo bar"},
		{prev: "(foo and bar)", in: "baz", out: "(foo and bar) baz"},
	}

	for i, v := range tt {
		expr := bpf.Expr(v.prev).Join(v.in)
		if string(expr) != v.out {
			t.Fatalf("%d: unexpected expression: expected \"%v\", found \"%v\"", i, v.out, expr)
		}
	}
}

func TestFromAddr(t *testing.T) {
	tt := []struct {
		addr string
		dir  bpf.Dir
		bpf  string
	}{
		// #0
		{
			addr: "",
			dir:  bpf.NODIR,
			bpf:  "",
		},
		{
			addr: "udp://localhost:1",
			dir:  bpf.SRC,
			bpf:  "udp and src host localhost and port 1",
		},
		// #2
		{
			addr: "tcp://left:1",
			dir:  bpf.DST,
			bpf:  "tcp and dst host left and port 1",
		},
		{
			addr: "tcp://localhost:3333",
			dir:  bpf.NODIR,
			bpf:  "tcp and host localhost and port 3333",
		},
		// #4
		{
			addr: "ip://172.31.11.33",
			dir:  bpf.NODIR,
			bpf:  "ip and host 172.31.11.33",
		},
		{
			addr: "tcp://*:57621",
			dir:  bpf.NODIR,
			bpf:  "tcp and port 57621",
		},
	}
	for i, v := range tt {
		addr := newAddr(v.addr)
		expr := bpf.FromAddr(v.dir, addr)
		if string(expr) != v.bpf {
			t.Fatalf("%d: expected \"%v\", found \"%v\"", i, v.bpf, expr)
		}
	}
}

func newAddr(s string) addr {
	url, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return addr{
		net:  url.Scheme,
		host: url.Host,
	}
}

type addr struct {
	net, host string
}

func (a addr) Network() string {
	return a.net
}
func (a addr) String() string {
	return a.host
}
