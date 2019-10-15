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

package netstat

import (
	"bytes"
	"reflect"
	"testing"
)

const netstatExample = `
Active Connections

  Proto  Local Address          Foreign Address        State           PID
  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING       748
  RpcSs
 [svchost.exe]
  TCP    0.0.0.0:445            0.0.0.0:0              LISTENING       4
 Can not obtain ownership information
  TCP    0.0.0.0:5357           0.0.0.0:0              LISTENING       4
 [svchost.exe]
  UDP    [::1]:62261            *:*                                    1036
`

func TestParseOutput(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBufferString(netstatExample)
	ll, err := ParseOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ll) != 4 {
		t.Fatalf("Unexpected ll length: wanted 4, found %d: %v", len(ll), ll)
	}
}

func TestParseActiveConnection(t *testing.T) {
	t.Parallel()
	line := "  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING       748"
	ac, err := ParseActiveConnection(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert(t, "TCP", ac.Proto)
	assert(t, "0.0.0.0:135", ac.SrcAddr.String())
	assert(t, "0.0.0.0:0", ac.DstAddr.String())
	assert(t, "LISTENING", ac.State)
	assert(t, 748, ac.Pid)
}

func assert(t *testing.T, exp, x interface{}) {
	if !reflect.DeepEqual(exp, x) {
		t.Fatalf("Assert failed: expected %v, found %v", exp, x)
	}
}
