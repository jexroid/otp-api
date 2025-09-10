package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/models"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/sirupsen/logrus"
)

func userinit(u *models.User) (*models.User, error) {
	u.UUID = pkg.Uuid()
	hashPass, err := utils.GenerateHash(u.Password)
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
	u.Password = hashPass
	return u, nil
}

// Signup godoc
// @Summary User registration
// @Description Register a new user with phone and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   request body models.User true "User Registration Data"
// @Success 200 {object} api.RegisterResponse
// @Failure 406 {object} api.RegisterResponse
// @Failure 500 {string} string "internal error"
// @Router /auth/signup [post]
func (db Database) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	if valid := utils.SignUpChecker(user); valid != nil {
		logrus.Error(valid)
		w.WriteHeader(http.StatusNotAcceptable)
		jsonResponse, _ := json.Marshal(&api.RegisterResponse{Ok: false})
		w.Write(jsonResponse)
		return
	}

	_, err := userinit(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}

	result := db.DB.Table("User").Create(&user)
	logrus.Error(result.Error)
	if nil != result.Error {
		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "uni_User_phone" (SQLSTATE 23505)` {
			jr, _ := json.Marshal(&api.RegisterResponse{
				Ok:        true,
				UserExist: true,
				Message:   "user with same phone number exist, use login",
			})
			w.Write(jr)
			return
		}
		logrus.Error(result.Error)
	}

	jsonResponse, _ := json.Marshal(&api.RegisterResponse{
		Ok:      true,
		Message: "User successfully registered! use OTP to verify",
	})
	w.Write(jsonResponse)
	return
}
