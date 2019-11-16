package cmd

import (
	"os"
	"testing"

	"github.com/shilangyu/driveignore/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_globalArg(t *testing.T) {
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
				args: []string{},
			},
			wantErr: nil,
		},
		{
			name: "Bad",
			args: args{
				cmd:  nil,
				args: []string{"something"},
			},
			wantErr: errNoArg,
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
	const p = "driveignore_Test_globalRun_file"
	globalRunConstructed := globalRun(p)

	// First time run: no global driveignore
	_, err := os.Stat(p)
	req.True(os.IsNotExist(err))
	out, _ := utils.CatchOutput(func() {
		req.NoError(globalRunConstructed(nil, []string{}))
	})
	req.Equal(p+"\n", out)

	// Second time run: global driveignore exists
	out, _ = utils.CatchOutput(func() {
		req.NoError(globalRunConstructed(nil, []string{}))
	})
	req.Equal(p+"\n", out)

	os.Remove(p)
}
