// Copyright © 2019 Marcin Wojnarowski xmarcinmarcin@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/shilangyu/driveignore/utils"
	"github.com/spf13/cobra"
)

// cleanCmd represents the upload command
var cleanCmd = &cobra.Command{
	Use:   "clean [path to clean]",
	Short: "Cleans your drive sync folder from old files",
	Long: `Will look through the drive sync folder and 
remove files that do not exist in your source files.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vPrint := utils.VPrintWrapper(verbose)

		// remove legacy files
		err := utils.Walker(args[0], func(currPath string, info os.FileInfo, relativePath string) error {
			// check if file/directory exists in source folder
			sourcePath := filepath.Join(cleanInput, relativePath)
			sourceStat, err := os.Stat(sourcePath)
			if os.IsNotExist(err) || (!os.SameFile(info, sourceStat) && !info.IsDir()) {
				os.Remove(currPath)
				vPrint("Removed:", relativePath)
			}
			return nil
		})

		return err
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("There should only be one argument")
		}
		fstat, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return errors.New("Passed path doesnt exist")
		}
		if !fstat.IsDir() {
			return errors.New("Passed path isnt a directory")
		}
		return nil
	},
}

var cleanInput string

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Local flags
	cleanCmd.Flags().StringVarP(&cleanInput, "input", "i", ".", "Input directory of source files")
}
