package jwt

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iam43x/interview-help-api/internal/config"
)

type RSAPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func ParseRSAPair(conf *config.Config) (*RSAPair, error) {
	private, err := os.ReadFile(conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %v: %w", conf.PrivateKey, err)
	}
	privatePEM, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки приватной части ключа для jwt: %w", err)
	}
	public, err := os.ReadFile(conf.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %v: %w", conf.PublicKey, err)
	}
	publicPEM, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки публичной части ключа для jwt: %w", err)
	}
	return &RSAPair{
		privateKey: privatePEM,
		publicKey:  publicPEM,
	}, nil
}
