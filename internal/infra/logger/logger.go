package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/marcodellorto/bitbucket-runner-autoscaler/internal/domain/ports"
)

type SlogLogger struct {
	log *slog.Logger
}

func NewLogger(loggingLevel int) ports.Logger {

	lvl := new(slog.LevelVar)
	lvl.Set(slog.Level(loggingLevel))

	sl := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	return &SlogLogger{log: sl}
}

func (zl *SlogLogger) Debug(msg string, keysAndValues ...interface{}) {
	zl.log.Debug(msg, keysAndValues...)
}

func (zl *SlogLogger) Info(msg string, keysAndValues ...interface{}) {
	zl.log.Info(msg, keysAndValues...)
}

func (zl *SlogLogger) Warn(msg string, keysAndValues ...interface{}) {
	zl.log.Warn(msg, keysAndValues...)
}

func (zl *SlogLogger) Error(msg string, keysAndValues ...interface{}) {
	zl.log.Error(msg, keysAndValues...)
}

func (zl *SlogLogger) Fatal(msg string, keysAndValues ...interface{}) {
	zl.log.Error(msg, keysAndValues...)
	log.Fatal(msg)
}
