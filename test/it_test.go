package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/janobono/captcha-service/generated/openapi"
	"github.com/janobono/captcha-service/generated/proto"
	"github.com/janobono/captcha-service/internal/config"
	"github.com/janobono/captcha-service/internal/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestIntegrationSomething(t *testing.T) {
	freePorts, err := getFreePorts(2)
	if err != nil {
		t.Fatalf("failed to get free ports: %v", err)
	}

	serverConfig := &config.ServerConfig{
		Prod:        false,
		GRPCAddress: (*freePorts)[0],
		HTTPAddress: (*freePorts)[1],
		ContextPath: "/api",
		AppConfig: &config.AppConfig{
			Characters:        "abcdefghijklmnopqrstuvwxyz0123456789",
			TextLength:        8,
			Width:             200,
			Height:            70,
			NoiseLines:        8,
			Font:              "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
			FontSize:          32,
			TokenIssuer:       "captcha",
			TokenExpiresIn:    time.Duration(30) * time.Minute,
			TokenJwkExpiresIn: time.Duration(720) * time.Minute,
		},
		CorsConfig: &config.CorsConfig{
			AllowedOrigins:   []string{"*"}, // Or restrict to specific domains
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	}

	s := server.NewServer(serverConfig)
	go s.Start()
	time.Sleep(500 * time.Millisecond)

	t.Run("gRPC: Generate and Validate Captcha", func(t *testing.T) {
		conn, err := grpc.NewClient(
			"localhost"+serverConfig.GRPCAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		captchaClient := proto.NewCaptchaClient(conn)
		result, err := captchaClient.Create(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("failed to create captcha in: %v", err)
		}
		t.Logf("captcha created: %v", result)

		valid, err := captchaClient.Validate(context.Background(), &proto.CaptchaData{
			Token: result.Token,
			Text:  "1234",
		})
		assert.NoError(t, err)
		assert.False(t, valid.Value)
	})

	t.Run("REST: livez", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost%s/api/livez", serverConfig.HTTPAddress))
		if err != nil {
			t.Fatalf("REST call failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected status: got %d", resp.StatusCode)
		}

		var result openapi.HealthStatus
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to parse REST response: %v", err)
		}

		if result.Status != "UP" {
			t.Errorf("wrong livez status: got %s", result.Status)
		}
	})

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	time.Sleep(1 * time.Second)
}

func getFreePorts(count int) (*[]string, error) {
	var ports []string
	for i := 0; i < count; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}
		defer l.Close()

		addr := l.Addr().(*net.TCPAddr)
		ports = append(ports, fmt.Sprintf(":%d", addr.Port))
	}
	return &ports, nil
}
