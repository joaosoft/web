package builder

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"os/exec"

	"io"
	"io/ioutil"

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

func NewBuilder(options ...BuilderOption) *Builder {
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

	return service
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
				b.logger.Info("received term signal")
				return
			case <-b.quit:
				b.logger.Info("received shutdown signal")
				return
			case <-time.After(time.Duration(b.reloadTime) * time.Second):
				b.logger.Info("watching changes")

				ev := <-b.event
				if ev.Operation == watcher.OperationChanges {
					b.build()
					b.start()
				}
			}
		}
	}()

	return nil
}

// build ...
func (b *Builder) build() error {
	b.logger.Info("executing build")
	cmd := exec.Command("go", "build", "-i", "-o", b.config.Destination, b.config.Source)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return b.logger.Errorf("error getting stderr pipe %s", err).ToError()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return b.logger.Errorf("error getting stdout pipe %s", err).ToError()
	}

	err = cmd.Start()
	if err != nil {
		return b.logger.Errorf("error executing build command %s", err).ToError()
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return b.logger.Errorf("error executing build %s", string(errBuf)).ToError()
	}
	b.logger.Info("build completed")

	return nil
}

// start ...
func (b *Builder) start() error {
	b.logger.Info("executing start")
	cmd := exec.Command("./" + b.config.Destination)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return b.logger.Errorf("error getting stderr pipe %s", err).ToError()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return b.logger.Errorf("error getting stdout pipe %s", err).ToError()
	}

	err = cmd.Start()
	if err != nil {
		return b.logger.Errorf("error executing restart command %s", err).ToError()
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return b.logger.Errorf("error executing restart %s", string(errBuf)).ToError()
	}
	b.logger.Info("start completed")

	return nil
}

// Start ...
func (b *Builder) Start(wg *sync.WaitGroup) error {
	b.started = true
	if wg != nil {
		defer wg.Done()
	}

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
	b.started = false
	if wg != nil {
		defer wg.Done()
	}

	if err := b.pm.Stop(); err != nil {
		return err
	}

	return nil
}
