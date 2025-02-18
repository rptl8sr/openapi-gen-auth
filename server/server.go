package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"openapi-gen-auth/api"
	"openapi-gen-auth/jwt"
	"openapi-gen-auth/middleware"
	"openapi-gen-auth/service"
)

var _ api.ServerInterface = (*Implementation)(nil)

type Implementation struct {
	jwtSecret string
}

func NewServerImplementation(jwtSecret string) *Implementation {
	return &Implementation{
		jwtSecret: jwtSecret,
	}
}

func (s *Implementation) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var req api.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := service.YourOwnGetOrCreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
	}

	token, err := jwt.GenerateToken(userID, s.jwtSecret)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
	}

	respondJSON(w, http.StatusOK, api.AuthResponse{Token: &token})
}

func (s *Implementation) GetHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Implementation) GetApiPrivate(w http.ResponseWriter, r *http.Request) {
	scopes := r.Context().Value(middleware.Scopes)
	if scopes == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userIDRaw := r.Context().Value(middleware.ContextKeyUserID)
	if userIDRaw == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("invalid user id: %v", userIDRaw))
	}

	respondJSON(w, http.StatusOK, api.PrivateResponse{Username: &userID})
}

func (s *Implementation) GetApiPublic(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
