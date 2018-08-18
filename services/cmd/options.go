package cmd

import (
	"db-migration/services"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// CmdServiceOption ...
type CmdServiceOption func(client *CmdService)

// Reconfigure ...
func (service *CmdService) Reconfigure(options ...CmdServiceOption) {
	for _, option := range options {
		option(service)
	}
}

// WithConfiguration ...
func WithConfiguration(config *services.DbMigrationConfig) CmdServiceOption {
	return func(client *CmdService) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) CmdServiceOption {
	return func(service *CmdService) {
		service.logger = logger
		service.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) CmdServiceOption {
	return func(service *CmdService) {
		service.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) CmdServiceOption {
	return func(service *CmdService) {
		service.pm = mgr
	}
}
