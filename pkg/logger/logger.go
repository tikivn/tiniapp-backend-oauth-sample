package logger

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	pkg_ctxutil "tiniapp-backend-oauth-sample/pkg/ctxutil"
)

type Config struct {
	Debug bool `yaml:"debug"`
}

type ILogger interface {
	Flush()
	Clone() ILogger

	WithContext(ctx context.Context) ILogger
	WithField(key string, value interface{}) ILogger
	WithPrefix(prefix string) ILogger

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger struct {
	config   *Config
	internal *zap.Logger
}

func (l *Logger) Flush() {
	err := l.internal.Sync()
	if err != nil {
		_ = fmt.Errorf("log could not flush, error: %+v", err)
	}
}

func (l *Logger) Clone() ILogger {
	return &Logger{
		config:   l.config,
		internal: l.internal,
	}
}

func (l *Logger) WithContext(ctx context.Context) ILogger {
	_, requestID := pkg_ctxutil.ExtractRequestID(ctx)
	return &Logger{
		internal: l.internal.With(
			zap.Any("request_id", requestID),
		),
	}
}

func (l *Logger) WithField(key string, value interface{}) ILogger {
	return &Logger{
		internal: l.internal.With(
			zap.Any(key, value),
		),
	}
}

func (l *Logger) WithPrefix(prefix string) ILogger {
	return l.WithField("prefix", prefix)
}

func (l *Logger) WithRequestID(requestID string) ILogger {
	return l.WithField("request_id", requestID)
}

func prepareArgs(args ...interface{}) string {
	items := []string{}
	for _, arg := range args {
		items = append(items, fmt.Sprintf("%+v", arg))
	}
	return strings.Join(items, " ")
}

func (l *Logger) Debug(args ...interface{}) {
	l.internal.Debug(prepareArgs(args...))
}

func (l *Logger) Debugln(args ...interface{}) {
	l.internal.Debug(prepareArgs(args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.internal.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.internal.Info(prepareArgs(args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.internal.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.internal.Warn(prepareArgs(args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.internal.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(args ...interface{}) {
	l.internal.Error(prepareArgs(args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.internal.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Panic(args ...interface{}) {
	l.internal.Panic(prepareArgs(args...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.internal.Panic(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.internal.Fatal(prepareArgs(args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.internal.Fatal(fmt.Sprintf(format, args...))
}

func NewLogger(config *Config) *Logger {
	zapConfig := zap.NewProductionConfig()
	if config.Debug {
		zapConfig = zap.NewDevelopmentConfig()
	}
	zapLogger, _ := zapConfig.Build()

	logger := &Logger{
		config:   config,
		internal: zapLogger,
	}

	return logger
}

func ProvideLogger(config *Config) (ILogger, func(), error) {
	logger := NewLogger(config)
	ReplaceGlobalLogger(logger)

	cleanup := func() {
		logger.Flush()
	}

	return logger, cleanup, nil
}
