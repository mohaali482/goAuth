package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohaali482/goAuth/auth"
)

func AuthMiddleware(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token", "")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "request does not contain an access token"})
		}

		_, err := s.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "access token is not valid"})
		}

		return c.Next()

	}

}
