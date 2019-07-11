// +build darwin linux

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

package lookup

import (
	"bytes"
	"regexp"

	"github.com/booster-proj/lsaddr/lookup/internal"
	"gopkg.in/pipe.v2"
)

// openNetFiles uses ``lsof'' to find the list of open network files. It
// then filters the result using "rgx": each line that does not match is
// discarded.
func openNetFiles(rgx *regexp.Regexp) ([]*internal.OpenFile, error) {
	p := pipe.Line(
		pipe.Exec("lsof", "-i", "-n", "-P"),
		pipe.Filter(func(line []byte) bool {
			return rgx.Match(line)
		}),
	)
	output, err := pipe.Output(p)
	if err != nil {
		return []*internal.OpenFile{}, err
	}

	buf := bytes.NewBuffer(output)
	return internal.DecodeLsofOutput(buf)
}
