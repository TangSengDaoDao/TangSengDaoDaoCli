package cmd

import "github.com/spf13/cobra"

type upgradeCMD struct {
	ctx *TangSengDaoDaoContext
}

func newUpgradeCMD(ctx *TangSengDaoDaoContext) *upgradeCMD {

	return &upgradeCMD{
		ctx: ctx,
	}
}

func (u *upgradeCMD) CMD() *cobra.Command {
	upgradeCMD := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade a TangSengDaoDao service.",
		RunE:  u.run,
	}

	return upgradeCMD
}

func (u *upgradeCMD) run(cmd *cobra.Command, args []string) error {

	err := u.pull("wukongim")
	if err != nil {
		return err
	}
	err = u.pull("tangsengdaodaoserver")
	if err != nil {
		return err
	}
	err = u.pull("tangsengdaodaoweb")
	if err != nil {
		return err
	}
	return nil
}

func (u *upgradeCMD) pull(service string) error {

	return u.ctx.DockerComposePull([]string{u.ctx.opts.dockerComposePath}, service)
}
