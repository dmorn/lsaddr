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

package internal_test

import (
	"bytes"
	"testing"

	"github.com/booster-proj/lsaddr/lookup/internal"
)

// Lsof

func TestUnmarshalLsofLine(t *testing.T) {
	t.Parallel()
	line := "Spotify   11778 danielmorandini  128u  IPv4 0x25c5bf09993eff03      0t0  TCP 192.168.0.61:51291->35.186.224.47:https (ESTABLISHED)"
	f, err := internal.UnmarshalLsofLine(line)
	if err != nil {
		t.Fatalf("Unexpcted error: %v", err)
	}

	assert(t, "Spotify", f.Command)
	assert(t, 11778, f.Pid)
	assert(t, "danielmorandini", f.User)
	assert(t, "128u", f.Fd)
	assert(t, "IPv4", f.Type)
	assert(t, "0x25c5bf09993eff03", f.Device)
	assert(t, "TCP", f.Node)
	assert(t, "192.168.0.61:51291->35.186.224.47:https", f.Name)
	assert(t, "(ESTABLISHED)", f.State)
}

const lsofExample = `Dropbox     614 danielmorandini  236u  IPv4 0x25c5bf09a4161583      0t0  TCP 192.168.0.61:58122->162.125.66.7:https (ESTABLISHED)
Dropbox     614 danielmorandini  247u  IPv4 0x25c5bf09a393d583      0t0  TCP 192.168.0.61:58282->162.125.18.133:https (ESTABLISHED)
postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051
`

func TestDecodeLsofOutput(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBufferString(lsofExample)
	ll, err := internal.DecodeLsofOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ll) != 3 {
		t.Fatalf("Unexpected ll length: wanted 3, found %d: %v", len(ll), ll)
	}
}

func TestUnmarshalName(t *testing.T) {
	t.Parallel()
	tt := []struct {
		node string
		name string
		src  string
		dst  string
		net  string
	}{
		{"TCP", "127.0.0.1:49161->127.0.01:9090", "127.0.0.1:49161", "127.0.01:9090", "tcp"},
		{"TCP", "127.0.0.1:5432", "127.0.0.1:5432", "", "tcp"},
		{"UDP", "192.168.0.61:50940->192.168.0.2:53", "192.168.0.61:50940", "192.168.0.2:53", "udp"},
		{"TCP", "[fe80:c::d5d5:601e:981b:c79d]:1024->[fe80:c::f9b9:5ecb:eeca:58e9]:1024", "[fe80:c::d5d5:601e:981b:c79d]:1024", "[fe80:c::f9b9:5ecb:eeca:58e9]:1024", "tcp"},
	}

	for i, v := range tt {
		f := internal.OpenFile{Node: v.node, Name: v.name}
		src, dst := f.UnmarshalName()
		if src.String() != v.src {
			t.Fatalf("%d: Unexpected src: wanted %s, found %s", i, v.src, src.String())
		}
		if dst.String() != v.dst {
			t.Fatalf("%d: Unexpected dst: wanted %s, found %s", i, v.dst, dst.String())
		}
		if src.Network() != v.net {
			t.Fatalf("%d: Unexpected net: wanted %s, found %s", i, v.net, src.Network())
		}
	}
}

// Netstat

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

func TestDecodeNetstatOutput(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBufferString(netstatExample)
	ll, err := internal.DecodeNetstatOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ll) != 4 {
		t.Fatalf("Unexpected ll length: wanted 4, found %d: %v", len(ll), ll)
	}
}

func TestUnmarshalNetstatLine(t *testing.T) {
	t.Parallel()
	line := "  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING       748"
	f, err := internal.UnmarshalNetstatLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert(t, "TCP", f.Node)
	assert(t, "0.0.0.0:135->0.0.0.0:0", f.Name)
	assert(t, "LISTENING", f.State)
	assert(t, 748, f.Pid)
}

// Tasklist

const tasklistExample = `
Image Name                     PID Session Name        Session#    Mem Usage
========================= ======== ================ =========== ============
System Idle Process              0 Services                   0          4 K
System                           4 Services                   0     15,376 K
smss.exe                       296 Services                   0      1,008 K
csrss.exe                      380 Services                   0      4,124 K
wininit.exe                    452 Services                   0      4,828 K
services.exe                   588 Services                   0      6,284 K
lsass.exe                      596 Services                   0     12,600 K
svchost.exe                    688 Services                   0     17,788 K
svchost.exe                    748 Services                   0      8,980 K
svchost.exe                    888 Services                   0     21,052 K
svchost.exe                    904 Services                   0     21,200 K
svchost.exe                    940 Services                   0     52,336 K
WUDFHost.exe                   464 Services                   0      6,128 K
svchost.exe                   1036 Services                   0     14,524 K
svchost.exe                   1044 Services                   0     27,488 K
svchost.exe                   1104 Services                   0     28,428 K
WUDFHost.exe                  1240 Services                   0      6,888 K
`

func TestDecodeTasklistOutput(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBufferString(tasklistExample)
	ll, err := internal.DecodeTasklistOutput(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ll) != 17 {
		t.Fatalf("Unexpected ll length: wanted 17, found: %d: %v", len(ll), ll)
	}
}

func TestUnmarshalTasklistLine(t *testing.T) {
	t.Parallel()
	tt := []struct {
		line  string
		segs  []int
		image string
		pid   int
	}{
		{
			line:  "svchost.exe                    940 Services                   0     52,336 K",
			segs:  []int{25, 8, 16, 11, 12},
			image: "svchost.exe",
			pid:   940,
		},
		{
			line:  "System Idle Process              0 Services                   0          4 K",
			segs:  []int{25, 8, 16, 11, 12},
			image: "System Idle Process",
			pid:   0,
		},
	}

	for i, v := range tt {
		task, err := internal.UnmarshalTasklistLine(v.line, v.segs)
		if err != nil {
			t.Fatalf("%d: unexpected error: %v", i, err)
		}
		assert(t, v.image, task.Image)
		assert(t, v.pid, task.Pid)
	}
}

// Plist

const infoExample = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>pico8</string>
	<key>CFBundleGetInfoString</key>
	<string>pico8</string>
	<key>CFBundleIconFile</key>
	<string>pico8.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.Lexaloffle.pico8</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>pico8</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>pico8</string>
	<key>CFBundleSignature</key>
	<string>????</string>
	<key>CFBundleVersion</key>
	<string>pico8</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.1</string>
</dict>
</plist>
`

func TestExtractAppName(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(infoExample)
	name, err := internal.ExtractAppName(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	exp := "pico8"
	if name != exp {
		t.Fatalf("Unexpected name: found %s, wanted %s", name, exp)
	}
}

// Private helpers

func assert(t *testing.T, exp, x interface{}) {
	if exp != x {
		t.Fatalf("Assert failed: expected %v, found %v", exp, x)
	}
}
