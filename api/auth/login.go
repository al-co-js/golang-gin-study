package auth

import (
	"command/configs/jwt"
	"command/db"
	"command/models"
	"crypto/sha512"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
)

func Login(c *gin.Context) {
	var data models.User
	c.BindJSON(&data)

	var user models.User
	err := db.Collection("users").FindOne(c, bson.M{"email": data.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	buf := strings.Split(user.Password, "|")
	encrypt := pbkdf2.Key([]byte(data.Password), []byte(buf[1]), 4096, 64, sha512.New)

	if buf[0] != string(encrypt) {
		c.JSON(http.StatusUnauthorized, bson.M{"message": "wrong password"})
		return
	}

	accessToken, err := jwt.CreateToken(user, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}
	refreshToken, err := jwt.CreateToken(user, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
	}

	c.JSON(http.StatusOK, bson.M{"message": "success", "data": bson.M{"access_token": *accessToken, "refresh_token": *refreshToken}})
}
