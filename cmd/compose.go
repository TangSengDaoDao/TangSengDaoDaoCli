package cmd

type ProjectOptions struct {
	ProjectName   string
	Profiles      []string
	ConfigPaths   []string
	WorkDir       string
	ProjectDir    string
	EnvFiles      []string
	Compatibility bool
}
