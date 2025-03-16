package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/iam43x/interview-help-api/internal/jwt"
	"github.com/iam43x/interview-help-api/internal/store"
	"github.com/iam43x/interview-help-api/internal/util"
)

type RegistrationAPI struct {
	issuer *jwt.Issuer
	store  *store.Store
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Invite   string `json:"invite"`
}

type RegistrationResponse struct {
	Token string `json:"token"`
}

func NewRegistrationAPI(i *jwt.Issuer, s *store.Store) *RegistrationAPI {
	return &RegistrationAPI{
		issuer: i,
		store:  s,
	}
}

func (rg *RegistrationAPI) RegistrationHttpHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendErrorResponse(w, http.StatusForbidden, err.Error())
		return
	}
	if err := rg.store.ExistsInvite(req.Invite); err != nil {
		log.Printf("ERROR %s", err.Error())
		util.SendErrorResponse(w, http.StatusBadRequest, "Invite corrupted")
		return
	}
	hp, err := util.HashedPass(req.Password)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while registration")
		return
	}
	u, err := rg.store.CreateUser(req.Login, req.Name, hp, req.Invite)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while registration")
		return
	}
	t, err := rg.issuer.GenerateJWT(u)
	if err != nil {
		util.SendErrorResponse(w, http.StatusServiceUnavailable, "error while generate jwt")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegistrationResponse{Token: t})
}
