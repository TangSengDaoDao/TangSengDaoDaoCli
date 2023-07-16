package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type doctorCMD struct {
	ctx *TangSengDaoDaoContext
}

func newDoctorCMD(ctx *TangSengDaoDaoContext) *doctorCMD {

	return &doctorCMD{
		ctx: ctx,
	}
}

func (d *doctorCMD) CMD() *cobra.Command {
	doctorCMD := &cobra.Command{
		Use:   "doctor",
		Short: "doctor a TangSengDaoDao service.",
		RunE:  d.run,
	}

	return doctorCMD
}

func (d *doctorCMD) run(cmd *cobra.Command, args []string) error {
	containers, err := d.ctx.DockerComposePs([]string{d.ctx.opts.dockerComposePath})
	if err != nil {
		return err
	}
	results := make([]doctorResult, 0, len(containers))
	serviceNames := d.ctx.DockerComposeServices()
	for _, serviceName := range serviceNames {
		exist := false
		for _, container := range containers {
			if container.Service == serviceName {
				results = append(results, doctorResult{
					service: serviceName,
					state:   container.State,
				})
				exist = true
			}
		}
		if !exist {
			results = append(results, doctorResult{
				service: serviceName,
				state:   "not exist",
			})
		}
	}

	okFlag := "[✓]"
	errFlag := "[✗]"

	for _, result := range results {
		if result.state != "running" {
			fmt.Printf("\x1B[31m%s %s is %s\x1b[0m\n", errFlag, result.service, result.state)
		} else {
			fmt.Printf("\x1B[32m%s %s is %s\x1b[0m\n", okFlag, result.service, result.state)
		}

	}
	return nil
}

type doctorResult struct {
	service string
	state   string
}
