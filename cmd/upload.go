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

	gitignore "github.com/monochromegane/go-gitignore"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [output path]",
	Short: "Upload files to your drive folder",
	Long: `Uploads files from the current directory (can be overwritten with --input flag) to a drive folder
It will ignore files that satisfy the .driveignore
The order of importance of a .driveignore file:
current folder > global config 
`,
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := os.Getwd()
		driveignore, _ := gitignore.NewGitIgnore("./.driveignore")
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			path, _ = filepath.Rel(root, path)
			if err != nil {
				panic(err)
			}

			// skip the folder itself
			if path == "." {
				return nil
			}

			// ignore .driveignore files/dirs
			if info.IsDir() && driveignore.Match(path, true) {
				return filepath.SkipDir
			} else if !info.IsDir() && driveignore.Match(path, false) {
				return nil
			}

			// link it doesnt yet exist
			goal := filepath.Join(args[0], path)
			if _, err := os.Stat(goal); os.IsNotExist(err) {
				if info.IsDir() {
					err = os.Mkdir(goal, os.ModePerm)
					if err != nil {
						panic(err)
					}
				} else {
					err = os.Link(path, goal)
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
	},
	Args: func(cmd *cobra.Command, args []string) error {
		cobra.ExactArgs(1)
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("passed path is invalid")
		}
		return nil
	},
}

var input string
var mergeIgnores bool

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Local flags
	uploadCmd.Flags().StringVarP(&input, "input", "i", "./", "Input directory of the files to be uploaded")
	uploadCmd.Flags().BoolVarP(&mergeIgnores, "mergeIgnores", "M", false, "Merges global and current dir .driveignore")
}
