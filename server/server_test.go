package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dev2choiz/api-skeleton/mocks/mockbusiness"
)

func getRequestBody(t *testing.T, data any) io.Reader {
	body, err := json.Marshal(data)
	assert.NoError(t, err)

	return bytes.NewReader(body)
}

func TestServer_Index(t *testing.T) {
	s := NewServer(mockbusiness.NewMockBusiness(t))
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	s.Index(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "OK")
}
