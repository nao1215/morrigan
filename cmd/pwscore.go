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
	"github.com/spf13/cobra"
)

var passwdScoreCmd = &cobra.Command{
	Use:   "pwscore USERNAME",
	Short: "check password strength",
	Long: `pwscore subcommand checks password strength.

pwscore checks the guessability of password using the dictionary.
It also outputs advice on how to increase password strength.`,
	Example: `  morrigan pwscore nao`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pwscore(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(passwdScoreCmd)
}

func pwscore(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("pwscore subcommand need one argument (=username)")
	}

	passwd, err := interactive.ReadPassword()
	if err != nil {
		return err
	}

	if err := validWeakPasswd(passwd); err != nil {
		return nil
	}
	validLength(passwd)
	validContainUserName(args[0], passwd)
	validContainNumber(passwd)
	validContainLowerAndUpper(passwd)
	validContainSymbol(passwd)
	print.Info("[WIP] Calculate entropy as a measure of password strength")

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
	upper := false
	lower := false
	for _, v := range passwd {
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
	re := regexp.MustCompile(`[[:punct:]]|\.|\|\?|-|\!|\,|\@|\,|\#|\$|\%|\^|\&|\*|\_|\~`)
	if !re.Match([]byte(passwd)) {
		print.Warn("[Contains symbol     ] NG (Better to include symbol)")
		return
	}
	print.Info("[Contains symbol     ] OK")
}
