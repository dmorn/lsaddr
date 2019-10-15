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
	"net"
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

// FetchAll retrieves the complete list of open network files. It does
// so using an external tool, `netstat` for windows and `lsof` for unix
// based systems.
func FetchAll() ([]ONF, error) {
	// fetchAll implementations may be found insiede the
	// runtime_*.go files.
	return fetchAll()
}

func Filter(set []ONF, pivot string) ([]ONF, error) {
	if pivot == "" || pivot == "*" {
		return set, nil
	}
	acc := make([]ONF, 0, len(set))
	return acc, fmt.Errorf("not implemented yet")
}

