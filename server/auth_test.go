package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/mocks/mockbusiness"
	"github.com/dev2choiz/api-skeleton/pkg/business"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

func TestServer_Register(t *testing.T) {
	body := registerPayload{
		Username: "user",
		Password: "secret",
	}

	tests := []struct {
		name     string
		want     entity.User
		reqBody  io.Reader
		busMock  func(t *testing.T) business.Business
		wantCode int
	}{
		{
			"bad formed body content",
			entity.User{},
			strings.NewReader("bad json"),
			func(t *testing.T) business.Business {
				return mockbusiness.NewMockBusiness(t)
			},
			http.StatusBadRequest,
		},
		{
			"error from business",
			entity.User{},
			getRequestBody(t, body),
			func(t *testing.T) business.Business {
				bus := mockbusiness.NewMockBusiness(t)
				bus.EXPECT().
					Register(mock.Anything, body.Username, body.Password).
					Return(entity.User{}, errapp.ErrAppBadRequest)

				return bus
			},
			http.StatusBadRequest,
		},
		{
			"register successfully",
			entity.User{Username: body.Username},
			getRequestBody(t, body),
			func(t *testing.T) business.Business {
				bus := mockbusiness.NewMockBusiness(t)
				bus.EXPECT().
					Register(mock.Anything, body.Username, body.Password).
					Return(entity.User{Username: body.Username}, nil)

				return bus
			},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := tt.busMock(t)

			srv := NewServer(bus)

			req, err := http.NewRequest("POST", "/register", tt.reqBody)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			srv.Register(rr, req)
			res := rr.Result()
			assert.Equal(t, res.StatusCode, tt.wantCode)

			got := entity.User{}
			err = json.NewDecoder(res.Body).Decode(&got)

			assert.NoError(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestServer_Aunthenticate(t *testing.T) {
	body := authPayload{
		Username: "user",
		Password: "secret",
	}

	tests := []struct {
		name     string
		want     authResponse
		reqBody  io.Reader
		busMock  func(t *testing.T) business.Business
		wantCode int
	}{
		{
			"bad formed body content",
			authResponse{},
			strings.NewReader("bad json"),
			func(t *testing.T) business.Business {
				return mockbusiness.NewMockBusiness(t)
			},
			http.StatusBadRequest,
		},
		{
			"error from business",
			authResponse{},
			getRequestBody(t, body),
			func(t *testing.T) business.Business {
				bus := mockbusiness.NewMockBusiness(t)
				bus.EXPECT().
					Authenticate(mock.Anything, body.Username, body.Password).
					Return(mock.Anything, errapp.ErrAppUnauthorized)

				return bus
			},
			http.StatusUnauthorized,
		},
		{
			"authenticate successfully",
			authResponse{Token: "tok"},
			getRequestBody(t, body),
			func(t *testing.T) business.Business {
				bus := mockbusiness.NewMockBusiness(t)
				bus.EXPECT().
					Authenticate(mock.Anything, body.Username, body.Password).
					Return("tok", nil)

				return bus
			},
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := tt.busMock(t)

			srv := NewServer(bus)

			req, err := http.NewRequest("POST", "/register", tt.reqBody)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			srv.Authenticate(rr, req)
			res := rr.Result()
			assert.Equal(t, res.StatusCode, tt.wantCode)

			got := authResponse{}
			err = json.NewDecoder(res.Body).Decode(&got)

			assert.NoError(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}
