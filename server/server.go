package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dev2choiz/api-skeleton/pkg/business"
	"github.com/dev2choiz/api-skeleton/pkg/httpx"
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
	httpx.WriteResponse(w, "OK", http.StatusOK)
}
