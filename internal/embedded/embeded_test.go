// Package embedded provides functions to read and manipulate
// files embedded in the morrigan command.
package embedded

import (
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/nao1215/morrigan/file"
)

func TestWeakPasswdList(t *testing.T) {
	list, err := file.ToList((path.Join("passwd", "weak.txt")))
	if err != nil {
		t.Fatal(err)
	}

	l := []string{}
	for _, v := range list {
		l = append(l, strings.ReplaceAll(v, "\n", ""))
	}

	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "success",
			want:    l,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WeakPasswdList()
			if (err != nil) != tt.wantErr {
				t.Errorf("WeakPasswdList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WeakPasswdList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargetLogList(t *testing.T) {
	list, err := file.ToList((path.Join("log-collect", "target-files.txt")))
	if err != nil {
		t.Fatal(err)
	}

	l := []string{}
	for _, v := range list {
		l = append(l, strings.ReplaceAll(v, "\n", ""))
	}

	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "success",
			want:    l,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TargetLogList()
			if (err != nil) != tt.wantErr {
				t.Errorf("TargetLogList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TargetLogList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLicense(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				pkg: "morrigan",
			},
			want: []string{
				"MIT License",
				"",
				"Copyright (c) 2022 CHIKAMATSU Naohiro",
				"",
				"Permission is hereby granted, free of charge, to any person obtaining a copy",
				"of this software and associated documentation files (the \"Software\"), to deal",
				"in the Software without restriction, including without limitation the rights",
				"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell",
				"copies of the Software, and to permit persons to whom the Software is",
				"furnished to do so, subject to the following conditions:",
				"",
				"The above copyright notice and this permission notice shall be included in all",
				"copies or substantial portions of the Software.",
				"",
				"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR",
				"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,",
				"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE",
				"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER",
				"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,",
				"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE",
				"SOFTWARE.",
			},
			wantErr: false,
		},
		{
			name: "no-exist package does not exist",
			args: args{
				pkg: "no-exist",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := License(tt.args.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("License() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("License() = %v, want %v", got, tt.want)
			}
		})
	}
}
