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

package lookup

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/booster-proj/lsaddr/lookup/internal"
	"gopkg.in/pipe.v2"
)

// OpenNetFiles uses ``lsof'' to find the network files that are currently open,
// filtering each line of ``lsof'' output using "s" as regular expression.
// Pass an empty string to return the entire output.
func OpenNetFiles(s string) ([]NetFile, error) {
	empty := []NetFile{}
	expr := s

	if strings.HasSuffix(s, ".app") {
		// we suppose that "s" points to the root directory
		// of an application.
		name, err := AppName(s)
		if err != nil {
			return empty, err
		}
		Logger.Printf("app name: %s, path: %s", name, s)

		// Find process identifier associated with this app.
		pids := Pids(name)
		if len(pids) == 0 {
			return empty, nil
		}

		// Now find the associated open files.
		expr = strings.Join(pids, "|")
	}

	rgx, err := regexp.Compile(expr)
	if err != nil {
		return empty, err
	}
	return openNetFiles(rgx)
}

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

// AppName finds the "BundeExecutable" identifier from "Info.plist" file
// contained in the "Contents" subdirectory in "path".
// "path" should point to the target app root directory.
func AppName(path string) (string, error) {
	info := filepath.Join(path, "Contents", "Info.plist")
	f, err := os.Open(info)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return internal.ExtractAppName(f)
}

// Pids returns the process identifiers associated with "app".
func Pids(app string) []string {
	p := pipe.Exec("pgrep", app)
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
