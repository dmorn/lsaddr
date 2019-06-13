// +build darwin

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

func TestUnmarshalLsofLine(t *testing.T) {
	line := "Spotify   11778 danielmorandini  128u  IPv4 0x25c5bf09993eff03      0t0  TCP 192.168.0.61:51291->35.186.224.47:https (ESTABLISHED)"
	f, err := internal.UnmarshalLsofLine(line)
	if err != nil {
		t.Fatalf("Unexpcted error: %v", err)
	}

	if f.Command != "Spotify" {
		t.Fatalf("Unexpected %v", f.Command)
	}
	if f.Pid != "11778" {
		t.Fatalf("Unexpected %v", f.Pid)
	}
	if f.User != "danielmorandini" {
		t.Fatalf("Unexpected %v", f.User)
	}
	if f.Fd != "128u" {
		t.Fatalf("Unexpected %v", f.Fd)
	}
	if f.Type != "IPv4" {
		t.Fatalf("Unexpected %v", f.Type)
	}
	if f.Device != "0x25c5bf09993eff03" {
		t.Fatalf("Unexpected %v", f.Device)
	}
	if f.Node != "TCP" {
		t.Fatalf("Unexpected %v", f.Node)
	}
	if f.Name != "192.168.0.61:51291->35.186.224.47:https" {
		t.Fatalf("Unexpected %v", f.Name)
	}
	if f.State != "(ESTABLISHED)" {
		t.Fatalf("Unexpected %v", f.State)
	}
}

const lsofExample = `Dropbox     614 danielmorandini  236u  IPv4 0x25c5bf09a4161583      0t0  TCP 192.168.0.61:58122->162.125.66.7:https (ESTABLISHED)
Dropbox     614 danielmorandini  247u  IPv4 0x25c5bf09a393d583      0t0  TCP 192.168.0.61:58282->162.125.18.133:https (ESTABLISHED)
postgres    676 danielmorandini   10u  IPv6 0x25c5bf0997ca88e3      0t0  UDP [::1]:60051->[::1]:60051
`

func TestDecodeLsofOutput(t *testing.T) {
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
