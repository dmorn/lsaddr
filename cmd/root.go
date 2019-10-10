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

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jecoz/lsaddr/onf"
	"github.com/spf13/cobra"
)

var debug, raw bool
var format string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lsaddr",
	Short: "Show a subset of all network addresses being used by your apps",
	Long:  usage,
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetPrefix("[lsaddr] ")
		if !debug {
			log.SetOutput(ioutil.Discard)
		}

		format = strings.ToLower(format)
		if err := validateFormat(format); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		set, err := onf.FetchAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		//fset, err := onf.Filter(set, s)
		//if err != nil {
		//	fmt.Fprintf(os.Stderr, "error: unable to filter with %s: %w\n", s, err")
		//	os.Exit(1)
		//}

		log.Printf("# of open network files: %d", len(set))
		for _, v := range set {
			fmt.Printf("%v\n", v)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "increment logger verbosity")
	rootCmd.PersistentFlags().BoolVarP(&raw, "raw", "r", false, "increment logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "csv", "choose format of output produced <csv|bpf>")
}

func validateFormat(format string) error {
	switch format {
	case "csv", "bpf":
		return nil
	default:
		return fmt.Errorf("unrecognised format option %s", format)
	}
}

const usage = `
'lsaddr' takes the entire list of open network files and filters it out
using the argument provided, which can either be:

- "*.app" (macOS): It will be recognised as the path leading to the root directory of
an Application. The tool will then:
	1. Extract the Application's CFBundleExecutable value
	2. Use it to find the list of Pids associated with the program
	3. Build a regular expression out of them

- a regular expression: which will be used to filter out the list of open files. Each line that
does not match against the regex will be discarded (e.g. "chrome.exe", "Safari", "104|405").
Check out https://golang.org/pkg/regexp/ to learn how to properly format your regex.

Using the "--format" or "-f" flag, it is possible to decide the format/encoding of the output
produced. Possible values are:
- "bpf": produces a Berkley Packet Filter expression, which, if given to a tool that supports
bpfs, will make it capture only the packets headed to/coming from the destination addresses
of the open network files collected.
- "csv": produces a CSV encoded table of the open network files collected.
`
