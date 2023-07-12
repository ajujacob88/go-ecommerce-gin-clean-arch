package repository

import (
	"context"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/domain"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/model/request"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateAdmin(t *testing.T) {

	tests := []struct {
		testName       string
		inputField     request.NewAdminInfo
		expectedOutput domain.Admin
		buildStub      func(mock sqlmock.Sqlmock)
		expectedError  error
	}{
		{ // test case for creating a new admin
			testName: "create admin succesfull",
			inputField: request.NewAdminInfo{
				UserName: "Rahul",
				Email:    "rahul@gmail.com",
				Phone:    "9496074716",
				Password: "password",
			},
			expectedOutput: domain.Admin{
				//ID:           1,
				UserName: "Rahul",
				Email:    "rahul@gmail.com",
				Phone:    "9496074716",
				//Password:     "password",
				IsSuperAdmin: false,
				IsBlocked:    false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},

			buildStub: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_name", "email", "phone", "password"}).
					AddRow("Rahul", "rahul@gmail.com", "9496074716", "password")

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
				 VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *`)).
					//mock.ExpectQuery("INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
					//VALUES(\\$1$, \\$2$, \\$3$,\\$4$, false, false, NOW(), NOW()) RETURNING *;").
					WithArgs("Rahul", "rahul@gmail.com", "9496074716", "password").
					WillReturnRows(rows)
			},

			/*
				buildStub: func(mock sqlmock.Sqlmock) {
					columns := []string{"user_name", "email", "phone", "password"}
					mock.ExpectQuery(`INSERT INTO admins\(user_name, email, phone, password\)
											VALUES\(\$1, \$2, \$3, \$4\)`).
						WithArgs("Rahul", "rahul@gmail.com", "9496074716", "password").
						WillReturnRows(sqlmock.NewRows(columns).FromCSVString("Rahul,sujith@gmail.com,7902638845,sujith@123"))
				},
			*/
			// buildStub: func(mock sqlmock.Sqlmock) {
			// 	columns := []string{"user_name", "email", "phone", "password", "is_super_admin", "is_blocked", "created_at", "updated_at"}
			// 	query := `INSERT INTO admins(user_name, email, phone, password, is_super_admin, is_blocked, created_at, updated_at) VALUES (?, ?, ?, ?, false, false, NOW(), NOW()) RETURNING *`
			// 	mock.ExpectQuery(query).
			// 		WithArgs("Rahul", "rahul@gmail.com", "9496074716", "password").
			// 		WillReturnRows(sqlmock.NewRows(columns).AddRow("Rahul", "rahul@gmail.com", "9496074716", "password", false, false, nil, nil))
			// },

			// buildStub: func(mock sqlmock.Sqlmock) {
			// 	columns := []string{"user_name", "email", "phone", "password", "is_super_admin", "is_blocked", "created_at", "updated_at"}
			// 	query := `INSERT INTO admins(user_name, email, phone, password, is_super_admin, is_blocked, created_at, updated_at) VALUES (?, ?, ?, ?, false, false, NOW(), NOW()) RETURNING *`
			// 	mock.ExpectQuery(query).
			// 		WithArgs("Rahul", "rahul@gmail.com", "9496074716", "password").
			// 		WillReturnRows(sqlmock.NewRows(columns).AddRow("Rahul", "rahul@gmail.com", "9496074716", "password", false, false, nil, nil))
			// },

			expectedError: nil,
		},

		/*
			{ // test case for creating a new admin with duplictae phone
				testName: "duplicate phone",
				inputField: request.NewAdminInfo{
					UserName: "Rahul2",
					Email:    "rahul2@gmail.com",
					Phone:    "9496074716",
					Password: "password",
				},
				expectedOutput: domain.Admin{},
				// buildStub: func(mock sqlmock.Sqlmock) {
				// 	mock.ExpectQuery("^INSERT INTO admins\\(.+\\)$").
				// 		WithArgs("Rahul2", "rahul2@gmail.com", "9496074716", "password").
				// 		WillReturnError(errors.New("phone number already exists"))

				// },


					buildStub: func(mock sqlmock.Sqlmock) {
						mock.ExpectQuery(`^INSERT INTO admins \(user_name, email, phone, password, is_super_admin, is_blocked, created_at, updated_at\) VALUES \(\\$1$, \\$2$, \\$3$, \\$4$, false, false, NOW\(\), NOW\(\)\) RETURNING \*$`).
							WithArgs("Rahul2", "rahul2@gmail.com", "9496074716", "password").
							WillReturnError(errors.New("phone number already exists"))
					},


				expectedError: errors.New("validation error: Field 'UserName' is required"),

				//expectedError: errors.New("phone number already exists"),
			},
		*/
	}

	for _, tt := range tests {

		t.Run(tt.testName, func(t *testing.T) {
			//New() method from sqlmock package create sqlmock database connection and a mock to manage expectations.
			db, mock, err := sqlmock.New()
			//db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			//close the mock db connection after testing.
			defer db.Close()

			//initialize a mock db session
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}

			//create NewUserRepository mock by passing a pointer to gorm.DB
			adminRepository := NewAdminRepository(gormDB)

			tt.buildStub(mock)

			actualOutput, actualError := adminRepository.CreateAdmin(context.TODO(), tt.inputField)

			if tt.expectedError == nil {
				assert.NoError(t, actualError)
			} else {
				assert.Equal(t, tt.expectedError, actualError)
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}
			// Check that all expectations were met
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

		})
	}

}

/* This is correct
func TestIsSuperAdmin(t *testing.T) {
	tests := []struct {
		testName       string
		inputField     int
		expectedOutput bool
		buildStub      func(mock sqlmock.Sqlmock)
		expectedError  error
	}{
		{
			testName:       "admin is superadmin",
			inputField:     5,
			expectedOutput: true,
			buildStub: func(mock sqlmock.Sqlmock) {
				//columns := []string{"id"}
				mock.ExpectQuery("SELECT is_super_admin	FROM admins	WHERE id = \\$1$").
					WithArgs(5).
					WillReturnRows(sqlmock.NewRows([]string{"is_super_admin"}).AddRow(true))
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			//initialize a mock db session
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}

			//create NewUserRepository mock by passing a pointer to gorm.DB
			adminRepository := NewAdminRepository(gormDB)

			// before we actually execute our function, we need to expect required DB actions
			tt.buildStub(mock)

			//call the actual method
			actualOutput, actualErr := adminRepository.IsSuperAdmin(context.TODO(), tt.inputField)

			// validate err is nil if we are not expecting to receive an error
			if tt.expectedError == nil {
				assert.NoError(t, actualErr)
			} else { //validate whether expected and actual errors are same
				assert.Equal(t, tt.expectedError, actualErr)
			}

			assert.Equal(t, tt.expectedOutput, actualOutput)

			// if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
			// 	t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			// }

			// Check that all expectations were met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

		})

	}
}

/* NO NEED - delee this

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
		wantErr    error
	}{
		{
			name: "fail create admin",
			args: args{
				ctx:          context.TODO(),
				newAdminInfo: request.NewAdminInfo{UserName: "Ganesh", Email: "ganesh@gmail.com", Phone: "1234567895", Password: "ajujacob"},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectExec(regexp.QuoteMeta(
					`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
						 	VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *;`,
				)).WithArgs("Ganesh", "ganesh@gmail.com", "1234567895", "ajujacob").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			want:    domain.Admin{ID: 1, UserName: "Ganesh", Email: "ganesh@gmail.com", Phone: "1234567895", Password: "ajujacob", IsSuperAdmin: false, IsBlocked: false},
			wantErr: nil,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			//db, mock, err := sqlmock.New()

			u := &adminDatabase{
				DB: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.CreateAdmin(tt.args.ctx, tt.args.newAdminInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.CreateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.CreateAdmin() = %v, want %v", got, tt.want)
			}
		})

	}


}

*/
