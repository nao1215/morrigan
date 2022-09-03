package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_zipPwcrack(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "crack success",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"testdata/zip-password-is-ab.zip"},
			},
			wantErr: false,
		},
		{
			name: "no argument",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			wantErr: true,
		},
		{
			name: "bad file path",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"no_exist_path"},
			},
			wantErr: true,
		},
		{
			name: "zip file with no password",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"testdata/zip-no-password.zip"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := zipPwcrack(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("zipPwcrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
