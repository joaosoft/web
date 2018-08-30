package logger

var internalLogger = NewLoggerEmpty(InfoLevel)

func SetLevel(level Level) {
	internalLogger.SetLevel(level)
}

func With(prefixes, tags, fields, sufixes map[string]interface{}) ILogger {
	return internalLogger.With(prefixes, tags, fields, sufixes)
}

func WithPrefixes(prefixes map[string]interface{}) ILogger {
	return internalLogger.WithPrefixes(prefixes)
}

func WithTags(tags map[string]interface{}) ILogger {
	return internalLogger.WithTags(tags)
}

func WithFields(fields map[string]interface{}) ILogger {
	return internalLogger.WithFields(fields)
}

func WithSufixes(sufixes map[string]interface{}) ILogger {
	return internalLogger.WithSufixes(sufixes)
}

func WithPrefix(key string, value interface{}) ILogger {
	return internalLogger.WithPrefix(key, value)
}

func WithTag(key string, value interface{}) ILogger {
	return internalLogger.WithTag(key, value)
}

func WithField(key string, value interface{}) ILogger {
	return internalLogger.WithField(key, value)
}

func WithSufix(key string, value interface{}) ILogger {
	return internalLogger.WithSufix(key, value)
}

func Print(message interface{}) IAddition {
	return internalLogger.Print(message)
}

func Debug(message interface{}) IAddition {
	return internalLogger.Debug(message)
}

func Info(message interface{}) IAddition {
	return internalLogger.Info(message)
}

func Warn(message interface{}) IAddition {
	return internalLogger.Warn(message)
}

func Error(message interface{}) IAddition {
	return internalLogger.Error(message)
}

func Panic(message interface{}) IAddition {
	return internalLogger.Panic(message)
}

func Fatal(message interface{}) IAddition {
	return internalLogger.Fatal(message)
}

func Printf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Printf(format, arguments)
}

func Debugf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Debugf(format, arguments)
}

func Infof(format string, arguments ...interface{}) IAddition {
	return internalLogger.Infof(format, arguments)
}

func Warnf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Warnf(format, arguments)
}

func Errorf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Errorf(format, arguments)
}

func Panicf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Panicf(format, arguments)
}

func Fatalf(format string, arguments ...interface{}) IAddition {
	return internalLogger.Fatalf(format, arguments)
}

func Reconfigure(options ...LoggerOption) {
	internalLogger.Reconfigure(options...)
}

func Get() ILogger {
	return internalLogger
}
