package golog

import (
	"io"
	"sync"
)

type Log interface {
	With(prefixes, tags, fields map[string]interface{}) Log
	WithPrefixes(prefixes map[string]interface{}) Log
	WithTags(tags map[string]interface{}) Log
	WithFields(fields map[string]interface{}) Log

	Debug(message interface{})
	Info(message interface{})
	Warn(message interface{})
	Error(message interface{})

	Debugf(format string, arguments ...interface{})
	Infof(format string, arguments ...interface{})
	Warnf(format string, arguments ...interface{})
	Errorf(format string, arguments ...interface{})
}

type FormatHandler func(level Level, message *Message) ([]byte, error)

type Entry struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Message struct {
	Prefixes map[string]interface{} `json:"prefixes"`
	Tags     map[string]interface{} `json:"tags"`
	Message  interface{}                 `json:"message"`
	Fields   map[string]interface{} `json:"fields"`
}

// GoLog ...
type GoLog struct {
	level         Level
	writer        io.Writer
	prefixes      map[string]interface{} `json:"prefixes"`
	tags          map[string]interface{} `json:"tags"`
	fields        map[string]interface{} `json:"fields"`
	formatHandler FormatHandler
	mux           *sync.Mutex
}
