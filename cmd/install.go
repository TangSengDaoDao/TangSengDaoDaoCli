package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/TangSengDaoDao/TangSengDaoDaoCli/pkg/util"
	"github.com/spf13/cobra"
)

type installCMD struct {
	ctx *TangSengDaoDaoContext

	externalIP string // 外网IP
	mysqlPwd   string // mysql密码
	minioPwd   string // minio密码
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
	installCmd.Flags().StringVar(&i.externalIP, "ip", "", "external ip （外网IP）")

	err := os.MkdirAll(i.configDir(), 0755)
	if err != nil {
		panic(err)
	}

	i.mysqlPwd = util.GenUUID()[:16]
	i.minioPwd = util.GenUUID()

	return installCmd
}

func (i *installCMD) run(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(i.externalIP) == "" {
		ips, _ := util.GetIntranetIP()
		if len(ips) > 0 {
			i.externalIP = ips[0]
		}
	}
	var err error
	// ==================== 下载docker-compose.yaml文件 ====================
	if !i.existDockerCompose() {
		err = i.downloadDockerCompose(cmd)
		if err != nil {
			cmd.Println("download docker-compose.yaml file error:", err)
			return err
		}
	}
	// ==================== 下载.env文件 ====================
	if !i.existDotEnv() {
		// 下载.env文件
		err = i.downloadDotEnv(cmd)
		if err != nil {
			cmd.Println("download .env file error:", err)
			return err
		}
	}

	// ==================== 下载悟空IM的配置文件 ====================
	if !i.existWkConfig() {
		// 下载悟空IM的配置文件
		err = i.downloadWkConfig(cmd)
		if err != nil {
			cmd.Println("download wk.yaml file error:", err)
			return err
		}
	}
	// ==================== 下载唐僧叨叨的配置文件 ====================
	if !i.existTsddConfig() {
		// 下载唐僧叨叨的配置文件
		err = i.downloadTsddConfig(cmd)
		if err != nil {
			cmd.Println("download tsdd.yaml file error:", err)
			return err
		}
	}
	// ==================== 下载Caddy的配置文件 ====================
	if !i.existCaddyConfig() {
		// 下载Caddy的配置文件
		err = i.downloadCaddyConfig(cmd)
		if err != nil {
			cmd.Println("download Caddy file error:", err)
			return err
		}
	}
	// ==================== 取代变量 ====================
	// 替换.env文件中的变量
	err = i.replaceDotEnvVarz()
	if err != nil {
		return err
	}
	// 替换悟空IM的配置文件中的变量
	err = i.replaceWkConfigVarz()
	if err != nil {
		return err
	}
	// 替换唐僧叨叨的配置文件中的变量
	err = i.replaceTsddConfigVarz()
	if err != nil {
		return err
	}
	// caddy配置
	err = i.replaceCaddyVarz()
	if err != nil {
		return err
	}

	fmt.Println("Install at ", i.ctx.opts.rootDir)
	fmt.Println("Install success! please run 'tsdd start' to start service.")

	return nil
}

// 取代.env文件中的变量
func (i *installCMD) replaceDotEnvVarz() error {
	dotEnvPath := i.dotEnvPath()
	content, err := ioutil.ReadFile(dotEnvPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	// minio
	contentStr = strings.ReplaceAll(contentStr, "#MINIO_ROOT_PASSWORD#", i.minioPwd)
	contentStr = strings.ReplaceAll(contentStr, "#MINIO_SERVER_URL#", "")
	contentStr = strings.ReplaceAll(contentStr, "#MINIO_BROWSER_REDIRECT_URL#", "")
	// mysql
	contentStr = strings.ReplaceAll(contentStr, "#MYSQL_ROOT_PASSWORD#", i.mysqlPwd)
	// web
	contentStr = strings.ReplaceAll(contentStr, "#API_URL#", fmt.Sprintf("http://%s:8090/", i.externalIP))
	err = ioutil.WriteFile(dotEnvPath, []byte(contentStr), 0644)
	return err
}

// 取代悟空IM的配置文件中的变量
func (i *installCMD) replaceWkConfigVarz() error {
	confPath := i.wkConfigPath()
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "#EXTERNAL_IP#", i.externalIP)

	err = ioutil.WriteFile(confPath, []byte(contentStr), 0644)

	return err
}

