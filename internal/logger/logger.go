package logger

import (
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

var (
	instance *Logger
	once     sync.Once
)

func Get() *Logger {
	once.Do(func() {
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		instance = &Logger{
			Logger: zapLogger,
			sugar:  zapLogger.Sugar(),
		}
	})
	return instance
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

func (l *Logger) Close() error {
	return l.Logger.Sync()
}
