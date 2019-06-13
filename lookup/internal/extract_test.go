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

