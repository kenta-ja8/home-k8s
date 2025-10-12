package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
)

var (
	slogLogger = slog.New(newJSONHandler())
)

func Init(cfg *entity.Config) {
	if cfg.IS_LOCAL {
		slogLogger = slog.New(newTextHandler())
	}

}

func newJSONHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
}

func newTextHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
}

func log(level slog.Level, format string, args ...any) {
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	slogLogger.Log(context.Background(), level, message)
}

func Error(format string, args ...any) {
	log(slog.LevelError, format, args...)
}

func Warn(format string, args ...any) {
	log(slog.LevelWarn, format, args...)
}

func Info(format string, args ...any) {
	log(slog.LevelInfo, format, args...)
}

func Debug(format string, args ...any) {
	log(slog.LevelDebug, format, args...)
}
