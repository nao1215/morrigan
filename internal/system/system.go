package system

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

// IsWindows returns whether the execution environment is Windows or not.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsRootUser returns whether the executing user has root privileges.
func IsRootUser() bool {
	return os.Geteuid() == 0
}

// ReadEtcPasswdFile return contents of /etc/passwd
func ReadEtcPasswdFile() ([]string, error) {
	bytes, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}

// ReadEtcShadowFile return contents of /etc/shadow
func ReadEtcShadowFile() ([]string, error) {
	if !IsRootUser() {
		return nil, errors.New("root privileges are required to read /etc/shadow")
	}

	bytes, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}
