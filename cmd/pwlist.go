// Package cmd manages the entry points for the subcommands that morrigan has.
package cmd

import (
	"fmt"
	"os"

	"github.com/nao1215/morrigan/internal/embedded"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

var passwdListCmd = &cobra.Command{
	Use:   "pwlist [FLAGS]",
	Short: "Generate password list files",
	Long: `Generate password list files that be used in morrigan command.
`,
	Example: `  morrigan pwlist`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pwlist(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	passwdListCmd.Flags().StringP("output", "o", "password-list", "specifies the output directory")
	rootCmd.AddCommand(passwdListCmd)
}

func pwlist(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("%s: %w", "can not parse command line argument (--output)", err)
	}
	return embedded.GeneratePasswdListFiles(output)
}
