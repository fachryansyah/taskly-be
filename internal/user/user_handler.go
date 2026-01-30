package user

import (
	"errors"
	"tasklybe/internal/dto"
	"tasklybe/internal/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HandleRegister godoc
// @Summary Register user
// @Description Register a new user
// @Tags users
// @Produce json
// @Param request body RegisterUserRequest true "Register User Request"
// @Success 200 {object} UserResponse
// @Router /users/register [post]
func HandleRegister(c *fiber.Ctx) error {
	var req RegisterUserRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	u, err := RegisterUser(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[UserResponse]{
		Data:    ToUserResponse(u),
		Success: true,
		Message: "Success! User registered.",
	})
}

// HandleLogin godoc
// @Summary Login user
// @Description Login a user
// @Tags users
// @Produce json
// @Param request body LoginUserRequest true "Login User Request"
// @Success 200 {object} UserResponse
// @Router /users/login [post]
func HandleLogin(c *fiber.Ctx) error {
	var req LoginUserRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	u, apiKey, err := LoginUser(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[UserResponse]{
		Data: &UserResponse{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
			Name:      u.Name,
			ApiKey:    apiKey,
		},
		Success: true,
		Message: "Success! Login success.",
	})
}

func HandleMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
			Success: false,
			Message: "unauthorized",
		})
	}

	u, err := GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "user not found",
			})
		}
		if errors.Is(err, ErrValidation) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[any]{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.ResponseWrapper[UserResponse]{
		Data:    ToUserResponse(u),
		Success: true,
		Message: "Success! user found.",
	})
}
