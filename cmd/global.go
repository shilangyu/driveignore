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
	"os"
	"path/filepath"
	"runtime"

	"github.com/shilangyu/driveignore/utils"
	"github.com/spf13/cobra"
)

// globalCmd represents the global command
var globalCmd = &cobra.Command{
	Use:   "global [path to .driveignore]",
	Short: "Set your global .driveignore",
	Long: `If you wish to have a global .driveignore you can set the content of to it here.
You can later decide if you want to use global, local or merged .driveignore.
Once you set your global .driveignore you can delete the file you pointed to.`,
	Run: func(cmd *cobra.Command, args []string) {
		vPrint := utils.VPrintWrapper(verbose)

		_, currFile, _, _ := runtime.Caller(0)
		cwd, _ := os.Getwd()
		absPath := filepath.Join(cwd, args[0])
		content, _ := ioutil.ReadFile(absPath)
		vPrint("loaded file contents")
		if shouldAppend {
			currContent, _ := ioutil.ReadFile(filepath.Join(currFile, "../../.global_driveignore"))
			ioutil.WriteFile(filepath.Join(currFile, "../../.global_driveignore"), []byte(string(content)+"\n"+string(currContent)), os.ModePerm)
			vPrint("appended the contents to the .global_driveignore")
		} else {
			ioutil.WriteFile(filepath.Join(currFile, "../../.global_driveignore"), content, os.ModePerm)
			vPrint("saved the contents to the .global_driveignore")
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("There should only be one argument")
		}
		fstat, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			return errors.New("Passed path doesnt exist")
		}
		if fstat.IsDir() {
			return errors.New("Passed path is a directory, not file")
		}
		return nil
	},
}

var shouldAppend bool

func init() {
	rootCmd.AddCommand(globalCmd)

	globalCmd.Flags().BoolVarP(&shouldAppend, "append", "a", false, "appends file contents to existing .global_driveignore")
}
