package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateRequest[T any](c *fiber.Ctx, model *T) error {
	if err := c.BodyParser(model); err != nil {
		return ErrParseJSON
	}
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}
