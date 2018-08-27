package service

const (
	CmdDependencyGet    CmdDependency = "get"
	CmdDependencyUpdate CmdDependency = "update"
	CmdDependencyReset  CmdDependency = "reset"

	GenImportFile  = "import-gen.yml"
	LockImportFile = "import-lock.yml"

	CacheRepository           = ".dependency/cache"
	CacheRepositoryConfigFile = "cache.yml"
)

var (
	excludedPaths = []string{
		"vendor",
	}

	movedPackages = []*PackageAction{}

	ignoredPackages = []string{
		"golang.org/x",
		"google.golang.org",
	}
)
