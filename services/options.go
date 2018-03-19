package golog

import (
	"io"
)

// GoLogOption ...
type GoLogOption func(golog *GoLog)

// Reconfigure ...
func (golog *GoLog) Reconfigure(options ...GoLogOption) {
	for _, option := range options {
		option(golog)
	}
}

// WithWriter ...
func WithWriter(writer io.Writer) GoLogOption {
	return func(golog *GoLog) {
		golog.writer = writer
	}
}
