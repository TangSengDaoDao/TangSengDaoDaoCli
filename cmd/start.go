package cmd

import (
	"github.com/spf13/cobra"
)

type startCMD struct {
	ctx *TangSengDaoDaoContext
}

func newStartCMD(ctx *TangSengDaoDaoContext) *startCMD {

	return &startCMD{
		ctx: ctx,
	}
}

func (s *startCMD) CMD() *cobra.Command {
	startCMD := &cobra.Command{
		Use:   "start",
		Short: "start a TangSengDaoDao service.",
		RunE:  s.run,
	}

	return startCMD
}

func (s *startCMD) run(cmd *cobra.Command, args []string) error {

	return s.ctx.DockerComposeUp([]string{s.ctx.opts.dockerComposePath})
}
