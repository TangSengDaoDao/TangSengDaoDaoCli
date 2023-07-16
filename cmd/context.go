package cmd

import (
	"context"
	"embed"
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
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/compose/v2/pkg/utils"
	"github.com/spf13/cobra"
)

type TangSengDaoDaoContext struct {
	opts              *Options
	w                 *TangSengDaoDao
	dockerCompose     api.Service
	DockerComposeYaml string // docker compose yaml 文件内容
	DotEnv            string // .env文件内容
	Configs           embed.FS
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
	t.dockerCompose = api.NewServiceProxy().WithService(compose.NewComposeService(dockerCli))
	return t.dockerCompose
}

func (t *TangSengDaoDaoContext) DockerComposeUp(configs []string) error {

	project, err := t.getDockerProject(configs)
	if err != nil {
		return err
	}

	return t.dockerComposeUp(project)
}

func (t *TangSengDaoDaoContext) DockerComposeStop(configs []string) error {
	project, err := t.getDockerProject(configs)
	if err != nil {
		return err
	}
	return t.dockerComposeStop(project)
}

func (t *TangSengDaoDaoContext) DockerComposeDown(configs []string) error {
	project, err := t.getDockerProject(configs)
	if err != nil {
		return err
	}
	return t.dockerComposeDown(project)
}

func (t *TangSengDaoDaoContext) DockerComposeRemove(configs []string) error {
	project, err := t.getDockerProject(configs)
	if err != nil {
		return err
	}
	return t.dockerComposeRemove(project)
}
func (t *TangSengDaoDaoContext) DockerComposePs(configs []string) ([]api.ContainerSummary, error) {
	project, err := t.getDockerProject(configs)
	if err != nil {
		return nil, err
	}
	return t.dockerComposePs(project)
}

func (t *TangSengDaoDaoContext) DockerComposeServices() []string {
	project, err := t.getDockerProject([]string{t.opts.dockerComposePath})
	if err != nil {
		return nil
	}
	return project.ServiceNames()
}

func (t *TangSengDaoDaoContext) DockerComposePull(configs []string, services ...string) error {
	project, err := t.getDockerProject(configs, services...)
	if err != nil {
		return err
	}
	return t.dockerComposePull(project)
}

func (t *TangSengDaoDaoContext) DockerComposeConfig(configs []string) ([]byte, error) {
	project, err := t.getDockerProject(configs)
	if err != nil {
		return nil, err
	}
	return t.dockerComposeConfig(project)
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
	// consumer := formatter.NewLogConsumer(ctx, os.Stdout, os.Stderr, true, true, true)
	err := t.DockerCompose().Up(ctx, project, api.UpOptions{
		Create: create,
		Start: api.StartOptions{
			Project:  project,
			AttachTo: attachTo.Elements(),
			// Attach:      consumer,
			CascadeStop: false, // 在容器停止时停止应用程序
			Wait:        false, // Wait函数将等待容器达到运行或健康状态，直到返回
			WaitTimeout: time.Minute * 5,
		},
	})
	return err
}

func (t *TangSengDaoDaoContext) dockerComposePs(project *types.Project) ([]api.ContainerSummary, error) {
	ctx := context.Background()
	return t.DockerCompose().Ps(ctx, t.opts.projectName, api.PsOptions{
		Project: project,
		All:     true,
	})
}

func (t *TangSengDaoDaoContext) getDockerProject(configs []string, services ...string) (*types.Project, error) {
	composeProjectOpts := composecmd.ProjectOptions{
		ProjectName: t.opts.projectName,
		ProjectDir:  t.opts.rootDir,
		ConfigPaths: configs,
	}

	project, err := composeProjectOpts.ToProject(services, cli.WithName(t.opts.projectName), cli.WithDotEnv, cli.WithOsEnv, cli.WithDiscardEnvFile, cli.WithResolvedPaths(true))
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (t *TangSengDaoDaoContext) dockerComposeStop(project *types.Project) error {
	ctx := context.Background()
	return t.DockerCompose().Stop(ctx, t.opts.projectName, api.StopOptions{
		Project: project,
	})
}

func (t *TangSengDaoDaoContext) dockerComposeRemove(project *types.Project) error {
	ctx := context.Background()
	return t.DockerCompose().Remove(ctx, t.opts.projectName, api.RemoveOptions{
		Project: project,
		Force:   true,
	})
}

func (t *TangSengDaoDaoContext) dockerComposeDown(project *types.Project) error {
	ctx := context.Background()
	return t.DockerCompose().Down(ctx, t.opts.projectName, api.DownOptions{
		Project: project,
	})
}

func (t *TangSengDaoDaoContext) dockerComposePull(project *types.Project) error {
	ctx := context.Background()
	return t.DockerCompose().Pull(ctx, project, api.PullOptions{
		Quiet: false,
	})
}

func (t *TangSengDaoDaoContext) dockerComposeConfig(project *types.Project) ([]byte, error) {
	ctx := context.Background()
	return t.DockerCompose().Config(ctx, project, api.ConfigOptions{
		Format: "yaml",
	})
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
