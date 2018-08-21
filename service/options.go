package service

import (
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// WatcherOption ...
type WatcherOption func(client *Watcher)

// Reconfigure ...
func (w *Watcher) Reconfigure(options ...WatcherOption) {
	for _, option := range options {
		option(w)
	}
}

// WithConfiguration ...
func WithConfiguration(config *WatcherConfig) WatcherOption {
	return func(client *Watcher) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WatcherOption {
	return func(service *Watcher) {
		service.logger = logger
		service.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WatcherOption {
	return func(service *Watcher) {
		service.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) WatcherOption {
	return func(service *Watcher) {
		service.pm = mgr
	}
}

// WithQuitChannel ...
func WithQuitChannel(quit chan int) WatcherOption {
	return func(service *Watcher) {
		service.quit = quit
	}
}

// WithReloadTime ...
func WithReloadTime(reloadTime time.Duration) WatcherOption {
	return func(service *Watcher) {
		service.reloadTime = reloadTime
	}
}

// WithEventChannel ...
func WithEventChannel(event chan *Event) WatcherOption {
	return func(service *Watcher) {
		service.event = event
	}
}
