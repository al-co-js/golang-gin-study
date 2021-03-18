package auth

import (
	jsonwebtoken "command/configs/jwt"
	"command/db"
	"command/models"
	"command/types"
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
)

func Login(c *gin.Context) {
	var data types.LoginRequest
	c.BindJSON(&data)

	var user models.User
	err := db.Collection("users").FindOne(c, bson.M{"email": data.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	buf := strings.Split(user.Password, "|")
	salt, err := base64.StdEncoding.DecodeString(buf[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
	}
	encrypt := pbkdf2.Key([]byte(data.Password), salt, 80903, 512, sha512.New)
	if buf[0] != base64.StdEncoding.EncodeToString(encrypt) {
		c.JSON(http.StatusUnauthorized, bson.M{"message": "wrong password"})
		return
	}

	user.Password = ""

	accessToken, err := jsonwebtoken.CreateToken(user, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}
	refreshToken, err := jsonwebtoken.CreateToken(user, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	_, err = db.Collection("refresh_tokens").InsertOne(c, bson.M{"token": refreshToken, "uuid": data.UUID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	c.JSON(http.StatusOK, bson.M{"message": "success", "data": bson.M{"access_token": accessToken}})
}
