package vcr

import (
	"bytes"
	assert "github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestModifyHTTPRequestBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/", bytes.NewBufferString("abc"))
	assert.Nil(t, err)
	assert.Equal(t, int64(3), req.ContentLength)

	ModifyHTTPRequestBody(req, func(input string) string {
		assert.Equal(t, input, "abc")
		return "foofoo"
	})

	assert.Equal(t, int64(6), req.ContentLength)
	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, "foofoo", string(body))
}

func TestModifyHTTPRequestBodyWithNilBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), req.ContentLength)

	ModifyHTTPRequestBody(req, func(input string) string {
		return "foofoo"
	})

	assert.Equal(t, int64(0), req.ContentLength)
	assert.Equal(t, req.Body, nil)
}
