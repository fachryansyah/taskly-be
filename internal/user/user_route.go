package user

import (
	"tasklybe/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	users := app.Group("/user")
	users.Post("/register", HandleRegister)
	users.Post("/login", HandleLogin)
	users.Get("/me", middleware.Auth(), HandleMe)
}
