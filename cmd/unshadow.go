package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/internal/system"
	"github.com/spf13/cobra"
)

var unshadowCmd = &cobra.Command{
	Use:   "unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE",
	Short: "combine password fields in /etc/passwd and /etc/shadow",
	Long: `unshadow subcommand combine password fields in /etc/passwd and /etc/shadow.

unshadow replaces the encrypted password written in the second field
of /etc/shadow with the second field of /etc/passwd. If you do not
specify two files (passwd, shadow) in the arguments, unshadow reads
the local system files.
`,
	Example: `  sudo morrigan unshadow
  sudo morrigan unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := unshadow(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(unshadowCmd)
}

func unshadow(cmd *cobra.Command, args []string) error {
	if len(args) == 1 || len(args) > 2 {
		return errors.New("incorrect argument specified. see help, $ morrigan unshadow --help")
	}
	if len(args) == 0 {
		return unshadowLocalFiles()
	}

	return unshadowUserSpecifiedFiles(args[0], args[1])
}

func unshadowLocalFiles() error {
	passwdList, err := system.ReadEtcPasswdFile()
	if err != nil {
		return err
	}

	shadowList, err := system.ReadEtcShadowFile()
	if err != nil {
		return err
	}

	unshadowList, err := system.Unshadow(passwdList, shadowList)
	if err != nil {
		return err
	}

	for _, v := range unshadowList {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}

func unshadowUserSpecifiedFiles(passwdFile, shadowFile string) error {
	passwdBytes, err := os.ReadFile(passwdFile)
	if err != nil {
		return fmt.Errorf("%s: %w", "can not read password file", err)
	}
	passwdList := strings.Split(string(passwdBytes), "\n")

	shadowBytes, err := os.ReadFile(shadowFile)
	if err != nil {
		return fmt.Errorf("%s: %w", "can not read shadow file", err)
	}
	shadowList := strings.Split(string(shadowBytes), "\n")

	unshadowList, err := system.Unshadow(passwdList, shadowList)
	if err != nil {
		return err
	}

	for _, v := range unshadowList {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}
