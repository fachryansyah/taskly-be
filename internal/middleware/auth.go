package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"tasklybe/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type authClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "invalid authorization header",
			})
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "invalid authorization header",
			})
		}

		secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
		if secret == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "jwt secret missing",
			})
		}

		claims := &authClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
					Success: false,
					Message: "token expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "invalid token",
			})
		}
		if token == nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "invalid token",
			})
		}
		if strings.TrimSpace(claims.Subject) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[any]{
				Success: false,
				Message: "invalid token",
			})
		}

		c.Locals("userId", claims.Subject)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

