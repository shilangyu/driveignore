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

package utils

import (
	"os"
	"path/filepath"
)

// Walker walks through a directory with some preset actions
func Walker(path string, walk func(string, os.FileInfo, string) error) error {
	return filepath.Walk(path, func(currPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip the folder itself
		if currPath == "." {
			return nil
		}

		relativePath, _ := filepath.Rel(path, currPath)
		// adding slash to directories for print clarity
		if temp, _ := os.Stat(currPath); temp.IsDir() {
			relativePath += "\\"
		}

		return walk(currPath, info, relativePath)
	})
}

// GlobalDriveignorePath returns absolute path to .global_driveignore
func GlobalDriveignorePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "/driveignore/.global_driveignore")
}
