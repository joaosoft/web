package web

import (
	"db-migration/services"

	"github.com/joaosoft/logger"
)

// WebServiceOption ...
type WebServiceOption func(client *WebService)

// Reconfigure ...
func (service *WebService) Reconfigure(options ...WebServiceOption) {
	for _, option := range options {
		option(service)
	}
}

// WithConfiguration ...
func WithConfiguration(config *services.DbMigrationConfig) WebServiceOption {
	return func(client *WebService) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WebServiceOption {
	return func(service *WebService) {
		service.logger = logger
		service.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WebServiceOption {
	return func(service *WebService) {
		service.logger.SetLevel(level)
	}
}
