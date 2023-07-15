package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version string
var Commit string
var CommitDate string

type TangSengDaoDao struct {
	rootCmd *cobra.Command
}

func NewTangSengDaoDao() *TangSengDaoDao {

	return &TangSengDaoDao{
		rootCmd: &cobra.Command{
			Use:   "tsdd",
			Short: "TangSengDaoDao is an open-source instant messaging system.",
			Long:  `TangSengDaoDao is an open-source instant messaging system., please refer to the documentation at https://tangsengdaodao.com`,
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}
}

func (l *TangSengDaoDao) addCommand(cmd CMD) {
	l.rootCmd.AddCommand(cmd.CMD())
}

func (l *TangSengDaoDao) Execute() {
	ctx := NewTangSengDaoDaoContext(l)
	l.addCommand(newContextCMD(ctx))
	l.addCommand(newInstallCMD(ctx))

	if err := l.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
