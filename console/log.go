package console

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
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
func Success(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Info(color.GreenString(fmt.Sprintf(msg, args...)))
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
func Fatal(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Error(fmt.Sprintf(msg, args...))
		os.Exit(1)
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
func SuccessS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Info(color.GreenString(msg), args...)
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
func FatalS(msg string, args ...interface{}) {
	withOutputLock(func() {
		defaultLog.Error(msg, args...)
		os.Exit(1)
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
