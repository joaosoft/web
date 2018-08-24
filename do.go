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

	"gopkg.in/yaml.v2"
)

func (d *Dependency) doGet(dir string, loadedImports map[string]bool, installedImports Imports, loadExcludedPaths bool) error {
	sync := sync{
		internalImports:  make(Imports),
		externalImports:  make(Imports),
		loadedImports:    loadedImports,
		installedImports: installedImports,
	}

	if _, ok := loadedImports[strings.Replace(dir, "vendor/", "", 1)]; !ok {

		// load imports from project
		if err := d.doLoadImports(dir, &sync, loadExcludedPaths); err != nil {
			return err
		}

		sync.loadedImports[strings.Replace(dir, "vendor/", "", 1)] = true
	}

	// load locked imports
	if lockImports, err := d.doLoadLockImports(); err != nil {
		return err
	} else {
		// merge imports with lock
		if err := d.doMergeWithLockImports(&sync, lockImports); err != nil {
			return err
		}
	}

	// download imports
	if err := d.doDownloadImports(&sync); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doReload(dir string, loadedImports map[string]bool, installedImports Imports, loadExcludedPaths bool) error {
	sync := sync{
		internalImports:  make(Imports),
		externalImports:  make(Imports),
		loadedImports:    loadedImports,
		installedImports: installedImports,
	}

	// load imports from project
	if err := d.doLoadImports(dir, &sync, loadExcludedPaths); err != nil {
		return err
	}

	// load locked imports
	if lockImports, err := d.doLoadLockImports(); err != nil {
		return err
	} else {
		// merge imports with lock
		if err := d.doMergeWithLockImports(&sync, lockImports); err != nil {
			return err
		}
	}

	// download imports
	if err := d.doDownloadImports(&sync); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doReset() error {
	if file, err := os.OpenFile(LockImportFile, os.O_RDWR, 0666); err != nil {
		d.logger.Infof("creating file [%s]", LockImportFile)

		newFile, err := os.Create(LockImportFile)
		if err != nil {
			return d.logger.Errorf("error creating file [%s] %s", LockImportFile, err).ToError()
		}
		newFile.Close()
	} else {
		defer file.Close()
		if err := file.Truncate(0); err != nil {
			return d.logger.Errorf("error cleaning [%s] file", LockImportFile).ToError()
		}
	}
	return nil
}

func (d *Dependency) doLoadImports(dir string, sync *sync, loadExcludedPaths bool) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}

	// ignore hidden folder/files
	if fileInfo.Name() != "." && strings.HasPrefix(fileInfo.Name(), ".") {
		return nil
	}

	// if it is a directory
	if fileInfo.IsDir() {

		if dir == d.oldVendor {
			return nil
		}

		// exclude validation for prefix
		if !loadExcludedPaths {
			for _, exclude := range excludedPaths {
				if strings.HasPrefix(dir, exclude) {
					d.logger.Infof("the import [%s] is on excluded paths", dir)
					return nil
				}
			}
		}

		for _, exclude := range excludedImports {
			// also allow to validate in inner vendor projects
			if strings.HasPrefix(strings.Replace(dir, "vendor/", "", 1), exclude) {
				d.logger.Infof("the import [%s] is on excluded imports list", dir)
				return nil
			}
		}

		// exclude validation for suffix
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
			if err := d.doLoadImports(nextDir, sync, loadExcludedPaths); err != nil {
				return err
			}
		}

		return nil
	}

	if !strings.HasSuffix(fileInfo.Name(), ".go") {
		return nil
	}

	d.logger.Debugf("loading file [%s]", fileInfo.Name())

	if err := d.doGetFileImports(dir, sync); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doLoadLockImports() (Imports, error) {
	d.logger.Debugf("executing Load Lock Imports")
	imports := make(map[string]*Import)

	if _, err := os.Stat(LockImportFile); err == nil {
		if bytes, err := ioutil.ReadFile(LockImportFile); err != nil {
			return imports, d.logger.Errorf("error reading file [%s] %s", LockImportFile, err).ToError()
		} else {
			if err := yaml.Unmarshal(bytes, &imports); err != nil {
				return nil, d.logger.Errorf("error unmarshal file [%s] %s", LockImportFile, err).ToError()
			}
			return imports, nil
		}
	} else {
		newFile, err := os.Create(LockImportFile)
		if err != nil {
			return nil, d.logger.Errorf("error creating file [%s] %s", LockImportFile, err).ToError()
		}
		newFile.Close()
	}

	return imports, nil
}

func (d *Dependency) doSaveImports(imports Imports) error {
	d.logger.Debugf("executing Save Imports")

	d.doDelete(GenImportFile)

	if bytes, err := yaml.Marshal(imports); err != nil {
		return d.logger.Errorf("error marshal imports %s", err).ToError()
	} else {
		if err := ioutil.WriteFile(GenImportFile, bytes, 0644); err != nil {
			return d.logger.Errorf("error writing file [%s] %s", GenImportFile, err).ToError()
		}
	}

	return nil
}

