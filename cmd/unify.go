// Copyright © 2019 Marcin Wojnarowski xmaricnmarcin@gmail.com
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

	"github.com/spf13/cobra"
)

// unifyCmd represents the unify command
var unifyCmd = &cobra.Command{
	Use:   "unify [output path]",
	Short: "Unifies 2 directories where input is the source",
	Long: `Uploads all files (with respect to .driveignores)
aswell as removes legacy files from the drive sync folder.

Its an alias for: 'driveignore upload [args] [flags] --force' + 'driveignore clean [args] [flags]'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set flags
		uploadForce = true
		uploadMergeIgnores = unifyMergeIgnores
		uploadInput = unifyInput
		cleanInput = unifyInput

		// call commands
		err := uploadCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = cleanCmd.RunE(cmd, args)
		if err != nil {
			return err
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

var unifyInput string
var unifyMergeIgnores bool

func init() {
	rootCmd.AddCommand(unifyCmd)

	// local flags
	unifyCmd.Flags().StringVarP(&unifyInput, "input", "i", "./", "Input directory of the files to be uploaded")
	unifyCmd.Flags().BoolVarP(&unifyMergeIgnores, "merge-ignores", "M", false, "Merges global and input dir .driveignore")
}
