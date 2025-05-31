package server

import (
	"github.com/janobono/captcha-service/internal/config"
	"log/slog"
	"os"
)

func initSlog(config *config.ServerConfig) {
	var handler slog.Handler
	if config.Prod {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	slog.SetDefault(slog.New(handler))
}
