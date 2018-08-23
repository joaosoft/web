package dependency

import (
	"fmt"
	"path/filepath"

	"os"
	"strings"

	"go/parser"
	"go/token"

	"strconv"

	"time"

	"io/ioutil"

	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

func (d *Dependency) doGet(dir string, executed Imports, loadExcludedPaths bool) error {
	intImports := make(Imports)
	extImports := make(Imports)

	// load imports from project
	if err := d.doLoadImports(dir, intImports, extImports, executed, loadExcludedPaths); err != nil {
		return err
	}

	// load locked imports
	if lockImports, err := d.doLoadLockImports(); err != nil {
		return err
	} else {
		// merge imports with lock
		if err := d.doMergeWithLockImports(extImports, lockImports); err != nil {
			return err
		}
	}

	// download imports
	if err := d.doDownloadImports(extImports, executed); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doReset() error {
	if file, err := os.OpenFile(LockImportFileName, os.O_RDWR, 0666); err != nil {
		d.logger.Infof("creating file [%s]", LockImportFileName)

		newFile, err := os.Create(LockImportFileName)
		if err != nil {
			return d.logger.Errorf("error creating file [%s] %s", LockImportFileName, err).ToError()
		}
		newFile.Close()
	} else {
		defer file.Close()
		if err := file.Truncate(0); err != nil {
			return d.logger.Errorf("error cleaning [%s] file", LockImportFileName).ToError()
		}
	}
	return nil
}

func (d *Dependency) doLoadImports(dir string, intImports Imports, extImports Imports, executed Imports, loadExcludedPaths bool) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if fileInfo.Name() != "." && strings.HasPrefix(fileInfo.Name(), ".") {
		return nil
	}

	// if it is a directory
	if fileInfo.IsDir() {

		// exclude validation for prefix
		if !loadExcludedPaths {
			for _, exclude := range excludedPaths {
				if strings.HasPrefix(dir, exclude) {
					return nil
				}
			}
		}

		for _, exclude := range excludedImports {
			if strings.HasPrefix(dir, exclude) {
				return nil
			}
		}

		// exclude validation for sufix
		for _, exclude := range excludedPaths {
			if strings.HasSuffix(dir, exclude) {
				return nil
			}
		}

		d.logger.Debugf("loading files on directory [%s]", dir)
		subDir, err := filepath.Glob(fmt.Sprintf("%s/*", dir))
		if err != nil {
			d.logger.Errorf("error reading directory %s", err)
			return err
		}
		for _, nextDir := range subDir {
			if err := d.doLoadImports(nextDir, intImports, extImports, executed, loadExcludedPaths); err != nil {
				return err
			}
		}

		return nil
	}

	if !strings.HasSuffix(fileInfo.Name(), ".go") {
		return nil
	}

	d.logger.Debugf("loading file [%s]", fileInfo.Name())

	if err := d.doGetFileImports(dir, intImports, extImports, executed); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doLoadLockImports() (Imports, error) {
	d.logger.Debugf("executing Load Lock Imports")
	imports := make(map[string]Import)

	if _, err := os.Stat(LockImportFileName); err == nil {
		if bytes, err := ioutil.ReadFile(LockImportFileName); err != nil {
			return imports, d.logger.Errorf("error reading file [%s] %s", LockImportFileName, err).ToError()
		} else {
			if err := yaml.Unmarshal(bytes, &imports); err != nil {
				return nil, d.logger.Errorf("error unmarshal file [%s] %s", LockImportFileName, err).ToError()
			}
			return imports, nil
		}
	} else {
		newFile, err := os.Create(LockImportFileName)
		if err != nil {
			return nil, d.logger.Errorf("error creating file [%s] %s", LockImportFileName, err).ToError()
		}
		newFile.Close()
	}

	return imports, nil
}

func (d *Dependency) doSaveImports(imports Imports) error {
	d.logger.Debugf("executing Save Imports")

	d.doDelete(GenImportFileName)

	if bytes, err := yaml.Marshal(imports); err != nil {
		return d.logger.Errorf("error marshal imports %s", err).ToError()
	} else {
		if err := ioutil.WriteFile(GenImportFileName, bytes, 0644); err != nil {
			return d.logger.Errorf("error writing file [%s] %s", GenImportFileName, err).ToError()
		}
	}

	return nil
}

func (d *Dependency) doGetFileImports(dir string, intImports Imports, extImports Imports, executed Imports) error {
	d.logger.Debugf("executing Get Imports for file %s", dir)

	parsedFile, err := parser.ParseFile(token.NewFileSet(), dir, nil, parser.ImportsOnly|parser.ParseComments)
	if err != nil {
		if os.IsPermission(err) {
			return nil
		}

		d.logger.Warnf("error when parsing golang file [%s] %s", dir, err)
		return nil
	}

	for _, imprt := range parsedFile.Imports {
		name, err := strconv.Unquote(imprt.Path.Value)
		if err != nil {
			return d.logger.Errorf("error unquoting [%s] on file [%s]", imprt.Path.Value, dir).ToError()
		}

		if !strings.Contains(imprt.Path.Value, ".") {
			d.logger.Debugf("adding internal dependency [%s]", name)

			intImports[name] = Import{}
		} else {
			d.logger.Debugf("adding external dependency [%s]", name)

			if host, user, project, ssh, https, vendor, err := d.doGetRepositoryInfo(name); err != nil {
				return err
			} else {
				if _, ok := executed[ssh]; !ok {
					extImports[ssh] = Import{
						Branch: "master",
						internal: Internal{
							host:    host,
							user:    user,
							project: project,
							vendor:  vendor,
							repo: Repo{
								ssh:   ssh,
								https: https,
							},
						},
					}
				}
			}
		}
	}

	return nil
}

func (d *Dependency) doLoadLockedImports() (Imports, error) {
	d.logger.Debugf("executing Load Locked Imports")

	if _, err := os.Stat(LockImportFileName); err != nil {
		if bytes, err := ioutil.ReadFile(LockImportFileName); err != nil {
			return nil, d.logger.Errorf("error reading file [%s] %s", LockImportFileName, err).ToError()
		} else {
			imports := make(map[string]Import)
			if err := yaml.Unmarshal(bytes, &imports); err != nil {
				return nil, d.logger.Errorf("error unmarshal file [%s] %s", LockImportFileName, err).ToError()
			}
			return imports, nil
		}
	}

	return nil, nil
}

func (d *Dependency) doMergeWithLockImports(extImports Imports, lockImports Imports) error {
	d.logger.Debugf("executing Merge With Lock Imports")

	for lockKey, lockValue := range lockImports {
		if _, ok := extImports[lockKey]; ok {
			d.logger.Debugf("replacing [%s] with locked", lockKey)
			extImports[lockKey] = lockValue
		}
	}

	return nil
}

func (d *Dependency) doDownloadImports(imports Imports, executed Imports) error {
	d.logger.Debugf("executing Download imports to vendor")

	for repository, info := range imports {
		executed[repository] = info
		d.logger.Infof("downloading repository with ssh protocol [%s]", info.internal.repo.ssh)

		if _, err := exec.Command("git", "ls-remote", "-h", info.internal.repo.ssh).CombinedOutput(); err != nil {
			d.logger.Infof("the repository doesn't exist [%s]", info.internal.repo.ssh)
			return nil
		}

		if stderr, err := exec.Command("git", "clone", "--recursive", "-v", "--progress", "--depth", "1", "--shallow-submodules", info.internal.repo.ssh, info.internal.vendor).CombinedOutput(); err != nil {
			d.logger.Errorf("error executing git clone command %s", string(stderr))

			d.logger.Infof("retrying download with https protocol [%s]", info.internal.repo.https)
			if stderr, err := exec.Command("git", "clone", "--recursive", "-v", "--progress", "--depth", "1", "--shallow-submodules", info.internal.repo.https, info.internal.vendor).CombinedOutput(); err != nil {
				d.logger.Errorf("error executing git clone command %s", string(stderr)).ToError()
				d.logger.Infof("ignoring download of repository [%s]", info.internal.repo.https).ToError()
				continue
			}
		}

		d.logger.Infof("git clone completed for [%s]", repository)

		if _, err := os.Stat(fmt.Sprintf("%s/vendor/", info.internal.vendor)); err != nil {
			d.logger.Infof("getting vendor [%s] imports", info.internal.vendor)
			if err := d.doGet(info.internal.vendor, executed, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Dependency) doGetRepositoryInfo(name string) (string, string, string, string, string, string, error) {
	var host string
	var user string
	var project string
	var ssh string
	var https string
	var vendor string

	// moved packages
	for old, new := range movedPackages {
		name = strings.Replace(name, old, new, 1)
	}

	// example [github.com/username/path1/path2] and should be [git@github.com:username/path1]
	if nSplit := strings.Split(name, "/"); len(nSplit) >= 3 {
		host = nSplit[0]
		user = nSplit[1]
		project = nSplit[2]

		ssh = fmt.Sprintf("git@%s:%s/%s", host, user, project)
		https = fmt.Sprintf("https://%s/%s", host, project)
		vendor = fmt.Sprintf("vendor/%s/%s/%s", host, user, project)
	} else if len(nSplit) == 2 {
		host = nSplit[0]
		project = nSplit[1]
		ssh = fmt.Sprintf("ssh://%s/%s", host, project)
		https = fmt.Sprintf("https://%s/%s", host, project)
		vendor = fmt.Sprintf("vendor/%s/%s", host, project)
	} else {
		return "", "", "", "", "", "", d.logger.Errorf("invalid import [%s]", name).ToError()
	}

	return host, user, project, ssh, https, vendor, nil
}

func (d *Dependency) doBackupVendor() (string, error) {
	newName := fmt.Sprintf("vendor_%s", time.Now().Format("20060102150405"))
	d.logger.Debugf("executing Backup Vendor to [%s]", newName)

	if _, err := os.Stat("vendor"); err == nil {
		os.Rename("vendor", newName)
	}
	return newName, nil
}

func (d *Dependency) doUndoBackupVendor(oldName string) error {
	d.logger.Debugf("executing Undo Backup Vendor to [%s]", oldName)

	if _, err := os.Stat(oldName); err == nil {
		os.Rename(oldName, "vendor")
	}
	return nil
}

func (d *Dependency) doDelete(dir string) error {
	d.logger.Debugf("executing delete of [%s]", dir)

	if _, err := os.Stat(dir); err != nil {
		os.Remove(dir)
	}
	return nil
}
