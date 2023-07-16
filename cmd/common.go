package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	terminal "golang.org/x/term"
)

type CMD interface {
	CMD() *cobra.Command
}

type Options struct {
	ServerAddr        string
	Description       string
	Token             string
	rootDir           string // 唐僧叨叨的根目录
	projectName       string
	dockerComposePath string // docker compose yaml 文件路径

	Version    string // version
	Commit     string // git commit id
	CommitDate string // git commit date
	TreeState  string // git tree state
}

func NewOptions() *Options {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opts := &Options{
		projectName: "tsdd",
		ServerAddr:  "http://127.0.0.1:5001",
		Description: "",
	}
	opts.rootDir = filepath.Join(homeDir, opts.projectName)
	// 创建唐僧叨叨的根目录
	err = os.MkdirAll(opts.rootDir, 0700)
	if err != nil {
		panic(err)
	}
	// 创建唐僧叨叨的配置目录
	err = os.MkdirAll(opts.ContextDir(), 0700)
	if err != nil {
		panic(err)
	}
	opts.dockerComposePath = filepath.Join(opts.rootDir, "docker-compose.yaml")
	return opts
}

func (o *Options) ContextPath(name string) (string, error) {
	if !validName(name) {
		return "", fmt.Errorf("invalid context name %q", name)
	}
	return filepath.Join(o.ContextDir(), name+".json"), nil
}

func (o *Options) ContextDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	if u.HomeDir == "" {
		return ""
	}
	return filepath.Join(o.rootDir, ".config", "context")

}

func (o *Options) Load() error {
	data, err := ioutil.ReadFile(o.metaFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	name := string(data)
	filen, err := o.ContextPath(name)
	if err != nil {
		return err
	}

	optionData, err := ioutil.ReadFile(filen)
	if err != nil {
		return err
	}
	var optionMap map[string]interface{}
	err = json.Unmarshal(optionData, &optionMap)
	if err != nil {
		return err
	}
	if optionMap == nil {
		return nil
	}
	if optionMap["url"] != nil {
		o.ServerAddr = optionMap["url"].(string)
	}
	if optionMap["description"] != nil {
		o.Description = optionMap["description"].(string)
	}
	if optionMap["token"] != nil {
		o.Token = optionMap["token"].(string)
	}
	return nil
}

func (o *Options) Save(name string) error {
	p, err := o.ContextPath(name)
	if err != nil {
		return err
	}
	j, err := json.MarshalIndent(map[string]interface{}{
		"url":         o.ServerAddr,
		"description": o.Description,
		"token":       o.Token,
	}, "", "  ")
	if err != nil {
		return err
	}

	pf, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	_, err = pf.Write(j)
	if err != nil {
		return err
	}

	mf, err := os.OpenFile(o.metaFile(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	_, err = mf.Write([]byte(name))
	return err
}

func (o *Options) metaFile() string {
	return filepath.Join(o.ContextDir(), "meta")
}

func progressWidth() int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	minWidth := 10

	if w-30 <= minWidth {
		return minWidth
	} else {
		return w - 30
	}
}

func move(oldPath, newPath string) error {
	srcFile, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	err = os.Remove(oldPath)
	if err != nil {
		return err
	}
	return nil
}
