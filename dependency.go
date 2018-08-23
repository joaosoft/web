package dependency

import (
	"fmt"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Dependency struct {
	config        *DependencyConfig
	isLogExternal bool
	pm            *manager.Manager
	logger        logger.ILogger
}

func NewDependency(options ...DependencyOption) *Dependency {
	pm := manager.NewManager(manager.WithRunInBackground(true))
	log := logger.NewLogDefault("dependency", logger.InfoLevel)

	service := &Dependency{
		pm:     pm,
		logger: log,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Dependency.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Dependency
	service.Reconfigure(options...)

	return service
}

// execute ...
func (d *Dependency) Get() error {
	d.logger.Debug("executing Get")
	var err error
	executed := make(Imports)
	allImports := make(Imports)
	var newVendor string

	defer func() {
		if err != nil {
			d.doUndoBackupVendor(newVendor)
		}
	}()

	// backup old vendor folder
	if newVendor, err = d.doBackupVendor(); err != nil {
		return err
	}

	if err = d.doGet(d.config.Path, executed, false); err != nil {
		return err
	} else {
		// save generated imports
		if err = d.doSaveImports(allImports); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dependency) Reset() error {
	d.logger.Debug("executing Reset")

	return d.doReset()
}
