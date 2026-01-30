package validation

import (
	"fmt"
	"reflect"
	"strings"
	"tasklybe/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]

		if name == "-" {
			return ""
		}
		return name
	})
}

func BindAndValidate[T any](c *fiber.Ctx, dst *T) error {
	if err := c.BodyParser(dst); err != nil {
		return err
	}
	return Validate.Struct(dst)
}

func FormatValidationError(err error) *[]dto.ResponseError {
	errors := []dto.ResponseError{}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			// errors[e.Field()] = e.Tag()
			errors = append(errors, dto.ResponseError{
				Field:   e.Field(),
				Value:   fmt.Sprint(e.Value()),
				Tag:     e.Tag(),
				Message: fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()),
				Target:  "/task",
			})
		}
	}

	return &errors
}
