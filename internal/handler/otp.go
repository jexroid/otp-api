package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg"
	"github.com/jexroid/gopi/pkg/models"
	"github.com/jexroid/gopi/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// OTP godoc
// @Summary Request OTP
// @Description Generate and send OTP to phone number for authentication
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   request body api.OTPRequest true "OTP Request"
// @Success 200 {object} api.OTPResponse
// @Failure 400 {object} api.OTPResponse
// @Failure 500 {object} api.OTPResponse
// @Router /auth/otp [post]
func (db Database) OTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	var form api.OTPRequest
	json.NewDecoder(r.Body).Decode(&form)

	if err := utils.OTPChecker(form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.OTPResponse{Ok: false, Message: "Invalid request"})
		return
	}

	otpCode, err := utils.GenerateOTP()
	if err != nil {
		logrus.Error("Failed to generate OTP:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.OTPResponse{Ok: false, Message: "Failed to generate OTP"})
		return
	}

	var user models.User
	result := db.DB.Where("phone = ?", form.Phone).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newUser := models.User{
				Phone:     form.Phone,
				Code:      otpCode,
				ExpiresAt: time.Now().Add(2 * time.Minute),
				Used:      false,
			}

			createResult := db.DB.Create(&newUser)
			if createResult.Error != nil {
				logrus.Error("Failed to create user with OTP:", createResult.Error)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(api.OTPResponse{Ok: false, Message: "Failed to save OTP"})
				return
			}
			user = newUser
		} else {
			logrus.Error("Database error:", result.Error)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.OTPResponse{Ok: false, Message: "Database error"})
			return
		}
	} else {
		updateResult := db.DB.Model(&user).Updates(map[string]interface{}{
			"code":       otpCode,
			"expires_at": time.Now().Add(2 * time.Minute),
			"used":       false,
		})

		if updateResult.Error != nil {
			logrus.Error("Failed to update OTP:", updateResult.Error)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.OTPResponse{Ok: false, Message: "Failed to save OTP"})
			return
		}
	}

	logrus.Infof("Sending OTP to %d", form.Phone)

	json.NewEncoder(w).Encode(api.OTPResponse{Ok: true, Message: "OTP sent successfully"})
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP code and authenticate user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   request body api.OTPVerifyRequest true "OTP Verification Request"
// @Success 200 {object} api.OTPVerifyResponse
// @Failure 400 {object} api.OTPVerifyResponse
// @Failure 401 {object} api.OTPVerifyResponse
// @Failure 500 {object} api.OTPVerifyResponse
// @Router /auth/verify-otp [post]
func (db Database) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req api.OTPVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "Invalid request"})
		return
	}

	if err := utils.OTPVerifyChecker(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "Invalid phone number"})
		return
	}

	if len(req.Code) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "Invalid OTP code"})
		return
	}

	var user models.User
	result := db.DB.Where("phone = ? AND code = ? AND used = ?", req.Phone, req.Code, false).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "Invalid OTP"})
			return
		}

		logrus.Error("Database error:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "Database error"})
		return
	}

	if utils.IsOTPExpired(user.ExpiresAt) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: false, Error: "OTP has expired"})
		return
	}

	db.DB.Model(&user).Update("used", true)

	jwt := pkg.CreateToken(user.UUID, user.Phone)

	cookie := &http.Cookie{
		Name:     "identity",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(api.OTPVerifyResponse{Ok: true, Token: jwt})
}

func (db Database) createUserFromPhone(phone int) (*models.User, error) {
	user := models.User{
		Phone: phone,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
