package config

import (
	"github.com/janobono/go-util/common"
	"github.com/joho/godotenv"
	"log"
	"strings"
	"time"
)

type ServerConfig struct {
	Prod        bool
	GRPCAddress string
	HTTPAddress string
	ContextPath string
	AppConfig   *AppConfig
	CorsConfig  *CorsConfig
}

type AppConfig struct {
	Characters        string
	TextLength        int
	Width             int
	Height            int
	NoiseLines        int
	Font              string
	FontSize          int
	TokenIssuer       string
	TokenExpiresIn    time.Duration
	TokenJwkExpiresIn time.Duration
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	ExposedHeaders   []string
	MaxAge           time.Duration
}

func InitConfig() *ServerConfig {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("No .env.local file found")
	}

	return &ServerConfig{
		Prod:        common.EnvBool("PROD"),
		GRPCAddress: common.Env("GRPC_ADDRESS"),
		HTTPAddress: common.Env("HTTP_ADDRESS"),
		ContextPath: common.Env("CONTEXT_PATH"),
		AppConfig: &AppConfig{
			Characters:        common.Env("CAPTCHA_CHARACTERS"),
			TextLength:        common.EnvInt("CAPTCHA_TEXT_LENGTH"),
			Width:             common.EnvInt("CAPTCHA_IMAGE_WIDTH"),
			Height:            common.EnvInt("CAPTCHA_IMAGE_HEIGHT"),
			NoiseLines:        common.EnvInt("CAPTCHA_NOISE_LINES"),
			Font:              common.Env("CAPTCHA_FONT"),
			FontSize:          common.EnvInt("CAPTCHA_FONT_SIZE"),
			TokenIssuer:       common.Env("CAPTCHA_TOKEN_ISSUER"),
			TokenExpiresIn:    time.Duration(common.EnvInt("CAPTCHA_TOKEN_EXPIRES_IN")) * time.Minute,
			TokenJwkExpiresIn: time.Duration(common.EnvInt("CAPTCHA_TOKEN_JWK_EXPIRES_IN")) * time.Minute,
		},
		CorsConfig: &CorsConfig{
			AllowedOrigins:   EnvSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods:   EnvSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders:   EnvSlice("CORS_ALLOWED_HEADERS"),
			ExposedHeaders:   EnvSlice("CORS_EXPOSED_HEADERS"),
			AllowCredentials: common.EnvBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           time.Duration(common.EnvInt("CORS_MAX_AGE")) * time.Hour,
		},
	}
}

func EnvSlice(key string) []string {
	value := common.Env(key)
	return strings.Split(value, ",")
}
