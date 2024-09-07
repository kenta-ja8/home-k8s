package logger

import (
	"log/slog"
	"os"
)

var slogLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func Error(msg string, args ...any) {
	slogLogger.Error(msg, args...)
}
func Warn(msg string, args ...any) {
	slogLogger.Warn(msg, args...)
}
func Info(msg string, args ...any) {
	slogLogger.Info(msg, args...)
}
func Debug(msg string, args ...any) {
	slogLogger.Debug(msg, args...)
}
