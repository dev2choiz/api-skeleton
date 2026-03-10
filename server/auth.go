package server

import (
	"encoding/json"
	"net/http"

	"github.com/dev2choiz/api-skeleton/pkg/httpx"
)

type registerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Server) Register(w http.ResponseWriter, r *http.Request) {
	pay := registerPayload{}

	err := json.NewDecoder(r.Body).Decode(&pay)
	if err != nil {
		httpx.ResponseJSON(w, nil, http.StatusBadRequest)
		return
	}

	user, err := s.business.Register(r.Context(), pay.Username, pay.Password)
	if err != nil {
		httpx.ResponseErr(w, err, nil)
		return
	}

	httpx.ResponseJSON(w, user, http.StatusOK)
}

type authPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (s Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	pay := authPayload{}

	err := json.NewDecoder(r.Body).Decode(&pay)
	if err != nil {
		httpx.ResponseJSON(w, nil, http.StatusBadRequest)
		return
	}

	tok, err := s.business.Authenticate(r.Context(), pay.Username, pay.Password)
	if err != nil {
		httpx.ResponseErr(w, err, nil)
		return
	}

	httpx.ResponseJSON(w, authResponse{Token: tok}, http.StatusOK)
}
