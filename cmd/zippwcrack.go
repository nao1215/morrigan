// Package cmd manages the entry points for the subcommands that morrigan has.
// Original author is ICHINOSE Shogo (https://github.com/shogo82148)
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
	"github.com/yeka/zip"
)

var zipPasswordCrackCmd = &cobra.Command{
	Use:     "zip-pwcrack FILE_PATH",
	Short:   "Crack zip file with password",
	Long:    `crack zip file with password`,
	Example: `  morrigan zip-pwcrack sample.zip`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := zipPwcrack(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(zipPasswordCrackCmd)
}

func zipPwcrack(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("zip-pwcrack subcommand need one argument (=zip file path)")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// TODO: use thread
	ch := make(chan string, 100)
	go genPassword(ch)
	for pass := range ch {
		fmt.Printf(".")
		if err := readZipBuffer(pass, bytes.NewReader(b)); err == nil {
			fmt.Println("")
			print.Info("password is " + pass)
			return nil
		} else if errors.Is(err, errNoPassword) {
			fmt.Println("")
			print.Info(err.Error())
			return nil
		}
	}
	fmt.Println("")
	return errors.New("can not find password")
}

var errNoPassword = errors.New("this zip file has no password")

func readZipBuffer(password string, buf *bytes.Reader) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	zipr, err := zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		return err
	}

	for _, z := range zipr.File {
		if z.IsEncrypted() {
			z.SetPassword(password)
		} else {
			return errNoPassword
		}

		rc, err := z.Open()
		if err != nil {
			return err
		}
		if _, err := io.Copy(io.Discard, rc); err != nil {
			return err
		}
		rc.Close()
	}
	return nil
}

func genPassword(chpass chan<- string) {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ*+!$%&^~-=|"
	count := int64(0)
	idx := make([]int, 1, 36)
	for {
		pass := make([]byte, 0, 36)
		for _, i := range idx {
			pass = append(pass, chars[i])
		}
		chpass <- string(pass)
		for i := 0; ; i++ {
			if i >= len(idx) {
				idx = idx[:i+1]
			}
			idx[i]++
			if idx[i] < len(chars) {
				break
			}
			idx[i] = 0
		}
		count++
	}
}
