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
	ctx     *TangSengDaoDaoContext
}

func NewTangSengDaoDao() *TangSengDaoDao {

	tsdd := &TangSengDaoDao{
		rootCmd: &cobra.Command{
			Use:   "tsdd",
			Short: "TangSengDaoDao is an open-source instant messaging system.",
			Long:  `TangSengDaoDao is an open-source instant messaging system., please refer to the documentation at https://tangsengdaodao.com`,
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}
	ctx := NewTangSengDaoDaoContext(tsdd)
	tsdd.ctx = ctx
	return tsdd
}

func (l *TangSengDaoDao) Context() *TangSengDaoDaoContext {
	return l.ctx
}

func (l *TangSengDaoDao) Options() *Options {
	return l.ctx.opts
}

func (l *TangSengDaoDao) addCommand(cmd CMD) {
	l.rootCmd.AddCommand(cmd.CMD())
}

func (l *TangSengDaoDao) Execute() {

	l.addCommand(newContextCMD(l.ctx))
	l.addCommand(newInstallCMD(l.ctx))
	l.addCommand(newStartCMD(l.ctx))
	l.addCommand(newStopCMD(l.ctx))
	l.addCommand(newUninstallCMD(l.ctx))
	l.addCommand(newUpgradeCMD(l.ctx))
	l.addCommand(newPsCMD(l.ctx))
	l.addCommand(newConfigCMD(l.ctx))
	l.addCommand(newDoctorCMD(l.ctx))
	l.addCommand(newVersionCMD(l.ctx))

	if err := l.rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
