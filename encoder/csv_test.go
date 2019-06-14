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

package encoder_test

import (
	"strings"
	"testing"

	"github.com/booster-proj/lsaddr/encoder"
)

func TestEncode_CSV(t *testing.T) {
	l := netFiles0
	var w strings.Builder
	if err := encoder.NewCSV(&w).Encode(l); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expOut := `COMMAND,NET,SRC,DST
foo,udp,192.168.0.61:54104,52.94.218.7:443
bar,udp,[::1]:60051,[::1]:60052
`
	if expOut != w.String() {
		t.Fatalf("Unexpected output: wanted \"%s\", found \"%s\"", expOut, w.String())
	}
}
