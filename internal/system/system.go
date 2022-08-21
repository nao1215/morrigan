package system

import (
	"os"
	"runtime"
)

// IsWindows returns whether the execution environment is Windows or not.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsRootUser returns whether the executing user has root privileges.
func IsRootUser() bool {
	return os.Geteuid() == 0
}
