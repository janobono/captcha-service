package server

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/janobono/captcha-service/generated/openapi"
	"github.com/janobono/captcha-service/internal/config"
	impl "github.com/janobono/captcha-service/internal/server/impl"
	"github.com/janobono/captcha-service/internal/service"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	config         *config.ServerConfig
	captchaService *service.CaptchaService
}

func NewHttpServer(config *config.ServerConfig, captchaService *service.CaptchaService) *HttpServer {
	return &HttpServer{config, captchaService}
}

func (s *HttpServer) Start() *http.Server {
	slog.Info("Starting http server...")

	handleFunctions := openapi.ApiHandleFunctions{
		CaptchaControllerAPI: impl.NewCaptchaController(s.captchaService),
		HealthControllerAPI:  impl.NewHealthController(),
	}
	router := impl.NewRouter(handleFunctions, s.config.ContextPath)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.CorsConfig.AllowedOrigins,
		AllowMethods:     s.config.CorsConfig.AllowedMethods,
		AllowHeaders:     s.config.CorsConfig.AllowedHeaders,
		ExposeHeaders:    s.config.CorsConfig.ExposedHeaders,
		AllowCredentials: s.config.CorsConfig.AllowCredentials,
		MaxAge:           s.config.CorsConfig.MaxAge,
	}))

	httpServer := &http.Server{
		Addr:    s.config.HTTPAddress,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to serve", "error", err)
			panic(err)
		}
	}()

	slog.Info("Http server started", "port", s.config.HTTPAddress)
	return httpServer
}
