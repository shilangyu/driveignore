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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	gitignore "github.com/monochromegane/go-gitignore"
)

// IgnoreType is an enum representing the returned ignorer
type IgnoreType int

const (
	// NoIgnore says no .driveignores were found
	NoIgnore IgnoreType = iota
	// LocalIgnore says the local .driveignore has been taken
	LocalIgnore
	// GlobalIgnore says the global .driveignore has been taken
	GlobalIgnore
	// MergedIgnore says the local and global .driveignore has been taken
	MergedIgnore
)

// DriveIgnore returns a gitignore matcher with merge or not merged .driveignores
func DriveIgnore(localPath string, mergeIgnores bool) (driveignore gitignore.IgnoreMatcher, ignorer IgnoreType) {
	localDI := filepath.Join(localPath, ".driveignore")
	_, currFile, _, _ := runtime.Caller(0)
	globalDI := filepath.Join(currFile, "../../.global_driveignore")

	_, err1 := os.Stat(localDI)
	_, err2 := os.Stat(globalDI)
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		ignorer = NoIgnore
	} else if (!os.IsNotExist(err1) && !mergeIgnores) || (os.IsNotExist(err2) && mergeIgnores) {
		driveignore, _ = gitignore.NewGitIgnore(localDI, "./")
		ignorer = LocalIgnore
	} else if (!os.IsNotExist(err2) && !mergeIgnores) || (os.IsNotExist(err1) && mergeIgnores) {
		driveignore, _ = gitignore.NewGitIgnore(globalDI, "./")
		ignorer = GlobalIgnore
	} else if !os.IsNotExist(err1) && !os.IsNotExist(err2) && mergeIgnores {
		globalContent, _ := ioutil.ReadFile(globalDI)
		localContent, _ := ioutil.ReadFile(localDI)

		file, err := ioutil.TempFile("./", "tmp.*.temp")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
		file.Write([]byte(string(globalContent) + "\n" + string(localContent)))

		driveignore, _ = gitignore.NewGitIgnore(file.Name(), "./")
		file.Close()
		ignorer = MergedIgnore
	}
	return
}
