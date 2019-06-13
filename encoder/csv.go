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

package encoder

import (
	"encoding/csv"
	"io"

	"github.com/booster-proj/lsaddr/lookup"
)

type csvEncoder struct {
	w *csv.Writer
}

func newCSVEncoder(w io.Writer) *csvEncoder {
	return &csvEncoder{
		w: csv.NewWriter(w),
	}
}

func (e *csvEncoder) Encode(l []lookup.NetFile) error {
	header := []string{"COMMAND", "NET", "SRC", "DST"}
	if err := e.w.Write(header); err != nil {
		return err
	}

	for _, v := range l {
		record := []string{v.Command, v.Src.Network(), v.Src.String(), v.Dst.String()}
		if err := e.w.Write(record); err != nil {
			return err
		}
	}

	e.w.Flush()
	return e.w.Error()
}
