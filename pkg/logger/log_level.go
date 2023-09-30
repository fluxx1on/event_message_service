package logger

import (
	"strings"

	"log/slog"
)

func GetLevel(s string) slog.Level {
	str := strings.ToUpper(s)

	switch str {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}

	return -4
}
