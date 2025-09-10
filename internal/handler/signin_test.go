package handler_test

import (
	"testing"

	"github.com/jexroid/gopi/pkg/models"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) UserMock() models.User {
	args := m.Called()
	result := args.Get(0)
	return result.(models.User)
}

func TestSignin(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDB)

	mockDB.On("MockUser").Return(models.User{UUID: 1, Phone: 989174330422, Firstname: "amirreza", Lastname: "farzan", Password: "password"})

	// for mocking test watch the Signup_test.go im quite bad at mocking. i was willing to create mock interface with mockery but due to lack of time ...

	//// Define expected behavior for First method
	//expectedUser := api.User{Phone: 1234567890, Password: "hashed_password"}
	//mockDB.On("First", &api.User{}, "phone = ?", 1234567890).Return(nil, nil, &expectedUser)
	//
	//// Create a valid login request
	//loginRequest := api.LoginRequest{Phone: 1234567890, Password: "valid_password"}
	//
	//// Create a router for testing
	//router := mux.NewRouter()
	//router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
	//	handler.Signin(handler.Database{DB: mockDB}, w, r)
	//})
	//
	//// Create a test request
	//requestBody, err := json.Marshal(loginRequest)
	//assert.NoError(t, err)
	//req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(requestBody))
	//req.Header.Set("Content-Type", "application/json")
	//
	//// Record the response
	//rec := httptest.NewRecorder()
	//router.ServeHTTP(rec, req)
	//
	//// Assert response status code
	//assert.Equal(t, http.StatusOK, rec.Code)
	//
	//// Assert response body
	//var response api.LoginResponse
	//err = json.Unmarshal(rec.Body.Bytes(), &response)
	//assert.NoError(t, err)
	//assert.True(t, response.Ok)
	//assert.True(t, response.Valid)
	//assert.NotEmpty(t, response.Token)
	//
	//// Verify mock interactions
	//mockDB.AssertExpectations(t)
	//utils.AssertExpectations(t)
}
