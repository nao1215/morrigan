// Package print defines functions to accept colored standard output and user input
package print

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInfo(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"INFO : test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Info(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

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
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"WARN : test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Warn(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

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
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestErr(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want: []string{"ERROR: test message", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			Err(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

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
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
func TestFatal(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name     string
		args     args
		want     []string
		exitcode int
	}{
		{
			name: "Print message",
			args: args{
				msg: "test message",
			},
			want:     []string{"FATAL: test message", ""},
			exitcode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := Stdout
			orgStderr := Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			Stdout = pw
			Stderr = pw

			orgOsExit := OsExit
			exitCode := 0
			OsExit = func(code int) {
				exitCode = code
			}
			defer func() { OsExit = orgOsExit }()

			Fatal(tt.args.msg)
			pw.Close()
			Stdout = orgStdout
			Stderr = orgStderr

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
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}

			if exitCode != tt.exitcode {
				t.Errorf("value is mismatch. want=%d got=%d", exitCode, tt.exitcode)
			}
		})
	}
}
