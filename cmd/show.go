// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows the current .global_driveignore",
	Run: func(cmd *cobra.Command, args []string) {
		_, currFile, _, _ := runtime.Caller(0)
		content, _ := ioutil.ReadFile(filepath.Join(currFile, "../../.global_driveignore"))
		fmt.Println(string(content))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("There should be zero arguments")
		}
		return nil
	},
}

func init() {
	globalCmd.AddCommand(showCmd)
}
