package dependency

import (
	"fmt"

	"os"

	"io/ioutil"

	"os/exec"

	"github.com/joaosoft/logger"
	yaml "gopkg.in/yaml.v2"
)

type Vcs struct {
	cache  *Cache
	logger logger.ILogger
}

func NewVcs(path string, config string, logger logger.ILogger) (*Vcs, error) {
	vcs := &Vcs{
		cache: &Cache{
			imports: make(Imports),
			path:    path,
			config:  fmt.Sprintf("%s/%s", path, config),
		},
		logger: logger,
	}

	if err := vcs.StartCache(); err != nil {
		return nil, err
	}

	return vcs, nil
}

func (v *Vcs) StartCache() error {
	v.logger.Debugf("executing Start Cache")

	if _, err := os.Stat(v.cache.config); err == nil {

		if bytes, err := ioutil.ReadFile(v.cache.config); err != nil {
			return v.logger.Errorf("error reading file [%s] %s", LockImportFile, err).ToError()
		} else {
			if err := yaml.Unmarshal(bytes, &v.cache.imports); err != nil {
				return v.logger.Errorf("error unmarshal file [%s] %s", LockImportFile, err).ToError()
			}
			return nil
		}
	} else {
		// create file
		if err := os.MkdirAll(v.cache.path, os.ModePerm); err != nil {
			return v.logger.Errorf("error creating folder [%s] %s", v.cache.path, err).ToError()
		}

		newFile, err := os.Create(v.cache.config)
		if err != nil {
			return v.logger.Errorf("error creating file [%s] %s", v.cache.config, err).ToError()
		}
		newFile.Close()
	}

	return nil
}

func (v *Vcs) ClearCache(dir string) error {
	v.logger.Debugf("executing Clear Cache of [%s]", dir)

	if _, err := os.Stat(dir); err != nil {
		os.Remove(dir)
	}

	return v.StartCache()
}

func (v *Vcs) Clone(imprt *Import, copyTo string) error {

	var gitArgs []string
	v.logger.Debugf("executing Clone for [%s]", imprt.internal.repo.path)

	// git checkout tags/v1.0
	version := imprt.Branch
	if imprt.Version != "" {
		version = imprt.Version
	}

	key := fmt.Sprintf("%s/%s", imprt.internal.repo.path, version)
	pathCachedRepo := fmt.Sprintf("%s/%s/%s", v.cache.path, imprt.internal.repo.path, version)

	if _, ok := v.cache.imports[key]; !ok {

		// remove cached temporary folder to prevent errors
		os.Remove(key)

		v.logger.Infof("downloading repository with ssh protocol [%s] to [%s]", imprt.internal.repo.ssh, pathCachedRepo)

		if _, err := exec.Command("git", "ls-remote", "-h", imprt.internal.repo.ssh).CombinedOutput(); err != nil {
			v.logger.Infof("the repository doesn't exist [%s]", imprt.internal.repo.ssh)
			return err
		}

		gitArgs = []string{
			"clone",
			"--recursive",
			"-v",
			"--progress",
			"--depth", "1",
			"--shallow-submodules",
			"--branch",
			imprt.Branch,
			imprt.internal.repo.ssh,
			pathCachedRepo,
		}
		if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
			v.logger.Errorf("error executing git clone command %s", string(stderr))

			os.Remove(key)
			v.logger.Infof("retrying download with https protocol [%s] to [%s]", imprt.internal.repo.https, pathCachedRepo)

			gitArgs = []string{
				"clone",
				"--recursive",
				"-v",
				"--progress",
				"--depth", "1",
				"--shallow-submodules",
				"--branch",
				imprt.Branch,
				imprt.internal.repo.https,
				pathCachedRepo,
			}
			if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
				os.Remove(key)
				return v.logger.Errorf("error executing git clone command %s", string(stderr)).ToError()
			}
		}

		// git checkout tags/v1.0
		if imprt.Version != "" {
			v.logger.Infof("checkout version [%s]", imprt.Version)
			gitArgs = []string{
				"checkout",
				imprt.Version,
			}
			if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
				v.logger.Errorf("error executing [git checkout tags/%s]  command %s", imprt.Version, string(stderr))
			}
		}

		v.logger.Infof("git clone completed for [%s]", imprt.internal.repo.path)
		v.cache.imports[key] = imprt

		v.SaveCache()
	} else {
		v.logger.Infof("pull version [%s]", version)
		gitArgs = []string{
			"pull",
		}
		if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
			v.logger.Errorf("error executing [git pull]  command %s", imprt.Version, string(stderr))
		}
	}

	v.logger.Infof("copying import [%s] from cache", imprt.internal.repo.path)

	if err := CopyDir(fmt.Sprintf("%s%s", pathCachedRepo, imprt.internal.packag), imprt.internal.vendor); err != nil {
		return v.logger.Errorf("error executing Copying import [%s] to vendor [%s] %s", imprt.internal.repo.path, imprt.internal.vendor, err).ToError()
	}

	return nil
}

func (v *Vcs) SaveCache() error {
	v.logger.Debugf("executing Save Cache")

	os.Remove(v.cache.config)

	if bytes, err := yaml.Marshal(v.cache.imports); err != nil {
		return v.logger.Errorf("error marshal imports %s", err).ToError()
	} else {
		if err := ioutil.WriteFile(v.cache.config, bytes, 0644); err != nil {
			return v.logger.Errorf("error writing file [%s] %s", v.cache.config, err).ToError()
		}
	}

	return nil
}
