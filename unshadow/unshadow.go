// Package unshadow is designed to analyze password management systems for
// UNIX-like operating systems. This package provides the function to combine
// /etc/passwd and /etc/shadow. Otherwise, it accepts password input while hiding user input.
package unshadow

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/nao1215/morrigan/internal/gocrypt"
	"github.com/nao1215/morrigan/internal/print"

	"golang.org/x/term"
)

// IsRootUser returns whether the executing user has root privileges.
func IsRootUser() bool {
	return os.Geteuid() == 0
}

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

// Crypt provides a wrapper around the glibc crypt_r() function.
// For the meaning of the arguments, refer to the package README.
func Crypt(passwd, salt string) (string, error) {
	hash, err := gocrypt.Crypt(passwd, salt)
	if err != nil {
		return "", fmt.Errorf("%s%s: %w", "can not generate hash from password=", passwd, err)
	}
	return hash, nil
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

// Unshadow replaces the second field (password) in /etc/passwd with
// the second field (encrypted password) in /etc/shadow.
func Unshadow(passwdList, shadowList []string) ([]string, error) {
	if len(passwdList) != len(shadowList) {
		return nil, errors.New("/etc/passwd and /etc/shadow have different line numbers")
	}

	if !validEtcPasswd(passwdList) {
		return nil, errors.New("not /etc/passwd file")
	}

	if !validEtcShadow(shadowList) {
		return nil, errors.New("not /etc/shadow file")
	}

	sort.Strings(passwdList)
	sort.Strings(shadowList)

	var unshadowList []string
	for i, v := range passwdList {
		if len(v) == 0 {
			continue
		}

		passwdfields := strings.Split(v, ":")
		shadowFields := strings.Split(shadowList[i], ":")
		if passwdfields[0] != shadowFields[0] {
			return nil, errors.New("users do not match in /etc/passwd and /etc/shadow")
		}

		unshadow := passwdfields[0] + ":" + shadowFields[1] + ":" + strings.Join(passwdfields[2:], ":")
		unshadowList = append(unshadowList, unshadow)
	}
	return unshadowList, nil
}

func validEtcPasswd(passwdList []string) bool {
	if len(passwdList) == 0 {
		return false
	}
	return len(strings.Split(passwdList[0], ":")) == 7
}

func validEtcShadow(shadowList []string) bool {
	if len(shadowList) == 0 {
		return false
	}
	return len(strings.Split(shadowList[0], ":")) == 9
}
