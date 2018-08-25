package service

import (
	"fmt"

	"os"

	"io/ioutil"

	"os/exec"

	"strings"

	"sync"

	"github.com/joaosoft/logger"
	"github.com/go-yaml/yaml"
)

type Vcs struct {
	cache  *Cache
	logger logger.ILogger
	mux    sync.Mutex
}

func NewVcs(path string, config string, logger logger.ILogger) (*Vcs, error) {
	vcs := &Vcs{
		cache: &Cache{
			imports: make(Imports),
			path:    path,
			config:  fmt.Sprintf("%s/%s", path, config),
		},
		logger: logger,
		mux:    sync.Mutex{},
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

func (v *Vcs) CopyDependency(imprt *Import, copyTo string, update bool) error {

	v.logger.Debugf("executing CopyDependency for [%s]", imprt.internal.repo.path)

	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	// locking for repository operations
	v.mux.Lock()
	defer v.mux.Unlock()

	if _, ok := v.cache.imports[imprt.internal.repo.save]; !ok {
		if err := v.Clone(imprt); err != nil {
			return err
		}
	} else {
		branch := imprt.Branch
		if imprt.Version != "" {
			branch = imprt.Version
		}

		if branch == "" {
			if latestVersion, err := v.GetLatestVersion(pathCachedRepo); err != nil {
				return err
			} else {
				imprt.Version = latestVersion
				if err := v.Checkout(imprt); err != nil {
					return err
				}
			}
		}

		if update {
			if err := v.Pull(imprt); err != nil {
				return err
			}
		}
	}

	if err := v.doUpdateRepositoryInfo(imprt); err != nil {
		return err
	}

	v.logger.Infof("copying import [%s%s] from cache", imprt.internal.repo.path, imprt.internal.repo.packag)

	pkgCopyTo := fmt.Sprintf("%s%s", imprt.internal.repo.vendor, imprt.internal.repo.packag)
	if _, err := os.Stat(pkgCopyTo); err == nil {
		v.logger.Infof("package already copied to vendor [%s]", pkgCopyTo)
		return nil
	}
	if err := CopyDir(fmt.Sprintf("%s%s", pathCachedRepo, imprt.internal.repo.packag), fmt.Sprintf("%s%s", imprt.internal.repo.vendor, imprt.internal.repo.packag)); err != nil {
		return v.logger.Errorf("error executing Copying import [%s] to vendor [%s] %s", imprt.internal.repo.path, imprt.internal.repo.vendor, err).ToError()
	}

	return nil
}

func (v *Vcs) Clone(imprt *Import) error {

	var gitArgs []string
	v.logger.Debugf("executing Clone for [%s]", imprt.internal.repo.path)

	branch := imprt.Branch
	if imprt.Version != "" {
		branch = imprt.Version
	}

	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	// remove cached temporary folder to prevent errors
	os.Remove(imprt.internal.repo.save)

	v.logger.Infof("ping repository with ssh protocol [%s]", imprt.internal.repo.ssh)
	if err := exec.Command("git", "ls-remote", "-h", imprt.internal.repo.ssh).Run(); err != nil {
		v.logger.Infof("the repository doesn't exist [%s]", imprt.internal.repo.ssh)
		return err
	}

	v.logger.Infof("downloading repository with ssh protocol [%s] to [%s]", imprt.internal.repo.ssh, pathCachedRepo)

	gitArgs = []string{
		"clone",
		"--recursive",
		"-v",
		"--progress",
		"--depth", "1",
		"--shallow-submodules",
	}
	if branch != "" {
		gitArgs = append(gitArgs, "--branch", branch)
	}
	gitArgs = append(gitArgs, imprt.internal.repo.https, pathCachedRepo)

	if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
		v.logger.Errorf("error executing git clone command %s", string(stderr))

		os.Remove(imprt.internal.repo.save)
		v.logger.Infof("retrying download with https protocol [%s] to [%s]", imprt.internal.repo.https, pathCachedRepo)

		gitArgs = []string{
			"clone",
			"--recursive",
			"-v",
			"--progress",
			"--depth", "1",
			"--shallow-submodules",
		}
		if branch != "" {
			gitArgs = append(gitArgs, "--branch", branch)
		}
		gitArgs = append(gitArgs, imprt.internal.repo.https, pathCachedRepo)

		if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
			os.Remove(imprt.internal.repo.save)
			return v.logger.Errorf("error executing git clone command %s", string(stderr)).ToError()
		}
	}

	if err := v.doUpdateRepositoryInfo(imprt); err != nil {
		return err
	}

	v.logger.Infof("git clone completed for [%s]", imprt.internal.repo.save)
	v.cache.imports[imprt.internal.repo.save] = imprt

	v.SaveCache()

	return nil
}

func (v *Vcs) Pull(imprt *Import) error {
	v.logger.Debugf("executing Pull")

	v.logger.Infof("updating repository [%s]", imprt.internal.repo.path)
	gitArgs := []string{
		"pull",
	}
	if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git pull] command %s", string(stderr)).ToError()
	}

	return nil
}

func (v *Vcs) Checkout(imprt *Import) error {
	v.logger.Debugf("executing Checkout")

	branch := imprt.Branch
	if imprt.Version != "" {
		branch = imprt.Version
	}

	v.logger.Infof("checkout repository [%s]", imprt.internal.repo.path)
	gitArgs := []string{
		"checkout",
		branch,
	}
	if stderr, err := exec.Command("git", gitArgs...).CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git checkout] command %s", string(stderr)).ToError()
	}

	if err := v.doUpdateRepositoryInfo(imprt); err != nil {
		return err
	}

	return nil
}

func (v *Vcs) doUpdateRepositoryInfo(imprt *Import) error {
	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	imprt.Branch, _ = v.GetBranch(pathCachedRepo)
	imprt.Revision, _ = v.GetRevision(pathCachedRepo)
	imprt.Branch, _ = v.GetVersion(pathCachedRepo)

	return nil
}

func (v *Vcs) GetBranch(path string) (string, error) {
	v.logger.Debugf("executing Get Branch")

	gitArgs := []string{
		"rev-parse",
		"--abbrev-ref",
		"--abbrev-ref",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		v.logger.Infof("error getting git branch: %s", string(stderr))
		return "", nil
	} else {
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) GetLatestVersion(path string) (string, error) {
	v.logger.Debugf("executing Get Latest Version")

	gitArgs := []string{
		"describe",
		"--tags",
		"`git rev-list --tags --max-count=1`",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		v.logger.Infof("returned [%s] getting latest version", string(stderr))
		return "", nil
	} else {
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) GetVersion(path string) (string, error) {
	v.logger.Debugf("executing Get Version")

	gitArgs := []string{
		"describe",
		"--tags",
		"--abbrev=0",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		v.logger.Infof("error getting git version: %s", string(stderr))
		return "", nil
	} else {
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) GetRevision(path string) (string, error) {
	v.logger.Debugf("executing Get Revision")

	gitArgs := []string{
		"rev-parse",
		"HEAD",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		v.logger.Infof("error getting git revision: %s", string(stderr))
		return "", nil
	} else {
		rtn := strings.TrimSpace(string(stderr))
		if rtn == "" {
			return "", nil
		}
		return strings.TrimSpace(string(stderr)), nil
	}
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
