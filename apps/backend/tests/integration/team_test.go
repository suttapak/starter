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

func TestTeamCreate(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and login
	testUser := testutil.CreateTestUser(t, ts.DB, "teamowner@example.com", "teamowner", "password123")
	require.NotNil(t, testUser)

	loginPayload := service.LoginDto{
		UserNameEmail: "teamowner",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		payload        service.CreateTeamDto
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful team creation",
			payload: service.CreateTeamDto{
				Name:        "Test Team",
				Username:    "test",
				Description: "This is a test team",
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Contains(t, data, "id")
				assert.NotEmpty(t, data["id"])
			},
		},
		{
			name: "missing required fields",
			payload: service.CreateTeamDto{
				Name: "", // Empty name
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name: "unauthorized - no token",
			payload: service.CreateTeamDto{
				Name:        "Unauthorized Team",
				Username:    "test",
				Description: "Should fail",
			},
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ts.Router, "POST", "/teams/", tt.payload, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestTeamGetTeamsMe(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and teams
	testUser := testutil.CreateTestUser(t, ts.DB, "member@example.com", "member", "password123")
	require.NotNil(t, testUser)

	team1 := testutil.CreateTestTeam(t, ts.DB, "Team One", "username1", "First team", int(testUser.ID))
	team2 := testutil.CreateTestTeam(t, ts.DB, "Team Two", "username2", "Second team", int(testUser.ID))
	require.NotNil(t, team1)
	require.NotNil(t, team2)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "member",
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
			name: "get user's teams successfully",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 2, "Should have at least 2 teams")
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
			w := testutil.MakeRequest(t, ts.Router, "GET", "/teams/me", nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestTeamGetTeamByTeamId(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and team
	testUser := testutil.CreateTestUser(t, ts.DB, "teamviewer@example.com", "teamviewer", "password123")
	require.NotNil(t, testUser)

	team := testutil.CreateTestTeam(t, ts.DB, "Viewable Team", "team_username", "Team to view", int(testUser.ID))
	require.NotNil(t, team)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "teamviewer",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		teamID         int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "get team by id successfully",
			teamID: int(team.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Viewable Team", data["name"])
			},
		},
		{
			name:   "team not found",
			teamID: 99999,
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name:           "unauthorized - no token",
			teamID:         int(team.ID),
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d", tt.teamID)
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

func TestTeamGetTeamMembers(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test users and team
	owner := testutil.CreateTestUser(t, ts.DB, "owner@example.com", "owner", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Member Team", "member", "Team with members", int(owner.ID))
	require.NotNil(t, team)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "owner",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		teamID         int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "get team members successfully",
			teamID: int(team.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 1, "Should have at least 1 member (owner)")
			},
		},
		{
			name:   "team not found",
			teamID: 99999,
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/members", tt.teamID)
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

func TestTeamGetMemberCount(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and team
	owner := testutil.CreateTestUser(t, ts.DB, "countowner@example.com", "countowner", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Count Team", "count team", "Team for counting", int(owner.ID))
	require.NotNil(t, team)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "countowner",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		teamID         int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "get member count successfully",
			teamID: int(team.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				// Count should be at least 1 (the owner)
				count := body["data"].(float64)
				assert.GreaterOrEqual(t, count, float64(1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/member-count", tt.teamID)
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

func TestTeamGetTeamsFilter(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and teams
	testUser := testutil.CreateTestUser(t, ts.DB, "filterer@example.com", "filterer", "password123")
	require.NotNil(t, testUser)

	team1 := testutil.CreateTestTeam(t, ts.DB, "Alpha Team", "alpha_team", "First team", int(testUser.ID))
	team2 := testutil.CreateTestTeam(t, ts.DB, "Beta Team", "beta_team", "Second team", int(testUser.ID))
	require.NotNil(t, team1)
	require.NotNil(t, team2)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "filterer",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		queryParams    string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:        "get all teams without filter",
			queryParams: "",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
			},
		},
		{
			name:        "filter teams by name",
			queryParams: "?name=Alpha",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/teams/" + tt.queryParams
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

func TestTeamUpdateTeamInfo(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and team
	owner := testutil.CreateTestUser(t, ts.DB, "updater@example.com", "updater", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Old Team Name", "old team", "Old description", int(owner.ID))
	require.NotNil(t, team)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "updater",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		teamID         int
		payload        service.UpdateTeamInfoRequest
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "update team info successfully",
			teamID: int(team.ID),
			payload: service.UpdateTeamInfoRequest{
				Name:        "New Team Name",
				Username:    "New Username",
				Description: "New description",
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name:   "update with empty name",
			teamID: int(team.ID),
			payload: service.UpdateTeamInfoRequest{
				Name:        "",
				Description: "New description",
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d", tt.teamID)
			w := testutil.MakeRequest(t, ts.Router, "PUT", url, tt.payload, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}
