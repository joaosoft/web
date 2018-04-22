package gomapper

import "github.com/joaosoft/go-log/service"

// mapperOption ...
type mapperOption func(mapper *Mapper)

// Reconfigure ...
func (mapper *Mapper) Reconfigure(options ...mapperOption) {
	for _, option := range options {
		option(mapper)
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) mapperOption {
	return func(mapper *Mapper) {
		log = logger
		mapper.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) mapperOption {
	return func(mapper *Mapper) {
		log.SetLevel(level)
	}
}
