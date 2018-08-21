package builder

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// BuilderOption ...
type BuilderOption func(builder *Builder)

// Reconfigure ...
func (b *Builder) Reconfigure(options ...BuilderOption) {
	for _, option := range options {
		option(b)
	}
}

// WithConfiguration ...
func WithConfiguration(config *BuilderConfig) BuilderOption {
	return func(builder *Builder) {
		builder.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) BuilderOption {
	return func(builder *Builder) {
		builder.logger = logger
		builder.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) BuilderOption {
	return func(builder *Builder) {
		builder.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) BuilderOption {
	return func(builder *Builder) {
		builder.pm = mgr
	}
}

// WithReloadTime ...
func WithReloadTime(reloadTime int64) BuilderOption {
	return func(builder *Builder) {
		builder.reloadTime = reloadTime
	}
}

// WithQuitChannel ...
func WithQuitChannel(quit chan int) BuilderOption {
	return func(builder *Builder) {
		builder.quit = quit
	}
}
