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
	p := pipe.Exec("pgrep", wrap(app, "\""))
	output, err := pipe.Output(p)
	if err != nil {
		return []string{}
	}

	s := make([]string, len(output))
	for i, _ := range output {
		s[i] = string(output[i])
	}
	return s
}

type NetFile struct {
	Command string // command owning the file
	Device string // ??
	Src net.IP // source address
	Dst net.IP // destination address
	L3Proto string // UDP, TCP, ...
}

func OpenNetFiles(filter []string) ([]NetFile, error) {
	grepArgs := []string{"-E", strings.Join(filter, "|")}
	p := pipe.Line(
		pipe.Exec("lsof", "-i", "-n"),
		pipe.Exec("grep", grepArgs...),
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

