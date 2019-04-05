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
	"bytes"
	"io"
	"os"
	"sync"
)

// CatchOutput will temporarily catch all stdout and stderr and return it
func CatchOutput(f func()) (out string, err string) {
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	outs := make(chan string)
	errs := make(chan string)
	wg := sync.WaitGroup{}

	outReader, outWriter, _ := os.Pipe()
	os.Stdout = outWriter
	wg.Add(1)
	go func() {
		var buffer bytes.Buffer
		wg.Done()
		io.Copy(&buffer, outReader)
		outs <- buffer.String()
	}()
	wg.Wait()

	errReader, errWriter, _ := os.Pipe()
	os.Stderr = errWriter
	wg.Add(1)
	go func() {
		var buffer bytes.Buffer
		wg.Done()
		io.Copy(&buffer, errReader)
		errs <- buffer.String()
	}()
	wg.Wait()

	f()
	outWriter.Close()
	errWriter.Close()
	out = <-outs
	err = <-errs
	return
}
