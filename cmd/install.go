package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type installCMD struct {
	ctx *TangSengDaoDaoContext
}

func newInstallCMD(ctx *TangSengDaoDaoContext) *installCMD {
	c := &installCMD{
		ctx: ctx,
	}

	return c
}

func (i *installCMD) CMD() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install a TangSengDaoDao service.",
		RunE:  i.run,
	}

	return installCmd
}

func (i *installCMD) run(cmd *cobra.Command, args []string) error {
	fmt.Println("install...")

	return i.ctx.DockerComposeUp([]string{"/Users/tt/work/projects/tangsengdaodao/golang/TangSengDaoDaoCli/docker-compose.yaml"})
}
