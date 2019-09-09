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

package internal_test

import (
	"testing"

	"github.com/booster-proj/lsaddr/lookup/internal"
)

func TestPidsFromTasks(t *testing.T) {
	t.Parallel()
	tasks := []*internal.Task{
		newTask("foo", "1"),
		newTask("foo", "2"),
		newTask("foo", "3"),
		newTask("bar", "21"),
		newTask("bar", "22"),
		newTask("baz", "31"),
	}

	assertList(t, []string{"1", "2", "3"}, internal.PidsFromTasks(tasks, "foo"))
	assertList(t, []string{"21", "22"}, internal.PidsFromTasks(tasks, "bar"))
	assertList(t, []string{"31"}, internal.PidsFromTasks(tasks, "baz"))
	assertList(t, []string{}, internal.PidsFromTasks(tasks, "invalid"))
}

// Private helpers

func newTask(image, pid string) *internal.Task {
	return &internal.Task{
		Image: image,
		Pid:   pid,
	}
}

func assertList(t *testing.T, exp, x []string) {
	if len(exp) != len(x) {
		t.Fatalf("Unexpected list length. Wanted %d, found %d", len(exp), len(x))
	}
	for i := range exp {
		if exp[i] != x[i] {
			t.Fatalf("Unexpected item in list. Wanted %v, found %v. Content: %v", exp[i], x[i], x)
		}
	}
}
