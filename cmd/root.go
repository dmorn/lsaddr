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

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/booster-proj/lsaddr/encoder"
	"github.com/booster-proj/lsaddr/lookup"
	"github.com/spf13/cobra"
)

var Logger = log.New(os.Stderr, "[lsaddr] ", 0)
var debug bool
var csv bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lsaddr",
	Short: "Show a subset of all network addresses being used by the system",
	Long:  usage,
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !debug {
			Logger = log.New(ioutil.Discard, "", 0)
			lookup.Logger = log.New(ioutil.Discard, "", 0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		onf, err := lookup.OpenNetFiles(s)
		if err != nil {
			Logger.Printf("unable to find open network files for %s: %v\n", s, err)
			os.Exit(1)
		}
		Logger.Printf("# of open files: %d", len(onf))

		filtered := make([]lookup.NetFile, 0, len(onf))
		for _, v := range onf {
			if v.Dst.String() == "" {
				Logger.Printf("skipping open file: %v", v)
				continue
			}
			filtered = append(filtered, v)
		}

		if err := encoder.NewCSV(os.Stdout).Encode(filtered); err != nil {
			Logger.Printf("unable to encode open network files: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "print debug information to stderr")
	rootCmd.Flags().BoolVarP(&csv, "csv", "", true, "enable csv output encoding")
}

const usage = `
'lsaddr' takes the entire list of currently open network connections and filters it out
using the argument provided, which can either be:

- "*.app" (macOS): It will be recognised as the path leading to the root directory of
an Application. The tool will then:
	1. Extract the Application's CFBundleExecutable value
	2. Use it to find the list of Pids associated with the program
	3. Build a regular expression out of them

- a regular expression: which will be used to filter out the list of open files. Each line that
does not match against the regex will be discarded. On macOS, the list of open files is fetched
using 'lsof -i -n -P'.
Check out https://golang.org/pkg/regexp/ to learn how to properly format your regex.
`
