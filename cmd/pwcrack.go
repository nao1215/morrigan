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

var passwdCrackCmd = &cobra.Command{
	Use:   "pwcrack USERNAME",
	Short: "crack local user password",
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

		if err := compareChecksums(fields[1]); err != nil {
			return err
		}
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

func compareChecksums(encryptedPasswd string) error {
	if strings.HasPrefix(encryptedPasswd, "$1$") {
		print.Info("[WIP] md5sum. not implement")
		return nil
	} else if strings.HasPrefix(encryptedPasswd, "$5$") {
		print.Info("[WIP] sha256sum.  not implement")
		return nil
	} else if strings.HasPrefix(encryptedPasswd, "$6$") {
		print.Info("[WIP] sha256sum.  not implement")
		return nil
	}
	return errors.New("encrypted password=" + encryptedPasswd + " is unknown checksum format")
}
