// Package embedded provides functions to read and manipulate
// files embedded in the morrigan command.
package embedded

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

//go:embed passwd/worst
var worstPasswordListDir embed.FS

// WorstPasswdList return popular and weakness password list.
func WorstPasswdList() ([]string, error) {
	passwdList := []string{}
	err := fs.WalkDir(worstPasswordListDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		target := filepath.Join("passwd/worst", d.Name())
		in, err := worstPasswordListDir.ReadFile(target)
		if err != nil {
			return fmt.Errorf("%s %s: %w", "can not read", target, err)
		}
		passwdList = append(passwdList, strings.Split(string(in), "\n")...)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return passwdList, err
}

//go:embed log-collect/target-files.txt
var targetLogListFile embed.FS

// TargetLogList return log file path list to be collected
func TargetLogList() ([]string, error) {
	in, err := targetLogListFile.ReadFile(path.Join("log-collect", "target-files.txt"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not open target log file path list to be collected", err)
	}

	return strings.Split(string(in), "\n"), nil
}

//go:embed license/*
var licenseDir embed.FS

// License return package license
func License(pkg string) ([]string, error) {
	in, err := licenseDir.ReadFile(path.Join("license", pkg+".LICENSE"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(in), "\n"), nil
}
