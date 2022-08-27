// Package embedded provides functions to read and manipulate
// files embedded in the morrigan command.
package embedded

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed passwd/weak.txt
var weekPasswordListFile embed.FS

// WeakPasswdList return weak (famous) password list.
func WeakPasswdList() ([]string, error) {
	in, err := weekPasswordListFile.ReadFile(path.Join("passwd", "weak.txt"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not open weak password list", err)
	}

	return strings.Split(string(in), "\n"), nil
}
