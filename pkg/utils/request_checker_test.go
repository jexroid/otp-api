package utils_test

import (
	"testing"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/models"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// You CANT run these test in `t.Parallel()` because of memory management confusion
// Channeling will work but, it's not necessary here

func TestSigningChecker(t *testing.T) {
	var err error

	err = utils.SignInChecker(api.LoginRequest{Phone: 989174330422, Password: "password"})
	assert.NoError(t, err)

	err = utils.SignInChecker(api.LoginRequest{Phone: 174330422, Password: "password"})
	assert.Error(t, err, "Shouldn't be able to sign in with invalid phone number: 174330422")
	assert.ErrorContains(t, err, "phone: must be no less than 989000000000")

	err = utils.SignInChecker(api.LoginRequest{Phone: 1000174330422, Password: "password"})
	assert.Error(t, err, "Shouldn't be able to sign in with invalid phone number: 174330422")
	assert.ErrorContains(t, err, "phone: must be no greater than 989999999999")

	err = utils.SignInChecker(api.LoginRequest{Phone: 989174330422, Password: "am"})
	assert.Error(t, err, "Shouldn't be able to sign in with invalid password am")
	assert.ErrorContains(t, err, "password: the length must be between 4 and 40")

	err = utils.SignInChecker(api.LoginRequest{Phone: 989174330422, Password: "amaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	assert.Error(t, err, "Shouldn't be able to sign in with invalid password am")
	assert.ErrorContains(t, err, "password: the length must be between 4 and 40")
}

func TestSignupCheker(t *testing.T) {
	var SecErr error

	SecErr = utils.SignUpChecker(models.User{UUID: pkg.Uuid(), Phone: 989174330422, Firstname: "amirreza", Lastname: "farzan", Password: "password"})
	assert.NoError(t, SecErr)

	SecErr = utils.SignUpChecker(models.User{UUID: pkg.Uuid(), Phone: 989174330422, Firstname: "غیر انگلیسی", Lastname: "farzan", Password: "password"})
	assert.NoError(t, SecErr)

	SecErr = utils.SignUpChecker(models.User{UUID: pkg.Uuid(), Phone: 9174330422, Firstname: "amirreza", Lastname: "farzan", Password: "password"})
	assert.Error(t, SecErr)
	assert.ErrorContains(t, SecErr, "phone: must be no less than 989000000000")

	SecErr = utils.SignUpChecker(models.User{Phone: 9174330422, Firstname: "amirreza", Lastname: "farzan", Password: "password"})
	t.Log(SecErr)
	assert.Error(t, SecErr)
	assert.ErrorContains(t, SecErr, "phone: must be no less than 989000000000")
}
