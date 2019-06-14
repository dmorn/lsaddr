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
	"errors"
	"io"

	"github.com/booster-proj/lsaddr/lookup"
)

var AllowedEncoders = []string{"csv", "bpf"}

// Encoder is a wrapper around the Encode function.
type Encoder interface {
	Encode([]lookup.NetFile) error
}

// NewCSV returns an Encoder implementation which encodes
// in CSV format.
func NewCSV(w io.Writer) *CSVEncoder {
	return newCSVEncoder(w)
}

// NewBPF returns an Encoder implementation which encodes
// in Berkeley Packet Filter format.
func NewBPF(w io.Writer) *BPFEncoder {
	return newBPFEncoder(w)
}

// ValidateType returns an error if "s" does not point to
// a valid encoder. See `AllowedEncoders` to find which values
// are allowed.
func ValidateType(s string) error {
	for _, v := range AllowedEncoders {
		if s == v {
			return nil
		}
	}
	return errors.New("unsupported encoding type " + s)
}
