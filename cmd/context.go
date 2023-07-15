package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/flags"
	composecmd "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/cmd/formatter"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/compose/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type TangSengDaoDaoContext struct {
	opts          *Options
	w             *TangSengDaoDao
	dockerCompose api.Service
}

func NewTangSengDaoDaoContext(w *TangSengDaoDao) *TangSengDaoDaoContext {
	c := &TangSengDaoDaoContext{
		opts: NewOptions(),
		w:    w,
	}
	err := c.opts.Load()
	if err != nil {
		panic(err)
	}
	return c
}

func (t *TangSengDaoDaoContext) DockerCompose() api.Service {
	if t.dockerCompose != nil {
		return t.dockerCompose
	}
	dockerSock, err := t.findDockerSock()
	if err != nil {
		panic(err)
	}
	opts := &flags.ClientOptions{Hosts: []string{dockerSock}}
	apiClient, err := command.NewAPIClientFromFlags(opts, &configfile.ConfigFile{})
	if err != nil {
		panic(err)
	}
	dockerCli, err := command.NewDockerCli(command.WithAPIClient(apiClient))
	if err != nil {
		panic(err)
	}
	err = dockerCli.Initialize(&flags.ClientOptions{
		Debug: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("zzzz-end")

	t.dockerCompose = compose.NewComposeService(dockerCli)
	return t.dockerCompose
}

func (t *TangSengDaoDaoContext) DockerComposeUp(configs []string) error {

	composeProjectOpts := composecmd.ProjectOptions{
		ProjectName: "tsdd",
		ConfigPaths: configs,
	}

	project, err := composeProjectOpts.ToProject(nil, cli.WithName("tsdd"))
	if err != nil {
		return err
	}

	return t.dockerComposeUp(project)
}

func (t *TangSengDaoDaoContext) dockerComposeUp(project *types.Project) error {

	ctx := context.Background()
	timeout := time.Duration(time.Second * 60)
	create := api.CreateOptions{
		IgnoreOrphans: false,
		RemoveOrphans: true,
		Recreate:      api.RecreateDiverged, // 配置改变就重新创建
		QuietPull:     false,                // QuietPull使拉取过程变得安静
		Timeout:       &timeout,
	}
	attachTo := utils.Set[string]{}
	consumer := formatter.NewLogConsumer(ctx, os.Stdout, os.Stderr, true, true, true)
	err := t.DockerCompose().Up(ctx, project, api.UpOptions{
		Create: create,
		Start: api.StartOptions{
			Project:     project,
			AttachTo:    attachTo.Elements(),
			Attach:      consumer,
			CascadeStop: false, // 在容器停止时停止应用程序
			Wait:        false, // Wait函数将等待容器达到运行或健康状态，直到返回
			WaitTimeout: time.Minute * 2,
		},
	})
	return err
}

func (t *TangSengDaoDaoContext) findDockerSock() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dockerSockPath := path.Join(homeDir, ".docker", "run", "docker.sock")
	if _, err := os.Stat(dockerSockPath); err != nil {
		dockerSockPath = "/var/run/docker.sock"
	}
	return fmt.Sprintf("unix://%s", dockerSockPath), nil
}

type contextCMD struct {
	cmd         *cobra.Command
	description string
	server      string
	token       string
	ctx         *TangSengDaoDaoContext
}

func newContextCMD(ctx *TangSengDaoDaoContext) *contextCMD {
	c := &contextCMD{
		ctx: ctx,
	}
	c.cmd = &cobra.Command{
		Use:   "context",
		Short: "Manage TangSengDaoDao configuration contexts",
	}
	return c
}

func (c *contextCMD) CMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage TangSengDaoDao configuration contexts",
	}
	c.initSubCMD(cmd)
	return cmd
}

func (c *contextCMD) initSubCMD(cmd *cobra.Command) {
	addCMD := &cobra.Command{
		Use:   "add",
		Short: "Update or create a context",
		RunE:  c.add,
	}
	addCMD.Flags().StringVar(&c.description, "description", c.ctx.opts.Description, "Context description")
	addCMD.Flags().StringVarP(&c.server, "server", "s", c.ctx.opts.ServerAddr, "Http  api server address")
	addCMD.Flags().StringVar(&c.token, "token", "", "Token for connect TangSengDaoDao")
	cmd.AddCommand(addCMD)
}

func (c *contextCMD) add(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		cmd.Help()
		return nil
	}
	name := args[0]
	if !validName(name) {
		return errors.New("invalid name")
	}
	c.ctx.opts.Description = c.description
	c.ctx.opts.ServerAddr = c.server
	c.ctx.opts.Token = c.token
	return c.ctx.opts.Save(name)
}

func validName(name string) bool {
	return name != "" && !strings.Contains(name, "..") && !strings.Contains(name, string(os.PathSeparator))
}
