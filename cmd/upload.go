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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	gitignore "github.com/monochromegane/go-gitignore"
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
		// decide on .driveignore
		var driveignore gitignore.IgnoreMatcher
		var tempFileName string

		localDI := filepath.Join(uploadInput, ".driveignore")
		_, currFile, _, _ := runtime.Caller(0)
		globalDI := filepath.Join(currFile, "../../.global_driveignore")

		_, err1 := os.Stat(localDI)
		_, err2 := os.Stat(globalDI)
		if os.IsNotExist(err1) && os.IsNotExist(err2) {
			return errors.New("No local nor global .driveignores found")
		} else if (!os.IsNotExist(err1) && !mergeIgnores) || (os.IsNotExist(err2) && mergeIgnores) {
			driveignore, _ = gitignore.NewGitIgnore(localDI, "./")
			vPrint("loaded local .driveignore")
		} else if (!os.IsNotExist(err2) && !mergeIgnores) || (os.IsNotExist(err1) && mergeIgnores) {
			driveignore, _ = gitignore.NewGitIgnore(globalDI, "./")
			vPrint("loaded global .driveignore")
		} else if !os.IsNotExist(err1) && !os.IsNotExist(err2) && mergeIgnores {
			globalContent, _ := ioutil.ReadFile(globalDI)
			localContent, _ := ioutil.ReadFile(localDI)

			file, err := ioutil.TempFile("./", "tmp.*.temp")
			if err != nil {
				log.Fatal(err)
			}
			tempFileName = file.Name()
			defer os.Remove(file.Name())
			file.Write([]byte(string(globalContent) + "\n" + string(localContent)))

			driveignore, _ = gitignore.NewGitIgnore(file.Name(), "./")
			file.Close()
			vPrint("loaded merged global and local .driveignore")
		}
		err := filepath.Walk(uploadInput, func(path string, info os.FileInfo, err error) error {
			goalPath, _ := filepath.Rel(uploadInput, path)
			if err != nil {
				panic(err)
			}

			// skip the folder itself and temp merge file
			if path == "." || path == tempFileName {
				return nil
			}

			// ignore .driveignore files/dirs
			if info.IsDir() && driveignore.Match(path, true) {
				vPrint("skipped directory:", path)
				return filepath.SkipDir
			} else if !info.IsDir() && driveignore.Match(path, false) {
				vPrint("skipped file:", path)
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
				} else {
					err = os.Link(path, goal)
					vPrint("created hard link:", path)
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
var mergeIgnores bool

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Local flags
	uploadCmd.Flags().StringVarP(&uploadInput, "input", "i", "./", "Input directory of the files to be uploaded")
	uploadCmd.Flags().BoolVarP(&mergeIgnores, "mergeIgnores", "M", false, "Merges global and current dir .driveignore")
}