func (i *installCMD) replaceTsddConfigVarz() error {
	confPath := i.tsddConfigPath()
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "#EXTERNAL_IP#", i.externalIP)
	contentStr = strings.ReplaceAll(contentStr, "#MYSQL_ROOT_PASSWORD#", i.mysqlPwd)
	contentStr = strings.ReplaceAll(contentStr, "#MINIO_ROOT_PASSWORD#", i.minioPwd)

	err = ioutil.WriteFile(confPath, []byte(contentStr), 0644)
	return err
}

func (i *installCMD) replaceCaddyVarz() error {
	confPath := i.caddyConfigPath()
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "#APIAddr#", fmt.Sprintf("%s:8090", i.externalIP))

	err = ioutil.WriteFile(confPath, []byte(contentStr), 0644)
	return err
}

// 下载docker-compose.yaml文件
func (i *installCMD) downloadDockerCompose(cmd *cobra.Command) error {
	// 下载文件
	destPath := i.dockerComposePath()
	return ioutil.WriteFile(destPath, []byte(i.ctx.DockerComposeYaml), 0644)
}
func (i *installCMD) existDockerCompose() bool {
	return i.existFile(i.dockerComposePath())
}

// 下载.env文件
func (i *installCMD) downloadDotEnv(cmd *cobra.Command) error {
	// 下载文件
	destPath := i.dotEnvPath()
	return ioutil.WriteFile(destPath, []byte(i.ctx.DotEnv), 0644)
}
func (i *installCMD) existDotEnv() bool {
	return i.existFile(i.dotEnvPath())
}

// 下载悟空IM的配置文件
func (i *installCMD) downloadWkConfig(cmd *cobra.Command) error {

	wkContentBytes, err := i.ctx.Configs.ReadFile("configs/wk.yaml")
	if err != nil {
		return err
	}
	destPath := path.Join(i.configDir(), "wk.yaml")

	return ioutil.WriteFile(destPath, []byte(wkContentBytes), 0644)
}
func (i *installCMD) existWkConfig() bool {
	return i.existFile(i.wkConfigPath())
}

// 下载唐僧叨叨的配置文件
func (i *installCMD) downloadTsddConfig(cmd *cobra.Command) error {
	tsddContentBytes, err := i.ctx.Configs.ReadFile("configs/tsdd.yaml")
	if err != nil {
		return err
	}
	destPath := i.tsddConfigPath()
	return ioutil.WriteFile(destPath, []byte(tsddContentBytes), 0644)
}

func (i *installCMD) existTsddConfig() bool {
	return i.existFile(i.tsddConfigPath())
}

// 下载Caddy配置文件
func (i *installCMD) downloadCaddyConfig(cmd *cobra.Command) error {

	wkContentBytes, err := i.ctx.Configs.ReadFile("configs/Caddyfile")
	if err != nil {
		return err
	}
	destPath := i.caddyConfigPath()

	return ioutil.WriteFile(destPath, []byte(wkContentBytes), 0644)
}
func (i *installCMD) existCaddyConfig() bool {
	return i.existFile(i.caddyConfigPath())
}

func (i *installCMD) dotEnvPath() string {
	return path.Join(i.ctx.opts.rootDir, ".env")
}
func (i *installCMD) dockerComposePath() string {
	return i.ctx.opts.dockerComposePath
}
func (i *installCMD) wkConfigPath() string {
	return path.Join(i.configDir(), "wk.yaml")
}
func (i *installCMD) tsddConfigPath() string {
	return path.Join(i.configDir(), "tsdd.yaml")
}
func (i *installCMD) caddyConfigPath() string {
	return path.Join(i.configDir(), "Caddyfile")
}

// 获取配置目录
func (i *installCMD) configDir() string {

	return path.Join(i.ctx.opts.rootDir, "configs")
}

func (i *installCMD) existFile(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
