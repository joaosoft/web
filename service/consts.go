package service

const (
	CmdDependencyGet    CmdDependency = "get"
	CmdDependencyUpdate CmdDependency = "update"
	CmdDependencyReset  CmdDependency = "reset"

	GenImportFile  = "import-gen.yml"
	LockImportFile = "import-lock.yml"

	CacheRepository           = ".dependency/cache"
	CacheRepositoryConfigFile = "cache.yml"

	RegexForVendorFiles = `^vendor_[0-9]{14}$`

	ProtocolSSH   Protocol = "ssh"
	ProtocolHTTPS Protocol = "https"
)

var (
	excludedPaths = []string{
		"vendor",
	}

	movedPackages = []*PackageAction{}

	ignoredPackages = []string{
		"golang.org/x",
		"google.golang.org",
		"../",
	}
)
