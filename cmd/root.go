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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jecoz/lsaddr/bpf"
	"github.com/jecoz/lsaddr/csv"
	"github.com/jecoz/lsaddr/onf"
	"github.com/spf13/cobra"
)

// Build information.
var (
	Version   = "N/A"
	Commit    = "N/A"
	BuildTime = "N/A"
)

// Flags.
var (
	verbose bool
	version bool
	format  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lsaddr",
	Short: "Show a subset of all network addresses being used by your apps",
	Long:  usage,
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetPrefix("[lsaddr] ")
		if !verbose {
			log.SetOutput(ioutil.Discard)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("Version: %s, Commit: %s, Built at: %s\n\n", Version, Commit, BuildTime)
			os.Exit(0)
		}
		w := bufio.NewWriter(os.Stdout)
		enc, err := newEncoder(w, format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		set, err := onf.FetchAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		pivot := "*"
		if len(args) > 0 {
			pivot = args[0]
		}
		set, err = onf.Filter(set, pivot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to filter with %s: %v\n", pivot, err)
			os.Exit(1)
		}

		log.Printf("# of open network files: %d", len(set))
		if err := enc.Encode(set); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to encode output: %v\n", err)
			os.Exit(1)
		}
		w.Flush()
		os.Exit(0)
	},
}

type Encoder interface {
	Encode([]onf.ONF) error
}

func newEncoder(w io.Writer, format string) (Encoder, error) {
	switch strings.ToLower(format) {
	case "csv":
		return csv.NewEncoder(w), nil
	case "bpf":
		return bpf.NewEncoder(w), nil
	default:
		return nil, fmt.Errorf("unrecognised format option %s", format)
	}
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increment logger verbosity.")
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "", false, "Print build information such as version, commit and build time.")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "csv", "Choose output format.")
}

const usage = `
'lsaddr'

TODO: describe

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
