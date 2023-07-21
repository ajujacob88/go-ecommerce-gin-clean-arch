package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/mock/usecaseMock"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminMockUseCase := usecaseMock.NewMockAdminUseCase(ctrl)

	//ctx := context.Background()

	testData := []struct {
		testName   string
		inputField request.NewAdminInfo
		adminID    int
		buildStub  func(adminUseCase *usecaseMock.MockAdminUseCase)
		//expectedCode     int
		expectedResponse response.Response
		//expectedData     domain.Admin
		expectedError error
	}{
		{
			testName: "create admin successful",
			inputField: request.NewAdminInfo{
				UserName: "Rahul",
				Email:    "rahul@gmail.com",
				Phone:    "9496074716",
				Password: "password",
			},
			adminID: 1,
			buildStub: func(adminUseCase *usecaseMock.MockAdminUseCase) {

				adminUseCase.EXPECT().CreateAdmin(gomock.Any(), request.NewAdminInfo{
					UserName: "Rahul",
					Email:    "rahul@gmail.com",
					Phone:    "9496074716",
					Password: "password",
				}, 1).
					Times(1). //how many times the CreateUser usecase should be called
					Return(domain.Admin{
						UserName:     "Rahul",
						Email:        "rahul@gmail.com",
						Phone:        "9496074716",
						IsSuperAdmin: false,
						IsBlocked:    false,
						CreatedAt:    time.Time{},
						UpdatedAt:    time.Time{},
					}, nil)
			},
			//expectedCode: 201,
			expectedResponse: response.SuccessResponse(200, "admin created successfully", domain.Admin{
				UserName:     "Rahul",
				Email:        "rahul@gmail.com",
				Phone:        "9496074716",
				IsSuperAdmin: false,
				IsBlocked:    false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			}),

			expectedError: nil,
		},

		{
			testName: "duplicate admin",
			inputField: request.NewAdminInfo{
				UserName: "Rahul2",
				Email:    "rahul@gmail.com",
				Phone:    "9496074716",
				Password: "password",
			},
			adminID: 2,
			buildStub: func(adminUseCase *usecaseMock.MockAdminUseCase) {

				adminUseCase.EXPECT().CreateAdmin(gomock.Any(), request.NewAdminInfo{
					UserName: "Rahul2",
					Email:    "rahul@gmail.com",
					Phone:    "9496074716",
					Password: "password",
				}, 2).
					Times(1). //how many times the CreateUser usecase should be called
					Return(domain.Admin{}, errors.New("user already exists"))
			},
			//expectedCode:     400,
			expectedResponse: response.ErrorResponse(400, "failed to create the admin", "user already exists", nil),

			expectedError: fmt.Errorf("failed to create the admin"),
		},
	}

	for _, tt := range testData {
		t.Run(tt.testName, func(t *testing.T) {
			tt.buildStub(adminMockUseCase)

			adminHandler := AdminHandler{
				adminUseCase: adminMockUseCase,
			}

			// gin.Default will create a new engine instance with logger middleware by default
			engine := gin.Default()
			//NewRecorder from httptest package will create a recorder which records the response
			recorder := httptest.NewRecorder()
			//create new route for testing
			engine.POST("/admin", adminHandler.CreateAdmin)

			// Since the CreateAdmin function has a Gin context parameter,
			// we need to create a dummy Gin context for testing purposes.
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/admin", nil)
			c.Set("adminID", strconv.Itoa(tt.adminID))

			// Marshal the inputField (test data) to JSON and set it as the request body
			jsonData, err := json.Marshal(tt.inputField)
			assert.NoError(t, err)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
			c.Request.ContentLength = int64(len(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateAdmin handler function
			adminHandler.CreateAdmin(c)

			//url for the test
			url := "/admin"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			//recorder will record the response, and req is the mock req that we created with test data
			engine.ServeHTTP(recorder, req)

			//actual will hold the actual response
			var actual response.Response
			//unmarshalling json data to response.Response format
			err = json.Unmarshal(recorder.Body.Bytes(), &actual)

			// Convert the expectedResponse to JSON
			expectedJSON, err := json.Marshal(tt.expectedResponse)
			assert.NoError(t, err)

			// Assert the response

			// assert.Equal(t, http.StatusCreated, w.Code) // Ensure that the status code is 201
			assert.Equal(t, tt.expectedResponse.StatusCode, w.Code)
			assert.JSONEq(t, string(expectedJSON), w.Body.String())

			// // without using assert
			// actualResponse := w.Body.String()
			// if string(expectedJSON) != actualResponse {
			// 	t.Errorf("expected response does not match actual response:\nExpected: %s\nActual: %s", string(expectedJSON), actualResponse)
			// }

		})
	}
}

//no need.. not correct, only partial
/*

func TestCreateAdmin(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminMockUseCase := usecaseMock.NewMockAdminUseCase(ctrl)

	ctx := context.Background()

	testData := []struct {
		testName       string
		inputField     request.NewAdminInfo
		adminID        int
		expectedOutput domain.Admin
		buildStub      func(adminUseCase *usecaseMock.MockAdminUseCase)
		expectedError  error
	}{
		{
			testName: "create admin successful",
			inputField: request.NewAdminInfo{
				UserName: "Rahul",
				Email:    "rahul@gmail.com",
				Phone:    "9496074716",
				Password: "password",
			},
			adminID: 1,
			expectedOutput: domain.Admin{
				UserName:     "Rahul",
				Email:        "rahul@gmail.com",
				Phone:        "9496074716",
				IsSuperAdmin: false,
				IsBlocked:    false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
			buildStub: func(adminUseCase *usecaseMock.MockAdminUseCase) {
				adminUseCase.EXPECT().CreateAdmin(ctx, gomock.Any(), 1).Return(domain.Admin{
					UserName:     "Rahul",
					Email:        "rahul@gmail.com",
					Phone:        "9496074716",
					IsSuperAdmin: false,
					IsBlocked:    false,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				}, nil)
			},

			expectedError: nil,
		},
		// {
		// 	testName: "duplicate admin",
		// 	inputField: request.NewAdminInfo{
		// 		UserName: "Rahul",
		// 		Email:    "rahul@gmail.com",
		// 		Phone:    "9496074716",
		// 		Password: "password",
		// 	},
		// 	adminID: 1,
		// 	expectedOutput: domain.Admin{
		// 		UserName:     "Rahul",
		// 		Email:        "rahul@gmail.com",
		// 		Phone:        "9496074716",
		// 		IsSuperAdmin: false,
		// 		IsBlocked:    false,
		// 		CreatedAt:    time.Time{},
		// 		UpdatedAt:    time.Time{},
		// 	},
		// 	buildStub: func(adminUseCase *usecaseMock.MockAdminUseCase) {
		// 		adminUseCase.EXPECT().CreateAdmin(ctx, gomock.Any(), 1).Return(domain.Admin{}, errors.New("failed to create the admin"))
		// 	},

		// 	expectedError: errors.New("failed to create the admin"),
		// },
	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			data.buildStub(adminMockUseCase)

			adminHandler := AdminHandler{
				adminUseCase: adminMockUseCase,
			}

			// Since the CreateAdmin function has a Gin context parameter,
			// we need to create a dummy Gin context for testing purposes.
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/admin", nil)
			c.Set("adminID", strconv.Itoa(data.adminID))

			// Marshal the inputField (test data) to JSON and set it as the request body
			jsonData, err := json.Marshal(data.inputField)
			assert.NoError(t, err)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
			c.Request.ContentLength = int64(len(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateAdmin handler function
			adminHandler.CreateAdmin(c)

			// Assert the response
			assert.Equal(t, http.StatusCreated, w.Code) // Ensure that the status code is 201
			//assert.JSONEq(t, `{"message": "admin created successfully", "data": {"id":0, "user_name":"Rahul","email":"rahul@gmail.com","phone_no":"9496074716","password":"","is_super_admin":false,"IsBlocked":false,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}}`, w.Body.String())
			//check if data is of type map[string]interface{}

			expectedResponse := `{"status_code":201,"message":"admin created successfully","data":[{"id":0,"user_name":"Rahul","email":"rahul@gmail.com","phone_no":"9496074716","password":"","is_super_admin":false,"IsBlocked":false,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}]}`
			assert.JSONEq(t, expectedResponse, w.Body.String())

			// without using assert
			// actualResponse := w.Body.String()
			// if expectedResponse != actualResponse {
			// 	t.Errorf("expected response does not match actual response:\nExpected: %s\nActual: %s", expectedResponse, actualResponse)
			// }
		})
	}
}
*/
