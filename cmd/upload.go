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

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [output path]",
	Short: "Upload a directory to your drive folder",
	Long: `Uploads files from the current directory (can be overwritten with --input flag) to a drive folder
It will ignore files that satisfy the .driveignore
The order of importance of a .driveignore file:
current folder > global config 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vPrint := utils.VPrintWrapper(verbose)

		driveignore, driveignoreType := utils.DriveIgnore(uploadInput, uploadMergeIgnores)

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

		err := filepath.Walk(uploadInput, func(path string, info os.FileInfo, err error) error {
			goalPath, _ := filepath.Rel(uploadInput, path)
			if err != nil {
				panic(err)
			}

			// skip the folder itself
			if path == "." {
				return nil
			}

			// ignore .driveignore files/dirs
			if info.IsDir() && driveignore.Match(path, true) {
				vPrint("skipped directory:", goalPath)
				return filepath.SkipDir
			} else if !info.IsDir() && driveignore.Match(path, false) {
				vPrint("skipped file:", goalPath)
				return nil
			}

			// if same name file already exists, check if its the same hardlink, then ignore
			// else if not, create hardlink
			// if its a directory, create one if doesnt yet exist
			goal := filepath.Join(args[0], goalPath)
			goalStat, err := os.Stat(goal)
			pathStat, _ := os.Stat(path)
			sameNameDiffFile := false
			if !info.IsDir() && !os.IsNotExist(err) && !os.SameFile(pathStat, goalStat) {
				os.Remove(goal)
				sameNameDiffFile = true
				vPrint("overwritting a file with same name")
			}
			if os.IsNotExist(err) || sameNameDiffFile {
				if info.IsDir() {
					err = os.Mkdir(goal, os.ModePerm)
					if err != nil {
						panic(err)
					}
					vPrint("created directory:", goalPath)
				} else {
					err = os.Link(path, goal)
					vPrint("created hard link:", goalPath)
					if err != nil {
						panic(err)
					}
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
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

var uploadInput string
var uploadMergeIgnores bool

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Local flags
	uploadCmd.Flags().StringVarP(&uploadInput, "input", "i", "./", "Input directory of the files to be uploaded")
	uploadCmd.Flags().BoolVarP(&uploadMergeIgnores, "mergeIgnores", "M", false, "Merges global and current dir .driveignore")
}
