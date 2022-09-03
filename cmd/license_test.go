package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func Test_licenseRun(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no argument: print license list",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			want: []string{
				"The morrigan command user many codes that forked from other projects.",
				"All packages and subcommands not listed in the table are MIT LICENSE.",
				"We pay our respects to the author of the original code.",
				"",
				"+---------------+--------------------+----------------------------------------+--------------+",
				"| PACKAGE (CMD) |  ORIGINAL AUTHOR   |              ORIGINAL URL              |   LICENSE    |",
				"+---------------+--------------------+----------------------------------------+--------------+",
				"| gocrypt       | Akshay Moghe       | https://github.com/amoghe/go-crypt     | MIT          |",
				"| netcat        | Vadym Fedorov      | https://github.com/vfedoroff/go-netcat | MIT          |",
				"| morrigan      | Naohiro CHIKAMATSU | https://github.com/nao1215/morrigan    | MIT          |",
				"| ping          | u-root project     | https://github.com/u-root/u-root       | BSD 3-Clause |",
				"| whris         | harakeishi         | https://github.com/harakeishi/whris    | MIT          |",
				"| zip           | Alex Mullins       | https://github.com/yeka/zip            | MIT          |",
				"+---------------+--------------------+----------------------------------------+--------------+",
				"",
				"If you want to read the detailed license terms, pass the package",
				"name given in the table as arguments.",
				"[e.g.]",
				" $ morrigan license morrigan",
				"",
			},
			wantErr: false,
		},
		{
			name: "two argument: zip, morrigan package",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"zip", "morrigan"},
			},
			want: []string{
				"[zip]",
				"The MIT License (MIT)",
				"",
				"Copyright (C) 2022 Naohiro CHIKAMATSU",
				"Copyright (C) 2015 Alex Mullins",
				"",
				"Permission is hereby granted, free of charge, to any person obtaining a copy of",
				"this software and associated documentation files (the \"Software\"), to deal in",
				"the Software without restriction, including without limitation the rights to",
				"use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of",
				"the Software, and to permit persons to whom the Software is furnished to do so,",
				"subject to the following conditions:",
				"",
				"The above copyright notice and this permission notice shall be included in all",
				"copies or substantial portions of the Software.",
				"",
				"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR",
				"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS",
				"FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR",
				"COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER",
				"IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN",
				"CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.",
				"---------------------------------------------------------------------------",
				"[morrigan]",
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
				"",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgStdout := os.Stdout
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = pw

			if err := licenseRun(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("licenseRun() error = %v, wantErr %v", err, tt.wantErr)
			}
			pw.Close()
			os.Stdout = orgStdout

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
