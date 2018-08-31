package main

import (
	"fmt"
	"path/filepath"

	"os"
	"strings"

	"go/parser"
	"go/token"

	"strconv"

	"io/ioutil"

	"time"

	"regexp"

	"gopkg.in/yaml.v2"
)

func (d *Dependency) doGet(dir string, loadedImports map[string]bool, installedImports Imports, isVendorPackage bool, update bool) error {
	var err error
	var path string

	sync := Memory{
		generatedImports: make(Imports),
		lockedImports:    make(Imports),
		internalImports:  make(Imports),
		externalImports:  make(Imports),
		loadedImports:    loadedImports,
		installedImports: installedImports,
		update:           update,
	}

	if _, _, _, _, _, _, path, _, _, err = d.doGetRepositoryInfo(dir); err != nil {
		d.logger.Infof("error getting repository information on [%s]", dir)
	}

	if _, ok := loadedImports[path]; !ok {
		key := dir
		if strings.HasPrefix(dir, "vendor/") {
			key = strings.SplitN(dir, "vendor/", 2)[1]
		}
		sync.loadedImports[key] = true

		// load imports from project
		if err := d.doLoadImports(dir, &sync, isVendorPackage); err != nil {
			return err
		}

		// load locked imports
		if err := d.doLoadLockedImports(dir, &sync); err != nil {
			return err
		}

		// load generated imports
		if err := d.doLoadGeneratedImports(dir, &sync); err != nil {
			return err
		}

		// merge with locked imports
		if err := d.doMergeWithLockedImports(&sync); err != nil {
			return err
		}

		// merge with generated imports
		if err := d.doMergeWithGeneratedImports(&sync); err != nil {
			return err
		}

		// download imports
		if err := d.doDownloadImports(&sync); err != nil {
			return err
		}
	} else {
		d.logger.Infof("directory already copied [%s]", dir)
		return nil
	}

	return nil
}

