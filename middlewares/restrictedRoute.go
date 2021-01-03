package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// ProtectedRoute -> is the middleware function which will then pass in each route request as the middleman to indicate whether have or haven't has the token or the right to visit the route
func ProtectedRoute() fiber.Handler{
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx,err error) error{
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized to visit this route",
			})
		},
		SigningKey: []byte("secureSecretText"),
	})
}