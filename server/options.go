package server

import (
	"github.com/joaosoft/logger"
	"webserver"
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
	return func(webserver *WebServer) {
		webserver.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WebServerOption {
	return func(webserver *WebServer) {
		webserver.logger = logger
		webserver.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WebServerOption {
	return func(webserver *WebServer) {
		webserver.logger.SetLevel(level)
	}
}

// WithAddress ...
func WithAddress(address string) WebServerOption {
	return func(webserver *WebServer) {
		webserver.address = address
	}
}

// WithMultiAttachmentMode ...
func WithMultiAttachmentMode(mode webserver.MultiAttachmentMode) WebServerOption {
	return func(webserver *WebServer) {
		webserver.multiAttachmentMode = mode
	}
}
