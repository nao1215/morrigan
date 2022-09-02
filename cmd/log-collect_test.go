package cmd

import (
	"os"
	"testing"

	"github.com/nao1215/morrigan/file"
)

func Test_collect(t *testing.T) {
	type args struct {
		srcList  []string
		dest     string
		copyFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				srcList:  []string{"testdata/log-collect-src.txt"},
				dest:     "testdata/dest",
				copyFile: "testdata/dest/testdata/log-collect-src.txt",
			},
			wantErr: false,
		},
		{
			name: "copy no exist file and exist file",
			args: args{
				srcList:  []string{"testdata/log-collect-src2.txt", "testdata/no_exist.txt"},
				dest:     "testdata/dest",
				copyFile: "testdata/dest/testdata/log-collect-src2.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := collect(tt.args.srcList, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("collect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.copyFile != "" && !file.IsFile(tt.args.copyFile) {
				t.Error("failed to copy: " + tt.args.srcList[0])
			}
		})
	}

	if err := os.RemoveAll("./testdata/dest"); err != nil {
		t.Fatal(err)
	}
}
