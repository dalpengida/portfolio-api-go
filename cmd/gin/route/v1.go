package route

import (
	"github.com/dalpengida/portfolio-api-go-mysql/cmd/gin/handler"
	"github.com/dalpengida/portfolio-api-go-mysql/cmd/gin/middleware"
	"github.com/gin-gonic/gin"
)

const default_group_route = "/api/v1"

func SetupV1Route(app *gin.Engine) {

	v1Group := app.Group(default_group_route, func(c *gin.Context) {
		c.Next()
	})

	v1Account := handler.Account{}
	v1Group.POST("/signup", v1Account.Signup)
	v1Group.POST("/login", v1Account.Login)

	v1Group.POST("/logout", middleware.RequireJWT(), v1Account.Logout)
}
