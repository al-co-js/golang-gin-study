package token

import (
	jsonwebtoken "command/configs/jwt"
	"command/db"
	"command/models"
	"command/types"
	"net/http"
	"os"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Refresh(c *gin.Context) {
	var data types.TokenRefreshRequest
	c.BindJSON(&data)
	refreshSecret := os.Getenv("REFRESH_SECRET")

	var refreshTokenCheck models.RefreshToken
	err := db.Collection("refresh_tokens").FindOne(c, bson.M{"uuid": data.UUID}).Decode(&refreshTokenCheck)
	if err != nil {
		c.JSON(http.StatusUnauthorized, bson.M{"message": "unauthorized token"})
		return
	}

	_, err = db.Collection("refresh_tokens").DeleteOne(c, bson.M{"uuid": data.UUID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	hs := jwt.NewHS256([]byte(refreshSecret))
	var payload types.Payload
	_, err = jwt.Verify([]byte(refreshTokenCheck.Token), hs, &payload)
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

	_, err = db.Collection("refresh_tokens").InsertOne(c, bson.M{"token": refreshToken, "uuid": data.UUID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"message": "server error"})
		return
	}

	c.JSON(http.StatusOK, bson.M{"message": "success", "data": bson.M{"access_token": accessToken, "refresh_token": refreshToken}})
}
