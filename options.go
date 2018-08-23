package dependency

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// DependencyOption ...
type DependencyOption func(builder *Dependency)

// Reconfigure ...
func (d *Dependency) Reconfigure(options ...DependencyOption) {
	for _, option := range options {
		option(d)
	}
}

// WithConfiguration ...
func WithConfiguration(config *DependencyConfig) DependencyOption {
	return func(builder *Dependency) {
		builder.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) DependencyOption {
	return func(builder *Dependency) {
		builder.logger = logger
		builder.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) DependencyOption {
	return func(builder *Dependency) {
		builder.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) DependencyOption {
	return func(builder *Dependency) {
		builder.pm = mgr
	}
}
