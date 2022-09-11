// Package embedded provides functions to read and manipulate
// files embedded in the morrigan command.
package embedded

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nao1215/morrigan/file"
	"github.com/nao1215/morrigan/internal/print"
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

//go:embed passwd
var passwordDir embed.FS

// GeneratePasswdListFiles generate password list files that used in morrigan command at target directory.
func GeneratePasswdListFiles(targetDir string) error {
	err := fs.WalkDir(passwordDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		src, err := passwordDir.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s %s: %w", "can not read", path, err)
		}

		targetDir := strings.TrimSuffix(targetDir, string(filepath.Separator))
		destDir := strings.Replace(filepath.Dir(path), "passwd", targetDir, 1)
		if !file.IsDir(destDir) {
			if !file.Exists(destDir) {
				err := os.MkdirAll(destDir, 0755)
				if err != nil {
					return err
				}
			}
		}

		destFile := filepath.Join(destDir, d.Name())
		dest, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = dest.WriteString(string(src))
		if err != nil {
			return err
		}

		print.Info("generate " + destFile)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
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
