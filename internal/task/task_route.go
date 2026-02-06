package task

import (
	"tasklybe/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterTaskRoutes(app *fiber.App) {
	tasks := app.Group("/task")
	tasks.Get("/", middleware.Auth(), HandleGetTasks)
	tasks.Get("/:id", middleware.Auth(), HandleGetTask)
	tasks.Post("/", middleware.Auth(), HandleCreateTask)
	tasks.Put("/:id", middleware.Auth(), HandleEditTask)
	tasks.Delete("/:id", middleware.Auth(), HandleDeleteTask)
}
