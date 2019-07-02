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
	"net"

	"github.com/booster-proj/lsaddr/bpf"
	"github.com/booster-proj/lsaddr/lookup"
)

func TestHost(t *testing.T) {
	src, dst := lookup.Hosts(netFiles0) // split src and destination addresses

	validateExpr(t, bpf.NewExpr().Host(src), "host 192.168.0.61 or ::1")
	validateExpr(t, bpf.NewExpr().Host(dst), "host 52.94.218.7 or ::1")
}

func validateExpr(t *testing.T, e *bpf.Expr, expected string) {
	if e.String() != expected {
		t.Fatalf("Unexpected bpf expression: wanted \"%s\", found \"%v\"", expected, e)
	}
}

var netFiles0 = []lookup.NetFile{
	{"foo", newUDPAddr("192.168.0.61:54104"), newUDPAddr("52.94.218.7:443")},
	{"bar", newUDPAddr("[::1]:60051"), newUDPAddr("[::1]:60052")},
}

func newUDPAddr(address string) net.Addr {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}
	return addr
}
