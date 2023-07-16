package cmd

import "github.com/spf13/cobra"

type stopCMD struct {
	ctx *TangSengDaoDaoContext
}

func newStopCMD(ctx *TangSengDaoDaoContext) *stopCMD {

	return &stopCMD{
		ctx: ctx,
	}
}
func (s *stopCMD) CMD() *cobra.Command {
	startCMD := &cobra.Command{
		Use:   "stop",
		Short: "stop a TangSengDaoDao service.",
		RunE:  s.run,
	}

	return startCMD
}

func (s *stopCMD) run(cmd *cobra.Command, args []string) error {

	return s.ctx.DockerComposeStop([]string{s.ctx.opts.dockerComposePath})
}
