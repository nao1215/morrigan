package interactive

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassword get password from terminal (stdin).
func ReadPassword() (string, error) {
	// Get Ctrl-C (Interrupt) signal
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	// Get terminal state before password input
	currentState, err := terminal.GetState(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("%s: %w", "can not restore terminal state", err)
	}

	go func() {
		<-signalChan
		terminal.Restore(int(syscall.Stdin), currentState)
		os.Exit(1)
	}()

	fmt.Printf("Enter password: ")
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		return "", fmt.Errorf("%s: %w", "can not read password from stdin", err)
	}

	return string(passwd), nil
}
