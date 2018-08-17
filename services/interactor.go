package services

import (
	"github.com/joaosoft/logger"
)

type IStorageDB interface {
	GetMigration(idMigration string) (*Migration, error)
	GetMigrations(values map[string][]string) (ListMigration, error)
	CreateMigration(newMigration *Migration) error
	DeleteMigration(idMigration string) error
	DeleteMigrations() error
	ExecuteMigration(migration string) error
}

type Interactor struct {
	logger    logger.ILogger
	storageDB IStorageDB
}

func NewInteractor(logger logger.ILogger, storageDB IStorageDB) *Interactor {
	return &Interactor{
		logger:    logger,
		storageDB: storageDB,
	}
}

func (interactor *Interactor) GetMigrations(values map[string][]string) (ListMigration, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "GetMigrations"})
	interactor.logger.Debug("getting migrations")
	if categories, err := interactor.storageDB.GetMigrations(values); err != nil {
		err = interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting migrations on storage database %s", err).ToError()
		return nil, err
	} else {
		return categories, nil
	}
}

func (interactor *Interactor) GetMigration(idMigration string) (*Migration, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "GetMigration"})
	interactor.logger.Debugf("getting migration %s", idMigration)
	if category, err := interactor.storageDB.GetMigration(idMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting migration %S on storage database %s", idMigration, err).ToError()
		return nil, err
	} else {
		return category, nil
	}
}

func (interactor *Interactor) CreateMigration(newMigration *Migration) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "CreateMigration"})

	interactor.logger.Debugf("creating migration with id %s", newMigration.IdMigration)
	if err := interactor.storageDB.CreateMigration(newMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error creating migration %s on storage database %s", newMigration.IdMigration, err).ToError()
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) DeleteMigration(idMigration string) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteMigration"})
	interactor.logger.Debugf("deleting migration %s", idMigration)
	if err := interactor.storageDB.DeleteMigration(idMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting migration %s on storage database %s", idMigration, err).ToError()
		return err
	}
	return nil
}

func (interactor *Interactor) DeleteMigrations() error {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteMigrations"})
	interactor.logger.Debug("deleting migrations")
	if err := interactor.storageDB.DeleteMigrations(); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error deleting migrations on storage database %s", err).ToError()
		return err
	}
	return nil
}

func (interactor *Interactor) ExecuteMigration(migration string) error {
	interactor.logger.WithFields(map[string]interface{}{"method": "ExecuteMigration"})
	interactor.logger.Debug("execute migration")
	if err := interactor.storageDB.ExecuteMigration(migration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error executing migration on storage database %s", err).ToError()
		return err
	}
	return nil
}
