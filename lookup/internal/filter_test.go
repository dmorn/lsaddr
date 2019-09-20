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

func TestFilterTasks(t *testing.T) {
	t.Parallel()
	tasks := []*internal.Task{
		newTask("foo", 1),
		newTask("foo", 2),
		newTask("foo", 3),
		newTask("bar", 21),
		newTask("bar", 22),
		newTask("baz", 31),
	}

	assertPids(t, []int{1, 2, 3}, internal.FilterTasks(tasks, "foo"))
	assertPids(t, []int{21, 22}, internal.FilterTasks(tasks, "bar"))
	assertPids(t, []int{31}, internal.FilterTasks(tasks, "baz"))
	assertPids(t, []int{}, internal.FilterTasks(tasks, "invalid"))
}

// Private helpers

func newTask(image string, pid int) *internal.Task {
	return &internal.Task{
		Image: image,
		Pid:   pid,
	}
}

func assertPids(t *testing.T, exp []int, prod []*internal.Task) {
	if len(exp) != len(prod) {
		t.Fatalf("Unexpected list length. Wanted %d, found %d", len(exp), len(prod))
	}
	for i := range exp {
		if exp[i] != prod[i].Pid {
			t.Fatalf("Unexpected item in list. Wanted %v, found %v. Content: %v", exp[i], prod[i].Pid, prod)
		}
	}
}
