package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type versionCMD struct {
	ctx *TangSengDaoDaoContext
}

func newVersionCMD(ctx *TangSengDaoDaoContext) *versionCMD {

	return &versionCMD{
		ctx: ctx,
	}
}

func (v *versionCMD) CMD() *cobra.Command {
	versionCMD := &cobra.Command{
		Use:   "version",
		Short: "version of tsdd cli",
		RunE:  v.run,
	}

	return versionCMD
}

func (v *versionCMD) run(cmd *cobra.Command, args []string) error {
	fmt.Printf(" tsdd version: %s\n", v.ctx.opts.Version)
	fmt.Printf(" tsdd commit: %s\n", v.ctx.opts.Commit)
	fmt.Printf(" tsdd date: %s\n", v.ctx.opts.CommitDate)
	fmt.Printf(" tsdd treeState: %s\n", v.ctx.opts.TreeState)
	return nil
}
