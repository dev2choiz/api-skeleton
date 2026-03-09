package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getQueryString(t *testing.T) {
	tests := []struct {
		name  string
		query string
		key   string
		want  *string
	}{
		{
			name:  "missing key",
			query: "",
			key:   "foo",
			want:  nil,
		},
		{
			name:  "empty value",
			query: "foo=",
			key:   "foo",
			want:  nil,
		},
		{
			name:  "valid value",
			query: "foo=bar",
			key:   "foo",
			want:  new("bar"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			got := getQueryString(req, tt.key)

			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_getQueryInt(t *testing.T) {
	tests := []struct {
		name  string
		query string
		key   string
		want  *int
	}{
		{
			name:  "missing key",
			query: "",
			key:   "foo",
			want:  nil,
		},
		{
			name:  "empty value",
			query: "foo=",
			key:   "foo",
			want:  nil,
		},
		{
			name:  "valid int",
			query: "foo=4",
			key:   "foo",
			want:  new(4),
		},
		{
			name:  "invalid int",
			query: "foo=abc",
			key:   "foo",
			want:  new(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			got := getQueryInt(req, tt.key)

			assert.Equal(t, got, tt.want)
		})
	}
}
