// Copyright © 2019 Jecoz
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

package onf

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

// ONF represents an open network file.
type ONF struct {
	Raw       string   // raw string that produced this result
	Cmd       string   // command associated with Pid
	Pid       int      // pid of the owner
	Src       net.Addr // source address
	Dst       net.Addr // destination address
	CreatedAt time.Time
}

func (f ONF) String() string {
	return fmt.Sprintf("{Cmd: %s, Pid: %d, Conn: %v->%v}", f.Cmd, f.Pid, f.Src, f.Dst)
}

// FetchAll retrieves the complete list of open network files. It does
// so using an external tool, `netstat` for windows and `lsof` for unix
// based systems.
func FetchAll() ([]ONF, error) {
	// fetchAll implementations may be found insiede the
	// runtime_*.go files.
	return fetchAll()
}

// Filter takes `pivot` and creates a compiled regex out of it. It then uses
// it to filter `set`, removing every open network file that do not match.
// If an error occurs, it is returned together with the original list.
func Filter(set []ONF, pivot string) ([]ONF, error) {
	if pivot == "" || pivot == "*" {
		return set, nil
	}

	rgx, err := regexp.Compile(pivot)
	if err != nil {
		return set, fmt.Errorf("unable to filter open network file set: %w", err)
	}
	acc := make([]ONF, 0, len(set))
	for _, v := range set {
		if !rgx.MatchString(v.Raw) {
			log.Printf("[DEBUG] filtering open network file: %v", v)
			continue
		}
		acc = append(acc, v)
	}
	return acc, nil
}
