package server

import (
	"github.com/janobono/captcha-service/generated/proto"
	"github.com/janobono/captcha-service/internal/config"
	"github.com/janobono/captcha-service/internal/server/impl"
	"github.com/janobono/captcha-service/internal/service"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GrpcServer struct {
	config         *config.ServerConfig
	captchaService *service.CaptchaService
}

func NewGrpcServer(config *config.ServerConfig, captchaService *service.CaptchaService) *GrpcServer {
	return &GrpcServer{config, captchaService}
}

func (s *GrpcServer) Start() *grpc.Server {
	slog.Info("Starting gRPC server...")

	lis, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterCaptchaServer(grpcServer, impl.NewCaptchaServer(s.captchaService))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to serve", "error", err)
			panic(err)
		}
	}()

	slog.Info("gRPC server started", "port", s.config.GRPCAddress)
	return grpcServer
}
