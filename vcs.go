package main

import (
	"fmt"

	"os"

	"io/ioutil"

	"os/exec"

	"strings"

	"sync"

	"time"

	"github.com/joaosoft/logger"
	"gopkg.in/yaml.v2"
)

type Vcs struct {
	cache    *Cache
	logger   logger.ILogger
	mux      sync.Mutex
	protocol Protocol
}

func NewVcs(path string, config string, protocol Protocol, logger logger.ILogger) (*Vcs, error) {
	vcs := &Vcs{
		cache: &Cache{
			imports: make(Imports),
			path:    path,
			config:  fmt.Sprintf("%s/%s", path, config),
		},
		protocol: protocol,
		logger:   logger,
		mux:      sync.Mutex{},
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

func (v *Vcs) CopyDependency(sync *Memory, imprt *Import, copyTo string, update bool) error {

	v.logger.Debugf("executing Copy Dependency for [%s]", imprt.internal.repo.path)

	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	// locking for repository operations
	v.mux.Lock()
	defer v.mux.Unlock()

	if _, ok := v.cache.imports[imprt.internal.repo.save]; !ok {
		if err := v.Clone(imprt); err != nil {
			return err
		}

		if imprt.Version == "" && imprt.Revision == "" {
			// update branch
			imprt.Branch, _ = v.GetBranch(pathCachedRepo)
		}
	} else {
		if v.cache.imports[imprt.internal.repo.save].Branch != imprt.Branch ||
			v.cache.imports[imprt.internal.repo.save].Version != imprt.Version ||
			v.cache.imports[imprt.internal.repo.save].Revision != imprt.Revision {

			// if don't have gen file and has cached content, get the default branch
			if imprt.Branch == "" && imprt.Version == "" && imprt.Revision == "" {
				if defaultBranch, err := v.GetDefaultBranch(pathCachedRepo); err != nil {
					// if i'm in a revision i can't get the default branch
					imprt.Branch = "master"
					if err := v.Checkout(pathCachedRepo, imprt); err != nil {
						return err
					}
				} else {
					// get the default branch
					defaultBranch, _ = v.GetDefaultBranch(pathCachedRepo)
					if imprt.Branch != defaultBranch {
						imprt.Branch = defaultBranch
						if err := v.Checkout(pathCachedRepo, imprt); err != nil {
							return err
						}
					}
				}

			} else {
				if imprt.Version != "" {
					if err := v.FetchAllTags(pathCachedRepo); err != nil {
						return err
					}
				}

				if err := v.Checkout(pathCachedRepo, imprt); err != nil {
					return err
				}
			}

			if update {
				if _, ok := sync.lockedImports[imprt.internal.repo.path]; !ok {
					if err := v.Pull(pathCachedRepo, imprt); err != nil {
						return err
					}
				}
			}
		}

	}

	// update import
	if err := v.doUpdateImportInfo(imprt); err != nil {
		return err
	}

	// update cache
	v.cache.imports[imprt.internal.repo.save] = imprt

	// save cache
	if err := v.SaveCache(); err != nil {
		return err
	}

	pkgCopyTo := fmt.Sprintf("%s%s", imprt.internal.repo.vendor, imprt.internal.repo.packag)
	if _, err := os.Stat(pkgCopyTo); err == nil {
		v.logger.Infof("package already copied to vendor [%s]", pkgCopyTo)
		return nil
	}

	pkgCopyFrom := fmt.Sprintf("%s%s", pathCachedRepo, imprt.internal.repo.packag)
	if _, err := os.Stat(pkgCopyFrom); err == nil {
		v.logger.Infof("copying import [%s%s] from cache", imprt.internal.repo.path, imprt.internal.repo.packag)

		if err := CopyDir(pkgCopyFrom, pkgCopyTo); err != nil {
			return v.logger.Errorf("error executing copy of import [%s] to vendor [%s] %s", imprt.internal.repo.path, imprt.internal.repo.vendor, err).ToError()
		}
	}

	return nil
}

func (v *Vcs) Clone(imprt *Import) error {

	branch := imprt.Branch
	if imprt.Version != "" {
		branch = imprt.Version
	} else if imprt.Revision != "" {
		branch = imprt.Revision
	}

	var gitArgs []string
	v.logger.Debugf("executing Clone for [%s]", imprt.internal.repo.path)

	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	// remove cached temporary folder to prevent errors
	os.Remove(imprt.internal.repo.save)

	var repo string
	if v.protocol == ProtocolSSH {
		repo = imprt.internal.repo.ssh
	} else {
		repo = imprt.internal.repo.https
	}

	v.logger.Infof("downloading repository with %s protocol [%s] to [%s]", v.protocol, repo, pathCachedRepo)
	gitArgs = []string{
		"clone",
		"--recursive",
		"--shallow-submodules",
	}
	if branch != "" {
		gitArgs = append(gitArgs, "--branch", branch)
	}
	gitArgs = append(gitArgs, repo, pathCachedRepo)

	cmd := exec.Command("git", gitArgs...)

	if stderr, err := cmd.CombinedOutput(); err != nil {
		os.Remove(imprt.internal.repo.save)
		return v.logger.Warnf("error executing [git clone] command %s", string(stderr)).ToError()
	}

	v.logger.Infof("git clone completed for [%s]", imprt.internal.repo.save)

	return nil
}

func (v *Vcs) Fetch(path string) error {
	v.logger.Debugf("executing Fetch")

	gitArgs := []string{
		"fetch",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git fetch] command %s", string(stderr)).ToError()
	}

	return nil
}

func (v *Vcs) FetchAllTags(path string) error {
	v.logger.Debugf("executing Fetch All Tags")

	v.logger.Infof("fetching all tags")
	gitArgs := []string{
		"fetch",
		"--all",
		"--tags",
		"--prune",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git --all --tags --prune] command %s", string(stderr)).ToError()
	}

	return nil
}

