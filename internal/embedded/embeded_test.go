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
