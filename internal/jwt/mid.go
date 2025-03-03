package jwt

import (
	"net/http"
	"strings"

	"github.com/iam43x/interview-help-4u/internal/util"
)

type JWTMiddleware struct {
	Issuer *Issuer
}

func NewJWTMiddleware(i *Issuer) *JWTMiddleware {
	return &JWTMiddleware{
		Issuer: i,
	}
}

func (m *JWTMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.SendErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Ожидаем заголовок вида "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization format")
			return
		}

		tokenString := parts[1]

		if err := m.Issuer.ValidateJWT(tokenString); err != nil {
			util.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Если все ок — передаем управление дальше
		next.ServeHTTP(w, r)
	})
}
