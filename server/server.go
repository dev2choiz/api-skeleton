package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dev2choiz/api-skeleton/pkg/business"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

type Server struct {
	business business.Business
}

func NewServer(b business.Business) *Server {
	return &Server{b}
}

func ApplyRoutes(ser *Server, pub, sec chi.Router) {
	pub.Get("/", ser.Index)
	sec.Get("/users", ser.GetUsers)
	pub.Post("/register", ser.Register)
	pub.Post("/authenticate", ser.Authenticate)
}

func (s Server) Index(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "OK", http.StatusOK)
}

func writeResponse(w http.ResponseWriter, dat any, code int) {
	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	_ = json.NewEncoder(w).Encode(dat)
}

func responseJSON(w http.ResponseWriter, dat any, code int) {
	w.Header().Set("Content-Type", "application/json")

	writeResponse(w, dat, code)
}

func getCodeFromErr(err error) int {
	switch {
	case errors.Is(err, errapp.ErrAppNotFound):
		return http.StatusNotFound
	case errors.Is(err, errapp.ErrAppBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errapp.ErrAppConflict):
		return http.StatusConflict
	case errors.Is(err, errapp.ErrAppUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errapp.ErrAppInternal):
		return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}

func responseErr(w http.ResponseWriter, err error, dat any) {
	responseJSON(w, dat, getCodeFromErr(err))
}
