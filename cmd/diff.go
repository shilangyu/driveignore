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

	"github.com/fatih/color"
	"github.com/shilangyu/driveignore/utils"
	"github.com/spf13/cobra"
)

// diffCmd represents the upload command
var diffCmd = &cobra.Command{
	Use:   "diff [drive sync folder path]",
	Short: "Compares your directory with the drive one",
	Long: `Prints out the difference in files between your source (input) and
drive sync folder ([drive sync folder path])

Red    - your drive sync folder is missing a file
Yellow - your drive sync folder has a file that doesnt exist in input
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vPrint := utils.VPrintWrapper(verbose)
		redPrint := color.New(color.FgRed).PrintlnFunc()
		yellowPrint := color.New(color.FgHiYellow).PrintlnFunc()

		driveignore, driveignoreType := utils.DriveIgnore(diffInput, diffMergeIgnores)

		switch driveignoreType {
		case utils.GlobalIgnore:
			vPrint("loaded global .driveignore")
		case utils.LocalIgnore:
			vPrint("loaded local .driveignore")
		case utils.MergedIgnore:
			vPrint("loaded merged global and local .driveignore")
		case utils.NoIgnore:
			return errors.New("No local nor global .driveignores found")
		}

		// search for missing files
		filepath.Walk(diffInput, func(currPath string, info os.FileInfo, err error) error {
			relativePath, _ := filepath.Rel(diffInput, currPath)
			if err != nil {
				panic(err)
			}

			// skip the folder itself
			if currPath == "." {
				return nil
			}

			// adding slash to directories for print clarity
			if temp, _ := os.Stat(currPath); temp.IsDir() {
				relativePath += "\\"
			}

			// ignore .driveignore files/dirs
			if info.IsDir() && driveignore.Match(currPath, true) {
				return filepath.SkipDir
			} else if !info.IsDir() && driveignore.Match(currPath, false) {
				return nil
			}

			// check if file/directory exists in drive sync folder
			goalPath := filepath.Join(args[0], relativePath)
			goalStat, err := os.Stat(goalPath)
			if os.IsNotExist(err) || (!os.SameFile(info, goalStat) && !info.IsDir()) {
				redPrint(relativePath)
			}
			return nil
		})

		// search for legacy files/directories
		filepath.Walk(args[0], func(currPath string, info os.FileInfo, err error) error {
			relativePath, _ := filepath.Rel(args[0], currPath)
			if err != nil {
				panic(err)
			}

			// adding slash to directories for print clarity
			if temp, _ := os.Stat(currPath); temp.IsDir() {
				relativePath += "\\"
			}

			// skip the folder itself
			if currPath == "." {
				return nil
			}

			// check if file exists in input folder
			inputPath := filepath.Join(diffInput, relativePath)
			goalStat, err := os.Stat(inputPath)
			if os.IsNotExist(err) || (!os.SameFile(info, goalStat) && !info.IsDir()) {
				yellowPrint(relativePath)
			}
			return nil
		})

		return nil
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

var diffInput string
var diffMergeIgnores bool

func init() {
	rootCmd.AddCommand(diffCmd)

	// Local flags
	diffCmd.Flags().StringVarP(&diffInput, "input", "i", ".", "Input directory of the files to be compared")
	diffCmd.Flags().BoolVarP(&diffMergeIgnores, "merge-ignores", "M", false, "Merges global and input dir .driveignore")
}
