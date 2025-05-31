package service

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/janobono/captcha-service/internal/config"
	"github.com/janobono/go-util/security"
	"sync"
	"time"
)

type JwtService struct {
	appConfig *config.AppConfig

	mutex     sync.Mutex
	jwtToken  *security.JwtToken
	keyBuffer []keyEntry
}

type keyEntry struct {
	kid       string
	publicKey *rsa.PublicKey
}

func NewJwtService(appConfig *config.AppConfig) *JwtService {
	return &JwtService{
		appConfig: appConfig,
		keyBuffer: make([]keyEntry, 0, 2),
	}
}

func (j *JwtService) getJwtToken() (*security.JwtToken, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	now := time.Now().UTC()

	if j.jwtToken != nil && now.Before(j.jwtToken.KeyExpiration()) {
		return j.jwtToken, nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	publicKey := &privateKey.PublicKey
	kid := j.randomKid()

	// Add to buffer
	j.addToKeyBuffer(kid, publicKey)

	token := security.NewJwtToken(
		jwt.SigningMethodRS256,
		privateKey,
		publicKey,
		kid,
		j.appConfig.TokenIssuer,
		j.appConfig.TokenExpiresIn,
		now.Add(j.appConfig.TokenJwkExpiresIn*time.Second),
		j.getPublicKey,
	)

	j.jwtToken = token
	return token, nil
}

func (j *JwtService) getPublicKey(kid string) (interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	for _, entry := range j.keyBuffer {
		if entry.kid == kid {
			return entry.publicKey, nil
		}
	}
	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func (j *JwtService) addToKeyBuffer(kid string, key *rsa.PublicKey) {
	entry := keyEntry{kid: kid, publicKey: key}

	// Append new key
	j.keyBuffer = append(j.keyBuffer, entry)

	// Keep only last 2
	if len(j.keyBuffer) > 2 {
		j.keyBuffer = j.keyBuffer[len(j.keyBuffer)-2:]
	}
}

func (j *JwtService) randomKid() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("kid-%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}
