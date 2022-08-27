// Package unshadow is designed to analyze password management systems for
// UNIX-like operating systems. This package provides the function to combine
// /etc/passwd and /etc/shadow. Otherwise, it accepts password input while hiding user input.
package unshadow

import (
	"reflect"
	"testing"
)

func TestCrypt(t *testing.T) {
	type args struct {
		passwd string
		salt   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate hash from password",
			args: args{
				passwd: "password",
				salt:   "$y$j9T$EK7BPw2KNXh5fakmSslBN0$",
			},
			want:    "$y$j9T$EK7BPw2KNXh5fakmSslBN0$y/0n3K0U33ibc8Cegt7Bl39AuolJzYRogdbnKxcqSYD",
			wantErr: false,
		},
		{
			name: "bad format salt",
			args: args{
				passwd: "",
				salt:   "$10$j9T$EK7BPw2KNXh5fakmSslBN0$",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Crypt(tt.args.passwd, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Crypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnshadow(t *testing.T) {
	type args struct {
		passwdFilePath string
		shadowFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "generate unshadow file",
			args: args{
				passwdFilePath: "./testdata/01-ok-passwd.txt",
				shadowFilePath: "./testdata/01-ok-shadow.txt",
			},
			want: []string{
				"morrie:$y$j9T$AWBxIYtBpRyJuzyHhC/4M.$LpGFc4mc0F8/f9w150QMsvqku7hofX4r6YIFiUFiEj1:1002:1002::/home/morrie:/bin/sh",
				"testuser:$y$j9T$EK7BPw2KNXh5fakmSslBN0$y/0n3K0U33ibc8Cegt7Bl39AuolJzYRogdbnKxcqSYD:1003:1003:,,,:/home/testuser:/bin/bash",
			},
			wantErr: false,
		},
		{
			name: "not exist passwd file",
			args: args{
				passwdFilePath: "./testdata/not_exist.txt",
				shadowFilePath: "./testdata/01-ok-shadow.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not exist shadow file",
			args: args{
				passwdFilePath: "./testdata/01-ok-passwd.txt",
				shadowFilePath: "./testdata/not_exist.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not match password file and shadow file",
			args: args{
				passwdFilePath: "./testdata/02-ng-passwd.txt",
				shadowFilePath: "./testdata/02-ng-shadow.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid passwd file",
			args: args{
				passwdFilePath: "./testdata/01-ok-shadow.txt",
				shadowFilePath: "./testdata/01-ok-shadow.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid shadow file",
			args: args{
				passwdFilePath: "./testdata/01-ok-passwd.txt",
				shadowFilePath: "./testdata/01-ok-passwd.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "password file is empty",
			args: args{
				passwdFilePath: "./testdata/03-empty.txt",
				shadowFilePath: "./testdata/01-ok-shadow.txt",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "shadow file is empty",
			args: args{
				passwdFilePath: "./testdata/01-ok-passwd.txt",
				shadowFilePath: "./testdata/03-empty.txt",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unshadow(tt.args.passwdFilePath, tt.args.shadowFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unshadow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unshadow() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestReadPassword(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			// --- FAIL: TestReadPassword (0.00s)
            // --- FAIL: TestReadPassword/read_password_from_stdin (0.00s)
            // home/nao/github/github.com/nao1215/morrigan/unshadow/unshadow_test.go:34:
			// ReadPassword() error = can not restore terminal state: inappropriate ioctl for device, wantErr false
			name:    "read password from stdin",
			want:    "P@ssw0rd",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcDefer, err := mockStdin(t, tt.want)
			if err != nil {
				t.Fatal(err)
			}
			defer funcDefer()

			got, err := ReadPassword()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mockStdin is a helper function that lets the test pretend dummyInput as os.Stdin.
// It will return a function for `defer` to clean up after the test.
func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	tmpFile, err := os.CreateTemp(t.TempDir(), "morrigan_")

	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpFile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpFile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpFile.Name())
	}, nil
}
*/
