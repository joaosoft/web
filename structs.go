package dependency

type CmdDependency string
type Imports map[string]*Import

type sync struct {
	internalImports  Imports
	externalImports  Imports
	loadedImports    map[string]bool
	installedImports Imports
}

type Import struct {
	Branch   string   `json:"branch,omitempty" yaml:"branch,omitempty"`
	Package  []string `json:"package,omitempty" yaml:"package,omitempty"`
	Revision string   `json:"revision,omitempty" yaml:"revision,omitempty"`
	Version  string   `json:"version,omitempty" yaml:"version,omitempty"`
	internal internal
}

type internal struct {
	host    string
	user    string
	project string
	packag  string
	repo    repo
	vendor  string
}

type repo struct {
	https string
	ssh   string
	path  string
}

type Cache struct {
	imports Imports
	path    string
	config  string
}
