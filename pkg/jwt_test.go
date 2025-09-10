package pkg_test

import (
	"github.com/jexroid/gopi/pkg"

	"github.com/stretchr/testify/assert"
	"testing"
)

var jwt string

func TestCreateToken(t *testing.T) {
	jwt = pkg.CreateToken(pkg.Uuid(), 989174330422)

	assert.Equal(t, len(jwt), 157)
	assert.Contains(t, jwt, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ", "header didn't generated correctly")
}

func TestValidateToken(t *testing.T) {
	valid, payload := pkg.ValidateToken(jwt)

	assert.True(t, valid)

	phone, ok := payload["phone"]
	if !ok {
		t.Errorf("Phone key not found in payload")
	}

	phoneInt := int64(phone.(float64))
	if phoneInt != 989174330422 {
		t.Errorf("Expected phone number 9174330422, got %d", phoneInt)
	}
}
