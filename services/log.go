package golog

import (
	"io"
	"os"
)

type Log interface {
	With(tags []string, fields map[string]interface{}) Log
	WithTags(tags ...string) Log
	WithFields(fields map[string]interface{}) Log
	WithData(data map[string]interface{}) Log

	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// GoLog ...
type GoLog struct {
	fields map[string]string
	tags   map[string]string
	writer io.Writer

	quit    chan int
	started bool
}

// NewLog ...
func NewLog() *GoLog {
	return &GoLog{
		fields: make(map[string]string),
		tags:   make(map[string]string),
		writer: os.Stdout,
	}
}
