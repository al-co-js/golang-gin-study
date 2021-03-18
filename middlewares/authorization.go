package middlewares

import (
	"command/types"
	"net/http"
	"os"
	"strings"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Authorization(c *gin.Context) {
	accessSecret := os.Getenv("ACCESS_SECRET")
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if auth == token {
		c.JSON(http.StatusOK, bson.M{"message": "could not find bearer token in authorization header"})
		return
	}

	hs := jwt.NewHS256([]byte(accessSecret))
	var payload types.Payload
	_, err := jwt.Verify([]byte(token), hs, &payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, bson.M{"message": "unauthorized token"})
		return
	}

	c.Next()
}
