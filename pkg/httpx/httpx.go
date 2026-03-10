package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func WriteResponse(w http.ResponseWriter, dat any, code int) {
	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	_ = json.NewEncoder(w).Encode(dat)
}

func ResponseJSON(w http.ResponseWriter, dat any, code int) {
	w.Header().Set("Content-Type", "application/json")

	WriteResponse(w, dat, code)
}

func codeFromErr(err error) int {
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

func ResponseErr(w http.ResponseWriter, err error, dat any) {
	ResponseJSON(w, dat, codeFromErr(err))
}
