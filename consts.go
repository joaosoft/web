package dependency

const (
	CmdDependencyGet   CmdDependency = "get"
	CmdDependencyReset CmdDependency = "reset"

	GenImportFile  = "import-gen.yml"
	LockImportFile = "import-lock.yml"

	CacheRepository           = "/tmp/dependency/cache"
	CacheRepositoryConfigFile = "cache.yml"
)

var (
	excludedPaths = []string{
		"vendor",
	}

	excludedImports = []string{
		"golang.org/x",
		"google.golang.org",
		"github.com/golang",
	}

	movedPackages = map[string]string{
		"golang.org/x":      "github.com/golang",
		"google.golang.org": "github.com/golang",
	}
)
