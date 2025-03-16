package api

import (
	"encoding/json"
	"net/http"

	"github.com/iam43x/interview-help-api/internal/jwt"
	"github.com/iam43x/interview-help-api/internal/store"
	"github.com/iam43x/interview-help-api/internal/util"
)

type LoginAPI struct {
	issuer *jwt.Issuer
	store  *store.Store
}

func NewLoginAPI(i *jwt.Issuer, s *store.Store) *LoginAPI {
	return &LoginAPI{
		issuer: i,
		store: s,
	}
}

type LoginRequest struct {
	Login    string `json:"login"`
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
	u, err := l.store.GetUserByLogin(req.Login)
	if err != nil {
		util.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := util.VerifyPassword(u.Pass, req.Password); err != nil {
		util.SendErrorResponse(w, http.StatusUnauthorized, "pass not valid")
		return
	}
	token, err := l.issuer.GenerateJWT(u)
	if err != nil {
		util.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
