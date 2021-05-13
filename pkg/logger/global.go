package logger

import (
	"sync"
)

var defaultConfig = &Config{
	Debug: true,
}

var once sync.Once
var globalLogger ILogger

func GetLogger() ILogger {
	once.Do(func() {
		globalLogger = NewLogger(defaultConfig)
	})
	return globalLogger.Clone()
}

func ReplaceGlobalLogger(logger ILogger) {
	if globalLogger != nil {
		globalLogger.Flush()
	}
	globalLogger = logger
}
