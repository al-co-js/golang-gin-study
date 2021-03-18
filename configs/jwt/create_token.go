package jwt

import (
	"command/models"
	"command/types"
	"os"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

func CreateToken(user models.User, access bool) (string, error) {
	var secret string
	now := time.Now()
	if access {
		secret = os.Getenv("ACCESS_SECRET")
		now.Add(time.Minute * 15)
	} else {
		secret = os.Getenv("REFRESH_SECRET")
		now.Add(time.Hour * 24 * 7)
	}

	payload := types.Payload{
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now),
		},
		Access: access,
		User:   user,
	}
	hs := jwt.NewHS256([]byte(secret))
	token, err := jwt.Sign(payload, hs)
	if err != nil {
		return "", err
	}

	return string(token), nil
} 
