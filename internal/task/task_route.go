package task

import "github.com/gofiber/fiber/v2"

func RegisterTaskRoutes(app *fiber.App) {
	tasks := app.Group("/tasks")
	tasks.Get("/", HandleGetTasks)
	tasks.Get("/:id", HandleGetTask)
	tasks.Post("/", HandleCreateTask)
	tasks.Put("/:id", HandleEditTask)
	tasks.Delete("/:id", HandleDeleteTask)
}
