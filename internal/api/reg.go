package api

import (
	"encoding/json"
	"net/http"

	"github.com/iam43x/interview-help-api/internal/jwt"
	"github.com/iam43x/interview-help-api/internal/store"
	"github.com/iam43x/interview-help-api/internal/util"
)

type RegistrationAPI struct {
	Issuer *jwt.Issuer
	Store  *store.Store
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Invite   string `json:"invite"`
}

type RegistrationResponse struct {
	Token string
}

func NewRegistrationAPI(i *jwt.Issuer, s *store.Store) *RegistrationAPI {
	return &RegistrationAPI{
		Issuer: i,
		Store:  s,
	}
}

func (rg *RegistrationAPI) RegistrationHttpHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}
	inviteParam := r.URL.Query().Get("invite")
	if inviteParam != "test" {
		util.SendErrorResponse(w, http.StatusBadRequest, "Invite corrupted")
		return
	}
	hp, err := util.HashedPass(req.Password)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while registration")
		return
	}
	u, err := rg.Store.CreateUser(req.Login, req.Name, hp, req.Invite)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while registration")
		return
	}
	t, err := rg.Issuer.GenerateJWT(u)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while generate jwt")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegistrationResponse{Token: t})
}
