package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// MakeRequest makes an HTTP request to the given router
func MakeRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, path, nil)
	}

	require.NoError(t, err)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// ParseResponse parses the response body into the given struct
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), v)
	require.NoError(t, err)
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(t, expectedStatus, w.Code, "Expected status %d, got %d. Body: %s", expectedStatus, w.Code, w.Body.String())
}

// SetupTestMode sets Gin to test mode
func SetupTestMode() {
	gin.SetMode(gin.TestMode)
}
