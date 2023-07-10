package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"gorm.io/gorm"
)

func TestCreateAdmin(t *testing.T) {

	type fields struct {
		db *gorm.DB
	}

	type args struct {
		ctx          context.Context
		newAdminInfo request.NewAdminInfo
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       domain.Admin
		wantErr    bool
	}{
		{
			name: "fail create admin",
			args: args{
				ctx:          context.TODO(),
				newAdminInfo: request.NewAdminInfo{UserName: "Ganesh", Email: "ganesh@gmail.com", Phone: "1234567895", Password: "ajujacob"},
			},
			beforeTest: func(mockSql sqlmock.Sqlmock) {
				mockSql.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
						 	VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *;`,
				)).WithArgs("Ganesh", "ganesh@gmail.com", "1234567895", "ajujacob").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			want:    domain.Admin{ID: 1, UserName: "Ganesh", Email: "ganesh@gmail.com", Phone: "1234567895", Password: "ajujacob", IsSuperAdmin: false, IsBlocked: false},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			db, mock, err := sqlmock.New()
		})
	}

}
