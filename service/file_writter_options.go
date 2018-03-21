package gowriter

import "time"

// FileWriterOption ...
type FileWriterOption func(fileWriter *FileWriter)

// Reconfigure ...
func (fileWriter *FileWriter) Reconfigure(options ...FileWriterOption) {
	for _, option := range options {
		option(fileWriter)
	}
}

// WithDirectory ...
func WithDirectory(directory string) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.directory = directory
	}
}

// WithFileName ...
func WithFileName(fileName string) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.fileName = fileName
	}
}

// WithFileMaxMegaByteSize ...
func WithFileMaxMegaByteSize(fileMaxSize int64) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.fileMaxSize = fileMaxSize * MB_IN_BYTE
	}
}

// WithFlushTime ...
func WithFlushTime(flushTime time.Duration) FileWriterOption {
	return func(fileWriter *FileWriter) {
		fileWriter.config.flushTime = flushTime
	}
}
