package token

import (
	jsonwebtoken "command/configs/jwt"
	"command/types"
	"net/http"
	"os"
	"strings"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Refresh(c *gin.Context) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if auth == token {
		c.JSON(http.StatusOK, bson.M{"message": "could not find bearer token in authorization header"})
		return
	}

	hs := jwt.NewHS256([]byte(refreshSecret))
	var payload types.Payload
	_, err := jwt.Verify([]byte(token), hs, &payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, bson.M{"message": "unauthorized token"})
		return
	}

	accessToken, err := jsonwebtoken.CreateToken(payload.User, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}
	refreshToken, err := jsonwebtoken.CreateToken(payload.User, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	c.JSON(http.StatusOK, bson.M{"message": "success", "data": bson.M{"access_token": accessToken, "refresh_token": refreshToken}})
}
