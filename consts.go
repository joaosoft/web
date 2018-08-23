package dependency

const (
	CmdDependencyGet   CmdDependency = "get"
	CmdDependencyReset CmdDependency = "reset"

	GenImportFileName = "import-gen.yml"
	LockImportFileName = "import-lock.yml"
)

var (
	excludedPaths = []string{
		"vendor",
	}

	excludedImports = []string{
		"golang.org/x",
		"google.golang.org",
	}

	movedPackages = map[string]string{
		"golang.org/x":      "github.com/golang",
		"google.golang.org": "github.com/golang",
	}
)
