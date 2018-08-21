package service

import (
	"fmt"
	"os"
	"sync"

	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Watcher struct {
	config        *WatcherConfig
	watch         []string
	excluded      []string
	extensions    []string
	files         map[string]map[string]FileInfo
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
	reloadTime    time.Duration
	quit          chan int
	event         chan *Event
}

func NewWatcher(options ...WatcherOption) (*Watcher, error) {
	service := &Watcher{
		watch:      make([]string, 0),
		excluded:   make([]string, 0),
		extensions: make([]string, 0),
		reloadTime: time.Duration(time.Second * 5),
		files:      make(map[string]map[string]FileInfo),
		pm:         manager.NewManager(manager.WithRunInBackground(true)),
		logger:     logger.NewLogDefault("watcher", logger.InfoLevel),
		event:      make(chan *Event),
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
		level, _ := logger.ParseLevel(appConfig.Watcher.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.config = &appConfig.Watcher

	// loading each configuration
	service.watch = append(service.watch, service.config.Dirs.Watch...)
	service.excluded = append(service.excluded, service.config.Dirs.Excluded...)
	service.extensions = append(service.extensions, service.config.Dirs.Extensions...)

	service.Reconfigure(options...)

	return service, nil
}

func (w *Watcher) AddWatch(watchs ...string) *Watcher {
	w.watch = append(w.watch, watchs...)
	return w
}

func (w *Watcher) AddExtension(extensions ...string) *Watcher {
	w.extensions = append(w.extensions, extensions...)
	return w
}

func (w *Watcher) AddExcluded(excluded ...string) *Watcher {
	w.excluded = append(w.excluded, excluded...)
	return w
}

// execute ...
func (w *Watcher) execute() error {
	w.logger.Debugf("executing watcher for watch %+v", w.watch)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	// load
	for _, dir := range w.watch {

		_, err := os.Stat(dir)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-termChan:
					w.logger.Info("received term signal")
					return
				case <-w.quit:
					w.logger.Info("received shutdown signal")
					return
				case <-time.After(w.reloadTime):
					w.logger.Info("reloading data")

					// copy before reload files
					oldFiles := w.files[dir]
					w.files[dir] = make(map[string]FileInfo)

					if err = w.doLoad(oldFiles, dir, dir); err != nil {
						w.quit <- 1
					}

					if err = w.doRemove(dir, oldFiles); err != nil {
						w.quit <- 1
					}
				}
			}
		}()
	}

	<-termChan

	return nil
}

// Start ...
func (w *Watcher) Start() error {
	if err := w.pm.Start(); err != nil {
		return err
	}

	if err := w.execute(); err != nil {
		return err
	}

	return nil
}

// Stop ...
func (w *Watcher) Stop() error {
	w.quit <- 1
	return w.pm.Stop()
}

func (w *Watcher) doLoad(oldFiles map[string]FileInfo, dir string, next string) error {
	fileInfo, err := os.Stat(next)
	if err != nil {
		return err
	}

	if strings.HasPrefix(fileInfo.Name(), ".") {
		return nil
	}

	// if it is a directory
	if fileInfo.IsDir() {

		// exclude validation
		for _, exclude := range w.excluded {
			if strings.HasPrefix(next, exclude) {
				return nil
			}
		}

		w.logger.Debugf("loading files on directory [%s]", next)

		subDir, err := filepath.Glob(fmt.Sprintf("%s/*", next))
		if err != nil {
			w.logger.Errorf("error reading directory %s", err)
			return err
		}
		for _, nextDir := range subDir {
			w.logger.Debugf("loading files on subdirectory [%s]", nextDir)
			w.doLoad(oldFiles, dir, nextDir)
		}
		return nil
	}

	// extension validation
	if index := strings.LastIndex(next, "."); index > 0 {
		allowed := false
		for _, extension := range w.extensions {
			if extension == next[index+1:] {
				allowed = true
			}
		}

		if !allowed {
			return nil
		}
	}

	// if it is a file
	w.files[dir][next] = FileInfo{
		FullName: next,
		Name:     fileInfo.Name(),
		Size:     fileInfo.Size(),
		ModTime:  fileInfo.ModTime(),
	}

	if oldFileInfo, ok := oldFiles[next]; !ok {
		// new file
		w.logger.Debugf("added a new file on directory [%s]", next)
		w.event <- &Event{
			File:      next,
			Operation: OperationCreate,
		}
	} else {
		if oldFileInfo.ModTime != fileInfo.ModTime() ||
			oldFileInfo.Size != fileInfo.Size() {
			// updated file
			w.logger.Debugf("changed file on directory [%s]", next)
			w.event <- &Event{
				File:      next,
				Operation: OperationUpdate,
			}
		}
	}

	return nil
}

func (w *Watcher) doRemove(dir string, oldFiles map[string]FileInfo) error {

	for fullName, _ := range oldFiles {
		if _, ok := w.files[dir][fullName]; !ok {
			delete(w.files[dir], fullName)
			w.logger.Debugf("deleted file on directory [%s]", dir)
			w.event <- &Event{
				File:      fullName,
				Operation: OperationDelete,
			}
		}
	}

	return nil
}
