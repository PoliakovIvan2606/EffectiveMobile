package logger

import (
	"log/slog"
	"os"
)

var LevelsLog = map[string]slog.Level{
    "DEBUG": slog.LevelDebug,
    "INFO":  slog.LevelInfo,
    "WARN":  slog.LevelWarn,
    "ERROR": slog.LevelError,
}

func InitLogger(level string) {
    opts := &slog.HandlerOptions{
        Level: LevelsLog[level], // минимальный уровень
    }
    handler := slog.NewTextHandler(os.Stdout, opts)
    slog.SetDefault(slog.New(handler))
}