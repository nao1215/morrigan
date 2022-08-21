package system

import (
	"errors"
	"os"
	"runtime"
	"sort"
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

		unshadow := passwdfields[0] + ":" + shadowFields[1] + strings.Join(passwdfields[2:], ":")
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
