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

func TestProductCreate(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and category
	owner := testutil.CreateTestUser(t, ts.DB, "productowner@example.com", "productowner", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Product Team", "product_team", "Team for products", int(owner.ID))
	require.NotNil(t, team)

	category := testutil.CreateTestProductCategory(t, ts.DB, "Test Category", int(team.ID))
	require.NotNil(t, category)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "productowner",
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
		payload        service.CreateProductRequest
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "create product successfully",
			teamID: int(team.ID),
			payload: service.CreateProductRequest{
				Name:        "Test Product",
				Description: "A test product",
				UOM:         "ชิ้น",
				Price:       9999, // Price in cents
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Test Product", data["name"])
				assert.NotEmpty(t, data["id"])
			},
		},
		{
			name:   "create product with missing required fields",
			teamID: int(team.ID),
			payload: service.CreateProductRequest{
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
			name:   "unauthorized - no token",
			teamID: int(team.ID),
			payload: service.CreateProductRequest{
				Name:  "Unauthorized Product",
				Price: 5000,
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
			url := fmt.Sprintf("/teams/%d/products", tt.teamID)
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

func TestProductGetProduct(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and product
	owner := testutil.CreateTestUser(t, ts.DB, "productget@example.com", "productget", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Get Team", "get team", "Team for getting", int(owner.ID))
	require.NotNil(t, team)

	product := testutil.CreateTestProduct(t, ts.DB, "p001", "Get Product", "Product to get", 9999, int(team.ID))
	require.NotNil(t, product)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "productget",
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
		productID      int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:      "get product successfully",
			teamID:    int(team.ID),
			productID: int(product.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "Get Product", data["name"])
			},
		},
		{
			name:      "product not found",
			teamID:    int(team.ID),
			productID: 99999,
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
			teamID:         int(team.ID),
			productID:      int(product.ID),
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/products/%d", tt.teamID, tt.productID)
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

func TestProductGetProducts(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and products
	owner := testutil.CreateTestUser(t, ts.DB, "productlist@example.com", "productlist", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "List Team", "list_team", "Team for listing", int(owner.ID))
	require.NotNil(t, team)

	product1 := testutil.CreateTestProduct(t, ts.DB, "p001", "Product Alpha", "First product", 10.00, int(team.ID))
	product2 := testutil.CreateTestProduct(t, ts.DB, "p002", "Product Beta", "Second product", 20.00, int(team.ID))
	require.NotNil(t, product1)
	require.NotNil(t, product2)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "productlist",
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
		queryParams    string
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:        "get all products for team",
			teamID:      int(team.ID),
			queryParams: "",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
				data := body["data"].([]interface{})
				assert.GreaterOrEqual(t, len(data), 2, "Should have at least 2 products")
			},
		},
		{
			name:        "filter products by name",
			teamID:      int(team.ID),
			queryParams: "?name=Alpha",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "data")
			},
		},
		{
			name:           "unauthorized - no token",
			teamID:         int(team.ID),
			queryParams:    "",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/teams/%d/products%s", tt.teamID, tt.queryParams)
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

func TestProductUpdate(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and product
	owner := testutil.CreateTestUser(t, ts.DB, "productupdate@example.com", "productupdate", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Update Team", "update_team", "Team for updating", int(owner.ID))
	require.NotNil(t, team)

	product := testutil.CreateTestProduct(t, ts.DB, "p002", "Old Name", "Old description", 50.00, int(team.ID))
	require.NotNil(t, product)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "productupdate",
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
		productID      int
		payload        service.UpdateProductRequest
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:      "update product successfully",
			teamID:    int(team.ID),
			productID: int(product.ID),
			payload: service.UpdateProductRequest{
				Code:        "p002",
				Name:        "New Name",
				Description: "New description",
				UOM:         "ชิ้น",
				Price:       10000, // $100.00
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
			name:      "update with invalid data",
			teamID:    int(team.ID),
			productID: int(product.ID),
			payload: service.UpdateProductRequest{
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
			name:      "product not found",
			teamID:    int(team.ID),
			productID: 99999,
			payload: service.UpdateProductRequest{
				Code:        "p002",
				Name:        "New Name",
				Description: "New description",
				UOM:         "ชิ้น",
				Price:       10000, // $100.00
			},
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
			url := fmt.Sprintf("/teams/%d/products/%d", tt.teamID, tt.productID)
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

func TestProductDelete(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Teardown(t)

	// Create test user, team, and product
	owner := testutil.CreateTestUser(t, ts.DB, "productdelete@example.com", "productdelete", "password123")
	require.NotNil(t, owner)

	team := testutil.CreateTestTeam(t, ts.DB, "Delete Team", "delete_team", "Team for deleting", int(owner.ID))
	require.NotNil(t, team)

	product := testutil.CreateTestProduct(t, ts.DB, "p002", "Delete Me", "Product to delete", 25.00, int(team.ID))
	require.NotNil(t, product)

	// Login
	loginPayload := service.LoginDto{
		UserNameEmail: "productdelete",
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
		productID      int
		headers        map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:      "delete product successfully",
			teamID:    int(team.ID),
			productID: int(product.ID),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Contains(t, body, "message")
			},
		},
		{
			name:      "product not found",
			teamID:    int(team.ID),
			productID: 99999,
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
			url := fmt.Sprintf("/teams/%d/products/%d", tt.teamID, tt.productID)
			w := testutil.MakeRequest(t, ts.Router, "DELETE", url, nil, tt.headers)
			testutil.AssertStatusCode(t, w, tt.expectedStatus)

			var response map[string]interface{}
			testutil.ParseResponse(t, w, &response)

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}
