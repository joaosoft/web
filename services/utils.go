package golog

import (
	"fmt"
	"strings"
)

func ParseLevel(level string) (Level, error) {
	switch strings.ToUpper(level) {
	case "panic":
		return PanicLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	default:
		return DefaultLevel, fmt.Errorf("invalid level: %s, set default level: %s", level, DefaultLevel)
	}
}

func (level Level) String() string {
	switch level {
	case PanicLevel:
		return "panic"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return "info"
	}
}
