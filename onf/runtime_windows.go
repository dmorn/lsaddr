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

// +build windows

package onf

import (
	"time"

	"github.com/jecoz/lsaddr/netstat"
)

func fetchAll() ([]ONF, error) {
	set, err := netstat.Run()
	if err != nil {
		return []ONF{}, err
	}
	mapped := make([]ONF, len(set))
	for i, v := range set {
		mapped[i] = ONF{
			Raw:       v.Raw,
			Pid:       v.Pid,
			Src:       v.SrcAddr,
			Dst:       v.DstAddr,
			CreatedAt: time.Now(),
		}
	}
	return mapped, nil
}
