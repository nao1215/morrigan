package interactive

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nao1215/morrigan/internal/print"
	"golang.org/x/term"
)

// ReadPassword get password from terminal (stdin).
func ReadPassword() (string, error) {
	// Get Ctrl-C (Interrupt) signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	// Get terminal state before password input
	currentState, err := term.GetState(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("%s: %w", "can not restore terminal state", err)
	}

	go func() {
		<-signalChan
		if err := term.Restore(int(syscall.Stdin), currentState); err != nil {
			print.Fatal(fmt.Errorf("%s: %w", "can not restore terminal state", err))
		}
		os.Exit(1)
	}()

	fmt.Printf("Enter password: ")
	passwd, err := term.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		return "", fmt.Errorf("%s: %w", "can not read password from stdin", err)
	}

	return string(passwd), nil
}
