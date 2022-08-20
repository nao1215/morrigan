package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nao1215/morrigan/internal/embedded"
	"github.com/nao1215/morrigan/internal/interactive"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/internal/system"
	"github.com/spf13/cobra"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd USERNAME",
	Short: "passwd subcommand checks password score or crack password",
	Long: `passwd subcommand checks password score or crack password.

It checks the guessability of passwords using a dictionary.
By default, passwd attempts to crack password for local accounts.
In crack mode, root privileges are required to read /etc/shadow.`,
	Example: `  morrigan passwd --score nao
  sudo morrigan passwd viktoriya`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := passwd(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	passwdCmd.Flags().BoolP("score", "s", false, "check password score")
	rootCmd.AddCommand(passwdCmd)
}

func passwd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("passwd subcommand need one argument (=username)")
	}

	scoreMode, err := cmd.Flags().GetBool("score")
	if err != nil {
		print.Fatal(err)
		return errors.New("can not parse command line argument (--score)")
	}

	if scoreMode {
		return score(args[0])
	}
	return crack(args[0])
}

func score(username string) error {

	passwd, err := interactive.ReadPassword()
	if err != nil {
		return err
	}

	list, err := embedded.WeakPasswdList()
	if err != nil {
		return err
	}
	for _, v := range list {
		if v == passwd {
			fmt.Println("out")
			break
		}
	}
	return nil
}

func crack(username string) error {
	if !system.IsRootUser() {
		return errors.New("passwd subcommand crack-mode (default) needs root privileges")
	}
	return nil
}
