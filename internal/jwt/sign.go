package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/iam43x/interview-help-api/internal/domain"
)

type Issuer struct {
	rSAPair        *RSAPair
	signingMethod  *jwt.SigningMethodRSA
	expirationTime time.Duration
}

func NewIssuer(r *RSAPair, t time.Duration) *Issuer {
	return &Issuer{
		rSAPair:        r,
		signingMethod:  jwt.SigningMethodRS256, // const
		expirationTime: t,
	}
}

func (i *Issuer) GenerateJWT(u *domain.User) (string, error) {
	expirationTime := time.Now().Add(i.expirationTime)
	// interface Claims
	claims := jwt.MapClaims{
		"sub": u.Login,
		"iss": "interview-help-api",
		"exp": expirationTime.Unix(),
		"u_name": u.Name,
	}
	token := jwt.NewWithClaims(i.signingMethod, claims)
	return token.SignedString(i.rSAPair.privateKey)
}

func (i *Issuer) ValidateJWT(value string) error {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		return i.rSAPair.publicKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("token forbidden: %v", err)
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("token exp time extraction error: %w", err)
	}

	if time.Now().After(exp.Time) {
		return fmt.Errorf("token expire after: %v", exp.Time)
	}

	return nil
}
