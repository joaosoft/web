package mapper

import logger "github.com/joaosoft/logger"

// mapperOption ...
type mapperOption func(mapper *Mapper)

// Reconfigure ...
func (mapper *Mapper) Reconfigure(options ...mapperOption) {
	for _, option := range options {
		option(mapper)
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) mapperOption {
	return func(mapper *Mapper) {
		log = logger
		mapper.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) mapperOption {
	return func(mapper *Mapper) {
		log.SetLevel(level)
	}
}
