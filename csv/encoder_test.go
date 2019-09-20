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

package csv_test

import (
	"net"
	"strings"
	"testing"

	"github.com/booster-proj/lsaddr/csv"
	"github.com/booster-proj/lsaddr/lookup"
)

func TestEncode_CSV(t *testing.T) {
	t.Parallel()
	l := netFiles0
	var w strings.Builder
	if err := csv.NewEncoder(&w).Encode(l); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expOut := `PID,CMD,NET,SRC,DST
101,foo,udp,192.168.0.61:54104,52.94.218.7:443
102,,udp,[::1]:60051,[::1]:60052
`
	if expOut != w.String() {
		t.Fatalf("Unexpected output: wanted\n\"%s\",\nfound\n\"%s\"", expOut, w.String())
	}
}

var netFiles0 = []lookup.NetFile{
	{"foo", 101, newUDPAddr("192.168.0.61:54104"), newUDPAddr("52.94.218.7:443")},
	{"", 102, newUDPAddr("[::1]:60051"), newUDPAddr("[::1]:60052")},
}

func newUDPAddr(address string) net.Addr {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}
	return addr
}
