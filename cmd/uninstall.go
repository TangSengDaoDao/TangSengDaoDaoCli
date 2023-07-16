package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type uninstallCMD struct {
	ctx *TangSengDaoDaoContext
	all bool
}

func newUninstallCMD(ctx *TangSengDaoDaoContext) *uninstallCMD {

	return &uninstallCMD{
		ctx: ctx,
	}
}

func (u *uninstallCMD) CMD() *cobra.Command {
	uninstallCMD := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall a TangSengDaoDao service.",
		RunE:  u.run,
	}
	uninstallCMD.Flags().BoolVar(&u.all, "all", false, "Remove all including data（移除所有包括数据）")

	return uninstallCMD
}

func (u *uninstallCMD) run(cmd *cobra.Command, args []string) error {
	err := u.ctx.DockerComposeDown([]string{u.ctx.opts.dockerComposePath})
	if err != nil {
		return err
	}
	if u.all {
		err = os.RemoveAll(u.ctx.opts.rootDir)
		if err != nil {
			return err
		}
	}
	return nil
}
