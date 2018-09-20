package server

import (
	"web/common"

	"github.com/joaosoft/logger"
)

// ServerOption ...
type ServerOption func(builder *Server)

// Reconfigure ...
func (w *Server) Reconfigure(options ...ServerOption) {
	for _, option := range options {
		option(w)
	}
}

// WithConfiguration ...
func WithConfiguration(config *ServerConfig) ServerOption {
	return func(webserver *Server) {
		webserver.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) ServerOption {
	return func(webserver *Server) {
		webserver.logger = logger
		webserver.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) ServerOption {
	return func(webserver *Server) {
		webserver.logger.SetLevel(level)
	}
}

// WithAddress ...
func WithAddress(address string) ServerOption {
	return func(webserver *Server) {
		webserver.address = address
	}
}

// WithMultiAttachmentMode ...
func WithMultiAttachmentMode(mode common.MultiAttachmentMode) ServerOption {
	return func(webserver *Server) {
		webserver.multiAttachmentMode = mode
	}
}
