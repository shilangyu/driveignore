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
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shilangyu/driveignore/utils"
)

// globalCmd represents the global command
var globalCmd = &cobra.Command{
	Use:   "global [path to .driveignore]",
	Short: "Set your global .driveignore",
	Long: `If you wish to have a global .driveignore you can set the content of to it here.
You can later decide if you want to use global, local or merged .driveignore.
Once you set your global .driveignore you can delete the file you pointed to.`,
	RunE: globalRun(utils.GlobalDriveignorePath()),
	Args: globalArg,
}

var (
	errRuntimeCaller    = errors.New("runtime.Caller not ok")
	errOneArg           = errors.New("There should only be one argument")
	errPathDoesntExist  = errors.New("Passed path doesnt exist")
	errPathIsADirectory = errors.New("Passed path is a directory, not file")
)

func globalRun(globalDriveignorePath string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		vPrint := utils.VPrintWrapper(verbose)
		_, currFile, _, ok := runtime.Caller(0)
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if !ok {
			return errRuntimeCaller
		}
		absPath := getAbsPathFrom(args[0], cwd)
		currDir, _ := path.Split(currFile)
		globalDriveignoreAbsPath := getAbsPathFrom(globalDriveignorePath, currDir)
		content, err := ioutil.ReadFile(absPath)
		if err != nil {
			return err
		}
		vPrint("loaded file contents")
		if shouldAppend {
			currContent, err := ioutil.ReadFile(globalDriveignoreAbsPath)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(globalDriveignoreAbsPath, []byte(string(content)+"\n"+string(currContent)), os.ModePerm); err != nil {
				return err
			}
			vPrint("appended the contents to the .global_driveignore")
		} else {
			if err := ioutil.WriteFile(globalDriveignoreAbsPath, content, os.ModePerm); err != nil {
				return err
			}
			vPrint("saved the contents to the .global_driveignore")
		}
		return nil
	}
}

func getAbsPathFrom(path, cwd string) string {
	if strings.HasPrefix(path, "/") {
		return path // already abs path
	}
	absPath := filepath.Join(cwd, path)
	return absPath
}

func globalArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errOneArg
	}
	fstat, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return errPathDoesntExist
	}
	if fstat.IsDir() {
		return errPathIsADirectory
	}
	return nil
}

var shouldAppend bool

func init() {
	rootCmd.AddCommand(globalCmd)

	globalCmd.Flags().BoolVarP(&shouldAppend, "append", "a", false, "appends file contents to existing .global_driveignore")
}
