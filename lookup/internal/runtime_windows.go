// +build windows

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
	"strings"

	"gopkg.in/pipe.v2"
)

var runtime = Runtime{
	LsofCmd:     pipe.Exec("netstat", "-ano"),
	LsofDecoder: DecodeNetstatOutput,
}

func prepareNFExpr(s string) string {
	if !strings.HasSuffix(s, ".exe") {
		// we're not able to use something that is not
		// an executable name.
		return s
	}
	// TODO: what if "s" is a path? Only its last component
	// is required.

	pids := pids(s)
	if len(pids) == 0 {
		Logger.Printf("cannot find any PID associated with %s", s)
		return s
	}

	return strings.Join(pids, "|")
}

func pids(image string) []string {
	p := pipe.Exec("tasklist")
	output, err := pipe.Output(p)
	if err != nil {
		Logger.Printf("unable to execute tasklist: %v", err)
		return []string{}
	}

	r := bytes.NewReader(output)
	tasks, err := DecodeTasklistOutput(r)
	if err != nil {
		Logger.Printf("unable to decode tasklist's output: %v", err)
		return []string{}
	}

	for _, v := range tasks {
		if v.Image == image {
			return v.Pids
		}
	}
	return []string{}
}
