package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	formatter2 "github.com/docker/cli/cli/command/formatter"
	"github.com/docker/compose/v2/cmd/formatter"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

type psCMD struct {
	ctx *TangSengDaoDaoContext
}

func newPsCMD(ctx *TangSengDaoDaoContext) *psCMD {

	return &psCMD{
		ctx: ctx,
	}
}

func (p *psCMD) CMD() *cobra.Command {
	psCMD := &cobra.Command{
		Use:   "ps",
		Short: "list  tangSengDaoDao service.",
		RunE:  p.run,
	}

	return psCMD
}

func (p *psCMD) run(cmd *cobra.Command, args []string) error {

	containers, err := p.ctx.DockerComposePs([]string{p.ctx.opts.dockerComposePath})
	if err != nil {
		return err
	}
	return formatter.Print(containers, "table", os.Stdout,
		writer(containers),
		"NAME", "IMAGE", "COMMAND", "SERVICE", "CREATED", "STATUS", "PORTS")
}

func writer(containers []api.ContainerSummary) func(w io.Writer) {
	return func(w io.Writer) {
		for _, container := range containers {
			ports := displayablePorts(container)
			createdAt := time.Unix(container.Created, 0)
			created := units.HumanDuration(time.Now().UTC().Sub(createdAt)) + " ago"
			status := container.Status
			command := formatter2.Ellipsis(container.Command, 20)
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", container.Name, container.Image, strconv.Quote(command), container.Service, created, status, ports)
		}
	}
}

func displayablePorts(c api.ContainerSummary) string {
	if c.Publishers == nil {
		return ""
	}

	ports := make([]types.Port, len(c.Publishers))
	for i, pub := range c.Publishers {
		ports[i] = types.Port{
			IP:          pub.URL,
			PrivatePort: uint16(pub.TargetPort),
			PublicPort:  uint16(pub.PublishedPort),
			Type:        pub.Protocol,
		}
	}

	return formatter2.DisplayablePorts(ports)
}