func (v *Vcs) Pull(path string, imprt *Import) error {
	v.logger.Debugf("executing Pull")

	v.logger.Infof("updating repository [%s]", imprt.internal.repo.path)
	gitArgs := []string{
		"pull",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git pull] command %s", string(stderr)).ToError()
	}

	return nil
}

func (v *Vcs) Checkout(path string, imprt *Import) error {
	v.logger.Debugf("executing Checkout")

	branch := imprt.Branch

	if imprt.Version != "" {
		branch = imprt.Version
	}

	if imprt.Revision != "" {
		branch = imprt.Revision
	}

	v.logger.Infof("checkout repository [%s] branch [%s]", imprt.internal.repo.path, branch)
	gitArgs := []string{
		"checkout",
		branch,
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return v.logger.Errorf("error executing [git checkout] command %s", string(stderr)).ToError()
	}

	return nil
}

func (v *Vcs) doUpdateImportInfo(imprt *Import) error {
	pathCachedRepo := fmt.Sprintf("%s/%s", v.cache.path, imprt.internal.repo.save)

	if imprt.Version != "" {
		imprt.Version, _ = v.GetVersion(pathCachedRepo)
	}

	imprt.Revision, _ = v.GetRevision(pathCachedRepo)

	return nil
}

func (v *Vcs) GetDefaultBranch(path string) (string, error) {
	v.logger.Debugf("executing Get Default Branch")

	gitArgs := []string{
		"symbolic-ref",
		"--short",
		"HEAD",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return "master", v.logger.Infof("error getting git default branch: %s", string(stderr)).ToError()
	} else {
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) GetBranch(path string) (string, error) {
	v.logger.Debugf("executing Get Branch")

	gitArgs := []string{
		"rev-parse",
		"--abbrev-ref",
		"HEAD",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return "", v.logger.Infof("error getting git branch: %s", string(stderr)).ToError()
	} else {
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) GetLatestTag(path string) (string, string, error) {
	v.logger.Debugf("executing Get Latest Version")

	gitArgs := []string{
		"describe",
		"--tags",
		"`git rev-list --tags --max-count=1`",
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = path

	if stderr, err := cmd.CombinedOutput(); err != nil {
		return "master", "", v.logger.Infof("returned [%s] getting latest version", string(stderr)).ToError()
	} else {
		return "", strings.TrimSpace(string(stderr)), nil
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
		return "", v.logger.Infof("error getting git version: %s", string(stderr)).ToError()
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
		return "", v.logger.Infof("error getting git revision: %s", string(stderr)).ToError()
	} else {
		rtn := strings.TrimSpace(string(stderr))
		if rtn == "" {
			return "", nil
		}
		return strings.TrimSpace(string(stderr)), nil
	}
}

func (v *Vcs) SaveCacheImport(imprt *Import) error {
	v.logger.Debugf("executing Save Cache")

	// update cache
	v.cache.imports[imprt.internal.repo.save] = imprt

	var file *os.File
	var err error
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if _, err = os.Stat(v.cache.config); err != nil {
		// create file
		if err := os.MkdirAll(v.cache.path, os.ModePerm); err != nil {
			return v.logger.Errorf("error creating folder [%s] %s", v.cache.path, err).ToError()
		}

		file, err = os.Create(v.cache.config)
		if err != nil {
			return v.logger.Errorf("error creating file [%s] %s", v.cache.config, err).ToError()
		}
	} else {
		file, err = os.OpenFile(v.cache.config, os.O_APPEND|os.O_WRONLY, 0644)
	}

	if bytes, err := yaml.Marshal(Imports{imprt.internal.repo.path: imprt}); err != nil {
		return v.logger.Errorf("error marshal imports %s", err).ToError()
	} else {
		if _, err = file.Write(bytes); err != nil {
			return v.logger.Errorf("error writing file [%s] %s", v.cache.config, err).ToError()
		}
	}

	return nil
}

func (v *Vcs) doBackupConfig() (string, error) {
	if _, err := os.Stat(v.cache.config); err == nil {
		bkName := fmt.Sprintf("%s_%s", v.cache.config, time.Now().Format("20060102150405"))
		v.logger.Debugf("executing Backup Vendor to [%s]", bkName)

		os.Rename(v.cache.config, bkName)

		return bkName, nil
	}
	return "", nil
}

func (v *Vcs) doUndoBackupConfig(bkName string) error {
	os.Remove(v.cache.config)
	if _, err := os.Stat(v.cache.config); err == nil {
		v.logger.Debugf("executing Undo Backup Vendor to [%s]", v.cache.config)
		os.Rename(bkName, v.cache.config)
	}
	return nil
}

func (v *Vcs) SaveCache() error {
	v.logger.Debugf("executing Save Cache")
	var file *os.File
	var err error
	defer file.Close()

	if file, err = os.OpenFile(v.cache.config, os.O_RDWR, 0666); err != nil {
		v.logger.Infof("creating file [%s]", v.cache.config)

		file, err = os.Create(v.cache.config)
		if err != nil {
			return v.logger.Errorf("error creating file [%s] %s", v.cache.config, err).ToError()
		}
	} else {
		if err = file.Truncate(0); err != nil {
			return v.logger.Errorf("error cleaning [%s] file", v.cache.config).ToError()
		}
	}

	if bytes, err := yaml.Marshal(v.cache.imports); err != nil {
		return v.logger.Errorf("error marshal imports %s", err).ToError()
	} else {
		if _, err = file.Write(bytes); err != nil {
			return v.logger.Errorf("error writing file [%s] %s", v.cache.config, err).ToError()
		}
	}

	return nil
}
