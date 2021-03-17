package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Authorization(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if auth == token {
		c.JSON(http.StatusOK, bson.M{"message": "could not find bearer token in authorization header"})
		return
	}

	// TODO verify token and decode token to user struct
}
