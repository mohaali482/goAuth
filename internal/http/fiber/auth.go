package fiber

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mohaali482/goAuth/auth"
	"github.com/mohaali482/goAuth/internal/http/fiber/errors"
	"github.com/mohaali482/goAuth/internal/http/fiber/middlewares"
)

func App(s auth.UserService) *fiber.App {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		TimeZone:   "Africa/Addis_Ababa",
		Format:     "[${time}] ${ip} ${latency} ${status} - ${method} ${path}\n",
	}))
	accountsGroup := app.Group("/accounts")
	{
		accountsGroup.Post("/login", Login(s))
		accountsGroup.Post("/signup", Signup(s))
		accountsGroup.Delete("/logout", Logout(s))
		accountsGroup.Get("/refresh", RefreshToken(s))
	}

	usersGroup := app.Group("/users").Use(middlewares.AuthMiddleware(s))
	{
		usersGroup.Get("", GetAll(s))
		usersGroup.Get("/:id", GetByID(s))
		usersGroup.Delete("/:id", Delete(s))
		usersGroup.Patch("/:id", Update(s))
		usersGroup.Post("", Create(s))
	}

	return app
}

func Login(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Login started")
		var userLogin auth.UserLogin
		err := c.BodyParser(&userLogin)
		if err != nil {
			log.Default().Println("Error binding json while trying to login. Error: ", err)
			return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(err)
		}
		user, err := s.Login(userLogin.Username, userLogin.Password)
		if err != nil {
			log.Default().Println("Error logging in. Error: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": auth.ErrWrongCredentials.Error()})
		}
		tokens, err := s.GenerateJWT(user)
		if err != nil {
			log.Default().Println("Error generating tokens. Error: ", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: tokens["refresh"], Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true})
		c.Cookie(&fiber.Cookie{Name: "access_token", Value: tokens["access"], Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true})

		log.Default().Println("Login successful")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	}
}

func Logout(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Logout started")
		c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: "", Expires: time.Now().Add(-time.Hour), HTTPOnly: true})
		c.Cookie(&fiber.Cookie{Name: "access_token", Value: "", Expires: time.Now().Add(-time.Hour), HTTPOnly: true})
		log.Default().Println("Logout successful")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	}
}

func Signup(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Signup started")
		var userForm auth.UserForm
		err := c.BodyParser(&userForm)
		if err != nil {
			log.Default().Println("Error binding json while trying to create user. Error: ", err)
			return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(err)
		}
		err = userForm.Validate()
		if err != nil {
			log.Default().Println("Error validating user while trying to create user. Error: ", err)
			return errors.ReturnErrorResponse(err, c)
		}

		user, err := s.Create(userForm.ToUserEntity())
		if err != nil {
			log.Default().Println("Error creating user while trying to create user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		log.Default().Println("User created successfully")
		return c.Status(fiber.StatusCreated).JSON(user)
	}

}

func RefreshToken(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token", "")
		if refreshToken != "" {
			log.Default().Println("Error getting refresh_token while trying to refresh token.")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token not found in cookie"})
		}

		tokens, err := s.RefreshToken(refreshToken)
		if err != nil {
			log.Default().Println("Error refreshing token while trying to refresh token. Error: ", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
		}
		c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: tokens["refresh"], Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true})
		c.Cookie(&fiber.Cookie{Name: "access_token", Value: tokens["access"], Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true})
		log.Default().Println("Token refreshed successfully")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	}
}

func Create(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Creating user started")
		var user auth.User
		err := c.BodyParser(&user)
		if err != nil {
			log.Default().Println("Error binding json while trying to create user. Error: ", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = user.Validate()
		if err != nil {
			log.Default().Println("Error validating user while trying to create user. Error: ", err)
			return errors.ReturnErrorResponse(err, c)
		}

		user, err = s.Create(user)
		if err != nil {
			log.Default().Println("Error creating user while trying to create user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		log.Default().Println("User created successfully")
		return c.Status(fiber.StatusCreated).JSON(user)
	}

}

func Delete(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Deleting user started")
		id, err := strconv.Atoi(c.Params("id", ""))
		if err != nil {
			log.Default().Println("Error converting id while trying to delete user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is not a valid id"})
		}
		_, err = s.GetByID(id)
		if err != nil {
			log.Default().Println("Error getting user by id while trying to delete user. Error: ", err)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		if err := s.Delete(id); err != nil {
			log.Default().Println("Error deleting user while trying to delete user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		log.Default().Println("User deleted successfully")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	}
}

func GetAll(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Getting all users started")
		users, err := s.GetAll()
		if err != nil {
			log.Default().Println("Error getting all users. Error: ", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		log.Default().Println("Users fetched successfully")
		return c.Status(http.StatusOK).JSON(users)
	}
}

func GetByID(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Getting user by id started")
		id, err := strconv.Atoi(c.Params("id", ""))
		if err != nil {
			log.Default().Println("Error converting id. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is not a valid id"})
		}
		user, err := s.GetByID(id)
		if err != nil {
			log.Default().Println("Error getting user by id. Error: ", err)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		log.Default().Println("User fetched successfully")
		return c.Status(fiber.StatusOK).JSON(user)
	}
}

func Update(s auth.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Println("Updating user started")
		var userForm auth.UserForm
		id, err := strconv.Atoi(c.Params("id", ""))
		if err != nil {
			log.Default().Println("Error converting id while trying to update user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is not a valid id"})
		}
		_, err = s.GetByID(id)
		if err != nil {
			log.Default().Println("Error getting user by id while trying to update user. Error: ", err)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		err = c.BodyParser(&userForm)
		if err != nil {
			log.Default().Println("Error binding json while trying to update user. Error: ", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		if userForm.Phone != "" {
			err = userForm.ValidatePhone()
			if err != nil {
				log.Default().Println("Error validating phone while trying to update user. Error: ", err)
				return errors.ReturnErrorResponse(err, c)
			}
		}

		user, err := s.Update(id, userForm.ToUserEntity())
		if err != nil {
			log.Default().Println("Error updating user while trying to update user. Error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		log.Default().Println("User updated successfully")
		return c.Status(fiber.StatusOK).JSON(user)
	}
}
