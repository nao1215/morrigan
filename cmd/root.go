package cmd

import (
	"os"

	"github.com/nao1215/morrigan/internal/completion"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "morrigan",
}

// Execute start command.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true
	completion.DeployShellCompletionFileIfNeeded(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		print.Err(err)
		os.Exit(1)
	}
}
