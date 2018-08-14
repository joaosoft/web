package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"db-migration/services"

	"io/ioutil"

	"sort"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type CmdService struct {
	config        *services.DbMigrationConfig
	interactor    *services.Interactor
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
}

func NewService(options ...CmdServiceOption) (*CmdService, error) {
	service := &CmdService{
		pm:     manager.NewManager(manager.WithRunInBackground(true)),
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

	simpleDB := manager.NewSimpleDB(&appConfig.DbMigration.Db)
	simpleDB.Start(&sync.WaitGroup{})
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	service.interactor = services.NewInteractor(service.logger, services.NewStoragePostgres(service.logger, simpleDB))

	return service, nil
}

// Execute ...
func (service *CmdService) Execute(option services.MigrationOption, number int) error {
	service.logger.Infof("executing migration with option '-%s %s'", services.MigrationCmd, option)

	// load
	executed, toexecute, err := service.load()
	if err != nil {
		return err
	}

	// validate
	if err := service.validate(executed, toexecute); err != nil {
		return err
	}

	// process
	if err := service.process(option, number, executed, toexecute); err != nil {
		return err
	}

	return nil
}

// load ...
func (service *CmdService) load() (executed []string, toexecute []string, err error) {

	// executed
	service.logger.Info("getting all executed migrations")
	migrations, er := service.interactor.GetMigrations(nil)
	if er != nil {
		service.logger.Error("error loading migrations from database").ToErr(er)
		return nil, nil, er
	}
	for _, migration := range migrations {
		executed = append(executed, migration.IdMigration)
	}

	// to execute
	service.logger.Info("getting all migrations from file system")
	dir, _ := os.Getwd()
	files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.sql", dir, service.config.Path))
	if err != nil {
		service.logger.Error("error loading migrations from file system").ToError(&err)
		return executed, nil, err
	}
	for _, file := range files {
		fileName := file[strings.LastIndex(file, "/")+1:]
		toexecute = append(toexecute, fileName)
	}

	return executed, toexecute, err
}

// validate ...
func (service *CmdService) validate(executed []string, toexecute []string) (err error) {
	service.logger.Info("validate migrations")
	for i, migration := range executed {
		if migration != toexecute[i] {
			service.logger.Errorf("the migrations are in a different order of the already executed migrations [%s] <-> [%s]", migration, toexecute[i]).ToError(&err)
			return err
		}
	}
	return nil
}

// process ...
func (service *CmdService) process(option services.MigrationOption, number int, executed []string, toexecute []string) error {
	var migrations []string

	if option == services.MigrationOptionUp {
		if len(toexecute) <= len(executed) {
			service.logger.Info("there are no migrations to execute!")
			return nil
		}

		if number > (len(toexecute) - len(executed)) {
			number = len(toexecute) - len(executed)
		}
		sort.Slice(toexecute, func(i, j int) bool {
			if toexecute[i] < toexecute[j] {
				return true
			}
			return false
		})

		if number > 0 {
			migrations = toexecute[len(executed) : len(executed)+number]
		} else {
			migrations = toexecute[len(executed):]
		}
	} else {
		if len(executed) == 0 {
			service.logger.Info("there are no migrations to execute!")
			return nil
		}
		toexecute = toexecute[:len(executed)]
		sort.Slice(toexecute, func(i, j int) bool {
			if toexecute[i] < toexecute[j] {
				return false
			}
			return true
		})

		if number == 0 {
			number = 1
		}

		if number > 0 {
			migrations = toexecute[:number]
		} else {
			migrations = toexecute
		}
	}

	// log migration information
	service.logger.Infof("migrations already executed %+v", executed)
	service.logger.Infof("migrations to execute %+v", migrations)

	for _, migration := range migrations {
		fileName := migration[strings.LastIndex(migration, "/")+1:]

		service.logger.Infof("Running sql migration: [%s]", fileName)

		dir, _ := os.Getwd()
		file, err := os.Open(fmt.Sprintf("%s/%s/%s", dir, service.config.Path, migration))
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(file)

		indexUp := strings.Index(string(data), string(services.TagMigrationUp))
		indexDown := strings.Index(string(data), string(services.TagMigrationDown))

		var migrationBody []string
		var migrationUp string
		var migrationDown string
		if indexUp < indexDown {
			migrationBody = strings.Split(string(data), string(services.TagMigrationDown))
			if len(migrationBody) > 0 {
				migrationUp = migrationBody[0]
			}
			if len(migrationBody) > 1 {
				migrationDown = migrationBody[1]
			}

		} else {
			migrationBody = strings.Split(string(data), string(services.TagMigrationUp))
			if len(migrationBody) > 0 {
				migrationDown = migrationBody[0]
			}
			if len(migrationBody) > 1 {
				migrationUp = migrationBody[1]
			}
		}

		var query string
		if option == services.MigrationOptionUp {
			query = migrationUp
			if migrationUp == "" {
				service.logger.Infof("empty migration up on migration %s", migration)
			}
		}

		if option == services.MigrationOptionDown {
			query = migrationDown
			if migrationDown == "" {
				service.logger.Infof("empty migration down on migration %s", migration)
			}
		}

		conn, err := service.config.Db.Connect()
		if err != nil {
			return err
		}
		defer conn.Close()

		tx, err := conn.Begin()
		if err != nil {
			return err
		}
		defer func(err error) {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}(err)

		_, err = tx.Exec(query)

		if option == services.MigrationOptionUp {
			if err == nil {
				if er := service.interactor.CreateMigration(&services.Migration{IdMigration: fileName}); er != nil {
					service.logger.Error("error adding migration to database")
					err = er
					return er
				}
			}
		} else {
			if err == nil {
				if er := service.interactor.DeleteMigration(migration); er != nil {
					service.logger.Error("error deleting migration to database")
					err = er
					return er
				}
			}
		}

		if err != nil {
			service.logger.Errorf("error executing the migration %s", fileName).ToError(&err)
			return err
		}
	}

	return nil
}

// Start ...
func (m *CmdService) Start() error {
	return m.pm.Start()
}

// Stop ...
func (m *CmdService) Stop() error {
	return m.pm.Stop()
}
