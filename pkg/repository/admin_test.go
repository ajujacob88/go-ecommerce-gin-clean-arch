package repository

import (
	"context"
	"errors"
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
				//ID:       1,
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

				// mock.ExpectQuery(`INSERT INTO admins\(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at\)
				//  VALUES\(\$1, \$2, \$3,\$4, false, false, NOW\(\), NOW\(\)\) RETURNING \*`).
				//actually above is correct without using quotemeta, regexp.QuoteMeta returns a string that escapes all regular expression metacharacters inside the argument text; the returned string is a regular expression matching the literal text. so https://pkg.go.dev/regexp#QuoteMeta, in the ofcicial documentation we can check

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
					 VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *`)).
					WithArgs("Rahul", "rahul@gmail.com", "9496074716", "password").
					WillReturnRows(rows)
			},

			expectedError: nil,
		},

		{ // test case for creating a new admin with duplictae phone
			testName: "duplicate phone",
			inputField: request.NewAdminInfo{
				UserName: "Rahul2",
				Email:    "rahul2@gmail.com",
				Phone:    "9496074716",
				Password: "password",
			},
			expectedOutput: domain.Admin{},

			buildStub: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
					 VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *`)).
					WithArgs("Rahul2", "rahul2@gmail.com", "9496074716", "password").
					WillReturnError(errors.New("phone number already exists"))

			},

			expectedError: errors.New("phone number already exists"),
		},

		{ // test case for creating a new admin with duplictae email
			testName: "duplicate phone",
			inputField: request.NewAdminInfo{
				UserName: "Rahul2",
				Email:    "rahul@gmail.com",
				Phone:    "7736832773",
				Password: "password",
			},
			expectedOutput: domain.Admin{},

			buildStub: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO admins(user_name, email, phone, password,is_super_admin, is_blocked, created_at, updated_at)
					 VALUES($1, $2, $3,$4, false, false, NOW(), NOW()) RETURNING *`)).
					WithArgs("Rahul2", "rahul@gmail.com", "7736832773", "password").
					WillReturnError(errors.New("email already exists- value violates unique constraint 'email'"))

			},

			expectedError: errors.New("email already exists- value violates unique constraint 'email'"),
		},
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

			actualOutput, actualError := adminRepository.CreateAdmin(context.Background(), tt.inputField)

			/* This is by using assert from testify package
			if tt.expectedError == nil {
				assert.NoError(t, actualError)
			} else {
				assert.Equal(t, tt.expectedError, actualError)
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}
			*/

			//without using testify assert package, using default testing package

			if tt.expectedError == nil {
				if actualError != nil {
					t.Errorf("expected no error, but got: %v", actualError)
				}
			} else {
				if tt.expectedError.Error() != actualError.Error() {
					t.Errorf("expected error: %v, but got: %v", tt.expectedError, actualError)
				}
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %+v, but want %+v", actualOutput, tt.expectedOutput)
			}

			// Check that all expectations were met
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

		})
	}

}

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
		{
			testName:       "admin is not a superadmin",
			inputField:     6,
			expectedOutput: false,
			buildStub: func(mock sqlmock.Sqlmock) {
				//columns := []string{"id"}
				mock.ExpectQuery("SELECT is_super_admin	FROM admins	WHERE id = \\$1$").
					WithArgs(6).
					WillReturnRows(sqlmock.NewRows([]string{"is_super_admin"}).AddRow(false))
			},
			expectedError: nil,
		},
		{
			testName:       "admin does not exist",
			inputField:     7,
			expectedOutput: false,
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT is_super_admin FROM admins WHERE id = \\$1$").
					WithArgs(7).
					WillReturnError(errors.New("admin not found"))
			},
			expectedError: errors.New("admin not found"),
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

func TestFindAdmin(t *testing.T) {

	tests := []struct {
		testName       string
		inputField     string
		expectedOutput domain.Admin
		buildStub      func(mock sqlmock.Sqlmock)
		expectedError  error
	}{
		{ // test case for finding a admin
			testName:   "Find admin by email succesfull",
			inputField: "rahul@gmail.com",
			expectedOutput: domain.Admin{
				//ID:       1,
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
				rows := sqlmock.NewRows([]string{"user_name", "email", "phone"}).
					AddRow("Rahul", "rahul@gmail.com", "9496074716")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT *	FROM admins	WHERE email = $1;`)).
					WithArgs("rahul@gmail.com").
					WillReturnRows(rows)
			},

			expectedError: nil,
		},

		{ // test case for finding a admin - case insensitive
			testName:   "Find admin by email succesfull - test case 2- case insensitive",
			inputField: "RAHUL@gmaiL.com",
			expectedOutput: domain.Admin{
				//ID:       1,
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
				rows := sqlmock.NewRows([]string{"user_name", "email", "phone"}).
					AddRow("Rahul", "rahul@gmail.com", "9496074716")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT *	FROM admins	WHERE email = $1;`)).
					WithArgs("RAHUL@gmaiL.com").
					WillReturnRows(rows)
			},

			expectedError: nil,
		},

		{ // test case for finding a admin - special characters
			testName:   "Find admin by email succesfull - test case 3- with special characters entered",
			inputField: "RAHUL@gmaiL.com",
			expectedOutput: domain.Admin{
				//ID:       1,
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
				rows := sqlmock.NewRows([]string{"user_name", "email", "phone"}).
					AddRow("Rahul", "rahul@gmail.com", "9496074716")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT *	FROM admins	WHERE email = $1;`)).
					WithArgs("RAHUL@gmaiL*com").
					WillReturnRows(rows)
			},

			expectedError: nil,
		},

		{ // test case
			testName:       "Find admin by email not succesfull - invalid email provided",
			inputField:     "aju@gmail.com",
			expectedOutput: domain.Admin{},

			buildStub: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT *	FROM admins	WHERE email = $1;`)).
					WithArgs("aju@gmail.com").
					WillReturnError(errors.New("Entered Email does not exist"))
			},
			expectedError: errors.New("Entered Email does not exist"),
		},
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

			actualOutput, actualError := adminRepository.FindAdmin(context.Background(), tt.inputField)

			/* This is by using assert from testify package
			if tt.expectedError == nil {
				assert.NoError(t, actualError)
			} else {
				assert.Equal(t, tt.expectedError, actualError)
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}
			*/

			//without using testify assert package, using default testing package

			if tt.expectedError == nil {
				if actualError != nil {
					t.Errorf("expected no error, but got: %v", actualError)
				}
			} else {
				if tt.expectedError.Error() != actualError.Error() {
					t.Errorf("expected error: %v, but got: %v", tt.expectedError, actualError)
				}
			}

			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %+v, but want %+v", actualOutput, tt.expectedOutput)
			}

			// Check that all expectations were met
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

		})
	}

}
