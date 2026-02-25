package server

import (
	"net/http"

	"github.com/dev2choiz/api-skeleton/entity"
)

func getUserFilters(r *http.Request) entity.UserFilters {
	filter := entity.UserFilters{}

	id := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")
	firstname := r.URL.Query().Get("firstname")
	lastname := r.URL.Query().Get("lastname")

	if id != "" {
		filter.ID = &id
	}
	if username != "" {
		filter.Username = &username
	}
	if firstname != "" {
		filter.Firstname = &firstname
	}
	if lastname != "" {
		filter.Lastname = &lastname
	}

	return filter
}

func (s Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.business.GetUsers(r.Context(), getUserFilters(r))
	if err != nil {
		responseErr(w, err, nil)
		return
	}

	responseJSON(w, users, http.StatusOK)
}
