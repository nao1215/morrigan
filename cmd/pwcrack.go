// Package cmd manages the entry points for the subcommands that morrigan has.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/nao1215/morrigan/internal/embedded"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/unshadow"
	"github.com/spf13/cobra"
)

var passwdCrackCmd = &cobra.Command{
	Use:   "pwcrack USERNAME",
	Short: "Crack local user password",
	Long: `pwcrack subcommand crack local user password (need root privirage).

pwcrack checks the guessability of passwords using a dictionary.
It attempts to crack password for local accounts. It does not yet
support password cracking of users over the network.
`,
	Example: `  sudo morrigan pwcrack viktoriya`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pwcrack(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(passwdCrackCmd)
}

func pwcrack(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("passwd subcommand need one argument (=username)")
	}

	return crack(args[0])
}

func crack(username string) error {
	unshadowList, err := unshadow.Unshadow(unshadow.PasswdFilePath, unshadow.ShadowFilePath)
	if err != nil {
		return fmt.Errorf("%s: %w", "can not generate unshadow from /etc/passwd and /etc/shadow", err)
	}

	if !existUser(unshadowList, username) {
		return errors.New(username + " does not exist in the system")
	}

	for _, v := range unshadowList {
		fields := strings.Split(v, ":")

		if fields[0] != username {
			continue
		}

		if fields[1] == "*" {
			print.Info(username + " is temporarily disabled user")
			break
		}

		if fields[1] == "" {
			print.Info(username + " has not set the password")
			break
		}

		password, err := compareChecksums(fields[1])
		if err != nil {
			return err
		}
		print.Info(username + "'s password is " + password)
		break
	}
	return nil
}

func existUser(unshadowList []string, username string) bool {
	for _, v := range unshadowList {
		fields := strings.Split(v, ":")
		if fields[0] == username {
			return true
		}
	}
	return false
}

func compareChecksums(encryptedPasswdWithSaltAndID string) (string, error) {
	lastIndex := strings.LastIndex(encryptedPasswdWithSaltAndID, "$")
	saltWithID := encryptedPasswdWithSaltAndID[:lastIndex]

	list, err := embedded.WorstPasswdList()
	if err != nil {
		return "", err
	}

	bar := pb.StartNew(len(list))
	for _, v := range list {
		hash, err := unshadow.Crypt(v, saltWithID)
		if err != nil {
			bar.Finish()
			return "", err
		}

		if hash == encryptedPasswdWithSaltAndID {
			bar.Finish()
			return v, nil
		}
		bar.Increment()
	}

	bar.Finish()
	return "", errors.New("can not carck password")
}
