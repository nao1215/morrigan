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

// Info print information message at STDOUT.
func Info(msg string) {
	fmt.Fprintf(Stdout, "%s: %s\n", color.GreenString("INFO "), msg)
}

// Warn print warning message at STDERR.
func Warn(err interface{}) {
	fmt.Fprintf(Stderr, ":%s: %v\n", color.YellowString("WARN "), err)
}

// Err print error message at STDERR.
func Err(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.HiYellowString("ERROR"), err)
}

// Fatal print dying message at STDERR.
func Fatal(err interface{}) {
	fmt.Fprintf(Stderr, "%s: %v\n", color.RedString("FATAL"), err)
	os.Exit(1)
}
