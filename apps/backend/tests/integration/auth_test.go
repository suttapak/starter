package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/internal/testutil"
)

func TestAuthRegister(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	tests := []struct {
		name           string
		payload        service.UserRegisterDto
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful registration",
			payload: service.UserRegisterDto{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "refresh_token")
				assert.NotEmpty(t, data["token"])
				assert.NotEmpty(t, data["refresh_token"])
			},
		},
		{
			name: "duplicate username",
			payload: service.UserRegisterDto{
				Username: "testuser", // Same as above
				Email:    "another@example.com",
				Password: "password123",
				FullName: "Another User",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "duplicate email",
			payload: service.UserRegisterDto{
				Username: "anotheruser",
				Email:    "test@example.com", // Same as first test
				Password: "password123",
				FullName: "Another User",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "missing required fields",
			payload: service.UserRegisterDto{
				Username: "incomplete",
				// Missing email, password, fullname
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "password too short",
			payload: service.UserRegisterDto{
				Username: "shortpass",
				Email:    "short@example.com",
				Password: "short", // Less than 8 characters
				FullName: "Short Pass User",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/register", tt.payload, nil)

			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestAuthLogin(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create a test user first
	testUser := testutil.CreateTestUser(t, ts.DB, "login@example.com", "loginuser", "password123")
	require.NotNil(t, testUser)

	tests := []struct {
		name           string
		payload        service.LoginDto
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful login with username",
			payload: service.LoginDto{
				UserNameEmail: "loginuser",
				Password:      "password123",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "refresh_token")
				assert.NotEmpty(t, data["token"])
				assert.NotEmpty(t, data["refresh_token"])
			},
		},
		{
			name: "successful login with email",
			payload: service.LoginDto{
				UserNameEmail: "login@example.com",
				Password:      "password123",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name: "wrong password",
			payload: service.LoginDto{
				UserNameEmail: "loginuser",
				Password:      "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "non-existent user",
			payload: service.LoginDto{
				UserNameEmail: "nonexistent",
				Password:      "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "missing credentials",
			payload: service.LoginDto{
				UserNameEmail: "loginuser",
				// Missing password
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", tt.payload, nil)

			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestAuthRefreshToken(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create a test user
	testUser := testutil.CreateTestUser(t, ts.DB, "refresh@example.com", "refreshuser", "password123")
	require.NotNil(t, testUser)

	// Login to get tokens
	loginPayload := service.LoginDto{
		UserNameEmail: "refreshuser",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)

	data := loginResponse["data"].(map[string]interface{})
	refreshToken := data["refresh_token"].(string)

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful token refresh",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", refreshToken),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "refresh_token")
				assert.NotEmpty(t, data["token"])
			},
		},
		{
			name:           "missing refresh token",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "invalid refresh token",
			headers: map[string]string{
				"Authorization": "Bearer invalid_token",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/refresh", nil, tt.headers)

			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestAuthLogout(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful logout",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/logout", nil, nil)

			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			// Check that the session cookie is cleared
			cookies := w.Result().Cookies()
			var sessionCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == "session" {
					sessionCookie = cookie
					break
				}
			}

			if sessionCookie != nil {
				assert.Equal(t, "", sessionCookie.Value)
				assert.True(t, sessionCookie.MaxAge < 0, "Session cookie should be expired")
			}
		})
	}
}

func TestAuthProtectedRoute(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create a test user and login
	testUser := testutil.CreateTestUser(t, ts.DB, "protected@example.com", "protecteduser", "password123")
	require.NotNil(t, testUser)

	loginPayload := service.LoginDto{
		UserNameEmail: "protecteduser",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)

	data := loginResponse["data"].(map[string]interface{})
	accessToken := data["token"].(string)

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name: "access with valid token",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "access without token",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "access with invalid token",
			headers: map[string]string{
				"Authorization": "Bearer invalid_token",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the protected route /auth/email/send-verify
			w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/email/send-verify", nil, tt.headers)

			testutil.AssertStatusCode(t, w, tt.expectedStatus)
		})
	}
}
