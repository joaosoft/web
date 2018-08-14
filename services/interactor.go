package services

import (
	errors "github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
)

type IStorageDB interface {
	GetMigration(idMigration string) (*Migration, *errors.Err)
	GetMigrations(values map[string][]string) (ListMigration, *errors.Err)
	CreateMigration(newMigration *Migration) *errors.Err
	DeleteMigration(idMigration string) *errors.Err
	DeleteMigrations() *errors.Err
	ExecuteMigration(migration string) *errors.Err
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

func (interactor *Interactor) GetMigrations(values map[string][]string) (ListMigration, *errors.Err) {
	interactor.logger.WithFields(map[string]interface{}{"method": "GetMigrations"})
	interactor.logger.Info("getting migrations")
	if categories, err := interactor.storageDB.GetMigrations(values); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting migrations on storage database %s", err).ToErr(err)
		return nil, err
	} else {
		return categories, nil
	}
}

func (interactor *Interactor) GetMigration(idMigration string) (*Migration, *errors.Err) {
	interactor.logger.WithFields(map[string]interface{}{"method": "GetMigration"})
	interactor.logger.Infof("getting migration %s", idMigration)
	if category, err := interactor.storageDB.GetMigration(idMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error getting migration %S on storage database %s", idMigration, err).ToErr(err)
		return nil, err
	} else {
		return category, nil
	}
}

func (interactor *Interactor) CreateMigration(newMigration *Migration) *errors.Err {
	interactor.logger.WithFields(map[string]interface{}{"method": "CreateMigration"})

	interactor.logger.Infof("creating migration with id %s", newMigration.IdMigration)
	if err := interactor.storageDB.CreateMigration(newMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error creating migration %s on storage database %s", newMigration.IdMigration, err).ToErr(err)
		return err
	} else {
		return nil
	}
}

func (interactor *Interactor) DeleteMigration(idMigration string) *errors.Err {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteMigration"})
	interactor.logger.Infof("deleting migration %s", idMigration)
	if err := interactor.storageDB.DeleteMigration(idMigration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting migration %s on storage database %s", idMigration, err).ToErr(err)
		return err
	}
	return nil
}

func (interactor *Interactor) DeleteMigrations() *errors.Err {
	interactor.logger.WithFields(map[string]interface{}{"method": "DeleteMigrations"})
	interactor.logger.Info("deleting migrations")
	if err := interactor.storageDB.DeleteMigrations(); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error deleting migrations on storage database %s", err).ToErr(err)
		return err
	}
	return nil
}

func (interactor *Interactor) ExecuteMigration(migration string) *errors.Err {
	interactor.logger.WithFields(map[string]interface{}{"method": "ExecuteMigration"})
	interactor.logger.Info("execute migration")
	if err := interactor.storageDB.ExecuteMigration(migration); err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
			Errorf("error executing migration on storage database %s", err).ToErr(err)
		return err
	}
	return nil
}
