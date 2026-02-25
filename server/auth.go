package server

import (
	"encoding/json"
	"net/http"
)

type registerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Server) Register(w http.ResponseWriter, r *http.Request) {
	pay := registerPayload{}

	err := json.NewDecoder(r.Body).Decode(&pay)
	if err != nil {
		responseJSON(w, nil, http.StatusBadRequest)
		return
	}

	user, err := s.business.Register(r.Context(), pay.Username, pay.Password)
	if err != nil {
		responseErr(w, err, nil)
		return
	}

	responseJSON(w, user, http.StatusOK)
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
		responseJSON(w, nil, http.StatusBadRequest)
		return
	}

	tok, err := s.business.Authenticate(r.Context(), pay.Username, pay.Password)
	if err != nil {
		responseErr(w, err, nil)
		return
	}

	responseJSON(w, authResponse{Token: tok}, http.StatusOK)
}
