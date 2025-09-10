package crud

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockDB) Read(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockDB) All(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockDB) Update(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockDB) Delete(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestCreateHandler(t *testing.T) {
	mockDB := new(MockDB)

	reqBody := []byte(`{"name": "Test Car"}`)
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	mockDB.On("Create", rec, req).Return(nil)

	handler := http.HandlerFunc(mockDB.Create)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestAllHandler(t *testing.T) {
	mockDB := new(MockDB)

	req, err := http.NewRequest("GET", "/cars", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	mockDB.On("All", rec, req).Return(nil)

	handler := http.HandlerFunc(mockDB.All)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteHandler(t *testing.T) {
	mockDB := new(MockDB)

	req, err := http.NewRequest("DELETE", "/cars/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	mockDB.On("Delete", rec, req).Return(nil)

	handler := http.HandlerFunc(mockDB.Delete)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

func TestUpdateHandler(t *testing.T) {
	mockDB := new(MockDB)

	reqBody := []byte(`{"name": "Updated Car"}`)
	req, err := http.NewRequest("PUT", "/cars/123", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	mockDB.On("Update", rec, req).Return(nil)

	handler := http.HandlerFunc(mockDB.Update)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestReadHandler(t *testing.T) {
	mockDB := new(MockDB)

	req, err := http.NewRequest("GET", "/cars/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	mockDB.On("Read", rec, req).Return(nil)

	handler := http.HandlerFunc(mockDB.Read)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// Write similar tests for other CRUD operations (Read, All, Update, Delete)
