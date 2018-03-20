package golog

// Level ...
type Level int

const (
	// DefaultLevel Level
	DefaultLevel = InfoLevel

	// PanicLevel, when there is no recover
	PanicLevel Level = iota
	// FatalLevel, when the error is fatal to the application
	FatalLevel
	// ErrorLevel, when there is a controlled error
	ErrorLevel
	// WarnLevel, when there is a warning
	WarnLevel
	// InfoLevel, when it is a informational message
	InfoLevel
	// DebugLevel, when it is a debugging message
	DebugLevel

	LEVEL = "{{LEVEL}}"
	TIME  = "{{TIME}}"
)
