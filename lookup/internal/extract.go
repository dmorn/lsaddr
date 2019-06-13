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

package internal

import (
	"bytes"
	"fmt"
	"io"

	"howett.net/plist"
)

// ExtractAppName is used to find the value of the "CFBundleExecutable" key.
// "r" is expected to be an ".plist" encoded file.
func ExtractAppName(r io.Reader) (string, error) {
	rs, ok := r.(io.ReadSeeker)
	if !ok {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return "", err
		}
		fmt.Printf("Buffer length: %d\n", buf.Len())
		rs = bytes.NewReader(buf.Bytes())
	}

	var data struct {
		Name string `plist:"CFBundleExecutable"`
	}
	if err := plist.NewDecoder(rs).Decode(&data); err != nil {
		return "", err
	}

	return data.Name, nil
}

