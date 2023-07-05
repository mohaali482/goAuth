package fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mohaali482/goAuth/auth"
)

func App(s auth.UserService) *fiber.App {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		TimeZone:   "Africa/Addis_Ababa",
		Format:     "[${time}] ${ip} ${latency} ${status} - ${method} ${path}\n",
	}))

	return app
}
