package route

import (
	"github.com/dalpengida/portfolio-api-go-mysql/cmd/fiber/handler"
	"github.com/dalpengida/portfolio-api-go-mysql/cmd/fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

const default_group_route = "/api/v1"

func SetupV1Route(app *fiber.App) {
	v1Group := app.Group(default_group_route, func(c *fiber.Ctx) error {
		return c.Next()
	})

	// account api
	v1Account := handler.Account{}
	v1Group.Post("/signup", v1Account.Signup)
	v1Group.Post("/login", v1Account.Login)

	v1Group.Post("/logout", middleware.RequireJWTWithClaims(), v1Account.Logout)
}
