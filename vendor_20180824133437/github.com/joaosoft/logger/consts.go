package logger

type Level int
type Prefix string

const (
	DefaultLevel = InfoLevel // DefaultLevel Level

	PanicLevel Level = iota // PanicLevel, when there is no recover
	FatalLevel              // FatalLevel, when the error is fatal to the application
	ErrorLevel              // ErrorLevel, when there is a controlled error
	WarnLevel               // WarnLevel, when there is a warning
	InfoLevel               // InfoLevel, when it is a informational message
	DebugLevel              // DebugLevel, when it is a debugging message
	PrintLevel              // PrintLevel, when it is a system message
	NoneLevel               // NoneLevel, when the logging is disabled

	// Special Prefixes
	LEVEL     Prefix = "{{LEVEL}}"     // Add the level value to the prefix
	TIMESTAMP Prefix = "{{TIMESTAMP}}" // Add the timestamp value to the prefix
	DATE      Prefix = "{{DATE}}"      // Add the date value to the prefix
	TIME      Prefix = "{{TIME}}"      // Add the time value to the prefix
	IP        Prefix = "{{IP}}"        // Add the client ip address
	TRACE     Prefix = "{{TRACE}}"     // Add the error trace
	PACKAGE   Prefix = "{{PACKAGE}}"   // Add the package name
	FILE      Prefix = "{{FILE}}"      // Add the file
	FUNCTION  Prefix = "{{FUNCTION}}"  // Add the function name
	STACK     Prefix = "{{STACK}}"     // Add the debug stack
)
