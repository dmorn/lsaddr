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
	"os"

	"github.com/spf13/cobra"
	"github.com/booster-proj/lsaddr/lookup"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lsaddr",
	Short: "Outputs IP addresses used by an application",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		fmt.Printf("App path: %v\n", path)
		name, err := lookup.AppName(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable find app name: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Name found: %v\n", name)

		// Find process identifier associated with this app.
		pids := lookup.Pids(name)

		onf, err := lookup.OpenNetFiles(pids)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to find open network files for %s: %v\n", name, err)
			os.Exit(1)
		}

		fmt.Printf("Network Files: %v\n", onf)
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

