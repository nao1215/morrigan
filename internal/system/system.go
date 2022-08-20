package system

import (
	"runtime"
)

// IsWindows returns whether the execution environment is Windows or not.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
