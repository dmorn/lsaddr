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
	"fmt"
	"log"
	"net"
	"os"
)

var Logger = log.New(os.Stderr, "[lookup] ", 0)

// NetFile contains some information obtained from a network file.
type NetFile struct {
	Command string   // command owning the file
	Src     net.Addr // source address
	Dst     net.Addr // destination address
}

func (f NetFile) String() string {
	return fmt.Sprintf("{%s %v->%v}", f.Command, f.Src, f.Dst)
}
