package users

import (
	"command/db"
	"command/models"
	"crypto/rand"
	"crypto/sha512"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
)

func Create(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	salt := make([]byte, 64)
	rand.Read(salt)
	encrypt := pbkdf2.Key([]byte(user.Password), salt, 4096, 64, sha512.New)
	user.Password = string(encrypt) + "|" + string(salt)

	_, err := db.Collection("users").InsertOne(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, bson.M{"message": "success", "data": user})
}
