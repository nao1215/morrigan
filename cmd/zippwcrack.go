// Package cmd manages the entry points for the subcommands that morrigan has.
// Original author is ICHINOSE Shogo (https://github.com/shogo82148)
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/zip"
	"github.com/spf13/cobra"
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

	// Before reading a ZIP file in goroutine,
	// verify that the file is zip format or has a password.
	noPasswd, err := isNoPassword(bytes.NewReader(b))
	if err != nil {
		return err
	}
	if noPasswd {
		print.Info("zip file has no password")
		return nil
	}

	type zipPasswordCrackResult struct {
		password string
	}
	ch := make(chan string, 100)
	resultCh := make(chan zipPasswordCrackResult, 100)

	go genPassword(ch)
	for i := 0; i < 10; i++ {
		go func(result chan<- zipPasswordCrackResult) {
			for pass := range ch {
				if err := readZipBuffer(pass, bytes.NewReader(b)); err == nil {
					result <- zipPasswordCrackResult{
						password: pass,
					}
					return
				}
				result <- zipPasswordCrackResult{}
			}
		}(resultCh)
	}

	print.Info("detect passwords by brute force...")
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Start()
	for r := range resultCh {
		if r.password != "" {
			s.Stop()
			print.Info("zip file's password = " + r.password)
			return nil
		}
	}
	s.Stop()
	return errors.New("can not find password")
}

func isNoPassword(buf *bytes.Reader) (bool, error) {
	zipr, err := zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		return false, err
	}

	for _, z := range zipr.File {
		if !z.IsEncrypted() {
			return true, nil
		}
	}
	return false, nil
}

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
		z.SetPassword(password)

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
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*+!$%&^~-=|@*[{]}+;:?_"
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
