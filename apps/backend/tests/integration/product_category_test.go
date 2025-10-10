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

func TestProductCategoryCreate(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user and team
	owner := testutil.CreateTestUser(t, ts.DB, "categoryowner@example.com", "categoryowner", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Category Team", "category_team", "Team for categories", int(owner.ID))
	require.NotNil(t, team)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "categoryowner",
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
		payload        service.CreateProductCategoryRequest
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "create category successfully",
			teamID: int(team.ID),
			payload: service.CreateProductCategoryRequest{
				TeamId: team.ID,
				Name:   "Electronics",
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
			name:   "missing required fields",
			teamID: int(team.ID),
			payload: service.CreateProductCategoryRequest{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/product_category", tt.teamID)
			w := testutil.MakeRequest(t, ts.Router, "POST", url, tt.payload, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestProductCategoryGet(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and category
	owner := testutil.CreateTestUser(t, ts.DB, "categoryget@example.com", "categoryget", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Get Category Team", "get_category_team", "Team for getting categories", int(owner.ID))
	require.NotNil(t, team)

	category := testutil.CreateTestProductCategory(t, ts.DB, "Get Category", int(team.ID))
	require.NotNil(t, category)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "categoryget",
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
		categoryID     int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:       "get category successfully",
			teamID:     int(team.ID),
			categoryID: int(category.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Get Category", data["name"])
			},
		},
		{
			name:       "category not found",
			teamID:     int(team.ID),
			categoryID: 99999,
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/product_category/%d", tt.teamID, tt.categoryID)
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

func TestProductCategoryGetAll(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and categories
	owner := testutil.CreateTestUser(t, ts.DB, "categorylist@example.com", "categorylist", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "List Category Team", "list_category_team", "Team for listing categories", int(owner.ID))
	require.NotNil(t, team)

	cat1 := testutil.CreateTestProductCategory(t, ts.DB, "Category 1", int(team.ID))
	cat2 := testutil.CreateTestProductCategory(t, ts.DB, "Category 2", int(team.ID))
	require.NotNil(t, cat1)
	require.NotNil(t, cat2)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "categorylist",
		Password:      "password123",
	}
	loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)
	require.Equal(t, http.StatusCreated, loginResp.Code)

	var loginResponse map[string]interface{}
	testutil.ParseResponse(t, loginResp, &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	url := fmt.Sprintf("/teams/%d/product_category", int(team.ID))
	w := testutil.MakeRequest(t, ts.Router, "GET", url, nil, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	})
	testutil.AssertStatusCode(t, w, http.StatusOK)

	var response map[string]interface{}
	testutil.ParseResponse(t, w, &response)

	assert.Contains(t, response, "data")
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2, "Should have at least 2 categories")
}
