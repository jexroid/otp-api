package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(uuid uint32, phone int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uuid":  uuid,
			"phone": phone,
			"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logrus.Error("JWT Failed to sign: ", err)
		logrus.Info("JWT signing error occurred")
		return ""
	}
	return tokenString
}

func ValidateToken(tokenString string) (
	bool,
	jwt.MapClaims,
) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return secretKey, nil
	})
	if err != nil {
		return false, nil
	}

	// validate the essential claims
	if !token.Valid {
		return false, nil
	}

	return true, token.Claims.(jwt.MapClaims)
}
