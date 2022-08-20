package cmd

import (
	"os"

	"github.com/nao1215/morrigan/internal/completion"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/internal/system"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "morrigan",
}

// Execute start command.
func Execute() {
	if system.IsWindows() {
		print.Err("not support windows")
		os.Exit(1)
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true
	completion.DeployShellCompletionFileIfNeeded(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
