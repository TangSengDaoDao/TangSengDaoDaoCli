package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/TangSengDaoDao/TangSengDaoDaoCli/pkg/util"
	"github.com/spf13/cobra"
)

type installCMD struct {
	ctx                  *TangSengDaoDaoContext
	dockerComposeFileUrl string // docker-compose.yaml文件的下载地址
	dotEnvFileUrl        string // .env文件的下载地址
	wkConfigUrl          string // 悟空IM的配置文件的下载地址
	tsddConfigUrl        string // 唐僧叨叨的配置文件的下载地址

	externalIP string // 外网IP
	mysqlPwd   string // mysql密码
	minioPwd   string // minio密码
}

func newInstallCMD(ctx *TangSengDaoDaoContext) *installCMD {
	baseURL := "https://gitee.com/TangSengDaoDao/TangSengDaoDaoCli/raw/main"
	c := &installCMD{
		ctx:                  ctx,
		dockerComposeFileUrl: fmt.Sprintf("%s/docker-compose.yaml", baseURL),
		dotEnvFileUrl:        fmt.Sprintf("%s/.env", baseURL),
		wkConfigUrl:          fmt.Sprintf("%s/configs/wk.yaml", baseURL),
		tsddConfigUrl:        fmt.Sprintf("%s/configs/tsdd.yaml", baseURL),
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

	i.mysqlPwd = util.GenUUID()[:16]
	i.minioPwd = util.GenUUID()

	return installCmd
}

func (i *installCMD) run(cmd *cobra.Command, args []string) error {
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

	// ==================== 启动服务 ====================
	return i.ctx.DockerComposeUp([]string{"docker-compose.yaml"})
}

// 取代.env文件中的变量
func (i *installCMD) replaceDotEnvVarz() error {
	dotEnvPath := path.Join(i.configDir(), ".env")
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
	err = ioutil.WriteFile(dotEnvPath, []byte(contentStr), 0644)
	return err
}

// 取代悟空IM的配置文件中的变量
func (i *installCMD) replaceWkConfigVarz() error {
	confPath := path.Join(i.configDir(), "configs", "wk.yaml")
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
	confPath := path.Join(i.configDir(), "configs", "tsdd.yaml")
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

// 下载docker-compose.yaml文件
func (i *installCMD) downloadDockerCompose(cmd *cobra.Command) error {
	cmd.Println("download docker-compose.yaml file from ", i.dockerComposeFileUrl)
	// 下载文件
	tmpDir, _ := ioutil.TempDir("", "tsddcli")
	tmpPath := path.Join(tmpDir, "docker-compose.yaml.tmp")
	destPath := path.Join(i.configDir(), "docker-compose.yaml")
	err := i.download(i.dockerComposeFileUrl, tmpPath)
	if err != nil {
		return err
	}
	return move(tmpPath, destPath)
}
func (i *installCMD) existDockerCompose() bool {
	return i.existFile(path.Join(i.configDir(), "docker-compose.yaml"))
}

// 下载.env文件
func (i *installCMD) downloadDotEnv(cmd *cobra.Command) error {
	cmd.Println("download .env file from ", i.dotEnvFileUrl)
	// 下载文件
	tmpDir, _ := ioutil.TempDir("", "tsddcli")
	tmpPath := path.Join(tmpDir, ".env.tmp")
	destPath := path.Join(i.configDir(), ".env")
	err := i.download(i.dotEnvFileUrl, tmpPath)
	if err != nil {
		return err
	}
	return move(tmpPath, destPath)
}
func (i *installCMD) existDotEnv() bool {
	return i.existFile(path.Join(i.configDir(), ".env"))
}

// 下载悟空IM的配置文件
func (i *installCMD) downloadWkConfig(cmd *cobra.Command) error {
	cmd.Println("download wk.yaml file from ", i.wkConfigUrl)
	// 下载文件
	tmpDir, _ := ioutil.TempDir("", "tsddcli")
	tmpPath := path.Join(tmpDir, ".wk.yaml.tmp")
	destPath := path.Join(i.configDir(), "wk.yaml")
	err := i.download(i.wkConfigUrl, tmpPath)
	if err != nil {
		return err
	}
	return move(tmpPath, destPath)
}
func (i *installCMD) existWkConfig() bool {
	return i.existFile(path.Join(i.configDir(), "wk.yaml"))
}

// 下载唐僧叨叨的配置文件
func (i *installCMD) downloadTsddConfig(cmd *cobra.Command) error {
	cmd.Println("download tsdd.yaml file from ", i.tsddConfigUrl)
	// 下载文件
	tmpDir, _ := ioutil.TempDir("", "tsddcli")
	tmpPath := path.Join(tmpDir, ".tsdd.yaml.tmp")
	destPath := path.Join(i.configDir(), "tsdd.yaml")
	err := i.download(i.tsddConfigUrl, tmpPath)
	if err != nil {
		return err
	}
	return move(tmpPath, destPath)
}
func (i *installCMD) existTsddConfig() bool {
	return i.existFile(path.Join(i.configDir(), "tsdd.yaml"))
}

// 获取配置目录
func (i *installCMD) configDir() string {

	return path.Join(i.ctx.opts.rootDir, "configs")
}

// 下载文件
func (i *installCMD) download(url string, destPath string) error {
	client := http.DefaultClient
	client.Timeout = 60 * 10 * time.Second
	reps, err := client.Get(url)
	if err != nil {
		return err
	}
	defer reps.Body.Close()
	if reps.StatusCode != http.StatusOK {
		return fmt.Errorf("http status[%d] is error", reps.StatusCode)
	}
	//保存文件
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close() //关闭文件
	return nil
}

func (i *installCMD) existFile(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
