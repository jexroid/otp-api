package handler

import (
	"net/http"
	"time"
)

// Logout godoc
// @Summary User logout
// @Description Logout user by clearing authentication cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string "cookie deleted"
// @Router /auth/logout [get]
// @Security BearerAuth
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

	var deleteCookie = &http.Cookie{
		Name:     "identity",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, deleteCookie)

	w.Write([]byte("cookie deleted"))
	return
}
