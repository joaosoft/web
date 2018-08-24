package dependency

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

	movedPackages = []*MovePackage{
		&MovePackage{old: "gopkg.in/yaml.v2", new: "github.com/go-yaml/yaml"},
		&MovePackage{old: "golang.org/x", new: "github.com/golang"},
		&MovePackage{old: "google.golang.org", new: "github.com/golang"},
		&MovePackage{old: "gopkg.in", new: "github.com/golang"},
	}
)
