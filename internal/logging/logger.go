package logging

import (
	"log/slog"
	"os"
)

func Init(level slog.Level, path string) (logger *slog.Logger) {
	defer func() { slog.SetDefault(logger) }()

	var handler slog.Handler

	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fallback to stderr if file opening fails
			hb := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			})
			logger := slog.New(hb)
			logger.Error("Failed to open log file, falling back to stderr", "error", err, "path", path)
			return logger
		}
		hb := slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: level,
		})
		handler = hb
	} else {
		hb := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		handler = hb
	}

	return slog.New(handler)
}
