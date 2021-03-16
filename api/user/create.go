package user

import (
	"command/db"
	"command/models"
	"crypto/rand"
	"crypto/sha512"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
)

func Create(c *gin.Context) {
	client := db.Connect()
	var user models.User
	c.BindJSON(&user)

	salt := make([]byte, 64)
	rand.Read(salt)
	encrypt := pbkdf2.Key([]byte(user.Password), salt, 4096, 64, sha512.New)
	user.Password = string(encrypt) + "|" + string(salt)

	_, err := client.Database("command").Collection("users").InsertOne(c, user)
	if err != nil {
		c.JSON(500, bson.M{"message": "server error"})
		return
	}

	user.Password = ""

	c.JSON(200, bson.M{"message": "success", "data": user})
}
