package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jexroid/gopi/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "identity",
		Value: "some_value",
		Path:  "/",
	})

	handler.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	expectedBody := "cookie deleted"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body to contain '%s', got '%s'", expectedBody, rr.Body.String())
	}

	// Check if the identity cookie  is deleted
	cookies := rr.Header()["Set-Cookie"]
	t.Log(cookies)
	if len(cookies) != 1 {
		t.Error("Expected 1 Set-Cookie header, got", len(cookies))
	}

	assert.Equal(t, rr.Code, 200, "response was not 200")
}
