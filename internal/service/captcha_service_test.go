package service

import (
	"context"
	"github.com/janobono/captcha-service/internal/config"
	"strings"
	"testing"
	"time"
)

func TestCaptchaService(t *testing.T) {
	appConfig := &config.AppConfig{
		Characters:        "abcdefghijklmnopqrstuvwxyz0123456789",
		TextLength:        8,
		Width:             200,
		Height:            70,
		Font:              "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		FontSize:          36,
		NoiseLines:        8,
		TokenIssuer:       "captcha-service",
		TokenExpiresIn:    time.Duration(30) * time.Minute,
		TokenJwkExpiresIn: time.Duration(720) * time.Minute,
	}

	captchaService := NewCaptchaService(appConfig)

	ctx := context.Background()

	// Step 1: Generate new CAPTCHA
	detail, err := captchaService.Create(ctx)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if detail.CaptchaToken == "" {
		t.Error("CaptchaToken is empty")
	}
	if !strings.HasPrefix(detail.CaptchaImage, "data:image/png;base64,") {
		t.Error("CaptchaImage is not a valid base64 image")
	}

	// Step 2: Validate wrong value
	wrong := &CaptchaData{
		CaptchaToken: detail.CaptchaToken,
		CaptchaValue: "WRONG",
	}
	if captchaService.Validate(ctx, wrong) {
		t.Error("Expected CAPTCHA validation to fail with incorrect input")
	}
}
