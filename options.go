package webserver

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// WebServerOption ...
type WebServerOption func(builder *WebServer)

// Reconfigure ...
func (w *WebServer) Reconfigure(options ...WebServerOption) {
	for _, option := range options {
		option(w)
	}
}

// WithConfiguration ...
func WithConfiguration(config *WebServerConfig) WebServerOption {
	return func(builder *WebServer) {
		builder.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WebServerOption {
	return func(builder *WebServer) {
		builder.logger = logger
		builder.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WebServerOption {
	return func(builder *WebServer) {
		builder.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) WebServerOption {
	return func(builder *WebServer) {
		builder.pm = mgr
	}
}

// WithPort ...
func WithPort(port int) WebServerOption {
	return func(builder *WebServer) {
		builder.port = port
	}
}
