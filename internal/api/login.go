package api

import (
	"encoding/json"
	"net/http"

	"github.com/iam43x/interview-help-4u/internal/domain"
	"github.com/iam43x/interview-help-4u/internal/jwt"
	"github.com/iam43x/interview-help-4u/internal/util"
)

type LoginAPI struct {
	Issuer *jwt.Issuer
}

func NewLoginAPI(i *jwt.Issuer) *LoginAPI {
	return &LoginAPI{Issuer: i}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (l *LoginAPI) LoginHttpHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}
	if req.Username != "test" && req.Password != "test" {
		util.SendErrorResponse(w, http.StatusUnauthorized, "Unknown user")
		return
	}
	token, err := l.Issuer.GenerateJWT(&domain.User{Name: req.Username})
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}