package cmdinfo

import (
	"fmt"
	"runtime/debug"
)

// Version value is set by ldflags
var Version string

// Name is command name
const Name = "morrigan"

// GetVersion return gup command version.
// Version global variable is set by ldflags.
func GetVersion() string {
	version := "unknown"
	if Version != "" {
		version = Version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		version = buildInfo.Main.Version
	}
	return fmt.Sprintf("%s version %s", Name, version)
}

