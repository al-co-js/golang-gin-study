package jwt

import (
	"command/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(user models.User, access bool) (*string, error) {
	var secret string
	var expires time.Duration
	if access {
		secret = os.Getenv("ACCESS_SECRET")
		expires = time.Minute * 15
	} else {
		secret = os.Getenv("REFRESH_SECRET")
		expires = time.Hour * 24 * 7
	}

	claims := jwt.MapClaims{}
	claims["access"] = access
	claims["user"] = user
	claims["exp"] = time.Now().Add(expires).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &token, nil
}
