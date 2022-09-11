// Package cmd manages the entry points for the subcommands that morrigan has.
package cmd

import (
	"os"
	"testing"

	"github.com/nao1215/morrigan/file"
	"github.com/spf13/cobra"
)

func Test_pwlist(t *testing.T) {
	t.Run("Success generate password list", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().StringP("output", "o", "password-list", "specifies the output directory")

		if err := cmd.Flags().Set("output", "testdata/passwd"); err != nil {
			t.Fatal(err)
		}

		if err := pwlist(cmd, []string{}); err != nil {
			t.Errorf("failed to generate password list: %v", err)
		}

		if !file.IsFile("testdata/passwd/worst/2020-worst-200.txt") {
			t.Errorf("failed to generate testdata/passwd/worst/2020-worst-200.txt")
		}

		if !file.IsFile("testdata/passwd/worst/worst-500.txt") {
			t.Errorf("failed to generate testdata/passwd/worst/worst-500.txt")
		}

		if err := os.RemoveAll("./testdata/passwd"); err != nil {
			t.Fatal(err)
		}
	})
}
