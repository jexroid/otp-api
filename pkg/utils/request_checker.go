package utils

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg/models"
)

func SignUpChecker(request models.User) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, validation.Min(989000000000), validation.Max(989999999999)),
		validation.Field(&request.Firstname, validation.Required, validation.Length(3, 65)),
		validation.Field(&request.Lastname, validation.Required, validation.Length(3, 65)),
		validation.Field(&request.Password, validation.Required, validation.Length(4, 40)),
	)
}

func SignInChecker(request api.LoginRequest) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, validation.Min(989000000000), validation.Max(989999999999)),
		validation.Field(&request.Password, validation.Required, validation.Length(4, 40)),
	)
}

func ValidateChecker(request api.ValidateRequest) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Token, validation.Required, validation.Length(110, 290)),
	)
}

func OTPChecker(request api.OTPRequest) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, validation.Min(989000000000), validation.Max(989999999999)),
	)
}

func OTPVerifyChecker(request api.OTPVerifyRequest) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, validation.Min(989000000000), validation.Max(989999999999)),
	)
}
