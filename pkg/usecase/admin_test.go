package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/mock/repositoryMock"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestCreateAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)

	adminMockRepo := repositoryMock.NewMockAdminRepository(ctrl)

	ctx := context.Background()

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
		{
			testName: "only superadmin can create new admins",
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
