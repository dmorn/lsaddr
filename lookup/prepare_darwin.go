// +build darwin

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
	"os"
	"path/filepath"
	"strings"

	"github.com/booster-proj/lsaddr/lookup/internal"
)

// prepareExpr returns "s" untouched if it does not end with ".app". In that case,
// "s" is used as the root directory of a macOS application. The "CFBundleExecutable"
// value of the app is searched, and used to build the an expression that will match
// each string that contains a process identifer owned by the "target" app.
func prepareExpr(s string) (string, error) {
	if !strings.HasSuffix(s, ".app") {
		return s, nil
	}

	// we suppose that "s" points to the root directory
	// of an application.
	name, err := appName(s)
	if err != nil {
		return "", err
	}
	Logger.Printf("app name: %s, path: %s", name, s)

	// Find process identifier associated with this app.
	pids := Pids(name)
	if len(pids) == 0 {
		return "", fmt.Errorf("cannot find any PID associated with %s", name)
	}

	return strings.Join(pids, "|"), nil
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
	return internal.ExtractAppName(f)
}
