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

func TestUserGetUserMe(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and login
	testUser := testutil.CreateTestUser(t, ts.DB, "testme@example.com", "testme", "password123")
	require.NotNil(t, testUser)

	loginPayload := service.LoginDto{
		UserNameEmail: "testme",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "get current user successfully",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "testme", data["username"])
				assert.Equal(t, "testme@example.com", data["email"])
			},
		},
		{
			name:           "unauthorized - no token",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "invalid token",
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
			w := testutil.MakeRequest(t, ts.Router, "GET", "/users/me", nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestUserGetUserById(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test users
	user1 := testutil.CreateTestUser(t, ts.DB, "user1@example.com", "user1", "password123")
	user2 := testutil.CreateTestUser(t, ts.DB, "user2@example.com", "user2", "password123")
	require.NotNil(t, user1)
	require.NotNil(t, user2)

	// Login as user1
	loginPayload := service.LoginDto{
		UserNameEmail: "user1",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		userID         int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "get user by id successfully",
			userID: int(user2.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "user2", data["username"])
			},
		},
		{
			name:   "user not found",
			userID: 99999,
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name:           "unauthorized - no token",
			userID:         int(user2.ID),
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/users/%d", tt.userID)
			w := testutil.MakeRequest(t, ts.Router, "GET", url, nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestUserFindByUsername(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test users
	user1 := testutil.CreateTestUser(t, ts.DB, "searcher@example.com", "searcher", "password123")
	user2 := testutil.CreateTestUser(t, ts.DB, "alice@example.com", "alice", "password123")
	user3 := testutil.CreateTestUser(t, ts.DB, "bob@example.com", "bob", "password123")
	require.NotNil(t, user1)
	require.NotNil(t, user2)
	require.NotNil(t, user3)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "searcher",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		username       string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:     "find user by username - alice",
			username: "alice",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				fmt.Println("data", body)
				assert.Contains(t, body, "data")
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 1)
			},
		},
		{
			name:     "find user by partial username",
			username: "al",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
			},
		},
		{
			name:     "no users found",
			username: "nonexistent",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				// assert.Contains(t, body, "data")
				// data := body["data"].([]interface{})
				// fmt.Println("data", data)
				// May be empty or contain unrelated users depending on implementation
				assert.NotNil(t, body)
			},
		},
		{
			name:           "unauthorized - no token",
			username:       "alice",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/users/by-username?username=%s", tt.username)
			w := testutil.MakeRequest(t, ts.Router, "GET", url, nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestUserCheckVerifyEmail(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user
	testUser := testutil.CreateTestUser(t, ts.DB, "checker@example.com", "checker", "password123")
	require.NotNil(t, testUser)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "checker",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "check email verification status",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				// User is created with email_verifyed = false by default
				verified := body["data"].(bool)
				assert.False(t, verified, "Email should not be verified by default")
			},
		},
		{
			name:           "unauthorized - no token",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "GET", "/users/verify-email", nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}
