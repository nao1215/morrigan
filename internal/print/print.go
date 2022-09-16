package print

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is new instance of Writer which handles escape sequence for stdout.
	Stdout = colorable.NewColorableStdout()
	// Stderr is new instance of Writer which handles escape sequence for stderr.
	Stderr = colorable.NewColorableStderr()
)

// Info print information message at STDOUT in green.
// This function is used to print some information (that is not error) to the user.
func Info(msg string) {
	fmt.Fprintf(Stdout, "%s: %s\n", color.GreenString("INFO "), msg)
}

// Warn print warning message at STDERR in yellow.
// This function is used to print warning message to the user.
func Warn(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.YellowString("WARN "), err)
}

// Err print error message at STDERR in yellow.
// This function is used to print error message to the user.
func Err(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.HiYellowString("ERROR"), err)
}

// OsExit is wrapper for  os.Exit(). It's for unit test.
var OsExit = os.Exit

// Fatal print dying message at STDERR in red.
// After print message, process will exit
func Fatal(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.RedString("FATAL"), err)
	OsExit(1)
}
