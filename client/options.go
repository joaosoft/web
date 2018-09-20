package client

import (
	"web/common"

	"github.com/joaosoft/logger"
)

// ClientOption ...
type ClientOption func(builder *Client)

// Reconfigure ...
func (c *Client) Reconfigure(options ...ClientOption) {
	for _, option := range options {
		option(c)
	}
}

// WithConfiguration ...
func WithConfiguration(config *ClientConfig) ClientOption {
	return func(WebClient *Client) {
		WebClient.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) ClientOption {
	return func(WebClient *Client) {
		WebClient.logger = logger
		WebClient.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) ClientOption {
	return func(WebClient *Client) {
		WebClient.logger.SetLevel(level)
	}
}

// WithMultiAttachmentMode ...
func WithMultiAttachmentMode(mode common.MultiAttachmentMode) ClientOption {
	return func(WebClient *Client) {
		WebClient.multiAttachmentMode = mode
	}
}
