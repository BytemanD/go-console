package console

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
)

var defaultLog *slog.Logger

func Debug(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Debug(fmt.Sprintf(msg, args...))
		pbrManager.Output()
	})
}
func Info(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Info(fmt.Sprintf(msg, args...))
		pbrManager.Output()
	})
}

func Warn(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Warn(fmt.Sprintf(msg, args...))
		pbrManager.Output()
	})
}

func Error(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Error(fmt.Sprintf(msg, args...))
		pbrManager.Output()
	})
}

func DebugS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Debug(msg, args...)
		pbrManager.Output()
	})
}
func InfoS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Info(msg, args...)
		pbrManager.Output()
	})
}

func WarnS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Warn(msg, args...)
		pbrManager.Output()
	})
}

func ErrorS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Error(msg, args...)
		pbrManager.Output()
	})
}

func DefultLog() *slog.Logger {
	return defaultLog
}

func EnableLogDebug() {
	defaultLog = slog.New(
		&ConsoleHandler{
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		},
	)
}

func SetLogFile(file string) {
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

type ConsoleHandler struct {
	*slog.TextHandler
}

func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	return slog.Default().Handler().Handle(ctx, r)
}

func init() {
	defaultLog = slog.New(
		&ConsoleHandler{
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		},
	)
}
