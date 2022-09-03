package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

func Test_pwscore(t *testing.T) {
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
			name: "no argument error",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pwscore(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("pwscore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_valid(t *testing.T) {
	type args struct {
		username string
		passwd   string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "weak password",
			args: args{
				username: "test",
				passwd:   "password",
			},
			want: []string{
				"WARN : [Weak password or not] NG (Included in the weak password list)",
				"WARN : [Length              ] NG (15 characters or more is recommended)",
				"INFO : [Not contain name    ] OK",
				"WARN : [Contains number     ] NG (Better to include number)",
				"WARN : [Contains upper&lower] NG (Better to include upper character)",
				"WARN : [Contains symbol     ] NG (Better to include symbol)",
				"INFO : [WIP] Calculate entropy as a measure of password strength",
				"",
			},
			wantErr: false,
		},
		{
			name: "good password",
			args: args{
				username: "test",
				passwd:   "pass123-!Buck-T",
			},
			want: []string{
				"INFO : [Weak password or not] OK",
				"INFO : [Length              ] OK",
				"INFO : [Not contain name    ] OK",
				"INFO : [Contains number     ] OK",
				"INFO : [Contains upper&lower] OK",
				"INFO : [Contains symbol     ] OK",
				"INFO : [WIP] Calculate entropy as a measure of password strength",
				"",
			},
			wantErr: false,
		},
		{
			name: "password contain username and only lower",
			args: args{
				username: "test",
				passwd:   "testuser",
			},
			want: []string{
				"INFO : [Weak password or not] OK",
				"WARN : [Length              ] NG (15 characters or more is recommended)",
				"WARN : [Not contain name    ] NG (Better not to contain user name)",
				"WARN : [Contains number     ] NG (Better to include number)",
				"WARN : [Contains upper&lower] NG (Better to include upper character)",
				"WARN : [Contains symbol     ] NG (Better to include symbol)",
				"INFO : [WIP] Calculate entropy as a measure of password strength",
				"",
			},
			wantErr: false,
		},
		{
			name: "only upper chara",
			args: args{
				username: "test",
				passwd:   "TEST",
			},
			want: []string{
				"INFO : [Weak password or not] OK",
				"WARN : [Length              ] NG (15 characters or more is recommended)",
				"WARN : [Not contain name    ] NG (Better not to contain user name)",
				"WARN : [Contains number     ] NG (Better to include number)",
				"WARN : [Contains upper&lower] NG (Better to include lower character)",
				"WARN : [Contains symbol     ] NG (Better to include symbol)",
				"INFO : [WIP] Calculate entropy as a measure of password strength",
				"",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := print.Stdout
			orgStderr := print.Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			print.Stdout = pw
			print.Stderr = pw

			err = valid(tt.args.username, tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
			}
			pw.Close()
			print.Stdout = orgStdout
			print.Stderr = orgStderr

			if err != nil {
				return
			}

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("User value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
