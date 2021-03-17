package routers

import (
	"command/api/auth"
	"command/api/users"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	authGroup := r.Group("/auth")
	authGroup.Use()
	{
		authGroup.POST("/login", auth.Login)
	}

	userGroup := r.Group("/user")
	userGroup.Use()
	{
		userGroup.POST("/", users.Create)
	}

	return r
}
