package builder

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/watcher"
)

type Builder struct {
	config        *BuilderConfig
	event         chan *watcher.Event
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
	reloadTime    int64
	quit          chan int
	started       bool
}

func NewBuilder(options ...BuilderOption) (*Builder, error) {
	pm := manager.NewManager(manager.WithRunInBackground(true))
	log := logger.NewLogDefault("builder", logger.InfoLevel)
	event := make(chan *watcher.Event)

	w := watcher.NewWatcher(watcher.WithLogger(log), watcher.WithManager(pm), watcher.WithEventChannel(event))
	pm.AddProcess("watcher", w)

	service := &Builder{
		event:      event,
		reloadTime: 1,
		pm:         pm,
		logger:     log,
		quit:       make(chan int),
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(service.logger))
	}

	// load configuration File
	appConfig := &AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Builder.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Builder

	service.Reconfigure(options...)

	return service, nil
}

// execute ...
func (b *Builder) execute() error {
	b.logger.Debug("executing builder")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for {
			select {
			case <-termChan:
				b.quit <- 1
				b.logger.Info("received term signal")
				return
			case <-b.quit:
				b.logger.Info("received shutdown signal")
				return
			case <-time.After(time.Duration(b.reloadTime) * time.Second):
				b.logger.Info("watching changes")

				ev := <-b.event
				fmt.Println(ev.Operation)
				if ev.Operation == watcher.OperationChanges {
					b.rebuild()
				}
			}
		}
	}()

	return nil
}

// execute ...
func (b *Builder) rebuild() error {
	b.logger.Debug("executing rebuild")

	return nil
}

// Start ...
func (b *Builder) Start(wg *sync.WaitGroup) error {
	b.started = true
	wg.Add(1)
	wg.Done()

	if err := b.pm.Start(); err != nil {
		return err
	}

	if err := b.execute(); err != nil {
		return err
	}

	return nil
}

// Started ...
func (b *Builder) Started() bool {
	return b.started
}

// Stop ...
func (b *Builder) Stop(wg *sync.WaitGroup) error {
	wg.Add(1)
	wg.Done()

	b.quit <- 1
	if err := b.pm.Stop(); err != nil {
		return err
	}

	b.started = false

	return nil
}
