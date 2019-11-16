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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shilangyu/driveignore/utils"
)

// globalCmd represents the global command
var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Get the path to your global .driveignore",
	Long: `If you wish to have a global .driveignore you can set the content of to it here.
You can later decide if you want to use global, local or merged .driveignore.`,
	Example: "vim $(driveignore global)",
	RunE:    globalRun(utils.GlobalDriveignorePath()),
	Args:    globalArg,
}

var (
	errNoArg = errors.New("There should only be no arguments")
)

func globalRun(globalDriveignorePath string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		vPrint := utils.VPrintWrapper(verbose)

		if _, err := os.Stat(globalDriveignorePath); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(globalDriveignorePath), os.ModePerm)
			err := ioutil.WriteFile(globalDriveignorePath, []byte{}, os.ModePerm)
			if err != nil {
				return nil
			}
			vPrint(".global_driveignore didnt exist, created a new one")
		}

		fmt.Println(globalDriveignorePath)
		return nil
	}
}

func globalArg(_ *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errNoArg
	}
	return nil
}

func init() {
	rootCmd.AddCommand(globalCmd)
}
