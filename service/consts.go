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

	movedPackages = []*PackageAction{
		{old: "gopkg.in/yaml.v2", new: "github.com/go-yaml/yaml"},
		{old: "golang.org/x", new: "github.com/golang"},
		{old: "google.golang.org", new: "github.com/golang"},
		{old: "gopkg.in", new: "github.com/golang"},
	}

	ignoredPackages = []string{
		"golang.org/x",
		"google.golang.org",
	}
)
