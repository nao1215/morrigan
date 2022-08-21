package cmd

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"

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

	if err := validWeakPasswd(passwd); err != nil {
		return nil
	}
	validLength(passwd)
	validContainUserName(username, passwd)
	validContainNumber(passwd)
	validContainLowerAndUpper(passwd)
	validContainSymbol(passwd)

	//TODO: Calculate entropy as a measure of password strength
	return nil
}

func validWeakPasswd(passwd string) error {
	list, err := embedded.WeakPasswdList()
	if err != nil {
		return err
	}
	for _, v := range list {
		if v == passwd {
			print.Warn("[Weak password or not] NG (Included in the weak password list)")
			return nil
		}
	}
	print.Info("[Weak password or not] OK")
	return nil
}

func validContainUserName(username, passwd string) {
	if strings.Contains(strings.ToLower(passwd), strings.ToLower(username)) {
		print.Warn("[Not contain name    ] NG (Better not to contain user name)")
		return
	}
	print.Warn("[Not contain name    ] OK")
}

func validContainNumber(passwd string) {
	re := regexp.MustCompile("[0-9]+")
	if !re.Match([]byte(passwd)) {
		print.Warn("[Contains number     ] NG (Better to include number)")
		return
	}
	print.Info("[Contains number     ] OK")
}

func validContainLowerAndUpper(passwd string) {
	r := []rune(passwd)

	upper := false
	lower := false
	for _, v := range r {
		if unicode.IsUpper(v) {
			upper = true
			continue
		}
		if unicode.IsLower(v) {
			lower = true
			continue
		}
	}

	if !lower {
		print.Warn("[Contains upper&lower] NG (Better to include lower character)")
		return
	}

	if !upper {
		print.Warn("[Contains upper&lower] NG (Better to include upper character)")
		return
	}
	print.Info("[Contains upper&lower] OK")
}

func validLength(passwd string) {
	if len(passwd) < 15 {
		print.Warn("[Length              ] NG (15 characters or more is recommended)")
		return
	}
	print.Info("[Length              ] OK")
}

func validContainSymbol(passwd string) {
	re := regexp.MustCompile("[[:punct:]]|\\.|\\|\\?|-|\\!|\\,|\\@|\\,|\\#|\\$|\\%|\\^|\\&|\\*|\\_|\\~")
	if !re.Match([]byte(passwd)) {
		print.Warn("[Contains symbol     ] NG (Better to include symbol)")
		return
	}
	print.Info("[Contains symbol     ] OK")
}

func crack(username string) error {
	if !system.IsRootUser() {
		return errors.New("passwd subcommand crack-mode (default) needs root privileges")
	}
	return nil
}
