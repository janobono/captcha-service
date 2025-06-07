package service

import (
	"bytes"
	"context"
	crypto "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/janobono/captcha-service/internal/config"
	"github.com/janobono/go-util/security"
	"image/png"
	"math/big"
	"math/rand"
	"strings"
)

const tokenKey = "encodedText"

type CaptchaDetail struct {
	CaptchaToken string
	CaptchaImage string
}

type CaptchaData struct {
	CaptchaToken string
	CaptchaValue string
}

type CaptchaService struct {
	appConfig       *config.AppConfig
	passwordEncoder *security.PasswordEncoder
	jwtService      *JwtService
}

func NewCaptchaService(appConfig *config.AppConfig) *CaptchaService {
	return &CaptchaService{
		appConfig:       appConfig,
		passwordEncoder: security.NewPasswordEncoder(10),
		jwtService:      NewJwtService(appConfig),
	}
}

func (s *CaptchaService) Create(ctx context.Context) (*CaptchaDetail, error) {
	randomText, err := s.randomString()
	if err != nil {
		return nil, err
	}

	captchaImage, err := s.generateImage(randomText)
	if err != nil {
		return nil, err
	}

	encodedText, err := s.passwordEncoder.Encode(randomText)
	if err != nil {
		return nil, err
	}

	jwtToken, err := s.jwtService.getJwtToken()
	if err != nil {
		return nil, err
	}

	captchaToken, err := jwtToken.GenerateToken(jwt.MapClaims{tokenKey: encodedText})
	if err != nil {
		return nil, err
	}

	return &CaptchaDetail{
		CaptchaToken: captchaToken,
		CaptchaImage: captchaImage,
	}, nil
}

func (s *CaptchaService) Validate(ctx context.Context, captchaData *CaptchaData) bool {
	if captchaData == nil {
		return false
	}

	jwtToken, err := s.jwtService.getJwtToken()
	if err != nil {
		return false
	}

	claims, err := jwtToken.ParseToken(ctx, captchaData.CaptchaToken)
	if err != nil {
		return false
	}

	encodedText := ((*claims)[tokenKey]).(string)

	return s.passwordEncoder.Compare(captchaData.CaptchaValue, encodedText) == nil
}

func (s *CaptchaService) randomString() (string, error) {
	var builder strings.Builder
	builder.Grow(s.appConfig.TextLength)
	for i := 0; i < s.appConfig.TextLength; i++ {
		num, err := crypto.Int(crypto.Reader, big.NewInt(int64(len(s.appConfig.Characters))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		builder.WriteByte(s.appConfig.Characters[num.Int64()])
	}
	return builder.String(), nil
}

func (s *CaptchaService) generateImage(text string) (string, error) {
	width := s.appConfig.Width
	height := s.appConfig.Height

	var rng = rand.New(rand.NewSource(rand.Int63()))

	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Draw a border around the image
	dc.SetRGB(0.8, 0.8, 0.8) // light gray
	dc.DrawRectangle(0, 0, float64(width-1), float64(height-1))
	dc.Stroke()

	// Optional noise lines
	for i := 0; i < s.appConfig.NoiseLines; i++ {
		dc.SetRGBA(rng.Float64(), rng.Float64(), rng.Float64(), 0.3)
		x1 := rng.Float64() * float64(width)
		y1 := rng.Float64() * float64(height)
		x2 := rng.Float64() * float64(width)
		y2 := rng.Float64() * float64(height)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	// Load font
	fontSize := float64(s.appConfig.FontSize)
	if fontSize > float64(height)*0.9 {
		fontSize = float64(height) * 0.9
	}
	if err := dc.LoadFontFace(s.appConfig.Font, fontSize); err != nil {
		return "", fmt.Errorf("failed to load font from %s: %w", s.appConfig.Font, err)
	}

	// Draw each character with distortion and color
	perChar := float64(width) / float64(len(text)+1)
	for i, c := range text {
		x := perChar*float64(i+1) + rng.Float64()*5 - 2.5
		y := float64(height)/2 + rng.Float64()*10 - 5
		angle := rng.Float64()*0.4 - 0.2 // rotate -0.2 to +0.2 radians

		dc.Push() // Save current state

		dc.RotateAbout(angle, x, y)
		dc.SetRGB(rng.Float64(), rng.Float64(), rng.Float64()) // random color
		dc.DrawStringAnchored(string(c), x, y, 0.5, 0.5)

		dc.Pop() // Restore state
	}

	// Encode image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Base64 + MIME prefix
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + b64, nil
}
