package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_globalArg(t *testing.T) {
	req := require.New(t)
	tmpFile, err := ioutil.TempFile("", "driveignore_Test_globalArg_file_*")
	req.NoError(err)

	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Good",
			args: args{
				cmd:  nil,
				args: []string{tmpFile.Name()},
			},
			wantErr: nil,
		},
		{
			name: "Bad1",
			args: args{
				cmd:  nil,
				args: []string{tmpFile.Name(), tmpFile.Name()},
			},
			wantErr: errOneArg,
		},
		{
			name: "Bad2",
			args: args{
				cmd:  nil,
				args: []string{tmpFile.Name() + "_some_salt"},
			},
			wantErr: errPathDoesntExist,
		},
		{
			name: "Bad3",
			args: args{
				cmd:  nil,
				args: []string{"."},
			},
			wantErr: errPathIsADirectory,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := globalArg(tt.args.cmd, tt.args.args); err != tt.wantErr {
				t.Errorf("globalArg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_globalRun(t *testing.T) {
	req := require.New(t)
	tmpGlobalDriveignore, tmpNewDriveignore1, tmpNewDriveignore2, tmpGlobalDriveignoreRelative, tmpNewDriveignoreRelative, closer := createTestFiles(req)
	defer closer()
	// region 1 - Simple add with shouldAppend
	globalRunConstructed := globalRun(tmpGlobalDriveignore.Name())
	shouldAppend = true
	req.NoError(globalRunConstructed(nil, []string{tmpNewDriveignore1.Name()}))
	req.NoError(globalRunConstructed(nil, []string{tmpNewDriveignore2.Name()}))

	content := fileContent(tmpGlobalDriveignore.Name(), req)

	req.Equal(
		[]string{
			"New2",
			"New1",
			"New0",
		},
		content,
	)
	// endregion
	// region 2 - Simple add no shouldAppend
	globalRunConstructed = globalRun(tmpGlobalDriveignore.Name())
	shouldAppend = false
	req.NoError(globalRunConstructed(nil, []string{tmpNewDriveignore1.Name()}))
	req.NoError(globalRunConstructed(nil, []string{tmpNewDriveignore2.Name()}))

	content = fileContent(tmpGlobalDriveignore.Name(), req)

	req.Equal(
		[]string{
			"New2",
		},
		content,
	)
	// endregion
	// region 3 - Relative paths
	shouldAppend = true
	globalRunConstructed = globalRun(getFileName(tmpGlobalDriveignoreRelative.Name()))
	req.NoError(globalRunConstructed(nil, []string{getFileName(tmpNewDriveignoreRelative.Name())}))
	req.NoError(globalRunConstructed(nil, []string{getFileName(tmpNewDriveignoreRelative.Name())}))

	content = fileContent(tmpGlobalDriveignoreRelative.Name(), req)

	req.Equal(
		[]string{
			"New4",
			"New4",
			"New3",
		},
		content,
	)
	// endregion
}

func createTestFiles(req *require.Assertions) (*os.File, *os.File, *os.File, *os.File, *os.File, func()) {
	tmpGlobalDriveignore, err := ioutil.TempFile("", "driveignore_Test_globalRun_file_g_*")
	req.NoError(err)
	tmpNewDriveignore1, err := ioutil.TempFile("", "driveignore_Test_globalRun_file_n1_*")
	req.NoError(err)
	tmpNewDriveignore2, err := ioutil.TempFile("", "driveignore_Test_globalRun_file_n2_*")
	req.NoError(err)
	tmpGlobalDriveignoreRelative, err := ioutil.TempFile(".", "driveignore_Test_globalRunRelative_file_g_*")
	req.NoError(err)
	tmpNewDriveignoreRelative, err := ioutil.TempFile(".", "driveignore_Test_DriveignoreRunRelative_file_n1_*")
	req.NoError(err)

	closer := func() {
		req.NoError(os.Remove(tmpGlobalDriveignoreRelative.Name()))
		req.NoError(os.Remove(tmpNewDriveignoreRelative.Name()))
		req.NoError(tmpGlobalDriveignore.Close())
		req.NoError(tmpNewDriveignore1.Close())
		req.NoError(tmpNewDriveignore2.Close())
		req.NoError(tmpGlobalDriveignoreRelative.Close())
		req.NoError(tmpNewDriveignoreRelative.Close())
	}

	_, err = tmpGlobalDriveignore.WriteString("New0")
	req.NoError(err)
	_, err = tmpNewDriveignore1.WriteString("New1")
	req.NoError(err)
	_, err = tmpNewDriveignore2.WriteString("New2")
	req.NoError(err)
	_, err = tmpGlobalDriveignoreRelative.WriteString("New3")
	req.NoError(err)
	_, err = tmpNewDriveignoreRelative.WriteString("New4")
	req.NoError(err)

	return tmpGlobalDriveignore,
		tmpNewDriveignore1,
		tmpNewDriveignore2,
		tmpGlobalDriveignoreRelative,
		tmpNewDriveignoreRelative,
		closer
}

func fileContent(path string, req *require.Assertions) []string {
	var res []string
	file, err := os.Open(path)
	req.NoError(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	req.NoError(file.Close())
	req.NoError(scanner.Err())
	return res
}

func getFileName(fullPath string) string {
	_, file := path.Split(fullPath)
	return file
}
