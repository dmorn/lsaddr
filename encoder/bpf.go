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
	"io"
	"net"
	"strings"

	"github.com/booster-proj/lsaddr/lookup"
)

// Flags used to configure which `lookup.NetFile` fields the encoder
// will use to encode the filter.
const (
	Fdst = 1 << iota
	Fsrc
	Fport
	FstdFields = Fdst | Fsrc | Fport // initial values for standard encoder
)

// BPFEncoder is an Encoder implementation which encodes `lookup.NetFiles`
// using the BPF format.
// The "Fields" field can be used to configure how the filter is composed.
type BPFEncoder struct {
	w      io.Writer
	Fields int
}

func newBPFEncoder(w io.Writer) *BPFEncoder {
	return &BPFEncoder{
		w: w,
	}
}

// Encode encodes "l" into a bpf. A new line is added at the end.
func (e *BPFEncoder) Encode(l []lookup.NetFile) error {
	if e.Fields == 0 {
		e.Fields = FstdFields
	}

	var builder bpfBuilder
	for _, v := range l {
		if err := builder.Or(v, e.Fields); err != nil {
			return err
		}
	}
	r := strings.NewReader(builder.b.String() + "\n")
	_, err := io.Copy(e.w, r)
	return err
}

type bpfBuilder struct {
	b strings.Builder
}

func (b *bpfBuilder) Or(f lookup.NetFile, fields int) error {
	l := []string{}
	if fields&Fsrc != 0 {
		i, err := b.buildAddr(f.Src.String(), fields)
		if err != nil {
			return err
		}
		l = append(l, i)
	}
	if fields&Fdst != 0 {
		i, err := b.buildAddr(f.Dst.String(), fields)
		if err != nil {
			return err
		}
		l = append(l, i)
	}

	cur := strings.Join(l, " or ")
	prev := b.b.String()
	if prev != "" {
		cur = strings.Join([]string{prev, cur}, " or ")
	}

	b.b.Reset()
	_, err := b.b.Write([]byte(cur))
	return err
}

func (e *bpfBuilder) buildAddr(addr string, fields int) (string, error) {
	l := []string{}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	l = append(l, "host", host)
	if fields&Fport != 0 {
		l = append(l, "and", "port", port)
	}
	return strings.Join(l, " "), nil
}
