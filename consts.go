package main

const (
	CmdDependencyGet    CmdDependency = "get"
	CmdDependencyUpdate CmdDependency = "update"
	CmdDependencyReset  CmdDependency = "reset"

	GenImportFile  = "import-gen.yml"
	LockImportFile = "import-lock.yml"

	CacheRepository           = ".dependency/cache"
	CacheRepositoryConfigFile = "cache.yml"

	RegexForVendorFiles = `^_vendor_[0-9]{14}$`

	ProtocolSSH   Protocol = "ssh"
	ProtocolHTTPS Protocol = "https"
)

var (
	excludedPaths = []string{
		"vendor",
	}

	packageActions = []*PackageAction{
		{old: "golang.org/x/net", new: "go.googlesource.com/net"},
		{old: "golang.org/x/exp", new: "go.googlesource.com/exp"},
	}

	ignoredPackages = []string{
		"../",
	}
)
