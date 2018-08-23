package dependency

type CmdDependency string
type Imports map[string]Import

type Import struct {
	Branch   string   `json:"branch,omitempty" yaml:"branch,omitempty"`
	Package  []string `json:"package,omitempty" yaml:"package,omitempty"`
	Revision string   `json:"revision,omitempty" yaml:"revision,omitempty"`
	Version  string   `json:"version,omitempty" yaml:"version,omitempty"`
	internal Internal
}

type Internal struct {
	host    string
	user    string
	project string
	repo    Repo
	vendor  string
}

type Repo struct {
	https string
	ssh   string
}
