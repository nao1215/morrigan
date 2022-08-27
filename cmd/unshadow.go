package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/unshadow"
	"github.com/spf13/cobra"
)

var unshadowCmd = &cobra.Command{
	Use:   "unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE",
	Short: "Combine password fields in /etc/passwd and /etc/shadow",
	Long: `unshadow subcommand combine password fields in /etc/passwd and /etc/shadow.

unshadow replaces the encrypted password written in the second field
of /etc/shadow with the second field of /etc/passwd. If you do not
specify two files (passwd, shadow) in the arguments, unshadow reads
the local system files.
`,
	Example: `  sudo morrigan unshadow
  sudo morrigan unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := unshadowRun(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unshadowCmd)
}

func unshadowRun(cmd *cobra.Command, args []string) error {
	if len(args) == 1 || len(args) > 2 {
		return errors.New("incorrect argument specified. see help, $ morrigan unshadow --help")
	}
	if len(args) == 0 {
		return unshadowLocalFiles()
	}

	return unshadowUserSpecifiedFiles(args[0], args[1])
}

func unshadowLocalFiles() error {
	unshadowList, err := unshadow.Unshadow(unshadow.PasswdFilePath, unshadow.ShadowFilePath)
	if err != nil {
		return err
	}

	for _, v := range unshadowList {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}

func unshadowUserSpecifiedFiles(passwdFile, shadowFile string) error {
	unshadowList, err := unshadow.Unshadow(passwdFile, shadowFile)
	if err != nil {
		return err
	}

	for _, v := range unshadowList {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}
