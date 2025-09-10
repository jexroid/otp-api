package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/models"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Signin godoc
// @Summary User signin
// @Description Authenticate user with phone and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   request body api.LoginRequest true "Login Request"
// @Success 200 {object} api.LoginResponse
// @Failure 406 {object} api.LoginResponse
// @Router /auth/signin [post]
func (db Database) Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	var form api.LoginRequest
	json.NewDecoder(r.Body).Decode(&form)

	if valid := utils.SignInChecker(form); valid != nil {
		logrus.Info(valid)
		w.WriteHeader(http.StatusNotAcceptable)
		jsonResponse, _ := json.Marshal(&api.LoginResponse{Ok: false})
		w.Write(jsonResponse)
		return
	}

	var user models.User
	result := db.DB.Table("User").First(&user, "phone = ?", strconv.Itoa(form.Phone))
	if nil != result.Error {
		if result.Error.Error() == "record not found" {
			jr, _ := json.Marshal(&api.RegisterResponse{
				Ok:        true,
				UserExist: false,
				Message:   "User does not exist. sign up",
			})
			w.Write(jr)
			return
		}
	}

	isPassValid, _ := utils.ComparePasswordAndHash(form.Password, user.Password)
	logrus.Info(isPassValid)

	if isPassValid {
		jwt := pkg.CreateToken(user.UUID, user.Phone)
		cookie := &http.Cookie{
			Name:     "identity",
			Value:    jwt,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)

		jsonResponse, _ := json.Marshal(&api.LoginResponse{Ok: true, Valid: true, Token: jwt})
		w.Write(jsonResponse)
		return
	}

	response := &api.LoginResponse{}
	jsonResponse, _ := json.Marshal(response)

	w.Write(jsonResponse)
	return
}
