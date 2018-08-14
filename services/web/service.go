package web

import (
	"fmt"

	"sync"

	"db-migration/services"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type WebService struct {
	config        *services.DbMigrationConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
}

// NewService ...
func NewService(options ...WebServiceOption) (*WebService, error) {
	service := &WebService{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("services-cmd", logger.InfoLevel),
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	// load configuration File
	appConfig := &services.AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", services.GetEnv()), appConfig); err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.DbMigration.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.DbMigration

	service.Reconfigure(options...)

	if service.config.Host == "" {
		service.config.Host = services.DefaultURL
	}

	simpleDB := manager.NewSimpleDB(&appConfig.DbMigration.Db)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	web := manager.NewSimpleWebEcho(service.config.Host)
	controller := NewController(service.logger, services.NewInteractor(service.logger, services.NewStoragePostgres(service.logger, simpleDB)))
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (m *WebService) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *WebService) Stop() error {
	return m.pm.Stop()
}
