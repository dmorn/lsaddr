// Copyright © 2019 Jecoz
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

package bpf

import (
	"fmt"
	"io"

	"github.com/jecoz/lsaddr/onf"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(set []onf.ONF) error {
	var expr Expr
	for _, v := range set {
		src := string(FromAddr(NODIR, v.Src).Wrap())
		dst := string(FromAddr(NODIR, v.Dst).Wrap())
		expr = expr.Or(src).Or(dst)
	}
	if _, err := io.Copy(e.w, expr.NewReader()); err != nil {
		return fmt.Errorf("unable to encode open network files: %w", err)
	}
	return nil
}
