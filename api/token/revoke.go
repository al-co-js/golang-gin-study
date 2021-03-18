package token

import (
	"command/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Revoke(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")

	if auth == token {
		c.JSON(http.StatusPreconditionFailed, bson.M{"message": "could not find bearer token in authorization header"})
		return
	}

	_, err := db.Collection("refresh_tokens").DeleteOne(c, bson.M{"token": token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	c.JSON(http.StatusOK, bson.M{"message": "success"})
}
