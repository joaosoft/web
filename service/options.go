package gomapper

import "github.com/joaosoft/go-log/service"

// MapperOption ...
type MapperOption func(mapper *Mapper)

// Reconfigure ...
func (mapper *Mapper) Reconfigure(options ...MapperOption) {
	for _, option := range options {
		option(mapper)
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) MapperOption {
	return func(mapper *Mapper) {
		log = logger
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) MapperOption {
	return func(mapper *Mapper) {
		log.SetLevel(level)
	}
}
