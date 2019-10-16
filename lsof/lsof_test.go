// Copyright © 2019 Jecoz
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

package lsof

import (
	"bytes"
	"testing"
)

func TestParseOpenFile(t *testing.T) {
	t.Parallel()

	line := "Spotify   11778 danielmorandini  128u  IPv4 0x25c5bf09993eff03      0t0  TCP 192.168.0.61:51291->35.186.224.47:443 (ESTABLISHED)"
	of, err := ParseOpenFile(line)
	if err != nil {
		t.Fatalf("Unexpcted error: %v", err)
	}

	assert(t, "Spotify", of.Command)
	assert(t, 11778, of.Pid)
	assert(t, "danielmorandini", of.User)
	assert(t, "128u", of.Fd)
	assert(t, "IPv4", of.Type)
	assert(t, "0x25c5bf09993eff03", of.Device)
	assert(t, "192.168.0.61:51291", of.SrcAddr.String())
	assert(t, "35.186.224.47:443", of.DstAddr.String())
	assert(t, "(ESTABLISHED)", of.State)
}

func assert(t *testing.T, exp, x interface{}) {
	if exp != x {
		t.Fatalf("Assert failed: expected %v, found %v", exp, x)
	}
}

const lsofExample = `Dropbox     614 danielmorandini  236u  IPv4 0x25c5bf09a4161583      0t0  TCP 192.168.0.61:58122->162.125.66.7:https (ESTABLISHED)
Dropbox     614 danielmorandini  247u  IPv4 0x25c5bf09a393d583      0t0  TCP 192.168.0.61:58282->162.125.18.133:https (ESTABLISHED)
postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051
`

func TestParseOutput(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBufferString(lsofExample)
	onfset, err := ParseOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(onfset) != 3 {
		t.Fatalf("Unexpected onfset length: wanted 3, found %d: %v", len(onfset), onfset)
	}
}

func TestParseName(t *testing.T) {
	t.Parallel()

	tt := []struct {
		node string
		name string
		src  string
		dst  string
		net  string
	}{
		{"TCP", "127.0.0.1:49161->127.0.0.1:9090", "127.0.0.1:49161", "127.0.0.1:9090", "tcp"},
		{"TCP", "127.0.0.1:5432", "127.0.0.1:5432", "", "tcp"},
		{"UDP", "192.168.0.61:50940->192.168.0.2:53", "192.168.0.61:50940", "192.168.0.2:53", "udp"},
		{"TCP", "[fe80:c::d5d5:601e:981b:c79d]:1024->[fe80:c::f9b9:5ecb:eeca:58e9]:1024", "[fe80:c::d5d5:601e:981b:c79d]:1024", "[fe80:c::f9b9:5ecb:eeca:58e9]:1024", "tcp"},
	}

	for _, v := range tt {
		src, dst, err := ParseName(v.node, v.name)
		if err != nil {
			t.Fatal(err)
		}
		assert(t, v.src, src.String())
		assert(t, v.dst, dst.String())
		assert(t, v.net, src.Network())
	}
}
