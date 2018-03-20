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

// WithLevel ...
func WithLevel(level Level) GoLogOption {
	return func(golog *GoLog) {
		golog.level = level
	}
}

// WithFormatHandler ...
func WithFormatHandler(formatHandler FormatHandler) GoLogOption {
	return func(golog *GoLog) {
		golog.formatHandler = formatHandler
	}
}
