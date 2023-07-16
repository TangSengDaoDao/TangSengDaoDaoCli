package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type configCMD struct {
	ctx *TangSengDaoDaoContext
}

func newConfigCMD(ctx *TangSengDaoDaoContext) *configCMD {

	return &configCMD{
		ctx: ctx,
	}
}

func (c *configCMD) CMD() *cobra.Command {
	configCMD := &cobra.Command{
		Use:   "config",
		Short: "config a TangSengDaoDao service.",
		RunE:  c.run,
	}

	return configCMD
}

func (c *configCMD) run(cmd *cobra.Command, args []string) error {

	content, err := c.ctx.DockerComposeConfig([]string{c.ctx.opts.dockerComposePath})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, string(content))
	return err
}
