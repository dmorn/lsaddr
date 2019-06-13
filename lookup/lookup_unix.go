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
	"strings"

	"github.com/booster-proj/lsaddr/lookup/internal"
	"gopkg.in/pipe.v2"
)

// buildRgx compiles a regular expression out of "s". Some manipulation
// may be performed on "s" before it is compiled, depending on the hosting
// operating system: on macOS for example, if "s" ends with ".app", it
// will be trated as the root path to an application.
func buildRgx(s string) (*regexp.Regexp, error) {
	expr, err := prepareExpr(s)
	if err != nil {
		return nil, err
	}
	rgx, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return rgx, nil
}

// openNetFiles uses ``lsof'' to find the list of open network files. It
// then filters the result using "rgx": each line that does not match is
// discarded.
func openNetFiles(rgx *regexp.Regexp) ([]NetFile, error) {
	p := pipe.Line(
		pipe.Exec("lsof", "-i", "-n", "-P"),
		pipe.Filter(func(line []byte) bool {
			return rgx.Match(line)
		}),
	)
	output, err := pipe.Output(p)
	if err != nil {
		return []NetFile{}, err
	}

	buf := bytes.NewBuffer(output)
	ll, err := internal.DecodeLsofOutput(buf)
	if err != nil {
		return []NetFile{}, err
	}

	onf := make([]NetFile, len(ll))
	for i, v := range ll {
		src, dst := v.UnmarshalName()
		onf[i] = NetFile{
			Command: v.Command,
			Src:     src,
			Dst:     dst,
		}
	}
	return onf, nil
}

// Pids returns the process identifiers of "proc".
func Pids(proc string) []string {
	p := pipe.Exec("pgrep", proc)
	output, err := pipe.Output(p)
	if err != nil {
		Logger.Printf("%v", err)
		return []string{}
	}

	var builder strings.Builder
	builder.Write(output)

	trimmed := strings.Trim(builder.String(), "\n")
	return strings.Split(trimmed, "\n")
}
