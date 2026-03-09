package server

import (
	"net/http"

	"github.com/dev2choiz/api-skeleton/entity"
)

func getUserFilters(r *http.Request) entity.UserFilters {
	return entity.UserFilters{
		ID:        getQueryString(r, "id"),
		Username:  getQueryString(r, "username"),
		Firstname: getQueryString(r, "firstname"),
		Lastname:  getQueryString(r, "lastname"),
		Limit:     getQueryInt(r, "limit"),
	}
}

func (s Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.business.GetUsers(r.Context(), getUserFilters(r))
	if err != nil {
		responseErr(w, err, nil)
		return
	}

	responseJSON(w, users, http.StatusOK)
}
