package utils_test

import (
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	password              = "password"
	FirstEncodedPassword  string
	secondEncodedPassword string
)

func TestGenerateHash(t *testing.T) {
	hash, err := utils.GenerateHash(password)

	assert.Nil(t, err, "there should be no error in hashing")

	assert.Contains(t, hash, "$argon2id$v=19$m=65536,t=3,p=2$", "argon2 hashing is not working properly")

	anotherHash, err := utils.GenerateHash("password")

	assert.Nil(t, err, "there should be no error in hashing")

	assert.NotEqual(t, hash, anotherHash, "salt is not working")

	FirstEncodedPassword = hash
	secondEncodedPassword = anotherHash
}

func TestComparePasswordAndHash(t *testing.T) {
	FirstMatch, err := utils.ComparePasswordAndHash(password, FirstEncodedPassword)

	assert.NoError(t, err, "there should be no error in comparing the password")

	assert.True(t, FirstMatch, "password should match!")

	SecondMatch, err := utils.ComparePasswordAndHash(password, secondEncodedPassword)

	assert.NoError(t, err, "there should be no error in comparing the password")

	assert.True(t, SecondMatch, "password should match!")
}
