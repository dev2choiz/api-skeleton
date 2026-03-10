package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

func Test_codeFromErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "not found error",
			err:  errapp.ErrAppNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "bad request error",
			err:  errapp.ErrAppBadRequest,
			want: http.StatusBadRequest,
		},
		{
			name: "conflict error",
			err:  errapp.ErrAppConflict,
			want: http.StatusConflict,
		},
		{
			name: "unauthorized error",
			err:  errapp.ErrAppUnauthorized,
			want: http.StatusUnauthorized,
		},
		{
			name: "internal error",
			err:  errapp.ErrAppInternal,
			want: http.StatusInternalServerError,
		},
		{
			name: "wrapped not found error",
			err:  fmt.Errorf("%w", errapp.ErrAppNotFound),
			want: http.StatusNotFound,
		},
		{
			name: "unknown error",
			err:  errors.New("something else"),
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := codeFromErr(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
