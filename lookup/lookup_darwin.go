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
	"os"
	"path/filepath"
	"strings"
	"net"
	"fmt"
	"regexp"

	"gopkg.in/pipe.v2"
	"github.com/booster-proj/lsaddr/lookup/internal"
)

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
		panic(err)
		return []string{}
	}

	var builder strings.Builder
	builder.Write(output)

	trimmed := strings.Trim(builder.String(), "\n")
	return strings.Split(trimmed, "\n")
}

// NetFile represents the decoded content of a ``lsof'' output line.
type NetFile struct {
	Command string // command owning the file
	Src net.IP // source address
	Dst net.IP // destination address
	L3Proto string // UDP, TCP, ...
}

// OpenNetFiles uses ``lsof'' to find the network files that are currently open,
// filtering each line of ``lsof'' output using "s" as regular expression.
// Pass an empty string to return the entire output.
func OpenNetFiles(s string) ([]NetFile, error) {
	rgx, err := regexp.Compile(s)
	if err != nil {
		return []NetFile{}, err
	}

	p := pipe.Line(
		pipe.Exec("lsof", "-i", "-n"),
		pipe.Filter(func(line []byte) bool {
			return rgx.Match(line)
		}),
	)
	output, err := pipe.Output(p)
	if err != nil {
		return []NetFile{}, err
	}


	fmt.Printf("DEBUG: OpenNetFiles output:\n%s\n", output)

	return []NetFile{}, fmt.Errorf("not implemented yet")
}

func wrap(word, s string) string {
	return fmt.Sprintf("%s%s%s", s, word, s)
}

