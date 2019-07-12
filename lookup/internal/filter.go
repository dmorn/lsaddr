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

package internal

import (
	"bytes"
	"os"
	"strings"
	"path/filepath"
	"regexp"

	"gopkg.in/pipe.v2"
)

// BuildNFFilter compiles a regular expression out of "s". Some manipulation
// may be performed on "s" before it is compiled, depending on the hosting
// operating system: on macOS for example, if "s" ends with ".app", it
// will be trated as the root path to an application.
func BuildNFFilter(s string) (*regexp.Regexp, error) {
	expr := runtime.PrepareNFExprFunc(s)
	rgx, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return rgx, nil
}

func PidsFromTasks(tasks []*Task, image string) []string {
	pids := []string{}
	for _, v := range tasks {
		if v.Image != image {
			continue
		}
		pids = append(pids, v.Pid)
	}
	return pids
}

// Private helpers

// Darwin helpers

// prepareExprDarwin returns "s" untouched if it does not end with ".app". In that case,
// "s" is used as the root directory of a macOS application. The "CFBundleExecutable"
// value of the app is searched, and used to build the an expression that will match
// each string that contains a process identifer owned by the "target" app.
func prepareNFExprDarwin(s string) string {
	if _, err := os.Stat(s); err != nil {
		// this is not a path
		return s
	}
	path := strings.TrimRight(s, "/")
	if !strings.HasSuffix(path, ".app") {
		// it is a path, but not one that we know how to handle.
		return s
	}

	// we suppose that "s" points to the root directory
	// of an application.
	name, err := appName(path)
	if err != nil {
		Logger.Printf("unable to find app name: %v", err)
		return s
	}
	Logger.Printf("app name: %s, path: %s", name, path)

	// Find process identifier associated with this app.
	pids := pids(name)
	if len(pids) == 0 {
		Logger.Printf("cannot find any PID associated with %s", name)
		return s
	}

	return strings.Join(pids, "|")
}

// appName finds the "BundeExecutable" identifier from "Info.plist" file
// contained in the "Contents" subdirectory in "path".
// "path" should point to the target app root directory.
func appName(path string) (string, error) {
	info := filepath.Join(path, "Contents", "Info.plist")
	f, err := os.Open(info)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return ExtractAppName(f)
}

// pids returns the process identifiers of "proc".
func pids(proc string) []string {
	p := pipe.Exec("pgrep", proc)
	output, err := pipe.Output(p)
	if err != nil {
		Logger.Printf("unable to find pids with pgrep: %v", err)
		return []string{}
	}

	var builder strings.Builder
	builder.Write(output)

	trimmed := strings.Trim(builder.String(), "\n")
	return strings.Split(trimmed, "\n")
}

// Windows helpers

func prepareNFExprWin(s string) string {
	if !strings.HasSuffix(s, ".exe") {
		// we're not able to use something that is not
		// an executable name.
		return s
	}

	// TODO: what if "s" is a path? Only its last component
	// is required.

	tasks := tasks(s)
	if len(tasks) == 0 {
		Logger.Printf("cannot find any task associated with %s", s)
		return s
	}
	pids := PidsFromTasks(tasks, s)
	if len(pids) == 0 {
		Logger.Printf("cannot find any PID associated with %s", s)
		return s
	}

	return strings.Join(pids, "|")
}

// tasks executes the ``tasklist'' command, which is only
// available on windows.
func tasks(image string) []*Task {
	empty := []*Task{}
	p := pipe.Exec("tasklist")
	output, err := pipe.Output(p)
	if err != nil {
		Logger.Printf("unable to execute tasklist: %v", err)
		return empty
	}

	r := bytes.NewReader(output)
	tasks, err := DecodeTasklistOutput(r)
	if err != nil {
		Logger.Printf("unable to decode tasklist's output: %v", err)
		return empty
	}
	return tasks
}
