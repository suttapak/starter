package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/internal/service"
)

func TestRegister(t *testing.T) {

	type (
		PepareAuhtService struct {
			output *response.AuthResponse
			err    error
		}
		TestCase struct {
			name              string
			input             dto.UserRegisterDto
			output            controller.Response[*response.AuthResponse]
			pepareAuthService PepareAuhtService
		}
	)
	testCases := []TestCase{
		{
			name: "TEST_SUCCESS",
			input: dto.UserRegisterDto{
				Email:     "email",
				Username:  "username",
				Password:  "password",
				FirstName: "firstname",
				LastName:  "lastname",
			},
			output: controller.Response[*response.AuthResponse]{
				Status:  http.StatusCreated,
				Message: "Success",
				Data: &response.AuthResponse{
					Token:        "token",
					RefreshToken: "retoken",
				},
			},
			pepareAuthService: PepareAuhtService{
				output: &response.AuthResponse{
					Token:        "token",
					RefreshToken: "retoken",
				},
			},
		},
		{
			name:  "TEST_ERROR_BADREQUEST",
			input: dto.UserRegisterDto{},
			output: controller.Response[*response.AuthResponse]{
				Status:  http.StatusBadRequest,
				Message: "Something went wrong",
				Data:    nil,
			},
			pepareAuthService: PepareAuhtService{
				output: &response.AuthResponse{
					Token:        "token",
					RefreshToken: "retoken",
				},
			},
		},
		{
			name: "TEST_ERROR_INTERNAL",
			input: dto.UserRegisterDto{
				Email:     "email",
				Username:  "username",
				Password:  "password",
				FirstName: "firstname",
				LastName:  "lastname",
			},
			output: controller.Response[*response.AuthResponse]{
				Status:  http.StatusBadRequest,
				Message: "ข้อมูลไม่ถูกต้อง",
				Data:    nil,
			},
			pepareAuthService: PepareAuhtService{
				output: nil,
				err:    errs.ErrBadRequest,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// service
			auth := service.NewAuthServiceMock()
			auth.On("Register").Return(tc.pepareAuthService.output, tc.pepareAuthService.err)

			// controller
			control := controller.NewAuth(auth)

			// router
			app := gin.Default()
			app.POST("/auth/register", control.Register)

			w := httptest.NewRecorder()
			dataJson, _ := json.Marshal(tc.input)

			req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(string(dataJson)))
			app.ServeHTTP(w, req)

			assert.Equal(t, tc.output.Status, w.Code)

			ex, _ := json.Marshal(tc.output)

			assert.Equal(t, string(ex), w.Body.String())
		})

	}
}
