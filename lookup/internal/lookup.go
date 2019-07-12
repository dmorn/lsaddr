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
	"io"
	"log"
	"regexp"

	"gopkg.in/pipe.v2"
)

var Logger *log.Logger

type Runtime struct {
	OFCmd pipe.Pipe // Open Files Command
	OFDecoder func(io.Reader) ([]*OpenFile, error) // Open Files Decoder
	PrepareNFExprFunc func(string) string
}

// OpenNetFiles uses ``lsof'' (or its platform dependent equivalent) to find
// the list of open network files. It then filters the result using "rgx":
// each line that does not match is discarded.
func OpenNetFiles(rgx *regexp.Regexp) ([]*OpenFile, error) {
	p := pipe.Line(
		runtime.OFCmd,
		pipe.Filter(func(line []byte) bool {
			return rgx.Match(line)
		}),
	)
	output, err := pipe.Output(p)
	if err != nil {
		return []*OpenFile{}, err
	}

	buf := bytes.NewBuffer(output)
	return runtime.OFDecoder(buf)
}
