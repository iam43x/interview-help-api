package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/iam43x/interview-help-4u/internal/domain"
)

type Issuer struct {
	RSAPair        *RSAPair
	SigningMethod  *jwt.SigningMethodRSA
	ExpirationTime time.Duration
}

func NewIssuer(r *RSAPair, t time.Duration) *Issuer {
	return &Issuer{
		RSAPair:       r,
		SigningMethod: jwt.SigningMethodRS256, // const 
		ExpirationTime: t,
	}
}

func (i *Issuer) GenerateJWT(u *domain.User) (string, error) {
	expirationTime := time.Now().Add(i.ExpirationTime)
	// interface Claims
	claims := jwt.MapClaims{
		"sub": u.Name,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(i.SigningMethod, claims)
	return token.SignedString(i.RSAPair.PrivateKey)
}

func (i *Issuer) ValidateJWT(value string) error {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		return i.RSAPair.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("token forbidden: %v", err)
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("token exp time extraction error: %w", err)
	}
	
	if time.Now().After(exp.Time) {
		return fmt.Errorf("token expire!")
	}

	return nil
}
