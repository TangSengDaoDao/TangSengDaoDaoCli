package main

import (
	"embed"
	_ "embed"

	"github.com/TangSengDaoDao/TangSengDaoDaoCli/cmd"
)

//go:embed docker-compose.yaml
var dockerComposeYaml string

//go:embed .env
var dotEnv string

//go:embed configs
var configs embed.FS

// go ldflags
var Version string    // version
var Commit string     // git commit id
var CommitDate string // git commit date
var TreeState string  // git tree state

func main() {
	tsdd := cmd.NewTangSengDaoDao()
	tsdd.Context().DockerComposeYaml = dockerComposeYaml
	tsdd.Context().Configs = configs
	tsdd.Context().DotEnv = dotEnv
	tsdd.Options().Version = Version
	tsdd.Options().Commit = Commit
	tsdd.Options().CommitDate = CommitDate
	tsdd.Options().TreeState = TreeState
	tsdd.Execute()
}
