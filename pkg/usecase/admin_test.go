package usecase

import (
	"context"
	"testing"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/mock/repositoryMock"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/response"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

/*
func TestCreateAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)

	adminMockRepo := repositoryMock.NewMockAdminRepository(ctrl)

	ctx := context.Background()

	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), 10)
	// hashedpwdstring := string(hashedPassword)
	// fmt.Println("hashed pwd 1", hashedpwdstring)

	testData := []struct {
		testName       string
		inputField     request.NewAdminInfo
		adminID        int
		expectedOutput domain.Admin
		buildStub      func(adminRepo *repositoryMock.MockAdminRepository)
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
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				adminRepo.EXPECT().IsSuperAdmin(ctx, 1).Return(true, nil)
				adminRepo.EXPECT().CreateAdmin(ctx, gomock.Any()).Return(domain.Admin{
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

		/*
			{
				testName: "create admin successful2",
				inputField: request.NewAdminInfo{
					UserName: "Rahul",
					Email:    "rahul@gmail.com",
					Phone:    "9496074716",
					Password: hashedpwdstring,
				},
				adminID: 1,
				expectedOutput: domain.Admin{
					UserName:     "Rahul",
					Email:        "rahul@gmail.com",
					Phone:        "9496074716",
					Password:     hashedpwdstring, // Update this with the expected hashed password
					IsSuperAdmin: false,
					IsBlocked:    false,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				},
				buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
					adminRepo.EXPECT().IsSuperAdmin(ctx, 1).Return(true, nil)
					// Hash the password in the CreateAdmin method and compare with the expected hashed password
					//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
					//fmt.Println("string of hashed pw is", string(hashedPassword))

					//	fmt.Println("hashed pwd 2", hashedpwdstring)
					fmt.Println("newadminpwd is")
					adminRepo.EXPECT().CreateAdmin(ctx, gomock.Any()).Do(func(ctx context.Context, newAdmin request.NewAdminInfo) {
						fmt.Println("nea admin is", newAdmin)
						// Hash the password in the CreateAdmin method and compare with the expected hashed password
						//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
						//fmt.Println("hashed pwd 3", hashedpwdstring, "and", newAdmin.Password)
						//assert.Equal(t, hashedpwdstring, newAdmin.Password)
					}).Return(domain.Admin{
						UserName:     "Rahul",
						Email:        "rahul@gmail.com",
						Phone:        "9496074716",
						Password:     hashedpwdstring, // Update this with the expected hashed password
						IsSuperAdmin: false,
						IsBlocked:    false,
						CreatedAt:    time.Time{},
						UpdatedAt:    time.Time{},
					}, nil)
				},
				expectedError: nil,
			},

		{
			testName: "non-existent superadmin",
			inputField: request.NewAdminInfo{
				UserName: "John",
				Email:    "john@gmail.com",
				Phone:    "1234567890",
				Password: "password",
			},
			adminID:        2,
			expectedOutput: domain.Admin{},
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				adminRepo.EXPECT().IsSuperAdmin(ctx, 2).Return(false, nil)
			},
			expectedError: fmt.Errorf("Only superadmin can create the new admins"),
		},

		{
			testName: "error while checking superadmin status",
			inputField: request.NewAdminInfo{
				UserName: "Jane",
				Email:    "jane@gmail.com",
				Phone:    "9876543210",
				Password: "password",
			},
			adminID:        4,
			expectedOutput: domain.Admin{},
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				adminRepo.EXPECT().IsSuperAdmin(ctx, 4).Return(false, fmt.Errorf("DB connection error"))
			},
			expectedError: fmt.Errorf("DB connection error"),
		},

		{
			testName: "error while creating admin in the repository",
			inputField: request.NewAdminInfo{
				UserName: "Jane",
				Email:    "jane@gmail.com",
				Phone:    "9876543210",
				Password: "password",
			},
			adminID:        5,
			expectedOutput: domain.Admin{},
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				adminRepo.EXPECT().IsSuperAdmin(ctx, 5).Return(true, nil)
				adminRepo.EXPECT().CreateAdmin(ctx, gomock.Any()).Return(domain.Admin{}, fmt.Errorf("Database error"))
			},
			expectedError: fmt.Errorf("Database error"),
		},
	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			data.buildStub(adminMockRepo)

			adminUseCase := adminUseCase{
				adminRepo: adminMockRepo,
			}

			result, err := adminUseCase.CreateAdmin(ctx, data.inputField, data.adminID)

			assert.Equal(t, data.expectedError, err)
			assert.Equal(t, data.expectedOutput, result)
		})
	}

	ctrl.Finish()
}

*/

/*
func TestAdminLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	adminMockRepo := repositoryMock.NewMockAdminRepository(ctrl)

	ctx := context.Background()

	testData := []struct {
		testName       string
		inputField     request.AdminLoginInfo
		adminInfo      domain.Admin
		buildStub      func(adminRepo *repositoryMock.MockAdminRepository)
		expectedToken  string
		expectedOutput response.AdminDataOutput
		expectedError  error
	}{
		{
			testName: "admin login successful",
			inputField: request.AdminLoginInfo{
				Email:    "admin@example.com",
				Password: "password",
			},
			adminInfo: domain.Admin{
				Email:    "admin@example.com",
				Password: "$2a$10$1UjMdeCD4/u846kZb1538e8wB3XEKg27bXzygb9VkBxke6MMkOp0G",
				// Include other relevant admin data here
			},
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				adminRepo.EXPECT().FindAdmin(ctx, "admin@example.com").Return(domain.Admin{
					Email:    "admin@example.com",
					Password: "$2a$10$1UjMdeCD4/u846kZb1538e8wB3XEKg27bXzygb9VkBxke6MMkOp0G",
					// Include other relevant admin data here
				}, nil)
			},
			expectedToken:  "jwt-token",
			expectedOutput: response.AdminDataOutput{
				// Set the expected output values based on the admin data
			},
			expectedError: nil,
		},
		// Include additional test cases for different scenarios

	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			data.buildStub(adminMockRepo)

			adminUseCase := adminUseCase{
				adminRepo: adminMockRepo,
			}

			token, output, err := adminUseCase.AdminLogin(ctx, data.inputField)

			assert.Equal(t, data.expectedError, err)
			assert.Equal(t, data.expectedToken, token)
			assert.Equal(t, data.expectedOutput, output)
		})
	}

	ctrl.Finish()
}
*/

func TestAdminLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	adminMockRepo := repositoryMock.NewMockAdminRepository(ctrl)

	ctx := context.Background()

	testData := []struct {
		testName       string
		inputField     request.AdminLoginInfo
		adminInfo      domain.Admin
		buildStub      func(adminRepo *repositoryMock.MockAdminRepository)
		expectedToken  string
		expectedOutput response.AdminDataOutput
		expectedError  error
	}{
		{
			testName: "admin login successful",
			inputField: request.AdminLoginInfo{
				Email:    "aju@gmail.com",
				Password: "amal@123",
			},
			adminInfo: domain.Admin{
				UserName:     "aju",
				Email:        "aju@gmail.com",
				Phone:        "1234567890",
				Password:     "amal@123", // Update this with the expected hashed password
				IsSuperAdmin: false,
				IsBlocked:    false,
			},
			buildStub: func(adminRepo *repositoryMock.MockAdminRepository) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("amal@123"), 10)
				if err != nil {
					t.Fatalf("failed to generate hash from password in build stub %v", err)
				}
				adminRepo.EXPECT().FindAdmin(ctx, "aju@gmail.com").Return(domain.Admin{
					UserName:     "aju",
					Email:        "aju@gmail.com",
					Phone:        "1234567890",
					Password:     string(hashedPassword), // Update this with the expected hashed password
					IsSuperAdmin: false,
					IsBlocked:    false,
				}, nil)
			},
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIwMTMyMTAsImlkIjowfQ.gjjQqDo0hhr34eoQ5BeXzKE_cdVYDKJjLnQ-w_qZRqA",
			expectedOutput: response.AdminDataOutput{
				ID:           0,
				UserName:     "aju",
				Email:        "aju@gmail.com",
				Phone:        "1234567890",
				IsSuperAdmin: false,
			},
			expectedError: nil,
		},
		// Include additional test cases for different scenarios

	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			data.buildStub(adminMockRepo)

			adminUseCase := adminUseCase{
				adminRepo: adminMockRepo,
			}

			token, actualOutput, actualError := adminUseCase.AdminLogin(ctx, data.inputField)

			//assert.Equal(t, data.expectedError, actualError)
			// //assert.Equal(t, data.expectedToken, token)
			//assert.Equal(t, data.expectedOutput, actualOutput)

			//using direct if condition instead of assert
			if actualError != nil {
				t.Errorf("Expected no error, but got: %v", actualError)
			}

			if token == "" {
				t.Error("Expected token, but got an empty string")
			}

			if actualOutput != data.expectedOutput {
				t.Errorf("Expected output: %v, but got: %v", data.expectedOutput, actualOutput)
			}
		})
	}

	ctrl.Finish()
}
