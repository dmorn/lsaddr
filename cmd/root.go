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

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/booster-proj/lsaddr/bpf"
	"github.com/booster-proj/lsaddr/encoder"
	"github.com/booster-proj/lsaddr/lookup"
	"github.com/spf13/cobra"
)

var debug bool
var output string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lsaddr",
	Short: "Show a subset of all network addresses being used by the system",
	Long:  usage,
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetPrefix("[lsaddr] ")
		if !debug {
			log.SetOutput(ioutil.Discard)
		}

		output = strings.ToLower(output)
		if err := validateOutput(output); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		ff, err := lookup.OpenNetFiles(s)
		if err != nil {
			fmt.Printf("unable to find open network files for %s: %v\n", s, err)
			os.Exit(1)
		}
		log.Printf("# of open files: %d", len(ff))

		w := bufio.NewWriter(os.Stdout)
		if err := writeOutputTo(w, output, ff); err != nil {
			fmt.Fprintf(os.Stderr, "unable to write output: %v", err)
			os.Exit(1)
		}
		w.Flush()
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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "print debug information to stderr")
	rootCmd.PersistentFlags().StringVarP(&output, "out", "o", "csv", "select output produced")
}

func writeOutputTo(w io.Writer, output string, ff []lookup.NetFile) error {
	switch output {
	case "csv":
		return encoder.NewCSV(w).Encode(ff)
	case "bpf":
		var expr bpf.Expr
		for _, v := range ff {
			expr = expr.Join(bpf.AND, v)
		}
		_, err := io.Copy(w, expr.NewReader())
		return err
	}
	return nil
}

func validateOutput(output string) error {
	switch output {
	case "csv", "bpf":
		return nil
	default:
		return fmt.Errorf("unrecognised output %s", output)
	}
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

It is possible to configure with the '-out' flag which output 'lsaddr' will produce. Possible
values are:
- "bpf": produces a Berkley Packet Filter expression, which, if given to a tool that supports
bpfs, will make it capture only the packets headed to/coming from the destination addresses
of the open network files collected.
- "csv": produces a CSV encoded table of the open network files collected.
`
