package server

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getRequestBody(t *testing.T, data any) io.Reader {
	body, err := json.Marshal(data)
	assert.NoError(t, err)

	return bytes.NewReader(body)
}
