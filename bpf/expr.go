// Copyright © 2019 booster authors 
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

//// BPFEncoder is an Encoder implementation which encodes `lookup.NetFiles`
//// using the BPF format.
//// The "Fields" field can be used to configure how the filter is composed.
//type BPFEncoder struct {
//	w io.Writer
//}
//
//func newBPFEncoder(w io.Writer) *BPFEncoder {
//	return &BPFEncoder{
//		w: w,
//	}
//}
//
//// Encode encodes "l" into a bpf. A new line is added at the end.
//func (e *BPFEncoder) Encode(l []lookup.NetFile) error {
//	// find all destinations
//	hosts := make(map[string]bool) // using a map to avoid duplicates
//	for _, v := range l {
//		host, _, err := net.SplitHostPort(v.Dst.String())
//		if err != nil {
//			return err
//		}
//
//		hosts[host] = true
//	}
//
//	hostsL := make([]string, 0, len(hosts))
//	for k := range hosts {
//		hostsL = append(hostsL, k)
//	}
//
//	var b strings.Builder
//	b.WriteString("host ")
//	b.WriteString(strings.Join(hostsL, " or "))
//	b.WriteString("\n")
//
//	r := strings.NewReader(b.String())
//	_, err := io.Copy(e.w, r)
//	return err
//}
