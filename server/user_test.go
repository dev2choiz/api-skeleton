package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/mocks/mockbusiness"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

func Test_getUserFilters(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  entity.UserFilters
	}{
		{
			"with empty query",
			"",
			entity.UserFilters{},
		},
		{
			"with all filter fields",
			"id=id&username=username&firstname=firstname&lastname=lastname&limit=2",
			entity.UserFilters{
				ID:        new("id"),
				Username:  new("username"),
				Firstname: new("firstname"),
				Lastname:  new("lastname"),
				Limit:     new(2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/users?"+tt.query, nil)
			assert.NoError(t, err)

			got := getUserFilters(req)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestServer_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		want     []entity.User
		busErr   error
		wantCode int
	}{
		{
			"business return an error",
			nil,
			errapp.ErrAppInternal,
			http.StatusInternalServerError,
		},
		{
			"business return a list of users",
			[]entity.User{{ID: "id", Username: "geralt"}},
			nil,
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := mockbusiness.NewMockBusiness(t)
			bus.EXPECT().GetUsers(context.Background(), mock.Anything).
				Return(tt.want, tt.busErr)

			srv := NewServer(bus)

			req, err := http.NewRequest("GET", "/users", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			srv.GetUsers(rr, req)
			res := rr.Result()
			assert.Equal(t, res.StatusCode, tt.wantCode)

			got := []entity.User{}
			err = json.NewDecoder(res.Body).Decode(&got)
			assert.NoError(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}
