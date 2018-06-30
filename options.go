package mapper

import logger "github.com/joaosoft/logger"

// MapperOption ...
type MapperOption func(mapper *Mapper)

// Reconfigure ...
func (mapper *Mapper) Reconfigure(options ...MapperOption) {
	for _, option := range options {
		option(mapper)
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) MapperOption {
	return func(mapper *Mapper) {
		log = logger
		mapper.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) MapperOption {
	return func(mapper *Mapper) {
		log.SetLevel(level)
	}
}
