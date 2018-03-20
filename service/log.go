package golog

import (
	"fmt"
	"os"
	"sync"
)

// NewLog ...
func NewLog(options ...GoLogOption) *GoLog {
	golog := &GoLog{
		writer:        os.Stdout,
		formatHandler: JsonFormatHandler,
		level:         InfoLevel,
		prefixes:      make(map[string]interface{}),
		tags:          make(map[string]interface{}),
		fields:        make(map[string]interface{}),
		mux:           &sync.Mutex{},
	}
	golog.Reconfigure(options...)

	return golog
}

func (log *GoLog) With(prefixes, tags, fields map[string]interface{}) Log {
	log.WithPrefixes(prefixes)
	log.WithTags(tags)
	log.WithFields(fields)
	return log
}

func (log *GoLog) WithPrefixes(prefixes map[string]interface{}) Log {
	log.prefixes = prefixes
	return log
}

func (log *GoLog) WithTags(tags map[string]interface{}) Log {
	log.tags = tags
	return log
}

func (log *GoLog) WithFields(fields map[string]interface{}) Log {
	log.fields = fields
	return log
}

func (log *GoLog) WithField(key string, value interface{}) Log {
	log.fields[key] = fmt.Sprintf("%s", value)
	return log
}

func (log *GoLog) Debug(message interface{}) {
	log.writeLog(DebugLevel, message)
}

func (log *GoLog) Info(message interface{}) {
	log.writeLog(InfoLevel, message)
}

func (log *GoLog) Warn(message interface{}) {
	log.writeLog(WarnLevel, message)
}

func (log *GoLog) Error(message interface{}) {
	log.writeLog(ErrorLevel, message)
}

func (log *GoLog) Debugf(format string, arguments ...interface{}) {
	log.writeLog(DebugLevel, fmt.Sprintf(format, arguments...))
}

func (log *GoLog) Infof(format string, arguments ...interface{}) {
	log.writeLog(InfoLevel, fmt.Sprintf(format, arguments...))
}

func (log *GoLog) Warnf(format string, arguments ...interface{}) {
	log.writeLog(WarnLevel, fmt.Sprintf(format, arguments...))
}

func (log *GoLog) Errorf(format string, arguments ...interface{}) {
	log.writeLog(ErrorLevel, fmt.Sprintf(format, arguments...))
}

func (log *GoLog) writeLog(level Level, message interface{}) {
	log.mux.Lock()
	defer log.mux.Unlock()
	if level > log.level {
		return
	}

	if bytes, err := log.formatHandler(level, &Message{Prefixes: log.prefixes, Tags: log.tags, Message: message, Fields: log.fields}); err != nil {
		return
	} else {
		log.writer.Write([]byte(fmt.Sprintf("%s\n", bytes)))
	}
}
