package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Validate godoc
// @Summary Validate user token
// @Description Validates a user's authentication token and returns user details if valid
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   validateRequest body api.ValidateRequest true "Validate Request"
// @Success 200 {object} api.ValidateResponse "Token is valid and user details returned"
// @Failure 400 {object} api.ValidateResponse "Invalid request body"
// @Failure 405 {object} api.ValidateResponse "Method not allowed"
// @Failure 406 {object} api.ValidateResponse "Token validation failed or invalid data"
// @Security ApiKeyAuth
// @Router /validate [post]
func (db Database) Validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	var user api.ValidateRequest
	json.NewDecoder(r.Body).Decode(&user)
	if valid := utils.ValidateChecker(user); valid != nil {
		logrus.Error(valid)
		w.WriteHeader(http.StatusNotAcceptable)
		jsonResponse, _ := json.Marshal(&api.RegisterResponse{Ok: false})
		w.Write(jsonResponse)
		return
	}

	valid, payload := pkg.ValidateToken(user.Token)
	if valid {
		var userPhone = strconv.Itoa(int(payload["phone"].(float64)))
		var userData api.Userdetail
		db.DB.Table("User").First(&userData, "phone = ?", userPhone)

		jsonResponse, _ := json.Marshal(&api.ValidateResponse{
			Ok: valid,
			User: api.Userdetail{
				Firstname: userData.Firstname,
				Lastname:  userData.Lastname,
				Phone:     userData.Phone,
			},
		})

		w.Write(jsonResponse)
		return
	} else {
		jsonResponse, _ := json.Marshal(&api.ValidateResponse{
			Ok: false,
		})
		w.WriteHeader(http.StatusNotAcceptable)

		w.Write(jsonResponse)
	}
}
