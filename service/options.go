package gomapper

import (
	logger "github.com/joaosoft/go-log/service"
)

// GoMapperOption ...
type GoMapperOption func(gomapper *GoMapper)

// Reconfigure ...
func (gomapper *GoMapper) Reconfigure(options ...GoMapperOption) {
	for _, option := range options {
		option(gomapper)
	}
}

// WithLogger ...
func WithLogger(logger logger.ILog) GoMapperOption {
	return func(gomapper *GoMapper) {
		log = logger
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) GoMapperOption {
	return func(gomapper *GoMapper) {
		log.SetLevel(level)
	}
}
