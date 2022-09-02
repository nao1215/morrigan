package file

import (
	"os"
	"path"
	"strings"
)

const (
	// Unknown : Don't use this bits
	Unknown os.FileMode = 1 << (9 - iota)
	// Readable : readable bits
	Readable
	// Writable : writable bits
	Writable
	// Executable : executable bits
	Executable
)

// IsFile reports whether the path exists and is a file.
func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (!stat.IsDir())
}

// Exists reports whether the path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

// IsDir reports whether the path exists and is a directory.
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (stat.IsDir())
}

// IsSymlink reports whether the path exists and is a symbolic link.
func IsSymlink(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}
	return false
}

// IsZero reports whether the path exists and is zero size.
func IsZero(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && (stat.Size() == 0)
}

// IsReadable reports whether the path exists and is readable.
func IsReadable(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && ((stat.Mode() & Readable) != 0)
}

// IsWritable reports whether the path exists and is writable.
func IsWritable(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && ((stat.Mode() & Writable) != 0)
}

// IsExecutable reports whether the path exists and is executable.
func IsExecutable(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil) && ((stat.Mode() & Executable) != 0)
}

// IsHiddenFile reports whether the path exists and is included hidden file.
func IsHiddenFile(filePath string) bool {
	_, file := path.Split(filePath)
	if IsFile(filePath) && strings.HasPrefix(file, ".") {
		return true
	}
	return false
}
