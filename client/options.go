package client

import (
	"web"

	"github.com/joaosoft/logger"
)

// WebClientOption ...
type WebClientOption func(builder *WebClient)

// Reconfigure ...
func (c *WebClient) Reconfigure(options ...WebClientOption) {
	for _, option := range options {
		option(c)
	}
}

// WithConfiguration ...
func WithConfiguration(config *WebClientConfig) WebClientOption {
	return func(webserver *WebClient) {
		webserver.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WebClientOption {
	return func(webserver *WebClient) {
		webserver.logger = logger
		webserver.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WebClientOption {
	return func(webserver *WebClient) {
		webserver.logger.SetLevel(level)
	}
}

// WithAddress ...
func WithAddress(address string) WebClientOption {
	return func(webserver *WebClient) {
		webserver.address = address
	}
}

// WithMultiAttachmentMode ...
func WithMultiAttachmentMode(mode web.MultiAttachmentMode) WebClientOption {
	return func(webserver *WebClient) {
		webserver.multiAttachmentMode = mode
	}
}