func (d *Dependency) doClearLock() error {
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

func (d *Dependency) doClearGen() error {
	if file, err := os.OpenFile(GenImportFile, os.O_RDWR, 0666); err != nil {
		d.logger.Infof("creating file [%s]", GenImportFile)

		newFile, err := os.Create(GenImportFile)
		if err != nil {
			return d.logger.Errorf("error creating file [%s] %s", GenImportFile, err).ToError()
		}
		newFile.Close()
	} else {
		defer file.Close()
		if err := file.Truncate(0); err != nil {
			return d.logger.Errorf("error cleaning [%s] file", GenImportFile).ToError()
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

	if regx, err := regexp.Compile(RegexForVendorFiles); err != nil {
		return d.logger.Errorf("error compiling regex for vendor folders: %s", err).ToError()
	} else if regx.MatchString(fileInfo.Name()) {
		return nil
	}

	// if it is a directory
	if fileInfo.IsDir() {

		if dir == d.bkVendor {
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

			// moved packages
			var newPackage string
			for _, action := range packageActions {
				if strings.HasPrefix(name, action.old) {
					d.logger.Infof("renaming package [%s] from [%s] to [%s]", name, action.old, action.new)
					newPackage = strings.Replace(name, action.old, action.new, 1)
					break
				}
			}

			if host, user, project, packag, ssh, https, path, vendor, save, err := d.doGetRepositoryInfo(name); err != nil {
				d.logger.Infof("repository ignored [%s]", name)
				return nil
			} else {
				if newPackage != "" {
					if _, _, _, _, ssh, https, _, _, _, err = d.doGetRepositoryInfo(newPackage); err != nil {
						d.logger.Infof("repository ignored [%s]", name)
						return nil
					}
				}
				if _, ok := sync.loadedImports[path]; !ok {
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

					sync.loadedImports[path] = true
				}
			}
		}
	next:
	}

	return nil
}

func (d *Dependency) doLoadLockedImports(dir string, sync *Memory) error {
	d.logger.Debugf("executing Load Lock Imports on [%s]", dir)
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
		for newKey, newValue := range newLockedImports {
			if _, ok := sync.lockedImports[newKey]; !ok {
				d.logger.Debugf("reading locked [%s] [%+v]", newKey, newValue).ToError()
				sync.lockedImports[newKey] = newValue
			}
		}
	} else {
		if !strings.HasPrefix(dir, "vendor") {
			newFile, err := os.Create(LockImportFile)
			if err != nil {
				return d.logger.Errorf("error creating file [%s] %s", LockImportFile, err).ToError()
			}
			newFile.Close()
		}
	}

	return nil
}

func (d *Dependency) doMergeWithLockedImports(sync *Memory) error {
	d.logger.Debugf("executing Merge With Locked Imports")
	for lockedKey, lockedValue := range sync.lockedImports {

		if externalValue, ok := sync.externalImports[lockedKey]; ok {
			lockedValue.internal = externalValue.internal
			d.logger.Debugf("replacing [%s] with locked [%+v]", lockedKey, lockedValue)
			sync.externalImports[lockedKey] = lockedValue
		}
	}

	return nil
}

func (d *Dependency) doLoadGeneratedImports(dir string, sync *Memory) error {
	d.logger.Debugf("executing Load Generated Imports on [%s]", dir)
	genImportFile := fmt.Sprintf("%s/%s", dir, GenImportFile)
	newGeneratedImports := make(Imports)

	if _, err := os.Stat(genImportFile); err == nil {
		if bytes, err := ioutil.ReadFile(genImportFile); err != nil {
			return d.logger.Errorf("error reading file [%s] %s", genImportFile, err).ToError()
		} else {
			if err := yaml.Unmarshal(bytes, &newGeneratedImports); err != nil {
				return d.logger.Errorf("error unmarshal file [%s] %s", genImportFile, err).ToError()
			}
		}

		for key, value := range newGeneratedImports {
			if _, ok := sync.generatedImports[key]; !ok {
				sync.generatedImports[key] = value
			}
		}
	}

	return nil
}

func (d *Dependency) doMergeWithGeneratedImports(sync *Memory) error {
	d.logger.Debugf("executing Merge With Generated Imports")
	for generatedKey, generatedValue := range sync.generatedImports {

		if _, ok := sync.lockedImports[generatedKey]; !ok {

			if externalValue, ok := sync.externalImports[generatedKey]; ok {
				generatedValue.internal = externalValue.internal
				d.logger.Debugf("replacing [%s] with generated [%+v]", generatedKey, generatedValue)
				sync.externalImports[generatedKey] = generatedValue
			}
		}
	}

	return nil
}

func (d *Dependency) doDownloadImports(sync *Memory) error {
	d.logger.Debugf("executing Download imports to vendor")

	for _, imprt := range sync.externalImports {

		if err := d.vcs.CopyDependency(sync, imprt, d.vendor, sync.update); err != nil {
			d.logger.Infof("repository ignored [%s]", imprt.internal.repo.ssh)
			continue
		}

		// to get inner vendor if it exists
		if _, err := os.Stat(imprt.internal.repo.vendor); err == nil {
			d.logger.Infof("getting vendor of [%s] import", imprt.internal.repo.vendor)
			if err := d.doGet(imprt.internal.repo.vendor, sync.loadedImports, sync.installedImports, true, sync.update); err != nil {
				return err
			}
		}

		sync.installedImports[imprt.internal.repo.path] = imprt
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

	// example [github.com/username/path1/path2] and should be [git@github.com:username/path1]
	if nSplit := strings.SplitN(name, "/", 4); len(nSplit) >= 3 {

		host = nSplit[0]
		user = nSplit[1]
		project = nSplit[2]

		if len(nSplit) > 3 {
			packag = fmt.Sprintf("/%s", nSplit[3])
		}

		ssh = fmt.Sprintf("git@%s:%s/%s.git", host, user, project)
		https = fmt.Sprintf("https://%s/%s/%s", host, user, project)
		path = fmt.Sprintf("%s/%s/%s", host, user, project)

	} else if len(nSplit) == 2 {

		host = nSplit[0]
		project = nSplit[1]

		ssh = fmt.Sprintf("git@%s:/%s.git", host, project)
		https = fmt.Sprintf("https://%s/%s", host, project)
		path = fmt.Sprintf("%s/%s", host, project)

	} else {
		return "", "", "", "", "", "", "", "", "", d.logger.Errorf("invalid import [%s]", name).ToError()
	}

	save = path
	vendor := fmt.Sprintf("%s/%s", d.vendor, save)

	return host, user, project, packag, ssh, https, path, vendor, save, nil
}

func (d *Dependency) doBackupVendor() error {
	if _, err := os.Stat(d.vendor); err == nil {
		d.bkVendor = fmt.Sprintf("_%s_%s", d.vendor, time.Now().Format("20060102150405"))
		d.logger.Debugf("executing Backup Vendor to [%s]", d.bkVendor)

		os.Rename(d.vendor, d.bkVendor)
	}
	return nil
}

func (d *Dependency) doUndoBackupVendor() error {
	os.Remove(d.vendor)
	if _, err := os.Stat(d.bkVendor); err == nil {
		d.logger.Debugf("executing Undo Backup Vendor to [%s]", d.bkVendor)
		os.Rename(d.bkVendor, d.vendor)
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
