package auth

import (
	"command/db"
	"command/models"
	"crypto/sha512"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
)

func Login(c *gin.Context) {
	client := db.Connect()
	var data models.User
	c.BindJSON(&data)

	var user models.User
	err := client.Database("command").Collection("users").FindOne(c, bson.M{"email": data.Email}).Decode(&user)
	if err != nil {
		c.JSON(500, bson.M{"message": "server error"})
		return
	}

	buf := strings.Split(user.Password, "|")
	encrypt := pbkdf2.Key([]byte(data.Password), []byte(buf[1]), 4096, 64, sha512.New)

	if buf[0] != string(encrypt) {
		c.JSON(401, bson.M{"message": "wrong password"})
		return
	}

	user.Password = ""

	c.JSON(200, bson.M{"message": "success", "data": user})
}
