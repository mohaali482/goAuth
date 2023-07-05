package errors

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ReturnErrorResponse(err error, c *fiber.Ctx) error {
	if err != nil {
		var validationErrors []ErrorResponse
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ErrorResponse{
				Field:   err.Field(),
				Message: SetValidationResult(err.Tag()),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
}

func SetValidationResult(tag string) string {
	switch tag {
	case "required":
		tag = "This field is required"
	case "email":
		tag = "This field is not a valid email"
	case "min":
		tag = "This field is too short"
	case "max":
		tag = "This field is too long"
	case "eqfield":
		tag = "This field is not equal to the other field"
	case "e164":
		tag = "Invalid phone number format. Please use +999999999999 format"
	}

	return tag
}
