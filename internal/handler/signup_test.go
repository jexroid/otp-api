package handler_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/models"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func (m *MockDB) MockUser() models.User {
	args := m.Called()
	result := args.Get(0)
	return result.(models.User)
}

type MockDatabase struct {
	DB *MockDB
}

func (m *MockDB) Create(user *models.User) (*gorm.DB, error) {
	args := m.Called(user)
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func userinit(u *models.User) (*models.User, error) {
	u.UUID = pkg.Uuid()
	hashpass, err := utils.GenerateHash(u.Password)
	if err != nil {
		logrus.Error("[...1...] error in making hash ", err)
		utils.Telegram("[...1...] error in hashing")
		anotherTry, anotherError := utils.GenerateHash(u.Password)
		if anotherError != nil {
			utils.Telegram("[...2...] error in hashing")
			logrus.Error("[...2...] error in making hash ", err)
			return u, err
		}
		logrus.Info("[..1..] second time was successfully")
		u.Password = anotherTry
		return u, nil
	}
	u.Password = hashpass
	return u, nil
}

func TestSignup(t *testing.T) {
	// Create a mock DB
	mockDB := &MockDB{}

	// Define mock behavior for Create (modify as needed)
	var response = gorm.DB{Error: nil, RowsAffected: 1}
	mockDB.On("Create", mock.AnythingOfType("*models.User")).Return(&response, nil)

	// Create a Database with the mock DB
	db := &MockDatabase{DB: mockDB}

	var user = models.User{UUID: pkg.Uuid(), Firstname: "amir", Lastname: "Farzan", Phone: 989174330422, Password: "password"}

	SuccessQuery, _ := db.DB.Create(&user)
	logrus.Error(SuccessQuery.Error)
	if nil != SuccessQuery.Error {
		if SuccessQuery.Error.Error() == "ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)" {
			_, err := json.Marshal(&api.RegisterResponse{
				Ok:        true,
				UserExist: true,
			})
			assert.NoError(t, err)
		}
	}

	var failed = gorm.DB{Error: errors.New("ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)"), RowsAffected: 0}
	mockDB.On("Create", mock.AnythingOfType("*api.User")).Return(&failed, nil)

	FailedQuery, _ := db.DB.Create(&user)
	logrus.Error(FailedQuery.Error)
	if nil != FailedQuery.Error {
		if FailedQuery.Error.Error() == "ERROR: duplicate key value violates unique constraint \"User_phone_key\" (SQLSTATE 23505)" {
			_, err := json.Marshal(&api.RegisterResponse{
				Ok:        true,
				UserExist: true,
			})
			assert.NoError(t, err)
		}
	}
	jwt := pkg.CreateToken(user.UUID, user.Phone)
	logrus.Info(jwt, user)

	jsonresponse, err := json.Marshal(&api.RegisterResponse{Ok: true, Message: "User successfully registered! use OTP to verify"})
	assert.Nil(t, err)
	t.Log(string(jsonresponse))
	assert.Contains(t, string(jsonresponse), `"Ok":true`)
	assert.Contains(t, string(jsonresponse), `"UserExist":false`)
}