func (d *Dependency) doGetFileImports(dir string, sync *sync) error {
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

			sync.internalImports[name] = &Import{}
		} else {
			d.logger.Debugf("adding external dependency [%s]", name)

			if host, user, project, packag, ssh, https, path, err := d.doGetRepositoryInfo(name); err != nil {
				d.logger.Infof("repository ignored [%s]", name)
				return nil
			} else {
				if _, ok := sync.loadedImports[ssh]; !ok {
					sync.externalImports[path] = &Import{
						Branch: "master",
						internal: internal{
							host:    host,
							user:    user,
							project: project,
							packag:  packag,
							repo: repo{
								ssh:   ssh,
								https: https,
								path:  path,
							},
							vendor: fmt.Sprintf("%s/%s", d.vendor, path),
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

	if _, err := os.Stat(LockImportFile); err != nil {
		if bytes, err := ioutil.ReadFile(LockImportFile); err != nil {
			return nil, d.logger.Errorf("error reading file [%s] %s", LockImportFile, err).ToError()
		} else {
			imports := make(map[string]*Import)
			if err := yaml.Unmarshal(bytes, &imports); err != nil {
				return nil, d.logger.Errorf("error unmarshal file [%s] %s", LockImportFile, err).ToError()
			}
			return imports, nil
		}
	}

	return nil, nil
}

func (d *Dependency) doMergeWithLockImports(sync *sync, lockImports Imports) error {
	d.logger.Debugf("executing Merge With Lock Imports")
	for lockKey, lockValue := range lockImports {

		if externalValue, ok := sync.externalImports[lockKey]; ok {
			d.logger.Debugf("replacing [%s] with locked", lockKey)
			lockValue.internal = externalValue.internal

			if lockValue.Branch == "" {
				lockValue.Branch = externalValue.Branch
			}

			sync.externalImports[lockKey] = lockValue
		}
	}

	return nil
}

func (d *Dependency) doDownloadImports(sync *sync) error {
	d.logger.Debugf("executing Download imports to vendor")

	for _, imprt := range sync.externalImports {

		if _, ok := sync.installedImports[imprt.internal.repo.path]; ok {
			continue
		}

		sync.installedImports[imprt.internal.repo.path] = imprt

		if err := d.vcs.Clone(imprt, d.vendor); err != nil {
			d.logger.Infof("repository ignored [%s]", imprt.internal.repo.ssh)
			continue
		}

		// to get inner vendor if it exists
		if _, err := os.Stat(fmt.Sprintf("%s/%s/", d.vendor, imprt.internal.vendor)); err != nil {
			d.logger.Infof("getting vendor of [%s] import", imprt.internal.vendor)
			if err := d.doGet(imprt.internal.vendor, sync.loadedImports, sync.installedImports, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Dependency) doGetRepositoryInfo(name string) (string, string, string, string, string, string, string, error) {
	var host string
	var user string
	var project string
	var packag string
	var ssh string
	var https string
	var path string

	// moved packages
	for old, new := range movedPackages {
		name = strings.Replace(name, old, new, 1)
	}

	// example [github.com/username/path1/path2] and should be [git@github.com:username/path1]
	if nSplit := strings.SplitN(name, "/", 4); len(nSplit) >= 3 {

		host = nSplit[0]
		user = nSplit[1]
		project = nSplit[2]

		if len(nSplit) > 3 {
			packag = fmt.Sprintf("/%s", nSplit[3])
		}

		ssh = fmt.Sprintf("git@%s:%s/%s", host, user, project)
		https = fmt.Sprintf("https://%s/%s/%s", user, host, project)
		path = fmt.Sprintf("%s/%s/%s", host, user, project)

	} else if len(nSplit) == 2 {

		host = nSplit[0]
		project = nSplit[1]

		ssh = fmt.Sprintf("git@%s:/%s", host, project)
		https = fmt.Sprintf("https://%s/%s", host, project)
		path = fmt.Sprintf("%s/%s", host, project)

	} else {
		return "", "", "", "", "", "", "", d.logger.Errorf("invalid import [%s]", name).ToError()
	}

	return host, user, project, packag, ssh, https, path, nil
}

func (d *Dependency) doBackupVendor() error {
	d.oldVendor = fmt.Sprintf("%s_%s", d.vendor, time.Now().Format("20060102150405"))
	d.logger.Debugf("executing Backup Vendor to [%s]", d.oldVendor)

	if _, err := os.Stat(d.vendor); err == nil {
		os.Rename(d.vendor, d.oldVendor)
	}
	return nil
}

func (d *Dependency) doUndoBackupVendor() error {
	d.logger.Debugf("executing Undo Backup Vendor to [%s]", d.oldVendor)

	if _, err := os.Stat(d.oldVendor); err == nil {
		os.Rename(d.oldVendor, d.vendor)
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
