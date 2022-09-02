package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/nao1215/morrigan/file"
	"github.com/nao1215/morrigan/internal/embedded"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/nao1215/morrigan/unshadow"
	"github.com/spf13/cobra"
)

var logCollectCmd = &cobra.Command{
	Use:   "log-collect",
	Short: "Collect the logs that exist in the system",
	Long: `Collect the logs that exist in the system.
`,
	Example: `  sudo morrigan log-collect`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := logCollect(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	logCollectCmd.Flags().StringP("output", "o", "log-collect", "copy destination directory")
	rootCmd.AddCommand(logCollectCmd)
}

func logCollect(cmd *cobra.Command, args []string) error {
	flag, err := parseLogCollectArgs(cmd, args)
	if err != nil {
		return err
	}

	if !unshadow.IsRootUser() {
		return errors.New("only use by root user")
	}

	logFileList, err := embedded.TargetLogList()
	if err != nil {
		return err
	}

	ids := &sudoIDs{}
	if err := ids.GetIDs(); err != nil {
		print.Warn(fmt.Errorf("%s: %w", "owner of the collected log files is root", err))
	}

	// Continue copying even if an error occurs
	for _, v := range logFileList {
		if !file.IsFile(v) {
			print.Info("skip. " + v + " does not exist")
			continue
		}

		destDir := filepath.Join(flag.outputDir, filepath.Dir(v))
		if !file.Exists(destDir) {
			err = os.MkdirAll(destDir, 0755)
			if err != nil {
				return err
			}
		}

		destFile := filepath.Join(destDir, filepath.Base(v))
		if err := file.Copy(v, destFile); err != nil {
			print.Err(fmt.Errorf("%s %s: %w", "can not copy", v, err))
			continue
		}
		print.Info("copy " + v + " to " + destFile)
	}

	if err := changeOwnerRecursive(flag.outputDir, ids.sudoUID, ids.sudoGID); err != nil {
		print.Warn(fmt.Errorf("%s: %w", "owner of the collected log files is root", err))
	}
	return nil
}

func changeOwnerRecursive(path string, uid int, gid int) error {
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := os.Chown(p, uid, gid); err != nil {
			return err
		}
		return nil
	})
	return err
}

type sudoIDs struct {
	sudoUID int
	sudoGID int
}

func (s *sudoIDs) GetIDs() error {
	sudoUIDStr := os.Getenv("SUDO_UID")
	sudoUID, err := strconv.Atoi(sudoUIDStr)
	if err != nil {
		return err
	}
	s.sudoUID = sudoUID

	user, err := user.LookupId(sudoUIDStr)
	if err != nil {
		return err
	}

	sudoGID, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}
	s.sudoGID = sudoGID

	return nil
}

type logCollectFlag struct {
	outputDir string
}

func parseLogCollectArgs(cmd *cobra.Command, args []string) (*logCollectFlag, error) {
	flag := logCollectFlag{}

	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--output)", err)
	}
	flag.outputDir = os.ExpandEnv(outputDir)

	return &flag, nil
}
