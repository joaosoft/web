package service

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

	yaml "gopkg.in/yaml.v2"
)

func (d *Dependency) doGet(dir string, loadedImports map[string]bool, installedImports Imports, isVendorPackage bool) error {
	sync := Memory{
		lockedImports:    make(Imports),
		internalImports:  make(Imports),
		externalImports:  make(Imports),
		loadedImports:    loadedImports,
		installedImports: installedImports,
		update:           false,
	}

	if _, ok := loadedImports[dir]; !ok {
		sync.loadedImports[dir] = true

		// load imports from project
		if err := d.doLoadImports(dir, &sync, isVendorPackage); err != nil {
			return err
		}
	} else {
		d.logger.Infof("directory already copied [%s]", dir)
		return nil
	}

	// load locked imports
	if err := d.doLoadLockImports(dir, &sync); err != nil {
		return err
	} else {
		// merge imports with lock
		if err := d.doMergeWithLockImports(&sync); err != nil {
			return err
		}
	}

	// download imports
	if err := d.doDownloadImports(&sync); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) doUpdate(dir string, loadedImports map[string]bool, installedImports Imports, isVendorPackage bool) error {
	sync := Memory{
		internalImports:  make(Imports),
		externalImports:  make(Imports),
		loadedImports:    loadedImports,
		installedImports: installedImports,
		update:           true,
	}

	if _, ok := loadedImports[dir]; !ok {
		sync.loadedImports[dir] = true

		// load imports from project
		if err := d.doLoadImports(dir, &sync, isVendorPackage); err != nil {
			return err
		}
	} else {
		d.logger.Infof("directory already copied [%s]", dir)
		return nil
	}

	// load locked imports
	if err := d.doLoadLockImports(dir, &sync); err != nil {
		return err
	} else {
		// merge imports with lock
		if err := d.doMergeWithLockImports(&sync); err != nil {
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

func (d *Dependency) doLoadImports(dir string, sync *Memory, isVendorPackage bool) error {
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
		if !isVendorPackage {
			for _, exclude := range excludedPaths {
				if strings.HasPrefix(dir, exclude) {
					d.logger.Infof("the import [%s] is on excluded paths", dir)
					return nil
				}
			}
		}

		// exclude validation for suffix
		for _, exclude := range excludedPaths {
			if strings.HasSuffix(dir, exclude) {
				d.logger.Infof("excluded path [%s]", exclude)
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
			if err := d.doLoadImports(nextDir, sync, isVendorPackage); err != nil {
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

func (d *Dependency) doLoadLockImports(dir string, sync *Memory) error {
	d.logger.Debugf("executing Load Lock Imports")
	lockImportFile := fmt.Sprintf("%s/%s", dir, LockImportFile)
	newLockedImports := make(Imports)

	if _, err := os.Stat(lockImportFile); err == nil {
		if bytes, err := ioutil.ReadFile(lockImportFile); err != nil {
			return d.logger.Errorf("error reading file [%s] %s", lockImportFile, err).ToError()
		} else {
			if err := yaml.Unmarshal(bytes, newLockedImports); err != nil {
				return d.logger.Errorf("error unmarshal file [%s] %s", lockImportFile, err).ToError()
			}
		}

		if !strings.Contains(dir, "vendor") {
			sync.lockedImports = newLockedImports
		} else {
			for newKey, newValue := range newLockedImports {
				if _, ok := sync.lockedImports[newKey]; !ok {
					sync.lockedImports[newKey] = newValue
				}
			}
		}
	} else {
		if !strings.Contains(dir, "vendor") {
			newFile, err := os.Create(LockImportFile)
			if err != nil {
				return d.logger.Errorf("error creating file [%s] %s", LockImportFile, err).ToError()
			}
			newFile.Close()
		}
	}

	return nil
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

	d.logger.Infof("configuration saved on [%s]", GenImportFile).ToError()

	return nil
}

func (d *Dependency) doGetFileImports(dir string, sync *Memory) error {
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

			// validate ignored packages
			for _, ignored := range ignoredPackages {
				if strings.Contains(name, ignored) {
					goto next
				}
			}

			if host, user, project, packag, ssh, https, path, vendor, save, err := d.doGetRepositoryInfo(name); err != nil {
				d.logger.Infof("repository ignored [%s]", name)
				return nil
			} else {
				if _, ok := sync.loadedImports[path]; !ok {
					if len(packag) > 0 {
						pkg := packag[1:]
						pkgs := strings.Split(pkg, "/")
						tmpPkg := ""
						for _, p := range pkgs {
							tmpPkg = fmt.Sprintf("%s/%s", tmpPkg, p)
							if _, ok := sync.loadedImports[fmt.Sprintf("%s%s", path, tmpPkg)]; ok {
								continue
							}
						}
					}

					sync.externalImports[path] = &Import{
						internal: Internal{
							repo: Repo{
								host:    host,
								user:    user,
								project: project,
								packag:  packag,
								ssh:     ssh,
								https:   https,
								path:    path,
								vendor:  vendor,
								save:    save,
							},
						},
					}

					sync.loadedImports[fmt.Sprintf("%s%s", path, packag)] = true
				}
			}
		}
	next:
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

func (d *Dependency) doMergeWithLockImports(sync *Memory) error {
	d.logger.Debugf("executing Merge With Lock Imports")
	for lockedKey, lockedValue := range sync.lockedImports {

		if externalValue, ok := sync.externalImports[lockedKey]; ok {
			lockedValue.internal = externalValue.internal
			d.logger.Debugf("replacing [%s] with locked [%+v]", lockedKey, lockedValue)
			sync.externalImports[lockedKey] = lockedValue
		}
	}

	return nil
}

func (d *Dependency) doDownloadImports(sync *Memory) error {
	d.logger.Debugf("executing Download imports to vendor")

	for _, imprt := range sync.externalImports {
		sync.installedImports[imprt.internal.repo.path] = imprt

		if err := d.vcs.CopyDependency(imprt, d.vendor, sync.update); err != nil {
			d.logger.Infof("repository ignored [%s]", imprt.internal.repo.ssh)
			continue
		}

		// to get inner vendor if it exists
		if _, err := os.Stat(fmt.Sprintf("%s/%s/", d.vendor, imprt.internal.repo.vendor)); err != nil {
			d.logger.Infof("getting vendor of [%s] import", imprt.internal.repo.vendor)
			if err := d.doGet(imprt.internal.repo.vendor, sync.loadedImports, sync.installedImports, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Dependency) doGetRepositoryInfo(name string) (string, string, string, string, string, string, string, string, string, error) {
	var host string
	var user string
	var project string
	var packag string
	var ssh string
	var https string
	var path string
	var save string

	save = name

	// moved packages
	for _, rename := range movedPackages {
		if strings.Contains(name, rename.old) {
			d.logger.Infof("renaming package [%s] from [%s] to [%s]", name, rename.old, rename.new)
			name = strings.Replace(name, rename.old, rename.new, 1)
			break
		}
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
		https = fmt.Sprintf("https://%s/%s/%s", host, user, project)
		path = fmt.Sprintf("%s/%s/%s", host, user, project)

	} else if len(nSplit) == 2 {

		host = nSplit[0]
		project = nSplit[1]

		ssh = fmt.Sprintf("git@%s:/%s", host, project)
		https = fmt.Sprintf("https://%s/%s", host, project)
		path = fmt.Sprintf("%s/%s", host, project)

	} else {
		return "", "", "", "", "", "", "", "", "", d.logger.Errorf("invalid import [%s]", name).ToError()
	}

	vendor := fmt.Sprintf("%s/%s", d.vendor, save)

	return host, user, project, packag, ssh, https, path, vendor, save, nil
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
